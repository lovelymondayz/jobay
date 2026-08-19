#!/bin/bash
# Hermes — push the newest backup off-site to GitHub.
#
# Usage:  ./backup-offsite.sh
# Cron:   30 3 * * * /root/hermes/scripts/backup-offsite.sh >> /var/log/hermes-backup.log 2>&1
#         (run AFTER backup-all.sh, which is at 03:00)
#
# Local backups do not survive losing the VPS or the disk. This pushes them to
# a PRIVATE GitHub repo so a total-loss recovery is possible.
#
# Contents include .env files and database dumps with customer PII — the repo
# MUST stay private. The script refuses to push if it is public.

set -uo pipefail

BACKUP_ROOT=/root/hermes-backups
REPO=BackupHermes
USER=lovelymondayz
PAT_FILE=/root/.github/pat
WORK=/tmp/hermes-offsite
KEEP=7            # keep the last N backups in the repo

log(){ echo "[$(date '+%F %T')] $*"; }

[ -f "$PAT_FILE" ] || { log "FATAL no PAT at $PAT_FILE"; exit 1; }
PAT=$(cat "$PAT_FILE")

NEWEST=$(ls -1dt "$BACKUP_ROOT"/2*/ 2>/dev/null | head -1)
[ -n "$NEWEST" ] || { log "FATAL no local backup to push — run backup-all.sh first"; exit 1; }
STAMP=$(basename "$NEWEST")

# Refuse to push secrets into a public repo.
VIS=$(curl -s -H "Authorization: token $PAT" \
  "https://api.github.com/repos/$USER/$REPO" | grep -oE '"private": *(true|false)' | head -1)
if echo "$VIS" | grep -q false; then
  log "FATAL $USER/$REPO is PUBLIC — refusing to push database dumps and .env files"
  exit 1
fi
if [ -z "$VIS" ]; then
  log "creating private repo $USER/$REPO"
  curl -s -X POST -H "Authorization: token $PAT" \
    https://api.github.com/user/repos \
    -d "{\"name\":\"$REPO\",\"private\":true,\"description\":\"Hermes VPS backups — databases and config. PRIVATE.\"}" \
    >/dev/null
fi

rm -rf "$WORK"; mkdir -p "$WORK"
log "cloning $REPO"
git clone --depth 1 "https://$USER:$PAT@github.com/$USER/$REPO.git" "$WORK" 2>&1 | tail -1

mkdir -p "$WORK/databases"
cp -r "$NEWEST" "$WORK/databases/$STAMP"
log "staged $STAMP ($(du -sh "$NEWEST" | cut -f1))"

# Retention inside the repo — keep history small enough to clone.
cd "$WORK/databases"
ls -1d 2*/ 2>/dev/null | sort -r | tail -n +$((KEEP+1)) | while read -r old; do
  log "pruning old remote backup $old"
  rm -rf "$old"
done

cd "$WORK"
cat > README.md <<EOF
# Hermes Backups — PRIVATE

Automated backups of every Docker project on the VPS.
**Contains database dumps and .env secrets. Never make this repo public.**

Latest: \`$STAMP\`
Updated: $(date '+%F %T %Z')

## Layout

\`\`\`
databases/<timestamp>/
  databases/*.sql.gz     pg_dump per project, gzipped
  config/<project>/      .env + docker-compose.yml
  MANIFEST.txt           what was captured, and how to restore
\`\`\`

## Restore

\`\`\`bash
# on a fresh VPS: clone each project repo, then
zcat databases/<stamp>/databases/qlio.sql.gz \\
  | docker exec -i qlio-db psql -U qlio -d qlio -v ON_ERROR_STOP=0
\`\`\`

Or use \`/root/hermes/scripts/restore-one.sh <project>\` on the VPS.

Retention: last $KEEP backups here, $((14)) days locally on the VPS.
EOF

git -c user.name="Hermes Bot" -c user.email="hermes-bot@arjism.com" add -A
if git diff --cached --quiet; then
  log "nothing changed"
  exit 0
fi
git -c user.name="Hermes Bot" -c user.email="hermes-bot@arjism.com" \
  commit -q -m "backup $STAMP"
git push origin HEAD 2>&1 | sed "s/$PAT/***/g" | tail -2
log "pushed $STAMP off-site"
rm -rf "$WORK"
