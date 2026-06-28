# 2026-06-25 Daily Learning Log

## Summary
Two sessions today: (1) Researching the `no-mistakes` git proxy tool, (2) Jakarta Munch project status update — swapping the logo to Toko Menteng branding and creating per-project update scripts.

---

## Session 1: no-mistakes Git Proxy Research
**Time:** 08:47 AM | **Type:** Research / Tool Discovery

### What was done
- kvinn shared a GitHub repo: [kunchenguid/no-mistakes](https://github.com/kunchenguid/no-mistakes)
- Investigated the tool by reading the README via curl (browser and gh CLI were blocked/slow)

### Key Learnings
- **no-mistakes** is a local git proxy that gates pushes through an AI-driven validation pipeline
- Workflow: `git push no-mistakes` → disposable worktree → review → test → docs → lint → push → PR → CI
- Agent-agnostic: works with Claude, Codex, Rovodev, OpenCode, Pi, ACP
- Human-in-the-loop: mechanical fixes auto-apply; intent-changing findings escalated for approval
- Install: `curl -fsSL https://raw.githubusercontent.com/kunchenguid/no-mistakes/main/docs/install.sh | sh`

### Mistakes Overcome
- `gh repo view` timed out (gh auth not configured for API calls)
- Browser crashed (Chrome DevTools port failure on VPS)
- Fell back to `curl` on raw README.md — worked fine

---

## Session 2: Jakarta Munch → Toko Menteng Rebrand + Update Scripts
**Time:** 02:59 AM (long session, 260 messages) | **Type:** Feature Update + DevOps

### What was done
1. **Logo swap:** Replaced Jakarta Munch text logo with Toko Menteng's actual logo downloaded from `tokomenteng.nl`
   - Downloaded logo to `/root/hermes/mentengdutch/public/toko-menteng-logo.png`
   - Patched `Header.tsx` to use `<img>` tag with the logo instead of text
   - Updated `index.html` to reference Toko Menteng favicon
2. **Created per-project update scripts** (`scripts/update.sh`) for all 3 projects:
   - `mentengdutch` (port 3002)
   - `wedding-invitation` (ports 3000 + 8080)
   - `digital-yearbook` (ports 3001 + 8081)
3. **Created `project-update-script` skill** in `~/.hermes/skills/devops/project-update-script/SKILL.md`

### Tech Stack
- React + TypeScript + Tailwind CSS (frontend)
- Docker Compose (deployment)
- Shell scripts (bash) for update automation
- Makefiles per project (`make deploy`)

### Key Learnings
- The project directory was renamed from `jakarta-munch` to `mentengdutch` — reflects the rebrand
- Update script pattern: `git fetch` → compare local/remote → rebuild containers if changed → show status
- Each project follows the same convention: `scripts/update.sh` + `Makefile` with `deploy` target

### Mistakes Overcome
- Browser kept crashing when trying to scrape tokomenteng.nl for logo URL
- Fell back to curl + grep to find the logo image URL in the HTML source
- Direct download of the PNG from `tokomenteng.nl/uploads/...` worked cleanly

---

## Projects Status
| Project | Status | Port(s) |
|---------|--------|---------|
| mentengdutch (formerly jakarta-munch) | ✅ Running, logo updated | 3002 |
| wedding-invitation | ✅ Has update script | 3000 + 8080 |
| digital-yearbook | ✅ Has update script | 3001 + 8081 |

## Tomorrow's Priorities
- [ ] Decide whether to install `no-mistakes` for our repos (kvinn to confirm)
- [ ] Test the update scripts actually work end-to-end
- [ ] Consider pushing mentengdutch to GitHub (currently only local, single "first commit")
