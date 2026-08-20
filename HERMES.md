# Hermes Fleet Documentation

## Overview

Hermes is the platform running on the VPS. Each project is an independent containerized application with its own GitHub repo, docker-compose.yml, and update scripts.

**Fleet count:** 10 projects (8 managed + 2 external)

## Master Script

```bash
/root/hermes/scripts/hermes.sh <command> [options]
```

| Command | What it does |
|---|---|
| `start` | Start all containers (or one with `--project`) |
| `stop` | Stop all containers |
| `restart` | Restart containers (keeps old image) |
| `rebuild` | Pull latest from GitHub + `--no-cache` build + `--force-recreate` |
| `status` | Show containers, ports, health checks |
| `logs` | Tail logs for all or one project |
| `validate` | Check all projects have required files |

**Examples:**
```bash
# Show everything
./hermes.sh status

# Rebuild one project (triggers its own update.sh)
./hermes.sh rebuild --project wedding-invitation

# Rebuild everything
./hermes.sh rebuild --all

# Rebuild even if GitHub hasn't changed
./hermes.sh rebuild --project pico --force
```

> ⚠️ `rebuild` does NOT delete data. It only recreates containers from the latest GitHub code. Database volumes are preserved.

## Projects

| # | Name | Repo | Frontend | Backend | DB Port | Domain |
|---|---|---|---|---|---|---|
| 1 | wedding | `lovelymondayz/wedding-invitation` | :3000 | :8080 | 5432 | wedding.arjism.com |
| 2 | yearbook | `lovelymondayz/digital-yearbook` | :3001 | :8081 | 5434 | yearbook.arjism.com |
| 3 | members | `lovelymondayz/members` | :3003 | :8082 | 5435 | members.arjism.com |
| 4 | noidk | `lovelymondayz/noidk` | :3005 | :8085 | 5439 | noidk.arjism.com |
| 5 | mentengdutch | `lovelymondayz/mentengdutch` | :3002 | — | — | menteng.arjism.com |
| 6 | pico | `lovelymondayz/pico` | :3005 | :8088 | 5436 | pico.arjism.com |
| 7 | qlio | `lovelymondayz/qlio-platform` | :3007 | :8087 | 5438 | qlio.arjism.com |
| 8 | sayless | `lovelymondayz/say-less` | :3006 | :8086 | 5437 | sayless.arjism.com |
| 9 | internet-storage | `lovelymondayz/internet-storage` | :3004 | :8084 | — | — |
| 10 | immich | *(external)* | :2283 | :2283 | — | storage.arjism.com |

## File Structure (per project)

```
project-name/
├── docker-compose.yml          # Container definition (build + env + ports)
├── scripts/
│   └── update.sh              # Manual update: fetch → check → build --no-cache → recreate
├── .env                       # Git-ignored secrets (DB passwords, JWT, etc.)
└── .git                       # GitHub repo (origin → lovelymondayz/<name>)
```

### update.sh standard pattern

Every project's `scripts/update.sh` follows this pattern:

1. `git fetch origin main`
2. Compare LOCAL vs REMOTE commit hashes
3. If different: `git pull origin main`
4. `docker compose build --no-cache` (fresh build, no stale layers)
5. `docker compose up -d --force-recreate` (new containers from latest image)
6. Health check curls on localhost ports

> ⚠️ `--no-cache` + `--force-recreate` ensures no stale images. Never use `docker compose restart` — it keeps old containers running.

## Auto-Update (Webhook)

When you push to GitHub, a centralized webhook receiver on `127.0.0.1:9000` validates the HMAC and runs the project's own `scripts/update.sh`.

```bash
push → GitHub webhook → cloudflared tunnel → deploy-webhook.py
→ validates HMAC → runs scripts/update.sh → build → recreate
```

### Deploy Config

`/root/hermes/scripts/deploy-config.json` maps GitHub repos → VPS directories:

```json
{
  "lovelymondayz/<repo>": {
    "dir": "/path/to/project",
    "update_script": "scripts/update.sh"
  }
}
```

### Non-managed services

| Service | Why no auto-update |
|---|---|
| **immich** | External OSS — v2.5.6 pinned, data bind-mounted at `/home/root/immich-data/` |
| **leadgen-localbiz** | Python scripts, not containers (cron-scheduled) |
| **potion-party** | Roblox Studio project |

## Immich Recovery

Immich is NOT a git repo. All data is bind-mounted at `/home/root/immich-data/`:
- `db/` = Postgres cluster (259MB)
- `upload/` = photo library (2.6GB)

If containers are deleted:

```bash
cd /root/immich
bash restore.sh
```

The script:
1. Checks `.env` and bind mounts exist
2. Pulls pinned images (v2.5.6)
3. Recreates containers
4. Waits for health check

Data survives because it's bind-mounted, not in container layers.

## Backup Strategy

| What | How | Schedule | Retention |
|---|---|---|---|
| DB dumps (7 DBs) | `pg_dump` via `backup-all.sh` | Daily 03:00 | 14 days local |
| .env files | Copied during backup | Daily 03:00 | 14 days local |
| Off-site | Push to `BackupHermes` private repo | Daily 03:30 | 7 snapshots |
| Immich photos (2.6GB) | Bind-mounted at `/home/root/immich-data/upload` | — | No off-site copy |

## Quick Start for New Projects

To add a new project:

1. Create repo on GitHub
2. Clone to VPS
3. Add `docker-compose.yml` with `build:` context
4. Create `scripts/update.sh` using the standard pattern
5. Add to `deploy-config.json`
6. Add webhook in GitHub → Settings → Webhooks → `https://deploy.client.arjism.com`

## Common Tasks

### Check what's running
```bash
/root/hermes/scripts/hermes.sh status
```

### Rebuild one project from latest GitHub
```bash
/root/hermes/scripts/hermes.sh rebuild --project <name>
```

### Tail logs for a project
```bash
/root/hermes/scripts/hermes.sh logs --project <name> --lines 100
```

### Force rebuild all
```bash
/root/hermes/scripts/hermes.sh rebuild --all --force
```

### Recreate everything from scratch (DANGEROUS — only with permission)
```bash
# 1. Stop everything
/root/hermes/scripts/hermes.sh stop --all

# 2. Delete all containers/images (destructive)
docker system prune -a --volumes -f

# 3. Rebuild from GitHub
/root/hermes/scripts/hermes.sh rebuild --all

# 4. Restore Immich (not a git repo)
cd /root/immich && bash restore.sh
```
