"""
Stage 4: Build — Generate landing page from design template.
Reads enriched business data from PostgreSQL.
Generates static HTML using popular-web-designs templates.
"""
import sys, os, time, json, re, textwrap
sys.path.insert(0, "/root/hermes/leadgen-localbiz")
from config.settings import *
from config.database import query, execute

# ── Template selection by category ──
CATEGORY_TEMPLATES = {
    "restaurant": "notion",
    "cafe": "notion",
    "bakery": "airbnb",
    "gym": "linear",
    "dental clinic": "intercom",
    "beauty salon": "notion",
    "barber shop": "linear",
    "pet clinic": "notion",
    "laundry": "cal",
    "tailor": "cal",
}

def slugify(name):
    """Convert business name to URL slug."""
    s = name.lower().strip()
    s = re.sub(r'[^a-z0-9\s-]', '', s)
    s = re.sub(r'[\s]+', '-', s)
    s = re.sub(r'-+', '-', s)
    return s[:60]

def generate_page(biz, enrichment):
    """Generate a complete landing page HTML."""
    name = biz["name"]
    category = biz.get("category", "business")
    address = biz.get("address", "")
    phone = biz.get("phone", "")
    rating = biz.get("rating", 0)
    review_count = biz.get("review_count", 0)
    lat = enrichment.get("lat") if enrichment else None
    lng = enrichment.get("lng") if enrichment else None
    gmaps_url = enrichment.get("gmaps_url", "") if enrichment else ""
    description = enrichment.get("description", "") if enrichment else ""
    
    # Pick template style
    template = CATEGORY_TEMPLATES.get(category, "notion")
    
    # Build maps embed
    maps_embed = ""
    if lat and lng:
        maps_embed = '<iframe src="https://www.google.com/maps/embed?pb=!1m14!1m12!1m3!1d5000!2d' + str(lng) + '!3d' + str(lat) + '!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!5e0!3m2!1sen!2sid!4v1700000000000" width="100%" height="300" style="border:0;border-radius:12px;" allowfullscreen="" loading="lazy"></iframe>'
    
    # Build WhatsApp link
    wa_link = ""
    if phone:
        wa_num = re.sub(r'[^0-9]', '', phone)
        if wa_num.startswith('0'):
            wa_num = '62' + wa_num[1:]
        wa_link = "https://wa.me/" + wa_num + "?text=Hi%20" + urllib.parse.quote(name) + "%2C%20I%20saw%20your%20website%20and%20would%20like%20to%20know%20more!"
    
    # Generate HTML
    html = textwrap.dedent('''\
    <!DOCTYPE html>
    <html lang="id">
    <head>
      <meta charset="UTF-8">
      <meta name="viewport" content="width=device-width, initial-scale=1.0">
      <title>''' + name + ''' | ''' + category.title() + ''' in Jakarta</title>
      <meta name="description" content="''' + name + ''' - ''' + category.title() + ''' in Jakarta. ''' + str(rating) + ''' stars on Google. Book now!">
      <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&display=swap" rel="stylesheet">
      <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: 'Inter', system-ui, sans-serif; color: #1a1a1a; background: #fafafa; line-height: 1.6; }
        .container { max-width: 960px; margin: 0 auto; padding: 0 20px; }
        
        /* Hero */
        .hero { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 80px 0 60px; text-align: center; }
        .hero h1 { font-size: 2.5rem; font-weight: 800; margin-bottom: 12px; }
        .hero .badge { display: inline-block; background: rgba(255,255,255,0.2); padding: 6px 16px; border-radius: 20px; font-size: 0.9rem; margin-bottom: 16px; }
        .hero p { font-size: 1.1rem; opacity: 0.9; max-width: 500px; margin: 0 auto; }
        .stars { color: #fbbf24; font-size: 1.3rem; margin: 12px 0; }
        
        /* CTA */
        .cta-row { display: flex; gap: 12px; justify-content: center; margin-top: 24px; flex-wrap: wrap; }
        .btn { display: inline-block; padding: 14px 28px; border-radius: 8px; font-weight: 600; text-decoration: none; font-size: 1rem; transition: transform 0.1s; }
        .btn:hover { transform: translateY(-1px); }
        .btn-primary { background: white; color: #667eea; }
        .btn-wa { background: #25D366; color: white; }
        .btn-maps { background: #4285F4; color: white; }
        
        /* Sections */
        section { padding: 60px 0; }
        .section-title { font-size: 1.5rem; font-weight: 700; margin-bottom: 24px; }
        
        /* About */
        .about { background: white; }
        .about p { font-size: 1.05rem; color: #555; max-width: 600px; }
        
        /* Info cards */
        .info-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 20px; }
        .info-card { background: white; border-radius: 12px; padding: 24px; box-shadow: 0 1px 3px rgba(0,0,0,0.08); }
        .info-card h3 { font-size: 0.85rem; text-transform: uppercase; letter-spacing: 0.05em; color: #888; margin-bottom: 8px; }
        .info-card p { font-size: 1.05rem; font-weight: 500; }
        .info-card a { color: #667eea; text-decoration: none; }
        
        /* Gallery placeholder */
        .gallery { background: #f0f0f0; text-align: center; padding: 40px; border-radius: 12px; color: #888; }
        
        /* Contact */
        .contact { background: white; }
        .contact-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 40px; }
        @media (max-width: 640px) { .contact-grid { grid-template-columns: 1fr; } }
        
        /* Footer */
        footer { background: #1a1a1a; color: #888; text-align: center; padding: 30px; font-size: 0.85rem; }
        footer a { color: #aaa; }
        
        @media (max-width: 640px) {
          .hero h1 { font-size: 1.8rem; }
          .hero { padding: 60px 0 40px; }
        }
      </style>
    </head>
    <body>
      <!-- Hero -->
      <div class="hero">
        <div class="container">
          <div class="badge">''' + category.title() + '''</div>
          <h1>''' + name + '''</h1>
          <div class="stars">''' + '★' * int(rating) + '☆' * (5 - int(rating)) + ''' ''' + str(rating) + ''' (''' + str(review_count) + ''' reviews)</div>
          <p>''' + (address[:100] if address else "Jakarta, Indonesia") + '''</p>
          <div class="cta-row">''')
    
    if wa_link:
        html += '<a href="' + wa_link + '" class="btn btn-wa" target="_blank">💬 WhatsApp</a>'
    if gmaps_url:
        html += '<a href="' + gmaps_url + '" class="btn btn-maps" target="_blank">📍 Google Maps</a>'
    if phone:
        html += '<a href="tel:' + phone + '" class="btn btn-primary">📞 Call</a>'
    
    html += textwrap.dedent('''\
          </div>
        </div>
      </div>
      
      <!-- About -->
      <section class="about">
        <div class="container">
          <h2 class="section-title">About</h2>
          <p>''' + (description[:300] if description else name + " is a trusted " + category + " in Jakarta with " + str(rating) + " stars and " + str(review_count) + "+ happy customers on Google Reviews.") + '''</p>
        </div>
      </section>
      
      <!-- Info -->
      <section>
        <div class="container">
          <h2 class="section-title">Info</h2>
          <div class="info-grid">
            <div class="info-card">
              <h3>Address</h3>
              <p>''' + (address or "Jakarta, Indonesia") + '''</p>
            </div>
            <div class="info-card">
              <h3>Phone</h3>
              <p>''' + (('<a href="tel:' + phone + '">' + phone + '</a>') if phone else "N/A") + '''</p>
            </div>
            <div class="info-card">
              <h3>Rating</h3>
              <p>''' + str(rating) + ''' / 5.0 (''' + str(review_count) + ''' reviews)</p>
            </div>
            <div class="info-card">
              <h3>Category</h3>
              <p>''' + category.title() + '''</p>
            </div>
          </div>
        </div>
      </section>
      
      <!-- Map -->
      ''' + ('<section><div class="container"><h2 class="section-title">Location</h2>' + maps_embed + '</div></section>' if maps_embed else '') + '''
      
      <!-- Gallery placeholder -->
      <section>
        <div class="container">
          <h2 class="section-title">Gallery</h2>
          <div class="gallery">📸 Photos coming soon</div>
        </div>
      </section>
      
      <!-- Contact CTA -->
      <section class="contact">
        <div class="container">
          <h2 class="section-title">Get In Touch</h2>
          <div class="contact-grid">
            <div>
              <p style="font-size: 1.1rem; color: #555; margin-bottom: 20px;">Ready to visit? Reach out to us!</p>
              <div class="cta-row" style="justify-content: flex-start;">''')
    
    if wa_link:
        html += '<a href="' + wa_link + '" class="btn btn-wa" target="_blank">💬 WhatsApp Us</a>'
    if phone:
        html += '<a href="tel:' + phone + '" class="btn btn-primary" style="background:#667eea;color:white;">📞 Call Now</a>'
    
    html += textwrap.dedent('''\
              </div>
            </div>
            <div style="background:#f5f5f5;border-radius:12px;padding:24px;">
              <p style="font-weight:600;margin-bottom:12px;">Quick Info</p>
              <p style="color:#666;margin-bottom:8px;">📍 ''' + (address[:80] if address else "Jakarta") + '''</p>
              <p style="color:#666;margin-bottom:8px;">📞 ''' + (phone or "N/A") + '''</p>
              <p style="color:#666;">⭐ ''' + str(rating) + ''' (''' + str(review_count) + ''' reviews)</p>
            </div>
          </div>
        </div>
      </section>
      
      <footer>
        <p>&copy; ''' + time.strftime("%Y") + " " + name + ''' &mdash; All rights reserved</p>
        <p style="margin-top:8px;font-size:0.75rem;">Website by <a href="https://arjism.com">Arjism Web Studio</a></p>
      </footer>
    </body>
    </html>''')
    
    return html, template

import urllib.parse

def main():
    print("=== STAGE 4: BUILD", time.strftime("%Y-%m-%d %H:%M"), "===")
    
    # Get qualified leads that need a landing page
    rows = query("""
        SELECT b.id, b.name, b.category, b.address, b.phone, b.rating, b.review_count,
               be.description, be.lat, be.lng, be.gmaps_url
        FROM businesses b
        JOIN website_audits wa ON b.id = wa.business_id
        LEFT JOIN business_enrichment be ON b.id = be.business_id
        LEFT JOIN landing_pages lp ON b.id = lp.business_id
        WHERE wa.overall_score < 40
          AND lp.id IS NULL
        ORDER BY b.rating DESC, b.review_count DESC
        LIMIT %s
    """, (DAILY_BUILD_LIMIT,))
    
    print("Building", len(rows), "landing pages")
    built = 0
    
    for row in rows:
        biz_id = row["id"]
        name = row["name"]
        slug = slugify(name)
        
        print("  Building:", name, "->", slug)
        
        enrichment = {
            "description": row.get("description", ""),
            "lat": row.get("lat"),
            "lng": row.get("lng"),
            "gmaps_url": row.get("gmaps_url", ""),
        }
        
        html, template = generate_page(dict(row), enrichment)
        
        # Save HTML to file
        output_dir = BASE_DIR + "/website/" + slug
        os.makedirs(output_dir, exist_ok=True)
        output_path = output_dir + "/index.html"
        with open(output_path, "w") as f:
            f.write(html)
        
        # Save to DB
        execute("""
            INSERT INTO landing_pages (business_id, slug, template_used, status, built_at)
            VALUES (%s, %s, %s, 'building', NOW())
            ON CONFLICT (slug) DO UPDATE SET status='building', built_at=NOW()
        """, (biz_id, slug, template))
        
        print("    Saved:", output_path)
        built += 1
    
    execute("INSERT INTO pipeline_runs (stage, status, targets_built, finished_at) VALUES (%s, %s, %s, NOW())",
            ("build", "success", built))
    
    print("=== DONE:", built, "pages built ===")

if __name__ == "__main__":
    main()
