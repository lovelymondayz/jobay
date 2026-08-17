# Lead Gen Pipeline — Architecture

## Overview
Automated 24/7 lead generation for local Jakarta businesses. Finds businesses with good Google reviews but no/bad website → builds them a free landing page preview → sends pitch email → converts to paying client.

## Tech Stack
- **Language:** Python 3.11+
- **Database:** PostgreSQL 16 (Docker: leadgen-postgres, port 5433)
- **Search:** DuckDuckGo (free), Yelp Fusion API (free tier 5k/day)
- **Scraping:** Firecrawl
- **Email:** Zoho Mail SMTP (free tier)
- **Hosting:** GitHub repos + Cloudflare Pages (free tier, 50 sites)
- **Notifications:** Discord webhook
- **Orchestration:** Hermes cron jobs

## Architecture

```
Hermes Cron (daily schedule)
  │
  ├─ 9:00am  Stage 1: SCOUT      → Find businesses (DuckDuckGo + Yelp)
  ├─ 9:30am Stage 2: EVALUATE    → Firecrawl sites, score quality (0-100)
  ├─10:00am Stage 3: ENRICH      → Scrape Google Business, geocode
  ├─12:00pm Stage 4: BUILD      → Generate landing page from template
  ├─12:00pm Stage 5: DEPLOY     → Push GitHub + Cloudflare Pages
  ├─ 3:00pm Stage 6: OUTREACH   → Send personalized pitch email
  └─ 8:00am Stage 7: TRACK      → Daily summary → Discord

All data stored in PostgreSQL (Docker: leadgen-postgres, port 5433)
```

## Directory Structure
```
leadgen-localbiz/
├── run_pipeline.py            # Master pipeline runner (all 7 stages)
├── Makefile
├── ARCHITECTURE.md
├── PLAN.md
├── docker-compose.yml         # PostgreSQL only
├── README.md                  # Setup checklist + costs
├── .venv/                     # Python virtual environment
├── config/                    # Configuration files
├── scraper/                   # Web scraping modules
├── scripts/                   # Stage scripts
│   ├── scout/                 # Stage 1: Business discovery
│   ├── evaluate/              # Stage 2: Website quality scoring
│   ├── enrich/                # Stage 3: Data enrichment
│   ├── build/                 # Stage 4: Landing page generation
│   ├── deploy/                # Stage 5: GitHub + Cloudflare deploy
│   ├── outreach/              # Stage 6: Email pitch sending
│   └── track/                 # Stage 7: Daily reporting
├── templates/                 # Landing page templates
├── email/                     # Email templates
├── assets/                    # Static assets
├── website/                   # Built landing pages (per business)
│   ├── business-name-slug/    # Each business gets a directory
│   │   └── index.html         # Generated landing page
│   └── ...
├── cron/                      # Cron schedule files
└── logs/                      # Pipeline run logs
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

## URL Pattern
```
https://client.arjism.com/business-name-slug
```

## Pipeline Stages

### Stage 1: SCOUT
- Search DuckDuckGo + Yelp for local Jakarta businesses
- Filter: good reviews (4+ stars), no/bad website
- Insert into `businesses` table

### Stage 2: EVALUATE
- Firecrawl each business website
- Score quality 0-100 (mobile, speed, design, content)
- Insert into `website_audits` table

### Stage 3: ENRICH
- Scrape Google Business Profile for photos, hours, services
- Geocode addresses
- Insert into `business_enrichment` table

### Stage 4: BUILD
- Generate landing page from template
- Populate with enriched data
- Save to `website/` directory

### Stage 5: DEPLOY
- Create/push GitHub repo
- Trigger Cloudflare Pages deploy
- Update `landing_pages` table

### Stage 6: OUTREACH
- Send personalized pitch email via Zoho SMTP
- Include link to preview landing page
- Track in `outreach` table

### Stage 7: TRACK
- Aggregate daily stats
- Send summary to Discord
- Log to `pipeline_runs` table

## Costs
| Item | Cost |
|------|------|
| PostgreSQL | Free (Docker on VPS) |
| Firecrawl | $0 (alpha/beta key) |
| DuckDuckGo search | Free |
| Yelp Fusion API | Free (5k/day) |
| GitHub | Free (public repos) |
| Cloudflare Pages | Free tier (50 sites) |
| Zoho Mail SMTP | Free (5 accounts) |
| **TOTAL** | **$0/month** |

## Capacity
- ~20-50 businesses scouted/day
- ~10-25 qualified/day (no/bad site)
- ~5-10 sites built/day (adjust via `DAILY_BUILD_LIMIT`)
- ~5-10 pitches/day (adjust via `DAILY_OUTREACH_LIMIT`)

## Design Decisions
- **Python over Go** — Faster iteration for scraping/scripting pipeline
- **PostgreSQL on port 5433** — Avoids conflict with yearbook DB on 5432
- **Template-based landing pages** — Simple HTML, no JS framework needed
- **Free-tier everything** — $0/month operating cost
- **Discord for tracking** — Real-time pipeline visibility
