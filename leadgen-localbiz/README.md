# Lead Gen Pipeline — Phase 1

Automated 24/7 lead generation for local Jakarta businesses.

## Goal
Find businesses with good Google reviews but no/bad website → build them a free landing page preview → send pitch email → convert to paying client.

## Architecture

```
Cron (daily schedule)
  │
  ├─ 9:00am  Stage 1: SCOUT      → Find businesses (DuckDuckGo + Yelp)
  ├─ 9:30am Stage 2: EVALUATE    → Firecrawl sites, score quality
  ├─10:00am Stage 3: ENRICH      → Scrape Google Business, geocode
  ├─12:00pm Stage 4: BUILD      → Generate landing page from template
  ├─12:00pm Stage 5: DEPLOY     → Push GitHub + Cloudflare Pages
  ├─ 3:00pm Stage 6: OUTREACH   → Send personalized pitch email
  └─ 8:00am Stage 7: TRACK      → Daily summary → Discord
  
All data stored in PostgreSQL (Docker: leadgen-postgres, port 5433)
```

## URL Pattern
```
https://client.arjism.com/business-name-slug
```

## Database Tables
| Table | Purpose |
|-------|---------|
| `businesses` | Raw scouted business data |
| `website_audits` | Website quality scores (0-100) |
| `business_enrichment` | Photos, hours, geocode, services |
| `landing_pages` | Built sites, GitHub repos, deploy status |
| `outreach` | Pitch email tracking |
| `lead_responses` | Client replies |
| `pipeline_runs` | Daily run logs |

## Setup Checklist
- [x] PostgreSQL schema created
- [x] Config files written
- [x] Stage 1-2, 6-7 scripts written
- [ ] Stage 3 (Enrich) script — TODO
- [ ] Stage 4 (Build) script — TODO
- [ ] Stage 5 (Deploy) script — TODO
- [ ] Configure SMTP credentials
- [ ] Configure Discord webhook
- [ ] Create `client.arjism.com` Cloudflare DNS + Pages
- [ ] Install cron: `crontab cron/schedule`

## Costs
| Item | Cost |
|------|------|
| PostgreSQL | Free (Docker on VPS) |
| Firecrawl | $0 (existing alpha/beta key?) |
| DuckDuckGo search | Free |
| Yelp Fusion API | Free (5k/day) |
| GitHub | Free (public repos) |
| Cloudflare Pages | Free tier (50 sites) |
| Zoho Mail SMTP | Free (5 accounts) |
| **TOTAL** | **$0/month** |

## Capacity
- ~20-50 businesses scouted/day
- ~10-25 qualified/day (no/bad site)
- ~5-10 sites built/day (adjust via DAILY_BUILD_LIMIT)
- ~5-10 pitches/day (adjust via DAILY_OUTREACH_LIMIT)
