"""
Stage 1.5b: Social Media Scout v2 — Targeted Instagram search for top no-website leads.
Strategy: Instagram only for 4.5+ star restaurants with no website.
"""
import sys, json, time, urllib.request, urllib.parse, re, os

sys.path.insert(0, "/root/hermes/leadgen-localbiz")
from config.settings import *
from config.database import query, execute

SERP_API_KEY = os.environ.get("SERP_API_KEY", "")


def serpapi_search_ig(name, city="Jakarta"):
    """Search for Instagram profile only."""
    if not SERP_API_KEY:
        return None

    query_str = f'site:instagram.com "{name}" {city}'
    params = urllib.parse.urlencode({
        "engine": "google",
        "q": query_str,
        "api_key": SERP_API_KEY,
        "num": 3,
    })
    url = f"https://serpapi.com/search?{params}"
    req = urllib.request.Request(url, headers={"User-Agent": "leadgen/1.0"})

    try:
        resp = urllib.request.urlopen(req, timeout=20)
        data = json.loads(resp.read())
        results = data.get("organic_results", [])

        for r in results:
            link = r.get("link", "")
            match = re.search(r"instagram\.com/([a-zA-Z0-9_.]{2,30})/?$", link)
            if match:
                username = match.group(1)
                bad = ["p", "reel", "reels", "stories", "explore", "accounts", "popular", "video"]
                if username.lower() not in bad:
                    return f"https://instagram.com/{username}"
        return None
    except Exception as e:
        print(f"    Error: {e}")
        return None


def main():
    print("=== SOCIAL MEDIA SCOUT v2 — Instagram Only ===\n")

    rows = query("""
        SELECT id, name, category, city, rating, review_count, phone
        FROM businesses
        WHERE has_website = FALSE
          AND (instagram_url IS NULL OR instagram_url = '')
          AND rating >= 4.5
          AND review_count >= 100
          AND category = 'restaurant'
        ORDER BY rating DESC, review_count DESC
    """)

    print(f"Target: {len(rows)} no-website restaurants (4.5+ stars, 100+ reviews)\n")

    searches = 0
    found = 0

    for row in rows:
        biz_id = row["id"]
        name = row["name"]
        city = row.get("city", "Jakarta")
        rating = row["rating"]
        reviews = row["review_count"]
        phone = row["phone"]

        print(f"  [{searches+1}] {name} ({rating} / {reviews} rev) | {phone}")

        ig = serpapi_search_ig(name, city)
        searches += 1

        if ig:
            found += 1
            print(f"    FOUND: {ig}")
            execute("UPDATE businesses SET instagram_url = %s WHERE id = %s", (ig, biz_id))
        else:
            print(f"    not found")

        time.sleep(1.5)

    print(f"\n=== DONE ===")
    print(f"Searches used: {searches}")
    print(f"IG found: {found}/{len(rows)}")


if __name__ == "__main__":
    main()
