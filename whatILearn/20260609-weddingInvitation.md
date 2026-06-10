---
date: 2026-06-09
title: Wedding Invitation App — Full-Stack Build
---

## What We Built Today

A complete premium digital wedding invitation website with:

### Backend (Go)
- REST API with Gin framework
- PostgreSQL database with pgx/v5
- JWT authentication (jwt/v5)
- Full CRUD for guests, RSVPs, wishes, gallery, music, schedule, love story, gift info
- CSV import/export for guests
- File upload endpoint
- Analytics dashboard endpoint

### Frontend (React + TypeScript)
- React 18 + Vite + TypeScript
- Tailwind CSS with custom cream/gold wedding theme
- Framer Motion animations throughout
- All public sections: Hero, Countdown, Wedding Info, Love Story Timeline, Schedule, Gallery, Video, Map, RSVP Form, Wishes, Gift
- Personalized invite popup with slug-based URLs
- Floating music player with play/pause/volume
- Floating bottom navigation
- Full admin dashboard: Overview stats, Guest Management, RSVPs, Wishes, Gallery, Music, Schedule, Love Story, Gift, Settings

### Deployment
- GitHub repo: https://github.com/lovelymondayz/wedding-invitation
- Notion journal entry created with full mistake documentation

## Mistakes & Learnings

1. **Go version incompatibility**: Go 1.18 can't use pgx/v5 or jwt/v5 (need Go 1.19+). Solution: upgrade Go to 1.22+ first.
2. **node_modules in git**: First commit included 5223 node_modules files. Solution: create .gitignore BEFORE git add -A.
3. **Subagent rate limits**: Delegating large code gen to subagents hits owl-alpha rate limits. Solution: use write_file directly.
4. **Sibling subagent file conflicts**: Subagents write files parent doesn't know about. Solution: always `find` before writing.
5. **TypeScript union types**: Radio button onChange returns string but state expects `'attending' | 'not_attending'`. Solution: use `as` cast.
6. **Tailwind @import order**: `@import url()` must come BEFORE `@tailwind` directives.
7. **API format mismatch**: Frontend services expected `{data: wrapper}` but backend returned raw objects. Solution: keep consistent.
8. **Missing React hook imports**: useState/useEffect not imported. Solution: include all hooks in import.

## Tech Stack Used
- Go 1.22.5 + pgx/v5 + Gin + JWT + PostgreSQL
- React 18 + TypeScript + Vite + Tailwind CSS + Framer Motion + React Router
- GitHub API for repo creation
- Notion API for journal entries

## Repositories Created/Updated
- lovelymondayz/wedding-invitation: https://github.com/lovelymondayz/wedding-invitation

## Infrastructure Setup Today
- Created `/root/.hermes/whatILearn/` folder for daily learning logs
- Created cron job "Daily whatILearn log" (ID: `e66555007b8c`) — runs at 19:15 UTC daily
- Format: `YYYY-MM-DD-{summary-title}.md`
- Cron checks for existing file first (won't duplicate), uses session_search to find today's work
- Consolidated memory: all projects now go in `/root/.hermes/` (not `/root/`)
- Updated `subagent-driven-development` skill to v1.3.0 with `frontend-build-pitfalls.md` reference

## Docker Setup

All 3 services containerized and running:

| Service | Container | Port |
|---------|-----------|------|
| PostgreSQL 16 | `wedding-db` | 5432 |
| Go Backend API | `wedding-backend` | 8080 |
| React Frontend (Nginx) | `wedding-frontend` | 3000 |

### Files Created
- `backend/Dockerfile` — multi-stage Go 1.25 → Alpine
- `frontend/Dockerfile` — Node 20 build → Nginx 1.25
- `frontend/nginx.conf` — SPA fallback + `/api` proxy to backend
- `docker-compose.yml` — 3 services, health checks, startup order, named volumes
- `.env` — DB credentials + JWT secret

### Mistakes & Learnings (Docker)
1. **go.mod `go 1.25` vs system Go 1.22**: System had Go 1.22.5 but go/mod dependencies required 1.25. Fix: use `golang:1.25-alpine` Docker image instead of `golang:1.22-alpine`.
2. **Migration SQL dollar-quoted strings**: Go runner splits on `$$` delimiters breaking `CREATE FUNCTION` blocks. Fix: run migrations via `docker exec ... psql` directly.
3. **PostgreSQL DATE/TIME to Go `*string`**: pgx can't auto-cast `DATE`/`TIME` columns to `*string`. Fix: add `::text` casts in SQL queries.
4. **`depends_on` with `condition: service_healthy`**: Compose v2 format requires service names to match exactly. Backend/frontend didn't auto-start with `up -d` — had to `docker start` manually.
5. **Duplicate seed data**: Running migrations twice (Go runner + psql) inserted duplicate love story, schedule, and gallery entries. Not critical but cleanup needed.

## Hermes Atlas Ecosystem Research

Studied hermesatlas.com — the community ecosystem map for Hermes Agent (168 repos, 12 categories).

### Key Takeaways for Our Setup

**Memory Architecture (3 layers):**
- Layer 1: Native MD files (MEMORY.md ~2200 chars, USER.md ~1375 chars) — always in prompt, no retrieval needed
- Layer 2: Pluggable MemoryProvider — 8 official options:
  - **Mem0** (58K⭐) — universal memory layer, $24M raised, managed cloud option
  - **Supermemory** (26K⭐) — sub-300ms recall, official provider
  - **OpenViking** (25K⭐) — context database, filesystem-based, official provider
  - **Hindsight** (16K⭐) — retain/recall/reflect workflows, learns from experience
  - **Honcho** (5K⭐) — stateful entities (people, groups, projects), official provider
  - **ByteRover** (4.8K⭐) — memory as git repo, 5-tier retrieval
  - **Mnemosyne** (1K⭐) — zero-dependency, sub-millisecond, fully local SQLite
  - **GBrain** (22K⭐) — knowledge graph + synthesis layer, markdown vault
- Layer 3: Community plugins — Mnemosyne (local, tiered), GBrain (world facts)

**Memory caps are deliberate** — force consolidation, prevent stake context. Increase via `hermes config set memory.memory_char_limit X` if needed (won't burn tokens — only affects stored size, not injected size).

**Skills Strategy:**
- 643 community skills exist — don't install all
- Start with: LLM Wiki (builtin), gstack (trusted), security-audit (official)
- Skills auto-generate after ~5 tool calls on a pattern or user correction
- Four trust tiers: Builtin → Official → Trusted → Community

**Multi-Agent Options:**
- **Mission Control** (5.2K⭐) — self-hosted fleet orchestration + spend monitoring
- **Agent of Empires** (2.5K⭐) — TUI/web control surface for multiple agents
- **SwarmClaw** (557⭐) — autonomous agent swarms with orchestration

**Deployment Pattern:**
- Docker + systemd is production standard (hermes-autonomous-server template)
- Helm chart available for Kubernetes
- Nix package available for reproducible deployments

**Hermes Learning Loop:**
- Retrospective after every task → writes skills/memory → compounds
- "Harness Engineering" framing: LLM is replaceable, real engineering is in the 5 layers around it (instruction, constraint, feedback, memory, orchestration)
- Skills > facts: Hermes stores *procedures*, not just *facts*

### Action Items
- [ ] Consider installing Mem0 or Hindsight as Layer 2 memory provider for richer recall
- [ ] Install gstack and security-audit skills
- [ ] Increase memory_char_limit to 5000 (requires manual edit to config.yaml)
- [ ] Explore Mission Control if we scale to multi-agent setup

## Next Steps
- Deploy backend to VPS with PostgreSQL
- Deploy frontend to GitHub Pages or Cloudflare
- Set up CI/CD pipeline
- Add SSL certificate for custom domain
- Implement Phase 2 features (photo gallery manager, theme customization, analytics)
- Clean up duplicate DB entries from double migration run
