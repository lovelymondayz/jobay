"""
Stage 2: Evaluate — Check website quality for each business.
Direct HTTP scrape (no API key needed). Scores 0-100.
If no site or score < 40 -> QUALIFY as lead.
"""
import sys, os, json, time, urllib.request, re
sys.path.insert(0, "/root/hermes/leadgen-localbiz")
from config.settings import *
from config.database import query, execute

def fetch_url(url, timeout=10):
    """Fetch URL content via direct HTTP."""
    if not url:
        return None, "no_url"
    # Normalize URL
    raw = url
    if not url.startswith("http://") and not url.startswith("https://"):
        url = "https://" + url
    req = urllib.request.Request(url, headers={
        "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
        "Accept": "text/html,application/xhtml+xml,*/*;q=0.8",
        "Accept-Language": "en-US,en;q=0.9",
    })
    try:
        resp = urllib.request.urlopen(req, timeout=timeout)
        content = resp.read().decode("utf-8", errors="ignore")
        return content, None
    except urllib.error.HTTPError as e:
        if e.code in [301, 302, 303, 307, 308]:
            return None, "redirect"
        return None, "http_" + str(e.code)
    except urllib.error.URLError as e:
        return None, "dns_error"
    except Exception as e:
        return None, "timeout"

def score_website(url, content=None):
    """Score 0-100. Returns (score, details_dict)."""
    if not url:
        return 0, {"has_site": False, "score": 0, "notes": "No website listed"}

    if content is None:
        return 5, {"has_site": True, "score": 5, "notes": "Website listed but could not fetch"}

    score = 50
    notes = []
    content_lower = content.lower()

    # Content length
    if len(content) > 3000:
        score += 10
    elif len(content) > 1000:
        score += 5
    else:
        score -= 15
        notes.append("Very thin content (" + str(len(content)) + " chars)")

    # Contact info
    has_contact = any(w in content_lower for w in ["contact", "whatsapp", "email", "phone", "hubungi", "kontak"])
    if has_contact:
        score += 10
    else:
        score -= 10
        notes.append("No contact info")

    # Menu/services
    has_menu = any(w in content_lower for w in ["menu", "service", "price", "harga", "layanan", "pelayanan"])
    if has_menu:
        score += 10
    else:
        score -= 5

    # Mobile-friendly (viewport meta)
    is_mobile = "viewport" in content_lower
    if is_mobile:
        score += 10
    else:
        score -= 10
        notes.append("Not mobile-friendly (no viewport)")

    # Modern framework signals
    modern = any(s in content for s in ["react", "vue", "next", "nuxt", "tailwind", "bootstrap"])
    if modern:
        score += 10

    # Old/outdated signals
    if "<table" in content and "<div" not in content:
        score -= 20
        notes.append("Table-based layout (outdated)")

    # Free hosting penalty
    free_hosts = ["wordpress.com", "wixsite.com", "jimdo.com", "godaddy.com", "webnode"]
    if any(h in url.lower() for h in free_hosts):
        score -= 15
        notes.append("Free/subdomain hosting")

    # HTTPS bonus
    if url.startswith("https://"):
        score += 5
    else:
        score -= 5
        notes.append("No HTTPS")

    # SEO basics
    has_title = "<title>" in content_lower
    has_h1 = "<h1" in content_lower
    if has_title and has_h1:
        score += 5
    else:
        score -= 5
        notes.append("Missing title or H1")

    # Clamp
    score = max(0, min(100, score))

    return score, {
        "has_site": True,
        "score": score,
        "is_mobile_friendly": is_mobile,
        "has_contact_form": has_contact,
        "has_menu": has_menu,
        "has_booking": "book" in content_lower or "reservation" in content_lower or "daftar" in content_lower,
        "is_modern_design": modern or score >= 60,
        "seo_score": min(100, max(0, score + (10 if has_title and has_h1 else -10))),
        "notes": "; ".join(notes) if notes else "OK",
        "qualified": score < 40,
    }

def main():
    print("=== STAGE 2: EVALUATE " + time.strftime("%Y-%m-%d %H:%M") + " ===")

    # Get businesses not yet audited
    rows = query("""
        SELECT b.id, b.name, b.website_url, b.has_website, b.category, b.rating, b.review_count
        FROM businesses b
        LEFT JOIN website_audits wa ON b.id = wa.business_id
        WHERE wa.id IS NULL
        ORDER BY b.rating DESC, b.review_count DESC
        LIMIT 50
    """)

    print("Auditing", len(rows), "businesses")
    qualified = 0

    for row in rows:
        biz_id = row["id"]
        name = row["name"]
        website = row["website_url"]
        has_site = row["has_website"]

        print("  " + name, end="")
        if website:
            print(" -> " + website[:60])
        else:
            print(" -> NO WEBSITE")

        if not has_site or not website:
            # No website at all — auto-qualify
            score = 0
            details = {"has_site": False, "score": 0, "notes": "No website", "qualified": True}
            qualified += 1
        else:
            content, err = fetch_url(website)
            if err:
                print("    Fetch failed:", err)
                score = 5
                details = {"has_site": True, "score": 5, "notes": "Could not fetch: " + err, "qualified": True}
                qualified += 1
            else:
                score, details = score_website(website, content)
                print("    Score:", score, "/100 —", details["notes"])

            if details.get("qualified"):
                qualified += 1
                print("    >>> QUALIFIED! <<<")

        # Save audit
        execute("""
            INSERT INTO website_audits
            (business_id, has_site, site_url, overall_score, seo_score,
             has_contact_form, has_menu, has_booking, is_mobile_friendly,
             is_modern_design, audit_notes, scraped_content, audited_at)
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, NOW())
        """, (
            biz_id,
            details.get("has_site", False),
            website or "",
            score,
            details.get("seo_score", score),
            details.get("has_contact_form", False),
            details.get("has_menu", False),
            details.get("has_booking", False),
            details.get("is_mobile_friendly", False),
            details.get("is_modern_design", False),
            (details.get("notes", "") or "")[:500],
            (content or "")[:3000] if content else None,
        ))

        time.sleep(1.5)  # Rate limit

    execute("INSERT INTO pipeline_runs (stage, status, targets_found, targets_qualified, started_at, finished_at) VALUES (%s, %s, %s, %s, NOW(), NOW())",
            ("evaluate", "success", len(rows), qualified))

    print("\n=== DONE:", len(rows), "audited,", qualified, "qualified ===")

if __name__ == "__main__":
    main()
