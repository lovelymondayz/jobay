# Jobay Project — AI Job Search Agent

Jobay is a Go + React web app that lets users upload a CV, then an AI agent searches for matching jobs, scores them, and attempts auto-apply.

## Architecture

```
User uploads CV → Go backend stores file
    ↓ (background goroutine)
CV parser (PDF/DOCX → text) via ledongthuc/pdf + docx
    ↓
AI extractor (text → structured profile) via 9router
    ↓
Job searcher (Firecrawl + RemoteOK) → job listings
    ↓
AI scorer (profile × job → 0-100 score) via 9router
    ↓
Auto-applier (best-effort: detects ATS, logs result)
    ↓
Jobs saved to jobay.json + WebSocket notify
```

## Deployment

- **Dir**: `/root/hermes/jobay`
- **Port**: `3010:8080` (NOT 3001 — that's digital-yearbook)
- **DB**: None (flat-file `jobay.json`)
- **Update script**: `/root/hermes/jobay/update.sh`

### Docker Compose Env Vars

| Var | Value |
|-----|-------|
| `PORT` | `8080` |
| `DB_PATH` | `/app/data/jobay.json` |
| `UPLOADS_DIR` | `/app/data/cvs` |
| `AI_BASE_URL` | `https://9router.nousresearch.com/v1` |
| `AI_MODEL` | `poolside/laguna-m.1:free` |
| `FIRECRAWL_API_KEY` | from `/root/.hermes/archive/firecrawl/.env.firecrawl` |

### Dockerfile

Uses **Go 1.24** (NOT 1.22). The `go get` of `github.com/ledongthuc/pdf` auto-upgraded `go.mod` from 1.22 → 1.24.1. If you add new Go deps and `go.mod` bumps, update the Dockerfile's `golang:` image to match.

```dockerfile
FROM golang:1.24-alpine AS builder  # must match go.mod directive
```

## Gotchas

### 1. Go Version Auto-Upgrade
Adding `github.com/ledongthuc/pdf` bumped `go.mod` from `go 1.22` to `go 1.24.1`. Always check `go.mod` after `go get` and update Dockerfile.

### 2. Firecrawl Key May Be Invalid
The key `fc-58316613a6a452575666905d76ed9ad0` has returned `Unauthorized: Invalid token`. The searcher now falls back to RemoteOK if Firecrawl fails, so the app stays functional. To fix Firecrawl: get a new key at https://firecrawl.dev and update `docker-compose.yml`.

### 3. Large Frontend Bundle
Build output: `dist/assets/index-CoESLpSK.js 537.26 kB │ gzip: 152.92 kB`. Vite warns "Some chunks are larger than 500 kB after minification." For a 2c/12GB VPS this loads fine, but if you add more deps consider code-splitting.

### 4. `docker compose up -d` Rejected by Terminal Tool
Hermes terminal tool rejects `docker compose up -d` as a long-lived process. Use `background=True`:
```python
terminal(command="cd /root/hermes/jobay && docker compose up -d", background=True)
```

### 5. Auto-Apply is Detection, Not Submission
The current applier only **detects** whether a job page has an apply form — it does NOT fill or submit forms. Real auto-apply needs either:
- A headless browser sidecar (Playwright)
- A Chrome extension that uses the user's logged-in session

### 6. Dual Search Pattern
The searcher uses both Firecrawl (if key is valid) and RemoteOK (always available). This is a graceful-degradation pattern: if the paid-tier source fails, the free-tier source still returns results. Useful pattern for any multi-source aggregator.

## API Endpoints

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/api/upload` | POST | Upload CV, create user, trigger agent |
| `/api/users/:slug` | GET | Get user profile |
| `/api/users/:slug/jobs` | GET | Get user's matched jobs |
| `/api/agent/run` | POST | Manually trigger agent for a user |
| `/api/jobs` | GET | List all jobs (admin) |

## Frontend Pages

| Path | Component |
|------|-----------|
| `/` | AdminDashboard |
| `/upload` | UploadPage |
| `/:slug` | UserDashboard |

## Per-User Dashboard URL

`https://jobay.arjism.com/<name-slug>` (e.g., `/arji-maulana`). Shows stats cards, job list with scores, activity log.
