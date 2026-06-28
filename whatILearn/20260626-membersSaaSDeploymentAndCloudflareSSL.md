# 2026-06-26: Members SaaS Deployment & Cloudflare SSL Custom Hostname

## Summary
Full SaaS membership platform ("members") built from scratch and deployed to VPS. Configured Cloudflare for SaaS with custom hostname SSL, HTTPS reverse proxy, and IP whitelist security. Also redesigned mentengdutch to Jakarta Munch style and researched Cloudflare for SaaS pricing.

---

## What Was Built

### 1. Members SaaS Platform (New Project)
Multi-tenant membership platform with 3 roles (Super Admin, Admin, Member). Super Admin creates Admin → auto-creates Store → Admin manages Members. Includes invoice + payment tracking.

**Tech Stack:** Go 1.25 + GORM + PostgreSQL (backend), React + Vite + TS + Tailwind (frontend), Docker, Nginx

**Features implemented:**
- Email/password (bcrypt) + Google OAuth2 + JWT auth
- RBAC middleware (super_admin, admin, member)
- 6 DB models: Role, User, Store, Member, Invoice, Payment
- Repository → Service → Handler architecture
- DTOs decoupling API contracts from DB models
- Admin onboarding flow (1 Admin = 1 Store)
- Member CRUD with store scoping
- Membership card with client-side QR code
- Invoice CRUD + Payment recording
- React frontend: Login, Dashboard, Members, MemberCard, Invoices pages

**Deployed:** db(5435), backend(8082), frontend(3003) — all containers Healthy

### 2. Cloudflare for SaaS + Custom SSL
- Added `members.arjism.com` as Custom Hostname in Cloudflare
- Configured nginx reverse proxy with IP whitelist (allow Cloudflare IPs only, deny all others by default 403)
- Set up cron job to auto-update Cloudflare IP whitelist daily at 04:00 UTC
- Troubleshot Terraform certificate validation records (TXT records for `_acme-challenge` and `_cf-custom-hostname`)

### 3. Mentengdutch Redesign
- Redesigned from dark theme to Jakarta Munch bright style
- Color palette: warm cream `#F5F2EB` bg, `#ED3F1C` accent, `#2f2f2f` text
- All components updated: Hero, Header, About, Reviews, Videos, Shop, Footer, etc.

### 4. Cloudflare for SaaS Pricing Research
- Free/Pro/Business: 100 custom hostnames included, $0.10/month per additional hostname
- Enterprise: custom pricing, unlimited hostnames, paid add-ons (Apex Proxying, BYOIP, Custom Metadata, mTLS)
- Billing is monthly per hostname (prorated for partial months)

---

## Key Learnings

### GORM Nested structs cause circular FK constraints
When GORM models embed other structs (e.g., `Role Role`, `Store Store` inside User), it tries to create circular foreign key relationships → migration fails. **Fix:** Use only FK ID fields (e.g., `RoleID uint`), do full joins manually via queries.

### PostCSS config: .js vs .cjs in ESM projects
When `package.json` has `"type": "module"`, Node treats `.js` files as ESM. `require()` in `postcss.config.js` fails. **Fix:** Rename to `postcss.config.cjs` with `module.exports`.

### Go build inside Docker is slow
First-time Go builds inside Docker take 5-6 minutes due to dependency downloads. Subsequent builds are faster due to layer caching. Always use multi-stage builds.

### Cloudflare Custom Hostname requires TXT validation
Adding a custom hostname to Cloudflare requires adding two DNS TXT records (`_acme-challenge` and `_cf-custom-hostname`) at your DNS provider. Without these, certificate stays pending and HTTPS fails.

### Cloudflare IP whitelist blocks localhost
When you restrict nginx to Cloudflare IPs only, tests from the host (localhost/127.0.0.1) get 403. You must explicitly `allow 127.0.0.1` and `::1` in the whitelist.

### Bcrypt hash truncation issue
When inserting bcrypt hashes via string concatenation in SQL, special characters can get truncated. The hash lost 2 chars (`$2a$10$caXPsv6...` was missing last 2 chars), causing all password checks to fail. **Lesson:** Always verify the stored hash is exactly 60 characters.

---

## Mistakes & Fixes

| Issue | Root Cause | Fix |
|-------|-----------|-----|
| GORM migration FK error | Nested struct fields in models created circular FKs | Removed embedded structs, kept FK ID fields |
| Package name mismatch | `package model` vs `models` imports | Renamed to `models` |
| PostCSS build error | `require()` in ESM context | Renamed to `.cjs` container |
| Login 401 (host) | Bcrypt hash truncated in DB | Verified stored hash length = 60 chars |
| Login 403 (host) | Cloudflare IP whitelist blocked localhost | Added `allow 127.0.0.1` |
| Bcrypt hash verification confusion | Double-hash verification needed | Used direct `bcrypt.CompareHashAndPassword` |

---

## Tech Stack Summary

```
members SaaS:
  Backend:  Go 1.25, GORM, PostgreSQL 16, JWT, bcrypt
  Frontend: React 18, Vite 5, Tailwind CSS 3, TypeScript
  Deploy:   Docker Compose (3 services), Nginx reverse proxy, Cloudflare DNS
  Ports:    db=5435, api=8082, fe=3003

Cloudflare security:
  - IP whitelist on nginx (Cloudflare IPs + localhost only)
  - Daily auto-update cron: /usr/local/bin/update-cloudflare-ips.sh
  - SSL: Custom Hostname edge certificate (pending TXT validation on 2026-06-26)

Mentengdutch redesign:
  - React + Vite + Tailwind
  - Theme: Jakarta Munch bright style
  - Live on port 3002 (jakarta-munch container)
```

---

## What's Next / Blocked

- [ ] TXT records for SSL certificate — waiting for kvinn to add them
- [ ] Cloudflare SSL mode change to "Full"
- [ ] Super Admin user officially set up on the platform
- [ ] Google Auth credentials — blank in `.env.production`
- [ ] Member card restore (store info removed with nested struct removal)
