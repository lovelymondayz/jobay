#!/usr/bin/env python3
"""
Hermes workspace backup script.
Zips /root/hermes and pushes to GitHub repo BackupHermes via the Contents API.
Keeps only the last 30 days of backups.
"""

import os, sys, json, zipfile, subprocess, ssl, http.client, base64, time
from datetime import datetime, timedelta, timezone

# Config
GITHUB_USER = "lovelymondayz"
REPO = "BackupHermes"
SOURCE_DIR = "/root/hermes"
BACKUP_DIR = "/tmp/backup-workspace"
REMOTE_FOLDER = "backupHermesDaily"
# Note: All projects now live under /root/hermes/ (pico, qlio, say-less, internet-storage moved here)
RETENTION_DAYS = 30
GIT_NAME = "Hermes Bot"
GIT_EMAIL = "hermes-bot@hermes.arjism.com"
TOKEN_FILE = "/root/.github/pat"

# Bangkok time (UTC+7)
BKK = timezone(timedelta(hours=7))
DATE = datetime.now(BKK).strftime("%Y-%m-%d")
ARCHIVE_NAME = f"backup-{DATE}.zip"
REMOTE_PATH = f"{REMOTE_FOLDER}/{ARCHIVE_NAME}"
LOCAL_ZIP = os.path.join(BACKUP_DIR, ARCHIVE_NAME)

log = lambda msg: print(f"[{datetime.now(BKK).strftime('%Y-%m-%d %H:%M:%S')}] {msg}", flush=True)

def gh_api(method, path, body=None, headers=None):
    """Make a GitHub API call using http.client (avoids shell escaping issues)."""
    ctx = ssl.create_default_context()
    conn = http.client.HTTPSConnection("api.github.com", context=ctx)
    hdrs = {
        "Authorization": f"Bearer {TOKEN}",
        "Accept": "application/vnd.github+json",
        "X-GitHub-Api-Version": "2022-11-28",
        "User-Agent": "Hermes-Backup/1.0"
    }
    if headers:
        hdrs.update(headers)
    if body and isinstance(body, str):
        body = body.encode()
    conn.request(method, path, body=body, headers=hdrs)
    resp = conn.getresponse()
    data = resp.read()
    conn.close()
    return resp.status, json.loads(data) if data else {}

def get_file_sha(path):
    """Get SHA of existing file in repo, or None."""
    status, data = gh_api("GET", f"/repos/{GITHUB_USER}/{REPO}/contents/{path}")
    if status == 200:
        return data.get("sha")
    return None

def upload_file(path, content_b64, message):
    """Create or update a file via GitHub Contents API."""
    sha = get_file_sha(path)
    body = {
        "message": message,
        "content": content_b64,
        "author": {"name": GIT_NAME, "email": GIT_EMAIL},
        "committer": {"name": GIT_NAME, "email": GIT_EMAIL}
    }
    if sha:
        body["sha"] = sha
    body_json = json.dumps(body)
    status, data = gh_api("PUT", f"/repos/{GITHUB_USER}/{REPO}/contents/{path}", body_json)
    if status in (200, 201):
        return True
    log(f"Upload failed {status}: {data.get('message', data)}")
    return False

def list_backups():
    """List backup files in the backup folder."""
    status, data = gh_api("GET", f"/repos/{GITHUB_USER}/{REPO}/contents/{REMOTE_FOLDER}")
    if status != 200:
        log(f"List failed: {data.get('message', data)}")
        return []
    return [f["name"] for f in data if f["name"].startswith("backup-") and f["name"].endswith(".zip")]

def delete_backup(filename):
    """Delete a file from the repo."""
    sha = get_file_sha(filename)
    if not sha:
        return
    body = json.dumps({
        "message": f"prune: remove {filename} (>30 days old)",
        "sha": sha,
        "author": {"name": GIT_NAME, "email": GIT_EMAIL}
    })
    status, data = gh_api("DELETE", f"/repos/{GITHUB_USER}/{REPO}/contents/{filename}", body)
    if status == 200:
        log(f"  Pruned: {filename}")
    else:
        log(f"  Failed to prune {filename}: {data.get('message', status)}")

def create_zip():
    """Create zip archive of source directory."""
    os.makedirs(BACKUP_DIR, exist_ok=True)
    if os.path.exists(LOCAL_ZIP):
        os.remove(LOCAL_ZIP)

    # Skip build artifacts and binaries
    skip_dirs = {'.git', 'node_modules', '__pycache__', '.venv', 'backupHermesDaily'}
    skip_files = {'.pyc', '.so', '.so.*'}  # Compiled binaries, shared libs

    with zipfile.ZipFile(LOCAL_ZIP, 'w', zipfile.ZIP_DEFLATED) as zf:
        for root, dirs, files in os.walk(SOURCE_DIR):
            dirs[:] = [d for d in dirs if d not in skip_dirs]
            for fname in files:
                # Skip compiled binaries and large build artifacts
                if any(fname.endswith(ext) for ext in skip_files):
                    continue
                # Skip dist directories (frontend builds - reproducible)
                if '/dist/' in root or root.endswith('/dist'):
                    continue
                # Skip Go server binaries
                if fname == 'server' or fname.endswith('.exe'):
                    continue
                fpath = os.path.join(root, fname)
                arcname = os.path.relpath(fpath, SOURCE_DIR)
                zf.write(fpath, arcname)

    size = os.path.getsize(LOCAL_ZIP)
    return size

def main():
    # Load token
    if not os.path.exists(TOKEN_FILE):
        log(f"ERROR: token file not found at {TOKEN_FILE}")
        sys.exit(1)
    with open(TOKEN_FILE, 'r') as f:
        global TOKEN
        TOKEN = f.read().strip()
    if not TOKEN:
        log("ERROR: token file is empty")
        sys.exit(1)

    log(f"Token loaded: {len(TOKEN)} chars")

    # Verify auth
    status, me = gh_api("GET", "/user")
    if status != 200:
        log(f"Auth failed {status}: {me.get('message', me)}")
        sys.exit(1)
    log(f"Authenticated as: {me.get('login')}")

    # Check source
    if not os.path.isdir(SOURCE_DIR) or not os.listdir(SOURCE_DIR):
        log("Source empty or missing — skipping")
        sys.exit(0)

    # Check if today's backup already exists
    existing = list_backups()
    if ARCHIVE_NAME in existing:
        log(f"{ARCHIVE_NAME} already exists — checking if content changed")
        # We'll upload anyway and let git detect no diff (sha match = no new commit needed)
        # Actually for Contents API always creates a commit, so skip
        # But user wants "no commit if nothing changed"
        # Compare zip content: create new zip, compare base64 with existing
        old_sha = get_file_sha(ARCHIVE_NAME)
        new_size = create_zip()
        with open(LOCAL_ZIP, 'rb') as f:
            new_b64 = base64.b64encode(f.read()).decode()
        # If old file exists and we can embed both, compare
        # For simplicity, re-upload (GitHub will create a new commit anyway)
        # To properly detect changes we'd need to download old, compare bytes
        # Let's do that
        status, old_data = gh_api("GET", f"/repos/{GITHUB_USER}/{REPO}/contents/{ARCHIVE_NAME}")
        if status == 200:
            old_b64 = old_data.get("content", "").replace("\n", "")
            if old_b64.strip() == new_b64.strip():
                log("No content changes — skipping commit")
                sys.exit(0)

    # Create zip
    log(f"Creating: {ARCHIVE_NAME}")
    size = create_zip()
    size_human = f"{size / 1024:.0f} KB" if size > 1024 else f"{size} B"
    log(f"Archive size: {size_human}")

    # Encode for upload
    with open(LOCAL_ZIP, 'rb') as f:
        b64 = base64.b64encode(f.read()).decode()

    # Upload to backupHermesDaily/ folder
    msg = f"backup: {DATE} - {ARCHIVE_NAME} ({size_human})"
    log(f"Uploading: {ARCHIVE_NAME} to {REMOTE_PATH}")
    if upload_file(REMOTE_PATH, b64, msg):
        log(f"Uploaded: {ARCHIVE_NAME}")
    else:
        log("Upload failed!")
        sys.exit(1)

    # Prune old backups
    log(f"Pruning backups older than {RETENTION_DAYS} days...")
    cutoff = datetime.now(BKK) - timedelta(days=RETENTION_DAYS)
    pruned = 0
    for fname in list_backups():
        # Parse date from filename
        try:
            file_date = datetime.strptime(fname.replace("backup-", "").replace(".zip", ""), "%Y-%m-%d")
            file_date = file_date.replace(tzinfo=BKK)
            if file_date < cutoff:
                delete_backup(fname)
                pruned += 1
        except ValueError:
            continue
    log(f"Pruned {pruned} old backup(s)")

    # Cleanup
    if os.path.exists(LOCAL_ZIP):
        os.remove(LOCAL_ZIP)
    if os.path.isdir(BACKUP_DIR):
        import shutil
        shutil.rmtree(BACKUP_DIR, ignore_errors=True)

    log(f"Done: {ARCHIVE_NAME} ({size_human})")

if __name__ == "__main__":
    main()