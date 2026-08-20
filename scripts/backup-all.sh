#!/bin/bash
# Hermes — backup every project's database and secrets.
#
# Usage:  ./backup-all.sh            back up everything
#         ./backup-all.sh --list     show what exists
#
# Cron:   0 3 * * * /root/hermes/scripts/backup-all.sh >> /var/log/hermes-backup.log 2>&1
#
# WHY THIS EXISTS
# Source code already lives in per-project GitHub repos, so code is recoverable.
# What is NOT recoverable is the data: Postgres volumes and the .env files that
# are gitignored on purpose. This script captures exactly that gap.
#
# Uses pg_dump, never a file copy of the volume — copying a live Postgres data
# directory produces a torn, often unrestorable snapshot.

set -uo pipefail

BACKUP_ROOT=/root/hermes-backups
RETENTION_DAYS=14
STAMP=$(date +%Y%m%d_%H%M%S)
OUT="$BACKUP_ROOT/$STAMP"

# container_name : db_user : db_name : project_dir
STACKS=(
  "wedding-db:wedding:wedding:/root/hermes/wedding-invitation"
  "yearbook-db:yearbook:yearbook:/root/hermes/digital-yearbook"
  "members-db:members:members:/root/hermes/members"
  "noidk-db:noidk:noidk:/root/hermes/noidk"
  "qlio-db:qlio:qlio:/root/hermes/qlio-platform"
  "sayless-db:postgres:sayless:/root/hermes/say-less"
  "pico-postgres:pico:pico:/root/hermes/pico"
)

# SQLite-based services (new revenue pipeline)
SQLITE_SERVICES=(
  "/root/hermes/lead-lists"
  "/root/hermes/content-packages"
  "/root/hermes/web-audits"
  "/root/hermes/email-sequences"
  "/root/hermes/competitor-intel"
  "/root/hermes/chatbot-deploy"
  "/root/hermes/content-calendars"
  "/root/hermes/invoice-extract"
  "/root/hermes/local-seo"
  "/root/hermes/resume-optimizer"
)

# Config-only projects (no database of their own)
CONFIG_ONLY=(
  "/root/hermes/mentengdutch"
  "/root/hermes/leadgen-localbiz"
  "/root/immich"
)

# Postgres volumes with no running container. pg_dump is impossible here, so we
# fall back to a tar of the data directory. That is only safe BECAUSE the
# cluster is stopped — never do this to a running Postgres.
OFFLINE_VOLUMES=(
  "leadgen_pgdata"
)

log(){ echo "[$(date '+%F %T')] $*"; }

if [ "${1:-}" = "--list" ]; then
  echo "Backups in $BACKUP_ROOT:"
  ls -1dt "$BACKUP_ROOT"/*/ 2>/dev/null | while read -r d; do
    printf "  %-22s %s\n" "$(basename "$d")" "$(du -sh "$d" 2>/dev/null | cut -f1)"
  done
  exit 0
fi

mkdir -p "$OUT"/{databases,config}
log "backup → $OUT"

OK=0; FAIL=0; SKIP=0

# ---------- databases ----------
for entry in "${STACKS[@]}"; do
  IFS=: read -r cname duser dname pdir <<< "$entry"

  if ! docker ps --format '{{.Names}}' 2>/dev/null | grep -qx "$cname"; then
    log "SKIP $cname (not running — cannot dump a stopped database)"
    SKIP=$((SKIP+1)); continue
  fi

  f="$OUT/databases/${dname}.sql.gz"
  # -i (not -T): `docker exec` has no -T flag — that belongs to `docker compose
  # exec`. Using -T here fails with "unknown shorthand flag" and, without the
  # size guard below, would silently write a 20-byte "backup" every night.
  if docker exec -i "$cname" pg_dump -U "$duser" -d "$dname" \
        --clean --if-exists --no-owner --no-privileges 2>/dev/null | gzip -9 > "$f"; then
    sz=$(stat -c%s "$f" 2>/dev/null || echo 0)
    # A gzipped dump of even an empty schema is ~400B (header + SET statements).
    # Below 300B means pg_dump wrote nothing at all — wrong user/db, or the
    # container rejected the command. Verified floor: say-less has a single
    # table and dumps to 867B, so do not raise this without checking.
    if [ "$sz" -lt 300 ]; then
      log "FAIL $dname (dump only ${sz}B — wrong credentials?)"
      rm -f "$f"; FAIL=$((FAIL+1)); continue
    fi
    if ! gzip -t "$f" 2>/dev/null; then
      log "FAIL $dname (gzip corrupt)"
      rm -f "$f"; FAIL=$((FAIL+1)); continue
    fi
    log "OK   $dname ($(numfmt --to=iec "$sz" 2>/dev/null || echo "${sz}B"))"
    OK=$((OK+1))
  else
    log "FAIL $dname (pg_dump error)"
    rm -f "$f"; FAIL=$((FAIL+1))
  fi
done

# ---------- offline volumes (tar fallback) ----------
for vol in "${OFFLINE_VOLUMES[@]}"; do
  vpath="/var/lib/docker/volumes/$vol/_data"
  [ -d "$vpath" ] || { log "SKIP $vol (volume absent)"; SKIP=$((SKIP+1)); continue; }

  # Refuse if something is actually using it — a tar of a live cluster is unsafe.
  if docker ps -q --filter "volume=$vol" 2>/dev/null | grep -q .; then
    log "SKIP $vol (in use — dump it via pg_dump instead)"
    SKIP=$((SKIP+1)); continue
  fi

  f="$OUT/databases/${vol}-coldcopy.tar.gz"
  if tar -czf "$f" -C "$(dirname "$vpath")" "$(basename "$vpath")" 2>/dev/null; then
    log "OK   $vol cold copy ($(du -h "$f" | cut -f1))"
    OK=$((OK+1))
  else
    log "FAIL $vol tar"
    rm -f "$f"; FAIL=$((FAIL+1))
  fi
done

# ---------- immich (bind-mounted, not a docker volume) ----------
# Immich's photo library is 2.6GB and lives at /home/root/immich-data/upload —
# far too large for the off-site git repo, so only the database is dumped here.
# The photos need a separate rclone target; see docs/BACKUP.md.
if docker ps --format '{{.Names}}' 2>/dev/null | grep -qx immich_postgres; then
  f="$OUT/databases/immich.sql.gz"
  if docker exec -i immich_postgres pg_dump -U immich -d immich \
        --clean --if-exists --no-owner --no-privileges 2>/dev/null | gzip -9 > "$f"; then
    sz=$(stat -c%s "$f" 2>/dev/null || echo 0)
    if [ "$sz" -lt 300 ]; then
      log "FAIL immich (${sz}B)"; rm -f "$f"; FAIL=$((FAIL+1))
    else
      log "OK   immich ($(numfmt --to=iec "$sz" 2>/dev/null || echo "${sz}B"))"
      OK=$((OK+1))
    fi
  else
    log "FAIL immich (pg_dump error)"; rm -f "$f"; FAIL=$((FAIL+1))
  fi
else
  log "SKIP immich (not running)"; SKIP=$((SKIP+1))
fi

# ---------- config / secrets ----------
# .env files are gitignored by design, so they exist nowhere else. Without them
# a restored dump cannot even be loaded (wrong DB password) and every JWT
# session is invalidated.
for entry in "${STACKS[@]}"; do
  IFS=: read -r _ _ _ pdir <<< "$entry"
  CONFIG_ONLY+=("$pdir")
done

# Also backup .env files for SQLite services
for pdir in "${SQLITE_SERVICES[@]}"; do
  CONFIG_ONLY+=("$pdir")
done

for pdir in "${CONFIG_ONLY[@]}"; do
  [ -d "$pdir" ] || continue
  n=$(basename "$pdir")
  mkdir -p "$OUT/config/$n"
  for f in .env .env.local docker-compose.yml; do
    [ -f "$pdir/$f" ] && cp "$pdir/$f" "$OUT/config/$n/" 2>/dev/null
  done
done
log "config captured for $(ls -1 "$OUT/config" | wc -l) projects"

# ---------- SQLite databases (new revenue pipeline) ----------
for pdir in "${SQLITE_SERVICES[@]}"; do
  [ -d "$pdir" ] || continue
  n=$(basename "$pdir")
  mkdir -p "$OUT/databases"
  db="$pdir/data/$n.db"
  if [ -f "$db" ]; then
    f="$OUT/databases/${n}.sql.gz"
    if sqlite3 "$db" .dump 2>/dev/null | gzip -9 > "$f"; then
      sz=$(stat -c%s "$f" 2>/dev/null || echo 0)
      if [ "$sz" -lt 50 ]; then
        log "FAIL $n SQLite dump empty (${sz}B)"; rm -f "$f"; FAIL=$((FAIL+1))
      else
        log "OK   $n SQLite ($(numfmt --to=iec "$sz" 2>/dev/null || echo "${sz}B"))"; OK=$((OK+1))
      fi
    else
      log "FAIL $n sqlite3 dump error"; rm -f "$f"; FAIL=$((FAIL+1))
    fi
  else
    log "SKIP $n (no data/${n}.db file)"; SKIP=$((SKIP+1))
  fi
done

# ---------- manifest ----------
{
  echo "Hermes backup — $STAMP"
  echo
  echo "## Databases"
  ls -lh "$OUT/databases" 2>/dev/null | tail -n +2 | awk '{print "  "$9"  "$5}'
  echo
  echo "## Running containers at backup time"
  docker ps --format '  {{.Names}}\t{{.Status}}' 2>/dev/null
  echo
  echo "## Restore"
  echo "  zcat databases/<name>.sql.gz | docker exec -i <container> psql -U <user> -d <db>"
  echo "  See docs/RESTORE.md in this directory."
} > "$OUT/MANIFEST.txt"

chmod -R 600 "$OUT"/config/*/.env 2>/dev/null || true

# ---------- retention ----------
find "$BACKUP_ROOT" -maxdepth 1 -type d -name '20*' -mtime +$RETENTION_DAYS \
  -exec rm -rf {} \; 2>/dev/null

log "done: $OK ok, $FAIL failed, $SKIP skipped — $(du -sh "$OUT" | cut -f1)"
[ "$FAIL" -gt 0 ] && exit 1
exit 0