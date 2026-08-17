# Lead Gen Pipeline — Implementation Plan

## Tech Stack
- **Language:** Python 3.11+
- **Database:** PostgreSQL 16 (Docker, port 5433)
- **Search:** DuckDuckGo (free), Yelp Fusion API (free tier)
- **Scraping:** Firecrawl
- **Email:** Zoho Mail SMTP (free tier)
- **Hosting:** GitHub repos + Cloudflare Pages (free tier)
- **Notifications:** Discord webhook
- **Orchestration:** Hermes cron jobs

## Database Schema (PostgreSQL)

```sql
CREATE TABLE businesses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    slug VARCHAR(200) UNIQUE NOT NULL,
    category VARCHAR(100),
    address TEXT,
    city VARCHAR(100) DEFAULT 'Jakarta',
    phone VARCHAR(20),
    website_url TEXT,
    google_place_id VARCHAR(100),
    google_rating DECIMAL(3,2),
    review_count INTEGER DEFAULT 0,
    yelp_id VARCHAR(100),
    source VARCHAR(50), -- 'duckduckgo' | 'yelp'
    status VARCHAR(20) DEFAULT 'scouted',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE website_audits (
    id SERIAL PRIMARY KEY,
    business_id INTEGER REFERENCES businesses(id),
    url TEXT,
    quality_score INTEGER, -- 0-100
    has_mobile BOOLEAN,
    has_ssl BOOLEAN,
    load_time_ms INTEGER,
    issues JSONB,
    audited_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE business_enrichment (
    id SERIAL PRIMARY KEY,
    business_id INTEGER REFERENCES businesses(id),
    photos JSONB, -- array of photo URLs
    hours JSONB, -- opening hours
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    services JSONB,
    enriched_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE landing_pages (
    id SERIAL PRIMARY KEY,
    business_id INTEGER REFERENCES businesses(id),
    repo_url TEXT,
    cloudflare_url TEXT,
    deploy_status VARCHAR(20) DEFAULT 'pending',
    built_at TIMESTAMPTZ,
    deployed_at TIMESTAMPTZ
);

CREATE TABLE outreach (
    id SERIAL PRIMARY KEY,
    business_id INTEGER REFERENCES businesses(id),
    email_to TEXT NOT NULL,
    subject TEXT,
    body TEXT,
    sent_at TIMESTAMPTZ,
    opened BOOLEAN DEFAULT FALSE,
    replied BOOLEAN DEFAULT FALSE
);

CREATE TABLE lead_responses (
    id SERIAL PRIMARY KEY,
    business_id INTEGER REFERENCES businesses(id),
    response_type VARCHAR(20), -- 'interested' | 'not_interested' | 'no_reply'
    message TEXT,
    received_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE pipeline_runs (
    id SERIAL PRIMARY KEY,
    run_date DATE NOT NULL,
    stage VARCHAR(50),
    businesses_scouted INTEGER DEFAULT 0,
    businesses_qualified INTEGER DEFAULT 0,
    sites_built INTEGER DEFAULT 0,
    emails_sent INTEGER DEFAULT 0,
    errors INTEGER DEFAULT 0,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ
);
```

## Pipeline Stages

### Stage 1: SCOUT
- Search DuckDuckGo for "best [category] in Jakarta"
- Query Yelp Fusion API for businesses with 4+ stars
- Filter out businesses with existing good websites
- Insert into `businesses` table
- **Status:** ✅ Implemented

### Stage 2: EVALUATE
- Firecrawl each business website URL
- Score quality 0-100 based on: mobile, SSL, load time, design, content
- Insert into `website_audits` table
- **Status:** ✅ Implemented

### Stage 3: ENRICH
- Scrape Google Business Profile for photos, hours, services
- Geocode addresses via OpenStreetMap/OSRM
- Insert into `business_enrichment` table
- **Status:** ❌ TODO

### Stage 4: BUILD
- Generate landing page from HTML template
- Populate with enriched data (name, photos, hours, contact)
- Save to `website/business-slug/index.html`
- **Status:** ❌ TODO

### Stage 5: DEPLOY
- Create GitHub repo (or push to existing)
- Trigger Cloudflare Pages deploy
- Update `landing_pages` table
- **Status:** ❌ TODO

### Stage 6: OUTREACH
- Send personalized pitch email via Zoho SMTP
- Include link to preview landing page
- Track in `outreach` table
- **Status:** ✅ Implemented

### Stage 7: TRACK
- Aggregate daily stats from all tables
- Send summary to Discord webhook
- Log to `pipeline_runs` table
- **Status:** ✅ Implemented

## Task Breakdown

### Phase 1: Foundation ✅
- [x] PostgreSQL schema created
- [x] Config files written
- [x] Docker Compose for PostgreSQL
- [x] Master pipeline runner (`run_pipeline.py`)
- [x] Stage 1 (Scout) script
- [x] Stage 2 (Evaluate) script
- [x] Stage 6 (Outreach) script
- [x] Stage 7 (Track) script

### Phase 2: Core Pipeline — TODO
- [ ] Stage 3 (Enrich) — Google Business scraping + geocode
- [ ] Stage 4 (Build) — Landing page template + generator
- [ ] Stage 5 (Deploy) — GitHub API + Cloudflare Pages integration

### Phase 3: Polish — TODO
- [ ] Configure SMTP credentials
- [ ] Configure Discord webhook
- [ ] Create `client.arjism.com` Cloudflare DNS + Pages
- [ ] Install cron: `crontab cron/schedule`
- [ ] Rate limiting + error handling
- [ ] Duplicate detection (already scouted businesses)
- [ ] Email template A/B testing

### Phase 4: Scale — TODO
- [ ] Multi-city support (Bogor, Depok, Tangerang, Bekasi)
- [ ] Category expansion (beyond salons/barbershops)
- [ ] Follow-up email sequences
- [ ] CRM integration
- [ ] Client onboarding flow

## Cron Schedule
```
# /root/hermes/leadgen-localbiz/cron/schedule
0 9 * * *   scout
30 9 * * *  evaluate
0 10 * * *  enrich
0 12 * * *  build
0 12 * * *  deploy
0 15 * * *  outreach
0 8 * * *   track (next day summary)
```

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
