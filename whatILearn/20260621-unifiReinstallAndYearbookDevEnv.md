# 2026-06-21 — UniFi Controller Reinstall & Yearbook Dev Environment Setup

## What Was Done

### 1. UniFi Controller — HTTP 400 Diagnosis & Fresh Reinstall

**Problem:** `unifi.arjism.com` returned HTTP 400 Bad Request. The UniFi controller was running (MongoDB healthy, Tomcat started) but rejecting all requests.

**Root cause:** Cloudflare Tunnel was proxying to `https://localhost:8443` (UniFi's self-signed HTTPS port). UniFi's embedded Tomcat has strict host-header checking and rejects requests from the tunnel. The error "This combination of host and port requires TLS" appeared when trying plain HTTP on port 8443.

**Diagnosis path:**
- Analyzed UniFi startup logs — MongoDB connections healthy, SCRAM-SHA-256 auth working, WiredTiger checkpoints normal
- Identified the Tomcat host-header mismatch as the root cause
- Explored fix options: `noTLSVerify`, `UNIFI_HOSTNAME` env var, HTTP port 8080 redirect

**Solution chosen:** Fresh reinstall with an nginx reverse proxy container (`unifi-proxy`) on port 8444 that proxies HTTP to UniFi's HTTPS on 8443. Cloudflare Tunnel points to `http://localhost:8444`.

**Key changes from old compose:**
- `mongo:7.0` → `mongo:4.4` (CPU lacks AVX, 7.0 crashes with SIGILL)
- `PID` → `PUID` (typo fix)
- `TZ=Etc/UTC` → `TZ=Asia/Jakarta`
- Added `MONGO_TLS=false`
- Added `unifi-proxy` nginx service on port 8444

**Result:** All 3 containers running (`unifi-db`, `unifi-network-application`, `unifi-proxy`). Verified with `curl -sk http://localhost:8444` → 302 to `/setup/` with valid HTML response.

**Remaining:** kvinn needs to update Cloudflare Tunnel ingress for `unifi.arjism.com` from `https://localhost:8443` to `http://localhost:8444` in the Cloudflare Zero Trust dashboard.

### 2. Digital Yearbook — Local Dev Environment DB Connection Fix

**Problem:** Running backend locally with `go run ./cmd/server` failed with:
```
ping database: failed to connect to `user=postgres database=yearbook`:
  [::1]:5434 (localhost): dial error: dial tcp [::1]:5434: connect: connection refused
  127.0.0.1:5434 (localhost): failed SASL auth: FATAL: password authentication failed for user "postgres"
```

**Root cause:** No `.env` file in the backend directory. The config fallback used `postgres://postgres:***@localhost:5432/yearbook` — placeholder password and wrong port. The Cloudflare tunnel for the DB is on port 5434 (`yearbook-pg.arjism.com` → `localhost:5434`), and the real DB user is `yearbook`/`yearbook_secret`.

**Fix:** Create `.env` file in `/root/hermes/digital-yearbook/backend/`:
```
DATABASE_URL=postgres://yearbook:yearbook_secret@localhost:5434/yearbook?sslmode=disable
```

**Also discussed:** Immich images not loading locally (likely `IMMICH_URL` and `IMMICH_API_KEY` not set in local `.env`), and a 404 on `/2027` route (frontend routing issue — status 200 but SPA fallback serving 404 page).

### 3. Yearbook Backend — Flipbook Page Sorting Fix

**Problem:** Flipbook pages were not sorted in the correct order.

**Fix:** Added sorting by `OriginalFileName` ascending in the Immich service layer (`immich.go`). Backend compiled cleanly, deployed via `docker cp` + `docker restart yearbook-backend`.

**Deploy:** `go build -o /tmp/server ./cmd/server` → `docker cp /tmp/server yearbook-backend:/app/server` → `docker restart yearbook-backend` → health check `curl -sf http://localhost:8080/api/health` → `{"status":"ok"}`

**Committed:** `fix: sort flipbook pages by original filename` → pushed to `lovelymondayz/digital-yearbook` (commit `89e7345`).

### 4. Hermes Workspace Backup (Cron)

**Scheduled cron ran at 19:01 UTC+7:**
- Script: `python3 /root/hermes/scripts/hermes-backup.py`
- Result: ✅ `backup-2026-06-22.zip` (25,558 KB) created and uploaded to GitHub as `lovelymondayz`
- 0 old backups pruned (within 30-day retention)

## Tech Stack Encountered
- **UniFi:** linuxserver/unifi-network-application, mongo:4.4, nginx:alpine, Cloudflare Tunnel
- **Yearbook:** Go backend, React 19 + R3F + Three.js frontend, PostgreSQL, Immich photo storage, Cloudflare Tunnel
- **Infrastructure:** Docker Compose, nginx reverse proxy, Cloudflared tunnels

## Key Learnings
1. **UniFi + Cloudflare Tunnel:** UniFi's self-signed HTTPS on 8443 causes 400 errors when proxied. Best pattern: add an nginx proxy container that handles the TLS termination cleanly, and point Cloudflare Tunnel at the nginx HTTP port.
2. **Mongo 7.0 requires AVX:** Older CPUs (like this VPS) crash with SIGILL on mongo:7.0. Always use mongo:4.4 for compatibility.
3. **Go godotenv fallback pattern:** When no `.env` exists, the Go config falls back to hardcoded defaults. The default had `***` as placeholder password — a trap for local dev. Always document the `.env` template.
4. **Yearbook deployment pattern:** `docker cp` binary into container + restart is the deploy workflow (no CI/CD pipeline for this project).

## Mistakes / Pitfalls
- UniFi compose had `PID` instead of `PUID` — linuxserver images require `PUID`/`PGID`
- mongo:7.0 was used initially — must check CPU AVX support before choosing mongo version
- The yearbook backend had no `.env.example` file, making local setup confusing
