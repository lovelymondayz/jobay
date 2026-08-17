# 2026-07-16 — Push-to-Deploy Pipeline (interrupted) + Hermes venv repair

**Sessions today:** 3
- `20260716_072629_e5a47291` (discord) — "Prioritizing Business Goals Discussion" — broken `hermes` CLI diagnosis
- `20260716_090117_a1613fa3` (discord) — "Model Delay Interrupts New Build" — push-to-deploy pipeline build (interrupted)
- `cron_4e68d6ab9ae2_20260716_190033` (cron) — daily Hermes workspace backup (routine, successful)

---

## What was built

### 1. GitHub Push-to-Deploy pipeline (PARTIAL — interrupted mid-build)
Picking up an earlier interrupted task: **git push → GitHub webhook → VPS receiver → pull + docker compose rebuild/redeploy** per project.

Files written to disk today (mtime 09:11):
- `/root/hermes/scripts/deploy-webhook.py` (4.7 KB) — stdlib-only HTTP receiver.
  - `ThreadingHTTPServer` on `127.0.0.1:9000`, only `/webhook` (POST) and `/health` (GET).
  - HMAC validation of `X-Hub-Signature-256` using `DEPLOY_WEBHOOK_SECRET`.
  - Responds `202` fast, runs `scripts/update.sh` in a background thread (900s timeout, deploy lock).
  - Handles GitHub `ping` event (answers, no deploy); ignores non-`main` refs and unknown repos.
  - Designed to sit behind the existing cloudflared tunnel (`deploy.client.arjism.com → 127.0.0.1:9000`), no public IP.
- `/root/hermes/scripts/deploy-config.json` — maps 4 repos to dirs + `scripts/update.sh`:
  - `lovelymondayz/members` → `/root/hermes/members`
  - `lovelymondayz/internet-storage` → `/root/internet-storage`
  - `lovelymondayz/wedding-invitation` → `/root/hermes/wedding-invitation`
  - `lovelymondayz/digital-yearbook` → `/root/hermes/digital-yearbook`
- `/root/hermes/scripts/.deploy_env` (chmod 600) — `DEPLOY_WEBHOOK_SECRET` = `openssl rand -hex 32`.

**NOT completed (build was interrupted at the `execute_code` step that was supposed to generate per-project `update.sh`):**
- `members/scripts/update.sh` — **missing** (config references it → would 404 on deploy).
- `internet-storage/scripts/update.sh` — **missing**.
- `mentengdutch` has an old (Jun 25) `update.sh` but is **not** in `deploy-config.json` (not wired).
- No `systemd` service unit for `deploy-webhook.py`.
- No GitHub webhook registered; no cloudflared route `deploy.client.arjism.com`.
- Receiver never started/tested. **Pipeline is NOT operational.**

> Note: the 3 `update.sh` files currently on disk (wedding-invitation Jun 25, mentengdutch Jun 25, digital-yearbook Jun 28) are PRE-EXISTING from June — they were NOT regenerated today. The generate step was the interruption point.

### 2. Hermes CLI venv repaired (between 07:26 breakage and 09:00)
- Symptom: `hermes` failed — `/root/.local/bin/hermes` shebang pointed at `/usr/local/lib/hermes-agent/venv/bin/python3`, a symlink whose target (a uv-managed Python) had been removed. System Python is 3.10.12.
- Attempted fix in-session (`python3 -m venv`) **failed**: `ensurepip is not available` because `python3.10-venv` is not installed on this box.
- **Actual fix (recovered by 09:00):** venv recreated against uv-managed **Python 3.11.15** (which ships pip/ensurepip). `hermes --version` now returns `Hermes Agent v0.18.2, Python: 3.11.15, install method: git`.

### 3. Daily backup (cron, routine)
- `backup-2026-07-17.zip` (13,860 KB) uploaded to `backupHermesDaily/`, 5 old backups (>30d) pruned. Exit 0.

---

## Key learnings
- **Recreate a broken Hermes venv with uv's Python, not system `python3 -m venv`.** System 3.10 here lacks `python3.10-venv` → `ensurepip` fails. uv's 3.11.15 already has pip.
- **Don't trust a session's "done" claim when it ends mid-`execute_code`.** The build session logged "generating update.sh" then was interrupted — verifying disk shows the generate step never wrote files. Always confirm artifacts exist.
- **Config that references missing scripts = silent deploy failure.** `deploy-config.json` points at `scripts/update.sh` for `members`/`internet-storage` that don't exist yet — receiver would log an error on first push. Gate the webhook go-live on all 4 `update.sh` existing.
- **Interruptions surface as `API call failed after 3 retries: [Errno 2] No such file or directory`.** Repeated "continue where you left off" kept failing; a fresh session is the reliable recovery path, not repeated resume.

## Tech stack
- Python stdlib only (`http.server`, `hmac`, `hashlib`, `subprocess`, `threading`) for the receiver — no extra deps.
- Docker Compose per project; `git pull origin main` + `docker compose build && up -d` redeploy pattern.
- Cloudflared tunnel (dashboard-managed, no local config file) planned as the public front for the receiver.
- GitHub webhooks (HMAC `X-Hub-Signature-256`) as the trigger.

## Mistakes / issues overcome
- **Model/infra delay** interrupted the build twice: first session stuck on `[Errno 2]` resume failures; second ("Model Delay Interrupts New Build") cut off at the `execute_code` generate step (48.4s wait). Neither self-recovered.
- **Broken venv** blocked the CLI entirely until recreated with uv Python.

## Status & recommended next steps (for kvinn)
1. Write `members/scripts/update.sh` and `internet-storage/scripts/update.sh` (mirror the wedding-invitation pattern).
2. Decide whether `mentengdutch` belongs in `deploy-config.json` (it has an update.sh but isn't wired).
3. Add a `systemd` unit for `deploy-webhook.py` (loads `.deploy_env`, binds `127.0.0.1:9000`).
4. Register the GitHub webhook (push event, `X-Hub-Signature-256`) → `https://deploy.client.arjism.com/webhook`; add the cloudflared route.
5. Smoke-test with a GitHub `ping`, then a real push to a low-risk repo.
