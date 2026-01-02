#!/bin/bash

# List versions for video/srs
# Usage: ./list-versions.sh [project] [app] [limit] [offset]

PROJECT=${1:-video}
APP=${2:-srs}
LIMIT=${3:-50}
OFFSET=${4:-0}
TOKEN="EpqQfkoT52yRNClG30KtokZljRxELHxopHYiYK5p_RE="
SERVER_URL="https://packages.slileisure.com"

echo "获取 ${PROJECT}/${APP} 的版本清单..."
echo ""

curl -s \
  -H "Authorization: Bearer ${TOKEN}" \
  "${SERVER_URL}/api/v1/projects/${PROJECT}/apps/${APP}/versions?limit=${LIMIT}&offset=${OFFSET}" \
  | python3 -m json.tool
