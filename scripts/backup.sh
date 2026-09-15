#!/bin/bash
# Zarrine Baft Backup Script (SRS 25)
# Backs up PostgreSQL database and product image uploads.
#
# Usage:
#   ./scripts/backup.sh
#   BACKUP_DIR=/custom/path ./scripts/backup.sh

set -euo pipefail

BACKUP_DIR="${BACKUP_DIR:-./backups}"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_PATH="${BACKUP_DIR}/${TIMESTAMP}"

# Docker container name
PG_CONTAINER="${PG_CONTAINER:-zarrine-postgres}"
PG_USER="${PG_USER:-postgres}"
PG_DB="${PG_DB:-zarrine_baft}"

# Uploads directory
UPLOADS_DIR="${UPLOADS_DIR:-./uploads}"

echo "=== Zarrine Baft Backup ==="
echo "Timestamp: ${TIMESTAMP}"
echo "Backup directory: ${BACKUP_PATH}"

mkdir -p "${BACKUP_PATH}"

# 1. Database backup
echo ""
echo "--- Backing up database ---"
if docker ps --format '{{.Names}}' | grep -q "^${PG_CONTAINER}$"; then
    docker exec "${PG_CONTAINER}" pg_dump -U "${PG_USER}" -d "${PG_DB}" --clean --if-exists \
        > "${BACKUP_PATH}/database.sql"
    echo "Database backup saved to: ${BACKUP_PATH}/database.sql"
else
    echo "WARNING: PostgreSQL container '${PG_CONTAINER}' is not running."
    echo "Attempting direct pg_dump (requires local psql)..."
    if command -v pg_dump &> /dev/null; then
        pg_dump -h 127.0.0.1 -U "${PG_USER}" -d "${PG_DB}" --clean --if-exists \
            > "${BACKUP_PATH}/database.sql"
        echo "Database backup saved to: ${BACKUP_PATH}/database.sql"
    else
        echo "ERROR: pg_dump not found. Skipping database backup."
    fi
fi

# 2. Uploads backup (local dev dir + compose named volume for prod)
echo ""
echo "--- Backing up uploads ---"
if [ -d "${UPLOADS_DIR}" ]; then
    tar -czf "${BACKUP_PATH}/uploads.tar.gz" -C "$(dirname "${UPLOADS_DIR}")" "$(basename "${UPLOADS_DIR}")"
    echo "Uploads (local dir) backup saved to: ${BACKUP_PATH}/uploads.tar.gz"
else
    echo "WARNING: Uploads directory '${UPLOADS_DIR}' not found. Skipping local backup."
fi

# Docker named volume used by docker-compose.yml in production
# (container path /app/uploads). Override with UPLOADS_VOLUME env if renamed.
UPLOADS_VOLUME="${UPLOADS_VOLUME:-zarrine-baft-backend_uploads_data}"
if docker volume inspect "${UPLOADS_VOLUME}" &> /dev/null; then
    # Host path in Windows form when running under Git Bash (cygpath),
    # plain PWD on Linux. NO_PATHCONV keeps container paths (/data, /backup)
    # from being mangled into Windows paths by MSYS2.
    HOST_DIR="${PWD}"
    if command -v cygpath &> /dev/null; then
        HOST_DIR="$(cygpath -m "${PWD}")"
    fi
    MSYS_NO_PATHCONV=1 docker run --rm \
        -v "${UPLOADS_VOLUME}:/data:ro" \
        -v "${HOST_DIR}/${BACKUP_PATH}:/backup" \
        alpine tar -czf /backup/uploads.volume.tar.gz -C /data .
    echo "Uploads (volume ${UPLOADS_VOLUME}) backup saved to: ${BACKUP_PATH}/uploads.volume.tar.gz"
else
    echo "WARNING: Docker volume '${UPLOADS_VOLUME}' not found. Skipping volume backup."
fi

# 3. Summary
echo ""
echo "=== Backup Complete ==="
echo "Location: ${BACKUP_PATH}"
ls -lh "${BACKUP_PATH}/"
echo ""
echo "To restore database:"
echo "  cat ${BACKUP_PATH}/database.sql | docker exec -i ${PG_CONTAINER} psql -U ${PG_USER} -d ${PG_DB}"
echo ""
echo "To restore uploads (local dir):"
echo "  tar -xzf ${BACKUP_PATH}/uploads.tar.gz -C ./"
echo ""
echo "To restore uploads (docker volume):"
echo "  docker run --rm -v ${UPLOADS_VOLUME}:/data -v \$(pwd)/${BACKUP_PATH}:/backup alpine tar -xzf /backup/uploads.volume.tar.gz -C /data"
