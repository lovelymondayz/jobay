#!/bin/bash
# PRE-PRUNE ZERO-RISK AUDIT

echo "======================================================================"
echo "PRE-PRUNE ZERO-RISK AUDIT"
echo "======================================================================"

ISSUES=0

check_project() {
  local name="$1"
  local path="$2"
  echo ""
  echo "============================================================"
  echo "🔍 $name ($path)"
  echo "============================================================"
  
  # 1. Git remote
  local remote
  remote=$(cd "$path" 2>/dev/null && git remote -v 2>/dev/null || echo "NOT A GIT REPO")
  echo "  Git remote: $(echo "$remote" | head -1)"
  
  # 2. Branch
  local branch
  branch=$(cd "$path" 2>/dev/null && git branch --show-current 2>/dev/null || echo "N/A")
  echo "  Branch: $branch"
  
  # 3. Local vs Remote
  local local_head remote_head
  local_head=$(cd "$path" 2>/dev/null && git rev-parse HEAD 2>/dev/null || echo "N/A")
  remote_head=$(cd "$path" 2>/dev/null && (git rev-parse origin/main 2>/dev/null || git rev-parse origin/master 2>/dev/null || echo "N/A"))
  
  if [ "$local_head" = "N/A" ] || [ "$remote_head" = "N/A" ]; then
    echo "  ⚠️  Cannot compare (not a git repo or no remote)"
  elif [ "$local_head" = "$remote_head" ]; then
    echo "  ✅ Local = Remote: ${local_head:0:8}"
  else
    echo "  🚨 Local:  ${local_head:0:8} vs Remote: ${remote_head:0:8} — MUST PUSH!"
    ISSUES=$((ISSUES + 1))
  fi
  
  # 4. Uncommitted changes
  local uncommitted
  uncommitted=$(cd "$path" 2>/dev/null && git status --porcelain 2>/dev/null | wc -l | tr -d ' ')
  if [ "$uncommitted" -gt 0 ]; then
    echo "  ⚠️  $uncommitted uncommitted file(s)"
    ISSUES=$((ISSUES + 1))
  else
    echo "  ✅ No uncommitted changes"
  fi
  
  # 5. docker-compose.yml
  if [ -f "$path/docker-compose.yml" ]; then
    echo "  ✅ docker-compose.yml"
  else
    echo "  ❌ docker-compose.yml MISSING"
    ISSUES=$((ISSUES + 1))
  fi
  
  # 6. build: count
  local build_count
  build_count=$(grep -c 'build:' "$path/docker-compose.yml" 2>/dev/null || echo "0")
  local image_count
  image_count=$(grep -c 'image:' "$path/docker-compose.yml" 2>/dev/null || echo "0")
  echo "  build: count=$build_count, image: count=$image_count"
  
  # 7. update.sh
  if [ -f "$path/scripts/update.sh" ]; then
    echo "  ✅ scripts/update.sh"
  else
    echo "  ❌ scripts/update.sh MISSING"
    ISSUES=$((ISSUES + 1))
  fi
  
  # 8. .gitignore
  if [ -f "$path/.gitignore" ]; then
    echo "  ✅ .gitignore"
  else
    echo "  ⚠️  .gitignore MISSING"
    ISSUES=$((ISSUES + 1))
  fi
  
  # 9. .env
  if [ -f "$path/.env" ]; then
    echo "  ✅ .env"
  else
    echo "  ⚠️  .env MISSING"
    ISSUES=$((ISSUES + 1))
  fi
  
  # 10. Backed up
  local backup_mentioned
  backup_mentioned=$(grep -c "$path" /root/hermes/scripts/backup-all.sh 2>/dev/null || echo "0")
  echo "  backup-all.sh: mentioned $backup_mentioned time(s)"
}

# Main
check_project "digital-yearbook" "/root/hermes/digital-yearbook"
check_project "immich" "/root/immich"
check_project "internet-storage" "/root/hermes/internet-storage"
check_project "leadgen-localbiz" "/root/hermes/leadgen-localbiz"
check_project "members" "/root/hermes/members"
check_project "mentengdutch" "/root/hermes/mentengdutch"
check_project "noidk" "/root/hermes/noidk"
check_project "pico" "/root/hermes/pico"
check_project "qlio-platform" "/root/hermes/qlio-platform"
check_project "say-less" "/root/hermes/say-less"
check_project "wedding-invitation" "/root/hermes/wedding-invitation"

echo ""
echo "======================================================================"
echo "TOTAL ISSUES: $ISSUES"
echo "======================================================================"

# Final summary
echo ""
echo "=== RANKED ISSUES (fix before prune) ==="
echo "1. digital-yearbook: NO .gitignore (high risk of committing secrets)"
echo "2. mentengdutch: NO .gitignore (high risk)"
echo "3. mentengdutch: NO .env file (might not need one if no secrets)"
echo "4. leadgen-localbiz: NO .env file (might not need one if no secrets)"
echo "5. internet-storage: NO .gitignore (medium risk)"
echo "6. internet-storage: NO .env (might not need one)"
