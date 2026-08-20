#!/bin/bash
# Hermes — restore one project's database from a backup.
#
# Usage:  ./restore-one.sh --list
#         ./restore-one.sh qlio                    restore qlio from the newest backup
#         ./restore-one.sh qlio 20260819_040000    restore from a specific backup
#
# DESTRUCTIVE. Takes a safety dump of the current state first, then requires you
# to type RESTORE.

set -uo pipefail
BACKUP_ROOT=/root/hermes-backups

# name : container : user : db
MAP=(
  "wedding:wedding-db:wedding:wedding"
  "yearbook:yearbook-db:yearbook:yearbook"
  "members:members-db:members:members"
  "noidk:noidk-db:noidk:noidk"
  "qlio:qlio-db:qlio:qlio"
  "sayless:sayless-db:postgres:sayless"
  "pico:pico-postgres:pico:pico"
)

# SQLite services (new revenue pipeline)
SQLITE_PROJECTS=(
  "lead-lists"
  "content-packages"
  "web-audits"
  "email-sequences"
  "competitor-intel"
  "chatbot-deploy"
  "content-calendars"
  "invoice-extract"
  "local-seo"
  "resume-optimizer"
)

if [ "${1:-}" = "--list" ] || [ $# -eq 0 ]; then
  echo "Projects: wedding yearbook members noidk qlio sayless pico"
  echo "SQLite:   lead-lists content-packages web-audits email-sequences competitor-intel chatbot-deploy content-calendars invoice-extract local-seo resume-optimizer"
  echo
  echo "Available backups:"
  ls -1dt "$BACKUP_ROOT"/*/ 2>/dev/null | head -20 | while read -r d; do
    printf "  %-20s %-8s  dbs: %s\n" "$(basename "$d")" \
      "$(du -sh "$d" 2>/dev/null | cut -f1)" \
      "$(ls -1 "$d/databases" 2>/dev/null | sed 's/\.sql\.gz//' | tr '\n' ' ')"
  done
  exit 0
fi

# Check if it's a SQLite project
SQLITE_NAME=""
for sp in "${SQLITE_PROJECTS[@]}"; do
  [ "$sp" = "$1" ] && SQLITE_NAME="$1"
done

if [ -n "$SQLITE_NAME" ]; then
  # SQLite restore
  STAMP="${2:-}"
  if [ -n "$STAMP" ]; then
    SRC="$BACKUP_ROOT/$STAMP/databases/${SQLITE_NAME}.sql.gz"
  else
    SRC=$(ls -1t "$BACKUP_ROOT"/*/databases/"${SQLITE_NAME}".sql.gz 2>/dev/null | head -1)
  fi
  [ -f "$SRC" ] || { echo "No backup found for $SQLITE_NAME."; exit 1; }

  PDIR="/root/hermes/$SQLITE_NAME"
  DB="$PDIR/data/${SQLITE_NAME}.db"
  [ -d "$PDIR" ] || { echo "Project dir $PDIR not found."; exit 1; }

  echo "Project   : $SQLITE_NAME (SQLite)"
  echo "Backup    : $SRC"
  echo "            $(du -h "$SRC" | cut -f1), $(date -r "$SRC" '+%F %T')"
  echo
  echo "This REPLACES the current $SQLITE_NAME database."
  printf "Type RESTORE to continue: "
  read -r confirm
  [ "$confirm" = "RESTORE" ] || { echo "Aborted."; exit 1; }

  SAFETY="$BACKUP_ROOT/pre-restore_${SQLITE_NAME}_$(date +%Y%m%d_%H%M%S).sql.gz"
  mkdir -p "$BACKUP_ROOT"
  echo "→ safety dump of current state: $SAFETY"
  sqlite3 "$DB" .dump 2>/dev/null | gzip -9 > "$SAFETY"
  echo "  $(du -h "$SAFETY" | cut -f1)"

  echo "→ restoring…"
  # Drop all tables first
  sqlite3 "$DB" "SELECT 'DROP TABLE ' || name || ';' FROM sqlite_master WHERE type='table';" 2>/dev/null | sqlite3 "$DB" 2>/dev/null
  zcat "$SRC" | sqlite3 "$DB" 2>/tmp/hermes_restore_$SQLITE_NAME.log
  echo "  tables now: $(sqlite3 "$DB" "SELECT COUNT(*) FROM sqlite_master WHERE type='table';" 2>/dev/null)"

  echo
  echo "→ recreating the app container for a clean state"
  ( cd "$PDIR" && docker compose up -d --force-recreate 2>&1 | tail -4 )
  echo
  echo "Done. Safety dump kept at: $SAFETY"
  exit 0
fi

NAME="$1"
STAMP="${2:-}"

ENTRY=""
for m in "${MAP[@]}"; do [ "${m%%:*}" = "$NAME" ] && ENTRY="$m"; done
[ -z "$ENTRY" ] && { echo "Unknown project '$NAME'. Run --list."; exit 1; }
IFS=: read -r _ CNAME DUSER DNAME <<< "$ENTRY"

if [ -n "$STAMP" ]; then
  SRC="$BACKUP_ROOT/$STAMP/databases/$DNAME.sql.gz"
else
  SRC=$(ls -1t "$BACKUP_ROOT"/*/databases/"$DNAME".sql.gz 2>/dev/null | head -1)
fi
[ -f "$SRC" ] || { echo "No backup found for $NAME."; exit 1; }

docker ps --format '{{.Names}}' | grep -qx "$CNAME" || {
  echo "Container $CNAME is not running. Start the stack first."; exit 1; }

echo "Project   : $NAME"
echo "Container : $CNAME  (db=$DNAME user=$DUSER)"
echo "Backup    : $SRC"
echo "            $(du -h "$SRC" | cut -f1), $(date -r "$SRC" '+%F %T')"
echo
echo "This REPLACES the current $DNAME database."
printf "Type RESTORE to continue: "
read -r confirm
[ "$confirm" = "RESTORE" ] || { echo "Aborted."; exit 1; }

SAFETY="$BACKUP_ROOT/pre-restore_${DNAME}_$(date +%Y%m%d_%H%M%S).sql.gz"
mkdir -p "$BACKUP_ROOT"
echo "→ safety dump of current state: $SAFETY"
docker exec -i "$CNAME" pg_dump -U "$DUSER" -d "$DNAME" \
  --clean --if-exists --no-owner --no-privileges 2>/dev/null | gzip -9 > "$SAFETY"
echo "  $(du -h "$SAFETY" | cut -f1)"

echo "→ restoring…"
# ON_ERROR_STOP=0: DROP ... IF EXISTS lines are expected to error on a fresh db.
zcat "$SRC" | docker exec -i "$CNAME" psql -U "$DUSER" -d "$DNAME" -v ON_ERROR_STOP=0 \
  > /tmp/hermes_restore_$NAME.log 2>&1

echo "  psql errors: $(grep -c '^ERROR' /tmp/hermes_restore_$NAME.log || echo 0)  (DROP-on-missing is normal)"
echo "  tables now : $(docker exec -i "$CNAME" psql -U "$DUSER" -d "$DNAME" -tAc \
  "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='public'" 2>/dev/null)"

echo
echo "→ recreating the app container for a clean connection pool"
PDIR=$(docker inspect "$CNAME" --format '{{index .Config.Labels "com.docker.compose.project.working_dir"}}' 2>/dev/null)
if [ -n "$PDIR" ] && [ -d "$PDIR" ]; then
  ( cd "$PDIR" && docker compose up -d --force-recreate 2>&1 | tail -4 )
else
  echo "  could not resolve project dir — recreate the stack manually"
fi

echo
echo "Done. Verify by logging into the app: if login works, password hashes"
echo "survived and the restore is genuinely good. Safety dump kept at:"
echo "  $SAFETY"
