# Wedding Invitation — Multi-Platform Music, Login Fixes & Countdown Repair

**Date:** August 16, 2026
**Project:** Wedding Invitation (`lovelymondayz/wedding-invitation`)
**Stack:** Go + GIN + pgx + PostgreSQL | React + Vite + TS + Tailwind + Framer Motion
**Sessions:** 1 major (337 messages) + backup cron

---

## What Was Built / Fixed

### 1. Couple Admin Login Fails (401)
**Problem:** Couple admin `Kambing` couldn't log in. Backend returned `401 Unauthorized`.

**Diagnosis:**
- Password was stored as plaintext "Kambing" (same as username) — not a valid bcrypt hash
- `couple_slug` lookup in login response used `_ =` (ignored errors) — could return empty slug

**Fix:**
- Generated proper bcrypt hash via Go script, updated DB
- Added distinct error codes: `invalid_credentials` vs `wrong_password`
- Frontend: animated inline error banner, red border on error fields, empty-field validation
- Added `/api/auth/me` endpoint that returns `couple_slug` from token for robust redirects

**Commit:** `293daa4`

---

### 2. Multi-Platform Music Support (Spotify, YouTube, SoundCloud)
**Problem:** Couples pasted Spotify links but nothing played. Root cause: `<audio>` element only supports direct MP3/WAV — not streaming URLs.

**Solution:** Built a source-aware music system that dispatches to the correct player.

| Source | Player | Conversion |
|--------|--------|------------|
| Direct MP3 | `<audio>` element | None needed |
| Spotify | Embed iframe | `open.spotify.com/track/ID` → `open.spotify.com/embed/track/ID` |
| YouTube | Embed iframe | Extract video ID → `youtube.com/embed/ID` |
| SoundCloud | Embed iframe | `w.soundcloud.com/player/?url=...` |
| Vimeo | Embed iframe | Standard vimeo embed |

**Backend:**
- New migration `003_music_source.sql` — adds `source` column to `music_tracks`
- `detectMusicSource(url string) string` — auto-detects from URL
- `AdminCreateMusicHandler` — accepts `source` field, stores it
- `GetActiveMusicHandler` — returns `source` in response

**Frontend:**
- **New `MusicEmbed.tsx`** — iframe-based player with postMessage play/pause
- **Updated `useMusic.ts`** — tracks `source` type, exposes to components
- **Updated all 3 templates** — dispatch to `MusicPlayer` (direct) or `MusicEmbed` (streaming)
- **Updated `MusicManagement`** — source selector dropdown, source labels with icons

**Commit:** `817805a`

---

### 3. Template C Countdown — Hardcoded Placeholders
**Problem:** Template C (Dark Luxe) showed `...` instead of actual countdown.

**Root Cause:** Template C was built with a static HTML mockup — four cards with literal `{ label: 'Days', value: '...' }` text. It never rendered the `CountdownSection` component.

**Fix:** Replaced static grid with `<CountdownSection countdown={data.countdown ?? undefined} />` — same as Templates A and B.

**Commit:** `6c88750`

---

### 4. Migration Not Executed (Silent Failure)
**Problem:** `003_music_source.sql` was created but never executed. Docker's `/docker-entrypoint-initdb.d/` only runs on **first container creation** — not on restarts. Since the DB volume already existed, the migration was skipped.

**Fix:**
```sql
ALTER TABLE music_tracks ADD COLUMN IF NOT EXISTS source VARCHAR(20) DEFAULT 'default';
```

**Lesson:** Always verify migrations ran, not just that the file exists. For existing containers, apply manually or use a migration runner.

---

## Key Learnings

1. **`<audio>` ≠ universal music player.** Streaming services (Spotify, YouTube) require embed iframes — not `<audio>` elements. Always detect the source URL and dispatch to the right player.

2. **Docker init scripts are one-shot.** `/docker-entrypoint-initdb.d/` only fires when the volume is first created. For existing databases, migrations must be applied manually or via a proper migration tool.

3. **bcrypt hash generation is straightforward.** A 10-line Go script with `golang.org/x/crypto/bcrypt` can generate/verify hashes for debugging auth issues.

4. **Static mockups become bugs.** Template C's countdown was built as a visual mockup with hardcoded values — it never connected to the data. When building templates, ensure every dynamic element uses the shared component, not a static duplicate.

5. **postMessage for iframe control.** Spotify and YouTube embeds use `iframe.contentWindow.postMessage()` for play/pause — not direct JS API calls.

---

## Commits Today

| Commit | Description |
|--------|-------------|
| `293daa4` | Fix couple admin login (password hash, error feedback, `/api/auth/me`) |
| `817805a` | feat: multi-platform music support (Spotify, YouTube, SoundCloud) |
| `6c88750` | fix: Template C countdown was hardcoded with '...' placeholders |

---

## Live At

- **Homepage:** `https://wedding.arjism.com`
- **Template A:** `wedding.arjism.com/brian-roro-20260815`
- **Template A:** `wedding.arjism.com/kambing-kambing-20260816`
- **Template C:** `wedding.arjism.com/jackie-chan-raisha-20260816`

---

## Backup Status

✅ `backup-2026-08-17.zip` (14,221 KB) uploaded to `backupHermesDaily/`. 35 old backups pruned.