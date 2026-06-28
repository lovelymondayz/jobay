---
date: 2026-06-23
title: Jakarta Munch Clone Build, Plagiarism Redesign, Obsidian Setup & Docker Status
---

## What Was Built / Fixed Today

### Session 1 (02:21 UTC) — Jakarta Munch Website Clone (Full Build)

**User request:** Build a clone of https://www.jakartamunch.com/ with exact styling.

**What was done:**
1. Analyzed the target site via `curl` (browser crashed due to Chrome DevTools issue — fell back to curl)
2. Identified the site as an **Astro + Tailwind CSS** static site for "Jakarta Munch" (Indonesian restaurant, NYC)
3. Key design system discovered:
   - Colors: `#F5F2EB` (warm cream), `#ED3F1C` (red), `#dfe7b4` (lime green), `#2f2f2f` (dark)
   - Font: Archivo (variable, 100-900 weight)
   - Sections: Hero, Banner, Tabbed Menu (Bowls/Signature/Desserts), About, Catering, Reviews, Location, Videos, Shop, Footer
   - Interactions: Order dropdown, mobile burger nav, menu tabs, video lightbox, scroll fade animations, client logo marquee
4. Scaffolded React + Vite + TypeScript + Tailwind project at `/root/hermes/jakarta-munch/`
5. Hit npm install timeout (60s) → retried with `--prefer-offline` → failed on `semver@7.8.5` not found
6. Resolved by letting npm resolve its own versions
7. Built all 10 sections as React components with exact styling
8. Dockerized: multi-stage build with `node:20` → `nginx:alpine`
9. Deployed as container `jakarta-munch` on port 3002
10. Verified live: `curl localhost:3002` → 200 OK

**Final state:** Container running, site live at port 3002.

**Mistakes & Learnings:**
1. **Browser crash on jakartamunch.com** — Chrome DevTools crashed with SIGILL. Fixed by falling back to `curl` for HTML extraction.
2. **npm install timeout** — `npm create vite` + `npm install` exceeded 60s timeout. Fixed with longer timeout (120s) and `--prefer-offline`.
3. **semver@7.8.5 not found** — npm registry had a transient issue. Fixed by retrying without the flag.
4. **Docker build warning** — `version` attribute in docker-compose.yml is obsolete (Compose v2 ignores it). Non-critical.

---

### Session 2 (03:04 UTC) — Docker Container Status Overview

**User request:** Run `docker ps` and post results.

**What was done:**
- Ran `docker ps --format "table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}"`
- Found 16 containers running, 6 marked healthy
- New: `jakarta-munch` (9 min old, port 3002) — no healthcheck yet
- All other services healthy: wedding (3), yearbook (3), unifi (3), immich (5), leadgen-postgres
- 3 PostgreSQL instances: wedding (5432), leadgen (5433), yearbook (5434)

---

### Session 3 (09:30 UTC) — Using Obsidian with whatILearn Logs

**User request:** Can Obsidian be installed on the VPS, or is whatILearn enough?

**What was done:**
1. Explained whatILearn is: flat folder of daily `.md` files, no backlinks/graph/plugins
2. Evaluated 3 options:
   - Install Obsidian on VPS → rejected (Electron GUI on headless server is wasteful)
   - Use locally, point at whatILearn → recommended
   - Self-hosted alternative (Trilium/Logseq) → offered
3. **Memory fix:** Found that goals/motto were in `fact_store` but never promoted to active memory
   - Searched fact_store, found the 4 pipeline goals from May 30
   - Promoted to memory: Motto "Burn VPS, Not Tokens" + 4 goals
   - Had to remove a less-critical memory entry (2,200 char limit hit)
   - Lesson: **fact_store ≠ memory**, probe it immediately when asked about stored facts

**Key learnings:**
- Memory is 2,200 chars max — must prune to make room for critical facts
- fact_store is a separate entity store; session_search doesn't cover it
- When user asks "do you remember X?" → check fact_store AND session_search AND memory

---

### Session 4 (10:50 UTC) — Website Redesign to Avoid Plagiarism (COMPLETED)

**User request:** "Change the colour or layout so it doesn't look much and avoid the plagiarism"

**Context:** kvinn noticed the Jakarta Munch clone looked too similar to the original. Wanted differentiation for client.arjism.com use (lead gen pipeline — Goal 1).

**What was done:**
1. Analyzed the original site's design DNA (colors, typography, layout, feel)
2. Planned differentiation strategy: change color palette, typography, layout structure, visual treatments, hero treatment
3. Asked user for color mood preference → chose "Surprise me"
4. Designed new "Toko Menteng" identity — Indonesian toko in Netherlands:
   - **Palette:** Deep charcoal `#0F0F0F` base + amber `#D4A843` accent + off-white `#F8F4ED`
   - **Typography:** Space Grotesk (geometric, modern) vs JM's Archivo
   - **Layout:** Split hero (text left, image right), list menu cards, sharp corners everywhere
   - **Visual treatment:** Thin 1px borders, grain texture overlay, no rounded elements, no marquee, no curved dividers
5. Rewrote README with full design system documentation
6. Build passed clean — dist generated

**Status:** ✅ Redesign COMPLETED. Build passed. **NOT YET REDEPLOYED** to Docker container.

**Note:** At this point the container `jakarta-munch` (port 3002) still serves the ORIGINAL Jakarta Munch clone. The redesigned "Toko Menteng" build exists in `/root/hermes/jakarta-munch/dist/` but hasn't been redeployed yet.

---

### Session 5 (11:16 UTC) — Wedding Invitation Project Recap

**User request:** "Remember our wedding-invitation project?"

**What was done:**
1. Searched for wedding-invitation sessions — found references in docker ps output
2. Located project at `/root/hermes/wedding-invitation/` (50 git objects, backend API at `wedding-api`)
3. Read existing whatILearn log from June 9 (`20260609-weddingInvitation.md`) — comprehensive doc
4. Summarized for kvinn:
   - Stack: Go + Gin + PostgreSQL (backend) / React 18 + Vite + Tailwind + Framer Motion (frontend)
   - Docker: 3 services (wedding-db:5432, wedding-backend:8080, wedding-frontend:3000)
   - All 3 containers running 44h+ at time of check
   - GitHub: `lovelymondayz/wedding-invitation`
   - 8 documented mistakes with solutions

---

## Tech Stack Used
- React 19 + TypeScript + Vite 8 + Tailwind CSS 3
- Docker (multi-stage: node:20 build → nginx:alpine)
- Astro (analyzed target site, not used in build)
- npm / node_modules management
- curl for site analysis (browser fallback)

## Infrastructure Events
- New container: `jakarta-munch` on port 3002 (React + Nginx)
- Total containers: 16 running
- No incidents or outages

## Mistakes Overcome
1. **Browser crash (Chrome SIGILL)** — fell back to curl for site analysis
2. **npm install timeout** — extended timeout to 120s
3. **semver registry issue** — retried without --prefer-offline
4. **Memory full (2,200 chars)** — pruned less-critical entry to make room for goals
5. **fact_store not searched initially** — now always check fact_store when user asks "do you remember?"

## Key Learnings
- `fact_store` is a separate knowledge entity store from `memory` — both must be checked
- Memory has hard 2,200 char limit — prioritize user goals/preferences > environment facts
- Browser is unreliable for heavy JS sites on this VPS — curl first, browser second
- Docker healthcheck: jakarta-munch has no healthcheck yet (unlike all other services)
- Astro sites can be perfectly replicated with React + Tailwind (same CSS classes, same DOM structure)

## Next Steps
- Add healthcheck to jakarta-munch container
- **Redeploy redesigned "Toko Menteng" build to Docker** (build exists in dist/, not yet deployed)
- Continue wedding-invitation Phase 2 features
- Address lead gen pipeline SerpApi rate limit (still blocked as of earlier sessions)
