#!/usr/bin/env python3
"""
GitHub Deploy Webhook Receiver
---------------------------------
Listens on 127.0.0.1:9000 for GitHub push webhooks.
Validates HMAC (X-Hub-Signature-256) using DEPLOY_WEBHOOK_SECRET.
On a verified push to refs/heads/main, runs that project's scripts/update.sh.

Exposed publicly via the existing cloudflared tunnel (e.g. deploy.client.arjism.com -> 127.0.0.1:9000).
No public IP exposure; origin stays hidden behind Cloudflare.

Stdlib only. Runs as a systemd service.
"""
import json
import hmac
import hashlib
import os
import subprocess
import threading
import logging
import sys
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer

CONFIG_PATH = "/root/hermes/scripts/deploy-config.json"
LOG_PATH = "/root/hermes/scripts/deploy.log"
SECRET = os.environ.get("DEPLOY_WEBHOOK_SECRET")

logging.basicConfig(
    filename=LOG_PATH,
    level=logging.INFO,
    format="%(asctime)s %(levelname)s %(message)s",
)
log = logging.getLogger("deploy")

deploy_lock = threading.Lock()


def load_config():
    with open(CONFIG_PATH) as f:
        return json.load(f)


def run_update(repo, proj):
    """Pull + rebuild + redeploy via the project's own update.sh."""
    with deploy_lock:
        dir_ = proj["dir"]
        script = proj.get("update_script", "scripts/update.sh")
        log.info("=== Deploy START: %s (dir=%s) ===", repo, dir_)
        try:
            r = subprocess.run(
                ["bash", script],
                cwd=dir_,
                capture_output=True,
                text=True,
                timeout=900,
            )
            out = (r.stdout or "")[-2500:]
            err = (r.stderr or "")[-2500:]
            log.info(
                "Deploy END %s rc=%s\n--- stdout ---\n%s\n--- stderr ---\n%s",
                repo, r.returncode, out, err,
            )
        except subprocess.TimeoutExpired:
            log.error("Deploy TIMEOUT for %s", repo)
        except Exception as e:
            log.error("Deploy ERROR %s: %s", repo, e)


class Handler(BaseHTTPRequestHandler):
    def _send(self, code, msg):
        self.send_response(code)
        self.send_header("Content-Type", "application/json")
        self.end_headers()
        self.wfile.write(json.dumps({"status": msg}).encode())

    def do_GET(self):
        if self.path == "/health":
            self._send(200, "ok")
        else:
            self._send(404, "not found")

    def do_POST(self):
        if self.path != "/webhook":
            self._send(404, "not found")
            return

        event = self.headers.get("X-GitHub-Event", "")
        length = int(self.headers.get("Content-Length", 0))
        body = self.rfile.read(length)

        # GitHub sends a ping on webhook creation — answer, don't deploy.
        if event == "ping":
            self._send(200, "pong")
            return

        if event != "push":
            self._send(200, "ignored event: %s" % event)
            return

        if not SECRET:
            log.error("DEPLOY_WEBHOOK_SECRET not set — refusing")
            self._send(500, "server misconfigured")
            return

        sig = self.headers.get("X-Hub-Signature-256", "")
        expected = "sha256=" + hmac.new(SECRET.encode(), body, hashlib.sha256).hexdigest()
        if not hmac.compare_digest(sig, expected):
            log.warning("Invalid signature from %s", self.client_address)
            self._send(403, "invalid signature")
            return

        try:
            payload = json.loads(body)
        except Exception:
            self._send(400, "bad json")
            return

        repo = payload.get("repository", {}).get("full_name")
        ref = payload.get("ref")
        config = load_config()
        proj = config.get(repo)

        if not proj:
            log.info("No deploy config for repo %s — skipping", repo)
            self._send(200, "no action for repo")
            return
        if ref != "refs/heads/main":
            log.info("Ignoring ref %s for %s", ref, repo)
            self._send(200, "ignored ref")
            return

        # Respond fast (GitHub times out at 10s); deploy in background.
        threading.Thread(target=run_update, args=(repo, proj), daemon=True).start()
        self._send(202, "deploy triggered")

    def log_message(self, *args):
        pass  # silence default request logging


if __name__ == "__main__":
    if not SECRET:
        print("DEPLOY_WEBHOOK_SECRET env required (see .deploy_env)", file=sys.stderr)
        raise SystemExit(1)
    server = ThreadingHTTPServer(("127.0.0.1", 9000), Handler)
    log.info("Deploy webhook listening on 127.0.0.1:9000")
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        pass
