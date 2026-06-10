"""
Stage 5: Deploy — Push HTML to GitHub repo.
You manually connect repo to Cloudflare Pages + DNS via dashboard.
"""
import sys, os, json, time, urllib.request, subprocess, shutil
sys.path.insert(0, "/root/hermes/leadgen-localbiz")
from config.settings import *
from config.database import query, execute

def github_create_repo(name, description=""):
    """Create a public GitHub repo via API."""
    data = json.dumps({
        "name": name,
        "description": description,
        "private": False,
        "auto_init": False,
    }).encode()
    req = urllib.request.Request(GITHUB_API + "/user/repos", data=data, method="POST")
    req.add_header("Authorization", "Bearer " + GITHUB_PAT)
    req.add_header("Content-Type", "application/json")
    req.add_header("Accept", "application/vnd.github.v3+json")
    try:
        resp = urllib.request.urlopen(req, timeout=15)
        result = json.loads(resp.read())
        return result.get("html_url"), result.get("clone_url")
    except Exception as e:
        err_body = ""
        if hasattr(e, "read"):
            err_body = e.read().decode()
        print("  GitHub create repo error:", e, err_body[:200])
        return None, None

def push_to_github(local_dir, repo_name, repo_url):
    """Push local directory content to GitHub repo."""
    try:
        git_dir = os.path.join(local_dir, ".git")
        if os.path.exists(git_dir):
            shutil.rmtree(git_dir)

        auth_url = repo_url.replace("https://", "https://" + GITHUB_PAT + "@")

        commands = [
            ["git", "init"],
            ["git", "config", "user.email", "bot@arjism.com"],
            ["git", "config", "user.name", "Leadgen Bot"],
            ["git", "add", "."],
            ["git", "commit", "-m", "Landing page for " + repo_name],
            ["git", "branch", "-M", "main"],
            ["git", "remote", "add", "origin", auth_url],
            ["git", "push", "-u", "origin", "main", "--force"],
        ]

        for cmd in commands:
            result = subprocess.run(cmd, cwd=local_dir, capture_output=True, text=True, timeout=60)
            if result.returncode != 0 and "nothing to commit" not in result.stderr:
                print("  Git error:", result.stderr[:200])
                return False

        return True
    except Exception as e:
        print("  Push error:", e)
        return False

def main():
    print("=== STAGE 5: DEPLOY", time.strftime("%Y-%m-%d %H:%M"), "===")

    rows = query("""
        SELECT lp.id as page_id, lp.slug, lp.business_id, lp.repo_name,
               b.name as biz_name
        FROM landing_pages lp
        JOIN businesses b ON lp.business_id = b.id
        WHERE lp.status = 'building'
          AND lp.github_pushed = false
        ORDER BY lp.built_at ASC
        LIMIT %s
    """, (DAILY_BUILD_LIMIT,))

    print("Deploying", len(rows), "pages")
    deployed = 0

    for row in rows:
        slug = row["slug"]
        biz_name = row["biz_name"]
        page_id = row["page_id"]
        repo_name = "leadgen-" + slug

        local_dir = BASE_DIR + "/website/" + slug

        if not os.path.exists(local_dir):
            print("  Missing local dir:", local_dir)
            continue

        print("  Deploying:", biz_name, "->", repo_name)

        # Create GitHub repo
        repo_html_url, repo_clone_url = github_create_repo(
            repo_name,
            "Landing page for " + biz_name
        )

        if repo_clone_url:
            print("    GitHub repo:", repo_html_url)

            if push_to_github(local_dir, repo_name, repo_clone_url):
                print("    Pushed to GitHub OK")
                execute("UPDATE landing_pages SET repo_name=%s, repo_url=%s, github_pushed=true WHERE id=%s",
                        (repo_name, repo_html_url, page_id))
            else:
                print("    Push failed!")
                continue
        else:
            # Repo may already exist — try pushing anyway
            print("    Repo may exist — trying push")
            existing_url = "https://github.com/" + GITHUB_USER + "/" + repo_name
            if push_to_github(local_dir, repo_name, existing_url + ".git"):
                repo_html_url = existing_url
                execute("UPDATE landing_pages SET repo_name=%s, repo_url=%s, github_pushed=true WHERE id=%s",
                        (repo_name, repo_html_url, page_id))
            else:
                print("    Skipped — push failed")
                continue

        # Mark as pushed — no Cloudflare API (you handle DNS + Pages manually)
        live_url = "https://" + slug + "." + CLIENT_BASE_DOMAIN
        execute("UPDATE landing_pages SET status='live', live_url=%s, deployed_at=%s WHERE id=%s",
                (live_url, time.strftime("%Y-%m-%d %H:%M:%S"), page_id))

        print("    Done:", repo_html_url)
        print("    Manual step: Connect repo to Cloudflare Pages + add CNAME", slug, "-> pages.dev")
        deployed += 1
        time.sleep(3)

    execute("INSERT INTO pipeline_runs (stage, status, targets_built, finished_at) VALUES (%s, %s, %s, NOW())",
            ("deploy", "success", deployed))

    print("=== DONE:", deployed, "pages pushed to GitHub ===")

if __name__ == "__main__":
    main()
