#!/usr/bin/env bash
set -euo pipefail

# Usage: ./deploy/deploy.sh user@host /opt/api

REMOTE=${1:?remote is required}
DEST=${2:-/opt/api}

rsync -avz --delete \
  bin/web-linux-amd64 \
  config.json \
  db/migrations \
  deploy/api.service \
  "${REMOTE}:${DEST}/"

ssh "${REMOTE}" "sudo cp ${DEST}/api.service /etc/systemd/system/api.service"
ssh "${REMOTE}" "sudo chmod +x ${DEST}/bin/web-linux-amd64"
ssh "${REMOTE}" "sudo systemctl daemon-reload"
ssh "${REMOTE}" "sudo systemctl enable --now api"
