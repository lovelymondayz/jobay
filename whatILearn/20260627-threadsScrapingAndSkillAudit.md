# 2026-06-27 — Threads Scraping Deep Dive & Full Skill Audit

## What Happened

### 1. Threads Post Content Investigation (Discord, 139 messages)
**Trigger:** kvinn shared a Threads URL (`@bevan.satriaa/post/DaDC2x4D033`) asking "what is this"

**Approaches attempted (all failed):**
- `curl` with browser User-Agent → returned only CSS custom properties + JS resource maps, zero readable text
- `browser_navigate` → Chrome crashed with stack trace on Threads specifically
- Threads API endpoints → all returned JS/CSS only, no JSON data
- Google search → captcha blocked
- DuckDuckGo HTML → empty results
- Nitter → not applicable (Threads, not Twitter)
- archive.org Wayback → no snapshots for this post

**What finally worked:**
- Used `curl` on a GitHub search result page to find the repo `Vanszs/qwencloud-generator`
- Extracted the post's meaning from the GitHub README + live dashboard screenshot embedded in the post

**Result:** The post was about a "Qwencloud API Key Generator" — a free tool that auto-generates Qwen (Alibaba AI) API keys. kvinn flagged it as sketchy.

**Key learning:** Threads (Meta) is fully client-side rendered. The SSR HTML contains only CSS variables, JS resource maps, and base64-encoded QPL (Quick Performance Logging) data blocks that are encrypted/obfuscated. No amount of curl/browser automation can extract post content without a working headless browser that doesn't crash.

### 2. Full Skill Audit (same session)
**Trigger:** kvinn asked to see all active skills and archived ones.

**Findings:**
- **21 active skill categories** including: apple, autonomous-ai-agents, creative, data-science, devops, diagramming, domain, email, gifs, github, inference-sh, mcp, media, mlops, note-taking, productivity, research, smart-home, social-media, software-development
- **5 archived skills:** xurl, r3f-flipbook, tailwind-vite-css, yuanbao, hermes-terminal-patterns
- The `web-content-extraction` skill (under `research/`) contains a `threads-scraping.md` reference documenting all known Threads limitations — last tested today

### 3. Cloudflare IP Whitelist Update (Cron, 04:00 AM)
- Script: `/usr/local/bin/update-cloudflare-ips.sh`
- Result: ✅ Success — nginx config test passed, IP ranges updated

### 4. Hermes Workspace Backup (Cron, 19:00 PM)
- Script: `python3 /root/hermes/scripts/hermes-backup.py`
- Result: ✅ `backup-2026-06-28.zip` (35,733 KB) uploaded successfully
- 0 old backups pruned (none older than 30 days)

### 5. Daily Learning Log Cron (Cron, 19:15 PM)
- This cron job itself — attempted but produced no output (the job that created *this* log was today's run on June 28)

## Tech Stack Used
- `curl` with various User-Agents for web scraping
- `grep`/`sed` for HTML parsing
- `browser_navigate` (failed — Chrome crash on Threads)
- DuckDuckGo HTML search endpoint
- GitHub search for cross-referencing
- Python (pytesseract + Pillow) for OCR (referenced from prior sessions)
- nginx config validation (`nginx -t`)
- Hermes backup script (Python, GitHub API via `lovelymondayz` token)

## Mistakes Overcome
1. **Browser crash on Threads** — pivoted to curl + search engine cross-referencing
2. **Google captcha** — switched to DuckDuckGo HTML endpoint
3. **grep lookbehind assertion error** — fixed regex pattern (Perl-style lookbehind not supported in all grep builds)
4. **Threads API returning only JS/CSS** — discovered there is no public Threads API; all content is CSR-only

## Key Learnings
- **Threads is un-scrapable** without a working headless browser. Documented in `research/web-content-extraction/references/threads-scraping.md`
- **QPL data blocks** in Threads SSR are encrypted — parsing JSON yields only component manifests, not user content
- **Skill library has 21 active categories** with 5 archived — good time to consider restoring `xurl` if URL unshortening is needed again
- **Backup size holding steady** at ~35.7 MB — no storage pressure
