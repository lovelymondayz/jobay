# 2026-07-01: Backup Script Fix and Large File Handling

## What Was Done Today

### 1. Hermes Daily Backup Failure Diagnosis & Fix
- **Problem:** Backup failed with HTTP 422 — archive too large (45 MB)
- **Root cause:** Workspace bloated with:
  - Go server binaries (20M each in members, wedding, yearbook)
  - `dist/` frontend builds in multiple projects
  - `.venv` directories with compiled libs
- **Fix:** Updated `hermes-backup.py` to skip:
  - `.venv`, `node_modules`, `__pycache__` (already skipped)
  - `dist/` directories (frontend build artifacts - reproducible)
  - `.so`, `.so.*`, `.pyc` files (compiled binaries)
  - `server` binaries (Go build artifacts)
- **Result:** Archive reduced to 13.8 MB, upload succeeded

### 2. whatILearn Cron Execution
- Daily learning log automation ran normally
- Captured all work done in prior session (backup fix)

### 3. Discord Status Check
- Morning check-in from kvinn ("yo", "yp")
- API connectivity issues reported (HTTP 503, 404 errors) - transient provider issues

## Key Learnings

| Lesson | Detail |
|--------|--------|
| GitHub API file limit | ~100MB for Contents API uploads. Exceeding causes 422 error. |
| Source vs build exclusion | Exclude compiled artifacts (`.so`, binaries, `dist/`) from backups — they're reproducible and bloat size by 3x. |
| Script verification | Running backup manually before automated cron catches size issues early. |

## Tech Stack

No changes. Existing stack unchanged:
- **Backup:** `hermes-backup.py` → GitHub (`BackupHermes` repo, `backupHermesDaily/` folder)
- **Cronjobs:** 2 active (backup 19:00 UTC, whatILearn 19:15 UTC)

## Mistakes Overcome

1. **Backup script didn't skip dist/** — Added `/dist/` path filtering to reduce 45MB → 13.8MB

## Context

- All prior work (Members Project, lead gen pipeline, etc.) stable and deployed
- Backup script now handles growing workspace size correctly