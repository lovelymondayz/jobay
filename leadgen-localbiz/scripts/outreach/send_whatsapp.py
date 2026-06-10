"""
Stage 6b: WhatsApp Outreach — Generate personalized WA deep links for all qualified leads.
Since 0% of leads have emails but 100% have phone numbers, WhatsApp is the primary channel.
Generates wa.me links with pre-filled messages. You send manually or via WA Business API.
"""
import sys, time, re, json, urllib.request, urllib.parse
sys.path.insert(0, "/root/hermes/leadgen-localbiz")
from config.settings import *
from config.database import query, execute

def format_wa_number(phone):
    """Convert Indonesian phone to international format for WhatsApp."""
    if not phone:
        return None
    # Strip everything except digits
    digits = re.sub(r'[^0-9]', '', phone)
    if not digits:
        return None
    # Indonesian numbers: 08xx -> 628xx
    if digits.startswith('0'):
        digits = '62' + digits[1:]
    if not digits.startswith('62'):
        digits = '62' + digits
    return digits

def craft_wa_message(name, category, rating, review_count, live_url):
    """Generate a personalized WhatsApp message."""
    cat_map = {
        "restaurant": "restoran",
        "gym": "gym/fitness center",
        "dental clinic": "klinik gigi",
        "beauty salon": "salon kecantikan",
        "barber shop": "barbershop",
        "cafe": "kafe",
        "bakery": "toko roti",
        "pet clinic": "klinik hewan",
        "laundry": "laundry",
        "tailor": "penjahit",
    }
    cat_id = cat_map.get(category, category)

    msg = f"""Halo {name}! 👋

Saya menemukan bisnis Anda di Google — {rating} bintang dengan {review_count}+ ulasan, keren banget!

Saya perhatikan website Anda bisa ditingkatkan, jadi saya buatkan preview landing page gratis untuk Anda:

{live_url}

Website ini mobile-friendly, cepat, dan menampilkan info kontak Anda di depan agar customer mudah menghubungi.

Kalau suka, saya bisa bikin website lengkapnya — harga terjangkau dan cepat. Preview ini gratis tanpa syarat.

Gimana menurut Anda? 😊"""
    return msg

def send_discord_summary(leads_data):
    """Send summary to Discord webhook."""
    if not DISCORD_WEBHOOK:
        return

    lines = [
        f"📱 **WhatsApp Outreach Ready** — {time.strftime('%Y-%m-%d %H:%M')}",
        f"",
        f"**{len(leads_data)} leads** ready for WhatsApp outreach.",
        f"",
        f"Top leads (copy-paste ready):",
    ]

    for i, lead in enumerate(leads_data[:10], 1):
        lines.append(f"{i}. **{lead['name']}** ({lead['category']}) — {lead['rating']}★")
        lines.append(f"   📞 {lead['phone']}")
        lines.append(f"   🔗 {lead['wa_link']}")
        lines.append(f"   🌐 {lead['live_url']}")
        lines.append("")

    lines.append(f"💡 **How to send:** Click each wa.me link → WhatsApp opens with pre-filled message → Send!")
    lines.append(f"📊 Full list saved to DB under outreach channel='whatsapp'")

    msg = "\n".join(lines)

    # Split if too long (Discord 2000 char limit)
    chunks = [msg[i:i+1900] for i in range(0, len(msg), 1900)]
    for chunk in chunks:
        data = json.dumps({"content": chunk}).encode()
        req = urllib.request.Request(DISCORD_WEBHOOK, data=data, method="POST")
        req.add_header("Content-Type", "application/json")
        try:
            urllib.request.urlopen(req, timeout=10)
        except Exception as e:
            print(f"  Discord error: {e}")
        time.sleep(1)

def main():
    print(f"=== STAGE 6b: WHATSAPP OUTREACH — {time.strftime('%Y-%m-%d %H:%M')} ===")

    # Get qualified leads with live pages, not yet contacted via WA
    rows = query("""
        SELECT b.id, b.name, b.category, b.rating, b.review_count, b.phone,
               lp.live_url, lp.slug
        FROM businesses b
        JOIN website_audits wa ON b.id = wa.business_id
        JOIN landing_pages lp ON b.id = lp.business_id
        LEFT JOIN outreach o ON b.id = o.business_id AND o.channel = 'whatsapp'
        WHERE wa.overall_score < 40
          AND lp.status = 'live'
          AND o.id IS NULL
        ORDER BY b.rating DESC, b.review_count DESC
    """)

    print(f"Found {len(rows)} leads ready for WhatsApp outreach")
    done = 0
    leads_data = []

    for row in rows:
        biz_id = row["id"]
        name = row["name"]
        category = row["category"]
        rating = row["rating"]
        review_count = row["review_count"]
        phone = row["phone"]
        live_url = row["live_url"]

        wa_num = format_wa_number(phone)
        wa_msg = craft_wa_message(name, category, rating, review_count, live_url)
        wa_link = f"https://wa.me/{wa_num}?text={urllib.parse.quote(wa_msg)}" if wa_num else None

        print(f"  {name} ({rating}★) — {phone}")
        if wa_link:
            print(f"    WA: {wa_link[:80]}...")
        else:
            print(f"    No valid phone")

        # Log outreach
        execute("""
            INSERT INTO outreach (business_id, channel, recipient_phone, subject, message_body, status, sent_at)
            VALUES (%s, %s, %s, %s, %s, %s, NULL)
        """, (biz_id, "whatsapp", phone, f"WhatsApp pitch for {name}", wa_msg, "ready"))

        leads_data.append({
            "name": name,
            "category": category,
            "rating": rating,
            "phone": phone,
            "wa_link": wa_link or "N/A",
            "live_url": live_url or "N/A",
        })

        done += 1

    # Send Discord summary
    if leads_data:
        send_discord_summary(leads_data)

    execute("""INSERT INTO pipeline_runs (stage, status, targets_contacted, finished_at)
               VALUES (%s, %s, %s, NOW())""", ("whatsapp_outreach", "success", done))

    print(f"\n=== DONE: {done} WhatsApp leads ready ===")
    print(f"Total outreach records: {len(rows)}")
    print(f"Next step: Click wa.me links to send messages manually, or connect WA Business API")

if __name__ == "__main__":
    main()
