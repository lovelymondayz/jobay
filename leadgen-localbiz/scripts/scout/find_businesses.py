"""
Stage 1: Scout — Find Jakarta businesses with good reviews.
Uses SerpApi Google Maps engine (200 searches/month — use wisely).
Strategy: 1 search per category per city, each returns up to 20 results.
"""
import sys, json, time, urllib.request, urllib.parse, os
sys.path.insert(0, "/root/hermes/leadgen-localbiz")
from config.settings import *
from config.database import query, execute

SERP_API_KEY = os.environ.get("SERP_API_KEY", "3c923a4adcde8464f9eb4c4cfaa6085d80839d051e943a0fcc4995e6f6633a0d")

def serpapi_search(query_str, limit=20):
    """Search SerpApi Google Maps engine."""
    if not SERP_API_KEY:
        print("  No SERP_API_KEY")
        return []
    params = urllib.parse.urlencode({
        "engine": "google_maps",
        "q": query_str,
        "type": "search",
        "api_key": SERP_API_KEY,
        "num": limit,
    })
    url = f"https://serpapi.com/search?{params}"
    req = urllib.request.Request(url, headers={"User-Agent": "leadgen/1.0"})
    try:
        resp = urllib.request.urlopen(req, timeout=20)
        data = json.loads(resp.read())
        return data.get("local_results", [])
    except Exception as e:
        print(f"  SerpApi error: {e}")
        return []

def save_business(biz):
    """Save business to DB, skip duplicates."""
    try:
        execute("""
            INSERT INTO businesses (name, category, city, address, phone, rating, review_count, has_website, website_url, source, raw_data)
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
            ON CONFLICT DO NOTHING
        """, (biz["name"], biz.get("category", ""), biz.get("city", "Jakarta"),
              biz.get("address", ""), biz.get("phone", ""), biz.get("rating", 0),
              biz.get("review_count", 0), biz.get("has_website", False),
              biz.get("website_url", ""), biz.get("source", ""),
              json.dumps(biz.get("raw_data", {}))))
        return True
    except Exception as e:
        print(f"  DB error: {e}")
        return False

def main():
    print(f"=== STAGE 1: SCOUT — {time.strftime('%Y-%m-%d %H:%M')} ===")
    print(f"Targets: {TARGET_CATEGORIES[:5]} x Jakarta/BSD")
    print(f"SerpApi budget: 200 searches/month\n")

    total_found = 0
    search_count = 0

    for city in ["Jakarta", "BSD Tangerang"]:
        for category in TARGET_CATEGORIES[:5]:  # 5 categories x 2 cities = 10 searches
            query_str = f"{category} in {city} Indonesia"
            print(f"  [{search_count+1}] Searching: {query_str}")

            results = serpapi_search(query_str, limit=20)
            search_count += 1
            found_here = 0

            for r in results:
                rating = float(r.get("rating", 0))
                reviews = int(r.get("reviews", 0) or 0)

                # Filter: minimum quality
                if rating < MIN_RATING or reviews < MIN_REVIEWS:
                    continue

                name = r.get("title", "").strip()
                if not name:
                    continue

                website = r.get("website", "") or ""
                phone = r.get("phone", "") or ""
                address = r.get("address", "") or ""

                biz = {
                    "name": name,
                    "category": category,
                    "city": city,
                    "rating": rating,
                    "review_count": reviews,
                    "phone": phone,
                    "address": address,
                    "has_website": bool(website),
                    "website_url": website,
                    "source": "serpapi",
                    "raw_data": {
                        "place_id": r.get("place_id", ""),
                        "place_id_search": r.get("place_id_search", ""),
                        "type": r.get("type", ""),
                        "thumbnail": r.get("thumbnail", ""),
                    },
                }

                if save_business(biz):
                    total_found += 1
                    found_here += 1
                    ws = "HAS site" if website else "NO site"
                    print(f"    + {name} ({rating}★ / {reviews} reviews) [{ws}]")

            print(f"    -> {found_here} saved")
            time.sleep(1)

    execute("INSERT INTO pipeline_runs (stage, status, targets_found, started_at, finished_at) VALUES (%s, %s, %s, NOW(), NOW())",
            ("scout", "success", total_found))

    print(f"\n=== DONE: {total_found} businesses saved | {search_count} SerpApi searches used ===")
    print(f"Remaining budget: {200 - search_count} searches this month")

if __name__ == "__main__":
    main()
