
"""
Stage 7: Track — Daily summary + Discord alerts.
Runs once per day, sends summary to Discord webhook.
"""
import sys, json, time
sys.path.insert(0, "/root/hermes/leadgen-localbiz")
from config.settings import *
from config.database import query, execute

def send_discord(msg):
    """Send message to Discord webhook."""
    if not DISCORD_WEBHOOK:
        print("  No Discord webhook configured")
        return
    data = json.dumps({"content": msg}).encode()
    req = urllib.request.Request(DISCORD_WEBHOOK, data=data, method="POST")
    req.add_header("Content-Type", "application/json")
    try:
        import urllib.request
        urllib.request.urlopen(req, timeout=10)
    except Exception as e:
        print(f"  Discord error: {e}")

def main():
    print(f"=== STAGE 7: TRACK — {time.strftime('%Y-%m-%d %H:%M')} ===")
    
    # Gather stats
    stats = query("SELECT COUNT(*) as cnt FROM businesses")[0]["cnt"]
    audited = query("SELECT COUNT(*) as cnt FROM website_audits")[0]["cnt"]
    qualified = query("SELECT COUNT(*) as cnt FROM website_audits WHERE overall_score < 40")[0]["cnt"]
    sites_built = query("SELECT COUNT(*) as cnt FROM landing_pages WHERE status = 'live'")[0]["cnt"]
    pitches_sent = query("SELECT COUNT(*) as cnt FROM outreach WHERE status = 'sent'")[0]["cnt"]
    replies = query("SELECT COUNT(*) as cnt FROM lead_responses")[0]["cnt"]
    
    # Top qualified leads
    top_leads = query("""
        SELECT b.name, b.category, b.rating, b.review_count, lp.live_url
        FROM businesses b
        JOIN website_audits wa ON b.id = wa.business_id
        JOIN landing_pages lp ON b.id = lp.business_id
        WHERE wa.overall_score < 40 AND lp.status = 'live'
        ORDER BY b.rating DESC, b.review_count DESC
        LIMIT 5
    """)
    
    # Build summary
    lines = [
        f"📊 **Lead Gen Pipeline Daily Report** — {time.strftime('%Y-%m-%d')}",
        f"",
        f"| Metric | Count |",
        f"|--------|-------|",
        f"| Total businesses scouted | {stats} |",
        f"| Websites audited | {audited} |",
        f"| Qualified leads (score <40) | {qualified} |",
        f"| Sites built & live | {sites_built} |",
        f"| Pitches sent | {pitches_sent} |",
        f"| Client replies | {replies} |",
        f"",
    ]
    
    if top_leads:
        lines.append("🔥 **Top Qualified Leads:**")
        for lead in top_leads:
            url_str = f" | {lead['live_url']}" if lead['live_url'] else ""
            lines.append(f"  • {lead['name']} ({lead['category']}) — {lead['rating']}★ {lead['review_count']} reviews{url_str}")
    
    summary = chr(10).join(lines)
    print(summary)
    
    send_discord(summary)
    
    # Log pipeline run
    execute("INSERT INTO pipeline_runs (stage, status, targets_found, targets_qualified, targets_built, targets_contacted, finished_at) VALUES (%s, %s, %s, %s, %s, %s, NOW())",
            ("track", "success", stats, qualified, sites_built, pitches_sent))
    
    print(f"\n=== DONE: Report sent ===")

if __name__ == "__main__":
    main()
