# 2026-06-28: Members Project Bug Fixes & Frontend Proxy Fix

## What Was Done Today

### 1. Members Project — 17 Bug Fixes & Full Redeploy
- Fixed **stores endpoint** (new `store.go` handler + `FindAll` repository method)
- Fixed **role-as-int bug** — Login/Register now return string roles (`"super_admin"/"admin"/"member"`) via new `roleName()` helper
- Fixed **invoice member display** — changed `inv.member?.user_name` → `inv.member_name` (flat field matching API)
- Fixed **invoice member dropdown** — changed `m.user?.name` → `m.user_name`
- Fixed **dashboard permissions** — frontend uses `/admin/dashboard` for super_admin, `/members` list for regular admin (avoids 403)
- Fixed **super admin store scoping** — super admin without `store_id` in JWT now returns ALL records (added `FindAll` to member + invoice repositories)
- Fixed **GetAdmins handler** — removed broken `Preload("Role")` call (no Role struct field after circular FK fix)
- Both backend (`go build`) and frontend (`tsc --noEmit`) compile clean

### 2. Critical Frontend Nginx Proxy Fix
- **Root cause:** Frontend container's nginx had no `/api/` proxy rule. Cloudflare hitting port 3003 directly got 405 on API calls.
- **Fix:** Added `location /api/` block to `frontend/nginx.conf` → `proxy_pass http://members-backend:8082`
- **Verified ALL endpoints working through port 3003:**
  - `POST /api/auth/login` → returns token + `role: "super_admin"`
  - `GET /api/members` → returns members with store_name
  - `GET /api/invoices` → returns invoices with member_name
  - `GET /api/stores` → returns all stores
  - `GET /api/admin/dashboard` → returns stats
  - `GET /api/admin/admins` → returns admins with role names
  - `GET /health` → `{"status":"ok"}`
- Pushed to GitHub (`lovelymondayz/members` commit `6eda14f`)

### 3. Cronjob Diagnostics
- Investigated why whatILearn cronjob stopped — found **OpenRouter 502 provider outage** (not a config/code issue)
- Triggered both cronjobs manually:
  - whatILearn log → re-triggered successfully
  - Hermes daily backup → ran successfully, `backup-2026-06-29.zip` (35.7 MB) uploaded to `backupHermesDaily/`

### 4. Cronjob Audit
- Confirmed only 2 active cronjobs (no rogue members cron):
  1. Hermes daily workspace backup — `0 19 * * *` — OK
  2. Daily whatILearn log — `15 19 * * *` — recovered from 502 error

## Key Learnings

| Lesson | Detail |
|--------|--------|
| **Nginx proxy in multi-container setups** | When frontend and backend are separate containers, the frontend nginx MUST proxy `/api/` to the backend. Cloudflare can hit any exposed port directly. |
| **GORM transient fields** | `gorm:\"-\"` fields are NOT auto-populated by Preload — repositories must manually populate them via separate queries. |
| **Circular FK crashes** | GORM AutoMigrate with struct relations (Member→Store→User→Member) causes FK creation order failures. Use transient fields + manual queries instead. |
| **Role as string vs int** | Frontend role checks break when backend returns int but frontend compares to string. Always return strings for roles. |
| **Docker build timeouts** | Go compilation in container takes ~5+ min. Build was killed at timeout. Retry works faster due to layer caching. |
| **OpenRouter 502** | Provider outages cause cronjob failures. Not a code issue — just retry. |

## Tech Stack
- **Backend:** Go 1.25 + Gin + GORM + Postgres 16
- **Frontend:** React + Vite + Tailwind + Zustand + React Query + Nginx
- **Infra:** Docker Compose, Cloudflare DNS, nginx reverse proxy
- **Repos:** `lovelymondayz/members` (separate from monorepo)

## Mistakes Overcome
1. **Frontend nginx missing /api proxy** — API calls from Cloudflare direct to port 3003 got 405. Fixed with proxy_pass block.
2. **Super admin store scoping** — super admin without `store_id` got "store_id required" error. Fixed by adding `FindAll` fallback for super admin role.
3. **GetAdmins broken after model refactor** — `Preload("Role")` referenced removed struct field. Fixed by removing the preload and using `roleName()` helper.
4. **Docker build killed by timeout** — Go compile too slow in container. Retried and succeeded with cached layers.

## Remaining Work
- WhatsApp OTP auth (Phase 4)
- Payment gateway integration (Midtrans/Xendit)
- Member self-service portal polish
- QR code image generation (currently returns string only)
- Rate limiting
