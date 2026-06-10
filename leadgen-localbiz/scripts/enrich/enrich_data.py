"""
Stage 3: Enrich — Geocode qualified business addresses via OpenStreetMap.
"""
import sys, json, time, urllib.request, urllib.parse
sys.path.insert(0, "/root/hermes/leadgen-localbiz")
from config.settings import *
from config.database import query, execute

def geocode_address(address):
    """Geocode via OpenStreetMap Nominatim (free, no key)."""
    if not address or len(address) < 5:
        return None, None
    url = "https://nominatim.openstreetmap.org/search?format=json&q=" + urllib.parse.quote(address)
    req = urllib.request.Request(url, headers={"User-Agent": "leadgen-pipeline/1.0"})
    try:
        resp = urllib.request.urlopen(req, timeout=10)
        data = json.loads(resp.read())
        if data:
            return float(data[0]["lat"]), float(data[0]["lon"])
    except Exception as e:
        print("    Geocode error:", e)
    return None, None

def main():
    print("=== STAGE 3: ENRICH " + time.strftime("%Y-%m-%d %H:%M") + " ===")

    # Get qualified businesses not yet enriched
    rows = query("""
        SELECT b.id, b.name, b.address, b.category, b.city
        FROM businesses b
        JOIN website_audits wa ON b.id = wa.business_id
        LEFT JOIN business_enrichment be ON b.id = be.business_id
        WHERE wa.overall_score < 40
          AND be.id IS NULL
        ORDER BY b.rating DESC, b.review_count DESC
        LIMIT 30
    """)

    print("Enriching", len(rows), "businesses")
    enriched = 0

    for row in rows:
        biz_id = row["id"]
        name = row["name"]
        address = row["address"]
        city = row.get("city", "Jakarta")

        print("  " + name)
        if address:
            print("    Address:", address[:80])

        lat, lng = None, None
        if address and len(address) > 5:
            lat, lng = geocode_address(address)
            if lat and lng:
                print("    Geocode:", lat, lng)
            else:
                print("    Geocode: not found")
            time.sleep(1.2)  # Nominatim rate limit

        # Build a simple description from what we have
        desc = name + " — " + row.get("category", "") + " di " + city
        if address:
            desc += ". " + address

        execute("""
            INSERT INTO business_enrichment (business_id, description, lat, lng, enriched_at)
            VALUES (%s, %s, %s, %s, NOW())
            ON CONFLICT DO NOTHING
        """, (biz_id, desc[:500], lat, lng))

        enriched += 1

    print("=== DONE:", enriched, "enriched ===")

if __name__ == "__main__":
    main()
