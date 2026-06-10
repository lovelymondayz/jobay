#!/bin/bash
# Hermes workspace backup script
# Backs up /root/hermes to GitHub repo BackupHermes
# Keeps only last 30 days of backups
set -euo pipefail

TOKEN_FILE="/root/.github/pat"
GITHUB_USER="lovelymondayz"
REPO="BackupHermes"
SOURCE_DIR="/root/hermes"
BACKUP_DIR="/tmp/backup-workspace"
DATE="$(date +%Y-%m-%d)"
ARCHIVE="backup-${DATE}.zip"
RETENTION_DAYS=30
GIT_NAME="Hermes Bot"
GIT_EMAIL="hermes-bot@hermes.arjism.com"

log() { echo "[$(date '+%Y-%m-%d %H:%M:%S')] $*"; }

if [ ! -f "$TOKEN_FILE" ]; then log "ERROR: no token file"; exit 1; fi
TOKEN="$(cat $TOKEN_FILE | tr -d '[:space:]')"
AUTH_URL="http://${GITHUB_USER}:${TOKEN}@github.com/${REPO}.git"

# Clone or update the backup repo
if [ -d "$BACKUP_DIR/.git" ]; then
    log "Updating mirror..."
    cd "$BACKUP_DIR"
    git restore . 2>/dev/null || true
    git clean -fd 2>/dev/null || true
    git pull --rebase origin main 2>&1 || true
else
    log "Cloning mirror..."
    rm -rf "$BACKUP_DIR"
    git clone "$AUTH_URL" "$BACKUP_DIR" 2>&1
fi

# Set git identity for this repo
cd "$BACKUP_DIR"
git config user.name "$GIT_NAME"
git config user.email "$GIT_EMAIL"

# Remove today's existing archive if any
rm -f "$ARCHIVE"

# Check source exists and is non-empty
if [ ! -d "$SOURCE_DIR" ] || [ -z "$(ls -A "$SOURCE_DIR" 2>/dev/null)" ]; then
    log "Source empty, skipping"; exit 0
fi

# Create zip archive (exclude noise dirs)
log "Creating: $ARCHIVE"
cd "$SOURCE_DIR"
zip -r "${BACKUP_DIR}/${ARCHIVE}" . \
    -x "*.git/*" \
    -x "*/node_modules/*" \
    -x "*/__pycache__/*" \
    -x "*/.venv/*" \
    2>/dev/null

ARCHIVE_SIZE=$(du -h "${BACKUP_DIR}/${ARCHIVE}" | cut -f1)
log "Size: $ARCHIVE_SIZE"

# Prune backups older than 30 days
log "Pruning >${RETENTION_DAYS} day backups..."
cd "$BACKUP_DIR"
for old in backup-*.zip; do
    [ -f "$old" ] || continue
    od="${old#backup-}"
    od="${od%.zip}"
    oe=$(date -d "$od" +%s 2>/dev/null) || continue
    now=$(date +%s)
    age=$(( (now - oe) / 86400 ))
    if [ "$age" -gt "$RETENTION_DAYS" ]; then
        log "  Remove: $old ($age days)"
        rm -f "$old"
    fi
done

# Commit and push only if there are changes
cd "$BACKUP_DIR"
git add -A
if git diff --cached --quiet -- .; then
    log "No changes - skipping commit"; exit 0
fi

git commit -m "backup: $DATE - $ARCHIVE ($ARCHIVE_SIZE)" --author="$GIT_NAME <$GIT_EMAIL>" 2>&1
git push origin main 2>&1
log "Backup pushed: $ARCHIVE"
