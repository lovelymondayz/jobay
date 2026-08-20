#!/bin/bash
# Hermes Fleet — Master Orchestration Script
#
# Usage: ./hermes.sh <command> [options]
#
# Commands:
#   start      Start all containers (or --project <name>)
#   stop       Stop all containers (or --project <name>)
#   restart    Restart all containers (or --project <name>)
#   rebuild    Rebuild from GitHub + recreate (runs scripts/update.sh)
#   status     Show running containers, ports, health
#   logs       Tail logs (or --project <name> --lines <n>)
#   validate   Check all projects have required files
#
# Options:
#   --project <name>   Target one project (e.g. --project wedding-invitation)
#   --all              Target all projects (default for most commands)
#   --lines <n>        Number of log lines (default 50)
#   --force            Force rebuild even if up to date
#
# Examples:
#   ./hermes.sh status
#   ./hermes.sh start --project wedding-invitation
#   ./hermes.sh rebuild --project noidk
#   ./hermes.sh rebuild --all
#   ./hermes.sh logs --project pico --lines 100

set -euo pipefail

# ─── Project Registry ────────────────────────────────────────
# Each project has its own docker-compose.yml and scripts/update.sh.
# This script ONLY orchestrates — it does NOT replace per-project files.
declare -A PROJECTS=(
  ["wedding-invitation"]="/root/hermes/wedding-invitation"
  ["digital-yearbook"]="/root/hermes/digital-yearbook"
  ["members"]="/root/hermes/members"
  ["noidk"]="/root/hermes/noidk"
  ["mentengdutch"]="/root/hermes/mentengdutch"
  ["pico"]="/root/hermes/pico"
  ["qlio-platform"]="/root/hermes/qlio-platform"
  ["say-less"]="/root/hermes/say-less"
  ["internet-storage"]="/root/hermes/internet-storage"
  ["immich"]="/root/immich"
)

# Ports (frontend/backend) for health checks
declare -A FRONTEND_PORTS=(
  ["wedding-invitation"]=3000
  ["digital-yearbook"]=3001
  ["members"]=3003
  ["noidk"]=3005
  ["mentengdutch"]=3002
  ["pico"]=3005
  ["qlio-platform"]=3007
  ["say-less"]=3006
  ["internet-storage"]=3004
  ["immich"]=2283
)

declare -A BACKEND_PORTS=(
  ["wedding-invitation"]=8080
  ["digital-yearbook"]=8081
  ["members"]=8082
  ["noidk"]=8085
  ["mentengdutch"]=""
  ["pico"]=8088
  ["qlio-platform"]=8087
  ["say-less"]=8086
  ["internet-storage"]=8084
  ["immich"]=2283
)

# Domains
declare -A DOMAINS=(
  ["wedding-invitation"]="wedding.arjism.com"
  ["digital-yearbook"]="yearbook.arjism.com"
  ["members"]="members.arjism.com"
  ["noidk"]="noidk.arjism.com"
  ["mentengdutch"]="menteng.arjism.com"
  ["pico"]="pico.arjism.com"
  ["qlio-platform"]="qlio.arjism.com"
  ["say-less"]="sayless.arjism.com"
  ["internet-storage"]=""
  ["immich"]="storage.arjism.com"
)

# ─── Helpers ──────────────────────────────────────────────────

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

get_projects() {
  if [[ -n "${TARGET_PROJECT:-}" ]]; then
    if [[ -z "${PROJECTS[$TARGET_PROJECT]+_}" ]]; then
      echo "❌ Unknown project: $TARGET_PROJECT"
      echo "   Known: ${!PROJECTS[*]}"
      exit 1
    fi
    echo "$TARGET_PROJECT"
  else
    for p in "${!PROJECTS[@]}"; do echo "$p"; done | sort
  fi
}

do_start() {
  local proj=$1 dir=${PROJECTS[$1]}
  echo "▶  Starting $proj..."
  (cd "$dir" && docker compose up -d)
}

do_stop() {
  local proj=$1 dir=${PROJECTS[$1]}
  echo "⏹  Stopping $proj..."
  (cd "$dir" && docker compose down)
}

do_restart() {
  local proj=$1 dir=${PROJECTS[$1]}
  echo "🔄 Restarting $proj..."
  (cd "$dir" && docker compose restart)
}

do_rebuild() {
  local proj=$1 dir=${PROJECTS[$1]}
  local update_script="$dir/scripts/update.sh"
  if [[ ! -f "$update_script" ]]; then
    echo "⚠️  $proj: no scripts/update.sh, skipping"
    return
  fi
  echo "🔨 Rebuilding $proj..."
  (cd "$dir" && bash "$update_script" ${FORCE:+--force})
}

do_status() {
  local proj=$1 dir=${PROJECTS[$1]}
  local fe_port=${FRONTEND_PORTS[$proj]:-""}
  local be_port=${BACKEND_PORTS[$proj]:-""}
  local domain=${DOMAINS[$proj]:-""}

  echo "┌─ $proj"
  echo "│  Dir:     $dir"
  [[ -n "$domain" ]] && echo "│  Domain:  $domain"

  # Container status
  local containers
  containers=$(cd "$dir" && docker compose ps --format '{{.Name}} {{.Status}}' 2>/dev/null || echo "  (none)")
  echo "│  Compose: $containers"

  # Health checks
  if [[ -n "$fe_port" ]]; then
    if curl -sf "http://localhost:${fe_port}/" > /dev/null 2>&1; then
      echo "│  Frontend: ✅ :$fe_port"
    else
      echo "│  Frontend: ❌ :$fe_port"
    fi
  fi
  if [[ -n "$be_port" && "$be_port" != "$fe_port" ]]; then
    if curl -sf "http://localhost:${be_port}/api/health" > /dev/null 2>&1 || \
       curl -sf "http://localhost:${be_port}/" > /dev/null 2>&1; then
      echo "│  Backend:  ✅ :$be_port"
    else
      echo "│  Backend:  ❌ :$be_port"
    fi
  fi
  echo "└─"
}

do_logs() {
  local proj=$1 dir=${PROJECTS[$1]}
  echo "📋 Logs for $proj (last ${LOG_LINES} lines):"
  (cd "$dir" && docker compose logs --tail="$LOG_LINES" 2>/dev/null || echo "  (no logs)")
}

do_validate() {
  echo "Validating all projects..."
  local ok=0 fail=0
  for proj in $(get_projects); do
    local dir=${PROJECTS[$proj]}
    local errors=""

    [[ ! -f "$dir/docker-compose.yml" ]] && errors+=" no-compose"
    [[ ! -f "$dir/scripts/update.sh" ]] && errors+=" no-update.sh"
    [[ ! -d "$dir/.git" ]] && errors+=" no-git"

    if [[ -z "$errors" ]]; then
      echo "  ✅ $proj"
      ok=$((ok + 1))
    else
      echo "  ❌ $proj:$errors"
      fail=$((fail + 1))
    fi
  done
  echo ""
  echo "Result: $ok OK, $fail issues"
  return $fail
}

# ─── Parse Args ───────────────────────────────────────────────
COMMAND=""
TARGET_PROJECT=""
ALL=false
LOG_LINES=50
FORCE=false

while [[ $# -gt 0 ]]; do
  case "$1" in
    start|stop|restart|rebuild|status|logs|validate)
      COMMAND="$1"
      shift
      ;;
    --project)
      TARGET_PROJECT="$2"
      shift 2
      ;;
    --all)
      ALL=true
      shift
      ;;
    --lines)
      LOG_LINES="$2"
      shift 2
      ;;
    --force)
      FORCE=true
      shift
      ;;
    -h|--help)
      sed -n '2,25p' "$0"
      exit 0
      ;;
    *)
      echo "Unknown arg: $1"
      exit 1
      ;;
  esac
done

if [[ -z "$COMMAND" ]]; then
  echo "Usage: $0 <command> [--project <name>|--all]"
  echo "Commands: start, stop, restart, rebuild, status, logs, validate"
  exit 1
fi

# ─── Execute ───────────────────────────────────────────────────
TARGET_PROJECT="${TARGET_PROJECT:-}"

case "$COMMAND" in
  start)    for p in $(get_projects); do do_start "$p"; done ;;
  stop)     for p in $(get_projects); do do_stop "$p"; done ;;
  restart)  for p in $(get_projects); do do_restart "$p"; done ;;
  rebuild)  for p in $(get_projects); do do_rebuild "$p"; done ;;
  status)   for p in $(get_projects); do do_status "$p"; echo; done ;;
  logs)     for p in $(get_projects); do do_logs "$p"; echo; done ;;
  validate) do_validate ;;
esac
