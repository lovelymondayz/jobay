"""
Stage 6: Outreach — Send personalized pitches via WhatsApp.
Reads from DB: qualified leads with landing pages deployed.
Generates WhatsApp deep links (wa.me) since we have phone numbers.
Falls back to logging for Instagram DMs.
"""
import sys, json, time, urllib.parse, urllib.request
sys.path.insert(0, "/root/hermes/leadgen-localbiz")
from config.settings import *
from config.database import query, execute

def craft_whatsapp_message(name, category, rating, review_count, live_url):
    """Generate a personalized WhatsApp pitch message."""
    category_display = {
        "restaurant": "restaurant",
        "gym": "gym/fitness center",
        "dental clinic": "dental clinic",
        "beauty salon": "beauty salon",
        "barber shop": "barbershop",
        "cafe": "cafe",
        "bakery": "bakery",
        "pet clinic": "pet clinic",
        "laundry": "laundry service",
        "tailor": "tailor shop",
    }.get(category, category)

    # Short, personal, Indonesian-style pitch
    msg = f"""Hi {name} team! 👋

I was looking for great {category_display}s in Jakarta and found yours on Google — {rating} stars with {review_count}+ reviews is impressive! ⭐

I noticed your business could use a better website, so I put together a free landing page preview:

{live_url}

It's mobile-friendly, loads fast, and has your WhatsApp button so customers can reach you easily.

If you like it, I can build the full site for you — affordable & quick. No strings attached on the preview! 🙏

Let me know what you think!
— {SMTP_FROM_NAME}"""
    return msg

def build_wa_link(phone, message):
    """Build WhatsApp deep link."""
    if not phone:
        return None
    # Normalize phone number
    wa_num = phone.replace("+", "").replace("-", "").replace(" ", "").replace("(", "").replace(")", "")
    if wa_num.startswith("0"):
        wa_num = "62" + wa_num[1:]
    elif not wa_num.startswith("62"):
        wa_num = "62" + wa_num
    encoded_msg = urllib.parse.quote(message)
    return f"https://wa.me/{wa_num}?text={encoded_msg}"

def send_discord_notification(msg):
    """Send notification to Discord webhook."""
    if not DISCORD_WEBHOOK:
        return
    data = json.dumps({"content": msg}).encode()
    req = urllib.request.Request(DISCORD_WEBHOOK, data=data, method="POST")
    req.add_header("Content-Type", "application/json")
    try:
        urllib.request.urlopen(req, timeout=10)
    except Exception as e:
        print(f"  Discord error: {e}")

def main():
    print(f"=== STAGE 6: OUTREACH — {time.strftime('%Y-%m-%d %H:%M')} ===")

    # Get qualified leads with deployed pages, not yet contacted
    rows = query("""
        SELECT b.id, b.name, b.category, b.rating, b.review_count, b.phone,
               b.instagram_url, b.facebook_url,
               lp.live_url, lp.slug, wa.audit_notes
        FROM businesses b
        JOIN website_audits wa ON b.id = wa.business_id
        JOIN landing_pages lp ON b.id = lp.business_id
        LEFT JOIN outreach o ON b.id = o.business_id
        WHERE wa.overall_score < 40
          AND lp.status = 'live'
          AND o.id IS NULL
        ORDER BY b.rating DESC, b.review_count DESC
        LIMIT %s
    """, (DAILY_OUTREACH_LIMIT,))

    print(f"Found {len(rows)} leads ready for outreach")
    sent = 0
    whatsapp_links = []

    for row in rows:
        name = row["name"]
        biz_id = row["id"]
        category = row["category"]
        rating = row["rating"]
        review_count = row["review_count"]
        phone = row["phone"]
        live_url = row["live_url"]
        ig_url = row.get("instagram_url")

        print(f"\n  📱 Pitching: {name} ({rating}★)")

        # Generate WhatsApp message
        message = craft_whatsapp_message(name, category, rating, review_count, live_url)
        wa_link = build_wa_link(phone, message)

        if wa_link:
            print(f"    WhatsApp: {phone}")
            print(f"    Link: {wa_link[:80]}...")
            whatsapp_links.append({
                "name": name,
                "phone": phone,
                "link": wa_link,
                "url": live_url
            })

            # Log outreach as 'wa_link_generated' — human clicks to send
            execute("""
                INSERT INTO outreach (business_id, channel, recipient_email, subject, message_body, status, sent_at)
                VALUES (%s, 'whatsapp', %s, %s, %s, 'wa_link_generated', NOW())
                ON CONFLICT DO NOTHING
            """, (biz_id, phone, f"WA pitch: {name}", message))
            sent += 1
        elif ig_url:
            # Fallback: Instagram
            print(f"    No phone, using Instagram: {ig_url}")
            execute("""
                INSERT INTO outreach (business_id, channel, recipient_email, subject, message_body, status, sent_at)
                VALUES (%s, 'instagram', %s, %s, %s, 'ig_logged', NOW())
                ON CONFLICT DO NOTHING
            """, (biz_id, ig_url, f"IG pitch: {name}", message))
            sent += 1
        else:
            print(f"    ❌ No contact info (no phone, no IG)")
            execute("""
                INSERT INTO outreach (business_id, channel, recipient_email, subject, message_body, status, sent_at)
                VALUES (%s, 'none', %s, %s, %s, 'no_contact', NOW())
                ON CONFLICT DO NOTHING
            """, (biz_id, "", f"No contact: {name}", ""))
            sent += 1

        time.sleep(2)  # Rate limit

    # Log pipeline run
    execute("INSERT INTO pipeline_runs (stage, status, targets_contacted, started_at, finished_at) VALUES (%s, %s, %s, NOW(), NOW())",
            ("outreach", "success", sent))

    # Print summary with clickable links
    print(f"\n{'='*60}")
    print(f"📱 WHATSAPP OUTREACH LINKS — Click to send:")
    print(f"{'='*60}")
    for item in whatsapp_links:
        print(f"\n  {item['name']} ({item['phone']})")
        print(f"  Site: {item['url']}")
        print(f"  WA:   {item['link']}")

    summary = f"📱 **Outreach Complete** — {sent} pitches ready\n" + \
              f"• WhatsApp links generated: {len(whatsapp_links)}\n" + \
              f"Click the links above to send via WhatsApp!"
    send_discord_notification(summary)

    print(f"\n=== DONE: {sent} pitches prepared ({len(whatsapp_links)} WhatsApp links) ===")

if __name__ == "__main__":
    main()
