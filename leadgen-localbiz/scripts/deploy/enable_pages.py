"""
Stage 5b: Enable GitHub Pages on all leadgen repos via API.
Deploys each repo's main branch to GitHub Pages, making sites live at:
  https://lovelymondayz.github.io/{repo-name}/
Then Cloudflare DNS can point client.arjism.com subdomains to these URLs.
"""
import sys, json, time, urllib.request, urllib.error
sys.path.insert(0, "/root/hermes/leadgen-localbiz")
from config.settings import *
from config.database import query, execute, get_conn

def enable_github_pages(repo_name):
    """Enable GitHub Pages on a repo, deploying from main branch."""
    pat = open("/root/.github/pat").read().strip()

    # First check if Pages is already enabled
    req = urllib.request.Request(f"https://api.github.com/repos/{GITHUB_USER}/{repo_name}/pages")
    req.add_header("Authorization", f"Bearer {pat}")
    req.add_header("Accept", "application/vnd.github.v3+json")
    try:
        resp = urllib.request.urlopen(req, timeout=10)
        existing = json.loads(resp.read())
        return "already_enabled", existing.get("html_url", "")
    except urllib.error.HTTPError as e:
        if e.code == 404:
            pass  # Not enabled yet, proceed
        else:
            return "error", str(e)

    # Enable Pages via API
    data = json.dumps({
        "source": {
            "branch": "main",
            "path": "/"
        },
        "build_type": "legacy"
    }).encode()

    req = urllib.request.Request(f"https://api.github.com/repos/{GITHUB_USER}/{repo_name}/pages", data=data, method="POST")
    req.add_header("Authorization", f"Bearer {pat}")
    req.add_header("Content-Type", "application/json")
    req.add_header("Accept", "application/vnd.github.v3+json")
    req.add_header("X-GitHub-Api-Version", "2022-11-28")

    try:
        resp = urllib.request.urlopen(req, timeout=15)
        result = json.loads(resp.read())
        return "enabled", result.get("html_url", "")
    except urllib.error.HTTPError as e:
        body = ""
        try:
            body = e.read().decode()[:300]
        except:
            pass
        return "error", f"{e.code}: {body}"
    except Exception as e:
        return "error", str(e)


def main():
    print(f"=== STAGE 5b: ENABLE GITHUB PAGES — {time.strftime('%Y-%m-%d %H:%M')} ===\n")

    # Get all repos that need Pages enabled
    rows = query("""
        SELECT id, slug, repo_name, repo_url, pages_enabled
        FROM landing_pages
        WHERE status = 'live'
          AND github_pushed = true
          AND (pages_enabled = false OR pages_enabled IS NULL)
        ORDER BY id ASC
    """)

    print(f"Repos to enable: {len(rows)}\n")

    enabled = 0
    skipped = 0
    failed = 0

    for i, row in enumerate(rows, 1):
        slug = row["slug"]
        repo_name = row["repo_name"]
        page_id = row["id"]

        print(f"  [{i}/{len(rows)}] {repo_name}...", end=" ", flush=True)

        status, url_or_err = enable_github_pages(repo_name)

        if status == "already_enabled":
            print(f"✅ Already enabled: {url_or_err}")
            # Update DB
            execute("UPDATE landing_pages SET pages_enabled=true, github_pages_url=%s WHERE id=%s",
                    (url_or_err, page_id))
            skipped += 1
        elif status == "enabled":
            print(f"✅ Enabled: {url_or_err}")
            execute("UPDATE landing_pages SET pages_enabled=true, github_pages_url=%s WHERE id=%s",
                    (url_or_err, page_id))
            enabled += 1
        else:
            print(f"❌ Error: {url_or_err}")
            failed += 1

        time.sleep(0.5)  # Rate limit

    print(f"\n=== SUMMARY ===")
    print(f"  Enabled:  {enabled}")
    print(f"  Skipped:  {skipped} (already on)")
    print(f"  Failed:   {failed}")
    print(f"\nGitHub Pages URLs will be:")
    print(f"  https://lovelymondayz.github.io/{{repo-name}}/")
    print(f"\nNext: Point client.arjism.com DNS to GitHub Pages")


if __name__ == "__main__":
    main()
