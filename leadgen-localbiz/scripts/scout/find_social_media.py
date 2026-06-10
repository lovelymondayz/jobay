"""
Stage 1.5: Social Media Scout — Find Instagram/Facebook/TikTok for existing leads.
Uses SerpApi Google search (counts toward 200/month budget).
Searches: site:instagram.com "restaurant_name" Jakarta
"""
import sys, json, time, urllib.request, urllib.parse, re, os

sys.path.insert(0, "/root/hermes/leadgen-localbiz")
from config.settings import *
from config.database import query, execute

SERP_API_KEY = os.environ.get("SERP_API_KEY", "")

def serpapi_search_site(site, name, city="Jakarta"):
    """Search for a specific site profile using SerpApi Google search."""
    if not SERP_API_KEY:
        return []
    
    query_str = f"site:{site} \"{name}\" {city}"
    params = urllib.parse.urlencode({
        "engine": "google",
        "q": query_str,
        "api_key": SERP_API_KEY,
        "num": 5,
    })
    url = f"https://serpapi.com/search?{params}"
    req = urllib.request.Request(url, headers={"User-Agent": "leadgen/1.0"})
    
    try:
        resp = urllib.request.urlopen(req, timeout=20)
        data = json.loads(resp.read())
        results = data.get("organic_results", [])
        return results
    except Exception as e:
        print(f"    SerpApi error: {e}")
        return []

def extract_social_url(results, platform):
    """Extract the best matching social URL from search results."""
    patterns = {
        "instagram": r"instagram\.com/([a-zA-Z0-9_.]+)",
        "facebook": r"facebook\.com/([a-zA-Z0-9_.]+)",
        "tiktok": r"tiktok\.com/@([a-zA-Z0-9_.]+)",
    }
    pattern = patterns.get(platform, "")
    
    for r in results:
        link = r.get("link", "")
        title = r.get("title", "").lower()
        snippet = r.get("snippet", "").lower()
        
        match = re.search(pattern, link)
        if match:
            username = match.group(1)
            # Filter out generic pages
            if username.lower() not in ["p", "reel", "stories", "explore", "accounts"]:
                return link.strip("/")
    
    return None

def find_social_media(name, city="Jakarta"):
    """Find all social media for a business."""
    found = {"instagram": None, "facebook": None, "tiktok": None}
    
    # Search Instagram
    ig_results = serpapi_search_site("instagram.com", name, city)
    ig_url = extract_social_url(ig_results, "instagram")
    if ig_url:
        found["instagram"] = ig_url
    
    # Search Facebook
    fb_results = serpapi_search_site("facebook.com", name, city)
    fb_url = extract_social_url(fb_results, "facebook")
    if fb_url:
        found["facebook"] = fb_url
    
    # Search TikTok
    tt_results = serpapi_search_site("tiktok.com", name, city)
    tt_url = extract_social_url(tt_results, "tiktok")
    if tt_url:
        found["tiktok"] = tt_url
    
    return found

def main():
    print("=== SOCIAL MEDIA SCOUT ===")
    print(f"Budget: SerpApi 200 searches/month\n")
    
    # Get businesses without social media data, prioritizing no-website restaurants
    rows = query("""
        SELECT id, name, category, city, rating, review_count, phone, has_website
        FROM businesses
        WHERE (social_media_found = FALSE OR social_media_found IS NULL)
          AND category = 'restaurant'
        ORDER BY has_website ASC, rating DESC, review_count DESC
    """)
    
    print(f"Found {len(rows)} restaurants to search\n")
    
    total_searches = 0
    found_count = 0
    search_budget = 200
    
    for row in rows:
        biz_id = row["id"]
        name = row["name"]
        city = row.get("city", "Jakarta")
        rating = row["rating"]
        reviews = row["review_count"]
        has_site = row["has_website"]
        
        site_label = "NO site" if not has_site else "has site"
        print(f"  [{total_searches+1}] {name} ({rating}★ / {reviews} rev) [{site_label}]")
        
        # Check budget (3 searches per business: IG + FB + TT)
        remaining = search_budget - total_searches
        if remaining < 3:
            print(f"  ⚠️  Budget exhausted ({total_searches}/{search_budget}). Stopping.")
            break
        
        social = find_social_media(name, city)
        total_searches += 3  # 3 searches per business
        
        ig = social.get("instagram")
        fb = social.get("facebook")
        tt = social.get("tiktok")
        
        has_social = any([ig, fb, tt])
        
        if has_social:
            found_count += 1
            parts = []
            if ig: parts.append(f"IG: {ig}")
            if fb: parts.append(f"FB: {fb}")
            if tt: parts.append(f"TT: {tt}")
            print(f"    ✅ {' | '.join(parts)}")
        else:
            print(f"    ❌ No social found")
        
        # Update DB
        execute("""
            UPDATE businesses 
            SET instagram_url = %s, facebook_url = %s, tiktok_url = %s, social_media_found = %s
            WHERE id = %s
        """, (ig, fb, tt, has_social, biz_id))
        
        time.sleep(1.5)  # Rate limit for SerpApi
    
    print(f"\n=== DONE ===")
    print(f"Searches used: {total_searches}")
    print(f"Social media found: {found_count}/{len(rows)}")
    print(f"Budget remaining: {search_budget - total_searches}")

if __name__ == "__main__":
    main()
