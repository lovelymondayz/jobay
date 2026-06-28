# 2026-06-22 — Yearbook Flipbook Routing, Student Sidebar & Styling Cleanup

## What Was Done

### 1. Route Restructure: `/:year` → `/yearbook/:year`

**Problem:** The flipbook was served at `/:year` (e.g. `/2027`), which conflicted with other potential routes and didn't match the desired URL structure `yearbook.arjism.com/yearbook/2027`.

**Changes:**
- `frontend/src/App.tsx`: Changed `<Route path="/:year">` → `<Route path="/yearbook/:year">`
- `frontend/src/pages/HomePage.tsx`: Updated all yearbook card links from `/${y.year}` → `/yearbook/${y.year}`
- Verified Navbar `isActive()` logic still works (it uses `startsWith`, so `/yearbook/2027` correctly highlights nav)

**Commit:** `15d2b80` — `feat: route /yearbook/:year, right sidebar with student list per page`

### 2. Backend API: Students by Page Endpoint

**Problem:** No API existed to fetch students grouped by page number for the sidebar.

**Changes:**
- `backend/internal/handler/flipbook.go`: Added `GetStudentsByPage` handler method
  - Route: `GET /api/v1/flipbook/{year}/students`
  - Looks up yearbook by year, then queries students with `page_number` set
  - Returns `{ yearbook_id, year, page_students: { [page]: [{ id, full_name, avatar_image_url, major }] } }`
- `backend/cmd/server/main.go`: Registered the new route under the flipbook public routes section

### 3. StudentSidebar Component (New)

**What:** A fixed right-side panel showing which students appear on each page of the flipbook.

**File:** `frontend/src/components/StudentSidebar.tsx` (new, ~150 lines)

**Features:**
- Fixed 288px right panel with glassmorphism styling (backdrop-blur, bg-white/5)
- **Top section — "On This Page":** Auto-expands when the user flips to a page that has students; shows avatar, full name, and major for each student
- **Scrollable accordion list:** All pages with students, grouped by page number; each section is collapsible
- **Collapse/expand toggle button:** Hides the sidebar to give full width to the flipbook
- Page indicator in FlipBookPage updated to "Page X of Y" format

### 4. FlipBookPage Updates

**Changes:**
- Added `currentPage` state tracked via the `FlipBookScene` component's page change events
- Integrated `StudentSidebar` component, passing `currentPage` and `studentsByPage` data
- Fetches `/api/v1/flipbook/{year}/students` on mount
- Header text shifts left to avoid overlapping the sidebar
- Page indicator repositioned (bottom-24 instead of bottom-7.5) to clear the sidebar

### 5. Styling Cleanup (Commit `4cb9145`)

**Changes:**
- Removed `framer-motion` from Navbar component (replaced with CSS transitions for performance)
- Fixed active Navbar state detection (was not highlighting correctly)
- Cleaned up unused imports

### 6. Deployment

**Deploy flow (no CI/CD for this project):**
```bash
# Backend
cd /root/hermes/digital-yearbook/backend
go build -o /tmp/server ./cmd/server
docker cp /tmp/server yearbook-backend:/app/server
docker restart yearbook-backend

# Frontend
cd /root/hermes/digital-yearbook/frontend
npm run build
docker cp dist/ yearbook-frontend:/usr/share/nginx/html/
docker restart yearbook-frontend
```

**Verification:**
- `curl -sf http://localhost:8080/api/health` → `{"status":"ok"}`
- Frontend rebuilt and served via nginx container
- All containers healthy after restart

### 7. Hermes Workspace Backup (Cron)

- Scheduled cron ran automatically
- `backup-2026-06-22.zip` created and uploaded to GitHub (`lovelymondayz` repo)
- 0 old backups pruned (within 30-day retention)

## Tech Stack

- **Frontend:** React 19, React Router DOM, Tailwind CSS, Framer Motion (being phased out), Lucide React icons
- **3D:** React Three Fiber, Three.js, `@react-three/drei`, `quick_flipbook` library
- **Backend:** Go 1.22+, Chi router, pgx/v5 PostgreSQL driver, Immich API integration
- **Infrastructure:** Docker Compose, nginx reverse proxy, Cloudflare Tunnel, PostgreSQL
- **Deployment:** Manual `docker cp` + `docker restart` (no CI/CD pipeline)

## Key Learnings

1. **Route design matters early.** Using `/:year` was too greedy — it could conflict with future routes. Moving to `/yearbook/:year` is more explicit and RESTful. Always consider future route expansion when designing URL structures.

2. **Student-page relationship is 1:many via `page_number`.** The `students` table has a `page_number` column that directly maps a student to a page in the yearbook. This is a simple but effective design — no join table needed for this use case.

3. **Lazy loading all non-critical routes.** Commit `7cfaf3f` (reviewed during this session) converted all routes except HomePage and FlipBookPage to `lazy()` with `Suspense` fallbacks. This reduces the initial bundle size significantly. The pattern:
   ```tsx
   const YearbookPage = lazy(() => import("./pages/YearbookPage"));
   // In route:
   <Suspense fallback={<LoadingFallback />}><YearbookPage /></Suspense>
   ```

4. **Sidebar + fullscreen canvas coexistence.** The flipbook uses a fullscreen Canvas (R3F). Adding a fixed sidebar required careful z-index management and shifting the header text left. The `h-full` (not `h-screen`) on the FlipBookPage container is critical when inside a flex Layout.

5. **Yearbook deployment is manual.** No CI/CD — `docker cp` binary into container + restart. This works for a small project but is error-prone. A future improvement would be a GitHub Actions pipeline that builds and deploys on push to main.

## Mistakes / Pitfalls

- **HomePage links used `/${y.year}`** — forgot to update these when the route changed. Had to catch and fix in the same commit.
- **Header overlap with sidebar:** The centered header text overlapped with the 288px right sidebar. Fixed by shifting the header container left when sidebar is visible.
- **Page indicator position:** Was at `bottom-7.5` which was too low after adding the sidebar. Moved to `bottom-24`.

## Files Changed (Today)

| File | Change |
|------|--------|
| `frontend/src/App.tsx` | Route `/:year` → `/yearbook/:year` |
| `frontend/src/pages/HomePage.tsx` | Card links updated to `/yearbook/{year}` |
| `frontend/src/pages/FlipBookPage.tsx` | Added currentPage state, StudentSidebar integration, page tracking |
| `frontend/src/components/StudentSidebar.tsx` | **New** — right sidebar with student list per page |
| `backend/internal/handler/flipbook.go` | Added `GetStudentsByPage` handler |
| `backend/cmd/server/main.go` | Registered `GET /api/v1/flipbook/{year}/students` route |

## Commits

```
15d2b80 feat: route /yearbook/:year, right sidebar with student list per page
4cb9145 fix: clean styling, remove framer-motion from navbar, fix active nav state
```

Both pushed to `lovelymondayz/digital-yearbook` on `origin/main`.
