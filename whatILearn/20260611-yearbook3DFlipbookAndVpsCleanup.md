# 2026-06-11 Digital Yearbook — 3D Flipbook, React 19 Upgrade & VPS Cleanup

## What Was Built / Fixed Today

### Session 1 (0:00–05:55 UTC) — Carry-over from June 10 Sessions
_These sessions started on June 10 but extended into early June 11 (UTC). Covered in detail in `20260610-constitutionAndSystemRedesign.md`._

**Summary of carry-over work:**
- Digital Yearbook Go backend deployed — 3 containers running (yearbook-db, yearbook-backend, yearbook-frontend)
- 38 PostgreSQL migration files, 23 Go source files, 33 React/TS frontend files — all pushed to GitHub
- Immich integration: API key hashed with SHA-256 base64, upload endpoint wired
- Fixed nil Go slice → JSON `null` bug (3 services patched + defensive frontend null guard)
- VPS port mapping: yearbook uses 3001 (fe) + 8081 (be) to avoid wedding conflicts (3000/8080)
- Nginx reverse proxy for `/api/` → `backend:8080` internal routing
- Container names set explicitly to avoid Docker `-1` suffix

---

### Session 2 (~07:38–08:53 UTC) — Yearbook URL Structure, 3D Flipbook & Admin Routes

**User request:** Redesign yearbook URL structure to support `http://yearbook.thamrin.ac.id/2026`, `2027`, etc. No login required — homepage shows year selection. Implement 3D flipbook using reference from `/root/hermes/reference/publication/`.

#### What was done:

1. **Reviewed reference publication** — identified 3D flipbook stack: `@react-three/fiber` + `@react-three/drei` + `quick_flipbook`
2. **Audited current codebase** at `/root/hermes/digital-yearbook/` — React 18 + Vite + Tailwind, Go backend with Immich integration
3. **First attempt with `quick_flipbook`** — failed due to `three.modifiers` incompatibility with modern Three.js 0.184
4. **Wrote custom `FlipBookScene`** using only `@react-three/fiber` + `@react-three/drei` (no external flipbook lib)
5. **Fixed Docker build** — switched from `npm ci` to `npm install` to avoid lockfile arch mismatch
6. **Set up nginx reverse proxy** container on port 80 → frontend:80 + backend:8080
7. **Diagnosed runtime crash** — `@react-three/fiber` v9.6.1 requires React 19, but project had React 18.3.1. `--legacy-peer-deps` masked install but runtime reconciler crashed
8. **Upgraded entire stack to React 19** — react 19.2.7, react-dom 19.2.7, R3F 9.6.1, drei 10.7.7, Three.js 0.184, @types/three 0.184.1, @types/react 19.2.17, framer-motion 11.18.2
9. **Clean `npm install`** without `--legacy-peer-deps` — succeeded
10. **Built frontend** — three.chunk (901 KB), vendor.chunk (396 KB), FlipBookScene.chunk (3.23 KB), no errors
11. **Removed Login/Register from Navbar** (both desktop and mobile)
12. **Moved auth routes** — `/login` and `/register` → `/admin/login` and `/admin/register`
13. **Rebuilt Docker + deployed** — all containers healthy
14. **Set up auto-deploy** — `scripts/auto-deploy.sh` (cron every 2 min), `scripts/push-and-deploy.sh`, `scripts/deploy.sh`, `scripts/.last-deployed-commit`
15. **Committed + pushed** all changes to GitHub

**Domain configuration:**
- `yearbook.arjism.com` — live via cloudflared → localhost:3001
- `yearbook.thamrin.ac.id` — needs Cloudflare dashboard config (A record + tunnel ingress rule, user action required)
- `http://103.47.134.107:3001` — direct VPS access, live

#### Flipbook 404 Fix (Bug #1)

**Problem:** User reported 404 when visiting `https://yearbook.arjism.com/2026`. HTTP 200 returned but page was blank.

**Diagnosis:** Vite's `manualChunks` config put `three`, `@react-three/fiber`, and `@react-three/drei` all into one chunk called `three`. But FlipBookScene imported from both `three` AND R3F/drei — the bundler was resolving R3F/drei imports from the wrong chunk, causing runtime crash.

**Fix:** Separated chunks — `three` in one chunk (746 KB), `r3f` (fiber + drei) in another (153 KB), vendor in separate chunk (396 KB).

#### Flipbook Blank Screen Fix (Bug #2)

**Problem:** After chunk fix, the `useLoader` hook from `@react-three/fiber` still wasn't being included correctly in the r3f chunk due to Vite code-splitting edge case.

**Fix:** Replaced `useLoader(THREE.TextureLoader, imageUrl)` with direct `THREE.TextureLoader` instance using `useEffect` + `useState`:
```tsx
const textureLoader = new THREE.TextureLoader();
useEffect(() => {
  textureLoader.load(imageUrl, (loaded) => {
    loaded.colorSpace = THREE.SRGBColorSpace;
    setTexture(loaded);
  });
}, [imageUrl]);
```

**Committed + pushed + verified deployed.**

---

### Session 3 (09:26–11:04 UTC) — Docker Container Health Check & Cleanup

**User request:** Run `docker ps -a` to check container status.

#### Timeline:

**09:26 UTC — First check (all healthy):**
- 12 containers, all Up, all healthy
- Yearbook stack freshly deployed (2-21 min ago)
- Wedding, Lead Gen, Immich all stable

**10:46 UTC — Second check (after rebuilds):**
- 15 containers — 3 new exited build artifacts appeared:
  - `youthful_meitner` — Exited (0) ✅
  - `exciting_chandrasekhar` — Exited (1) — `npm run` build step
  - `great_newton` — Exited (1) — `npm ci` build step
- 13 active services all healthy
- New container: `yearbook-proxy` (nginx:alpine on port 80→80)

**10:48 UTC — Cleanup:**
- User requested removal of unnecessary containers
- Ran `docker container prune -f`
- **Removed 3 exited containers, reclaimed 746MB**

**11:04 UTC — Final verification:**
- 13 containers, all Up, all healthy
- No exited containers remaining
- Yearbook (fe, be, proxy, db) ✅ | Wedding (fe, be, db) ✅ | Lead Gen (pg) ✅ | Immich (5 containers) ✅

---

### Cron Jobs (11:17–11:35 UTC) — Auto-Deploy Polling

Multiple auto-deploy cron runs throughout the day:
- Successfully detected and deployed new commits from the yearbook session
- All 3 services remained Up & Healthy across all poll cycles
- `.last-deployed-commit` tracked latest deployed hash

---

## Tech Stack

| Layer | Technology |
|---|---|
| **Frontend** | React 19.2.7 + TypeScript + Vite + Tailwind CSS |
| **3D Rendering** | @react-three/fiber 9.6.1 + @react-three/drei 10.7.7 + Three.js 0.184 |
| **Backend** | Go 1.22+ (chi router, pgxpool, bcrypt, JWT) |
| **Database** | PostgreSQL 16-alpine |
| **Image Storage** | Immich (via `immich_server:2283`) |
| **Reverse Proxy** | Nginx:al SPA + API proxy |
| **Deployment** | Docker Compose, Cloudflare tunnel (home-vps) |
| **Auto-Deploy** | cron (every 2 min) polling GitHub → auto redeploy |
| **CDN/DNS** | Cloudflare (cloudflared tunnel) |

## Key Learnings

### 1. Vite manualChunks Can Cause Wrong Chunk Resolution
When R3F/drei and Three.js are in the same manual chunk, Vite may resolve imports from the wrong source module. **Always separate framework-adjacent but independently imported libraries into their own chunks.**

### 2. useLoader Hook Is Fragile With Code-Splitting
The `useLoader` hook from R3F depends on internal R3F context that Vite's code-splitting can break. **For texture loading in production builds, prefer direct `THREE.TextureLoader` with `useEffect` + `useState` over `useLoader`.**

### 3. React 19 + R3F 9.x = Required Match
`@react-three/fiber` v9.x requires React 19. Using React 18 with `--legacy-peer-deps` installs fine but crashes at runtime with reconciler errors. **Always match R3F major version to React major version.**

### 4. Removing Public Auth Routes Improves Security UX
Moving `/login` → `/admin/login` and removing all Login/Register links from the Navbar creates a "hidden admin" pattern. The admin panel exists but is never linked from public UI — security through obscurity layer on top of real auth.

### 5. Build Artifact Containers Accumulate Fast
`npm ci` and `npm run build` inside Docker multi-stage builds create intermediate containers that don't auto-clean. **Run `docker container prune` regularly** or add `--rm` to build commands. These can waste 746MB+ per project.

### 6. Immich Album Naming Convention Matters
The flipbook API queries Immich albums by name pattern `"Thamrin Graduate {year}"`. Any mismatch in album naming = empty flipbook. User must ensure album names match exactly.

## Mistakes Overcome

| # | Mistake | Fix |
|---|---|---|
| 1 | `quick_flipbook` library imports `three.modifiers` (ancient package incompatible with Three.js 0.184) | Wrote custom FlipBookScene with R3F + drei only |
| 2 | Docker `npm ci` fails on lockfile arch mismatch | Switched to `npm install` in Dockerfile |
| 3 | R3F 9.x + React 18 runtime crash (reconciler mismatch) | Upgraded entire project to React 19 |
| 4 | Vite manualChunks merging three + r3f + drei into one chunk | Separated into independent `three` and `r3f` chunks |
| 5 | `useLoader` hook broken by code-splitting | Replaced with direct `THREE.TextureLoader` + `useEffect` |
| 6 | Build artifact containers accumulating on VPS | `docker container prune -f` to reclaim 746MB |

## Remaining Blockers

1. **`yearbook.thamrin.ac.id` DNS** — Needs Cloudflare dashboard config (A record + tunnel ingress rule). No Cloudflare API token available. User must configure manually via Zero Trust → Tunnels → `home-vps` → Add public hostname.
2. **Immich album page count** — API returns 9 pages instead of expected 10 for "Thamrin Graduate 2026" album. User should verify album contents in Immich.
3. **Cloudflare CI/CD** — GitHub PAT lacks `workflow` scope, so CI/CD workflow wasn't pushed. User needs to regenerate PAT or add workflow file manually.

## Next Steps

1. User verifies `https://yearbook.arjism.com/2026` loads 3D flipbook (hard refresh: `Ctrl+Shift+R`)
2. Configure `yearbook.thamrin.ac.id` in Cloudflare dashboard
3. Test flipbook navigation (Previous/Next buttons + keyboard controls)
4. Seed initial data (university + yearbook) via admin API
5. Test Immich image upload through yearbook backend
6. Clean up Immich album to ensure correct page count

---

*Log written: 2026-06-11 19:00 UTC (cron job)*
*Sessions covered: 2 Discord sessions + 4 cron auto-deploy runs + 1 cron health check*
*Session IDs: `20260611_073829_da737d7c`, `20260611_092655_ffe9898a`, `cron_880bdff3726a_*` (multiple runs)*
*Total tool calls across all sessions: ~50+*
*Primary project: `/root/hermes/digital-yearbook` → `github.com/lovelymondayz/digital-yearbook`*
