#!/bin/sh
set -eu

touch go.mod

PROJECT_NAME=suryakencana007
CURRENT_DIR=$(basename $(pwd))

CONTENT=$(cat <<-EOD
module github.com/${PROJECT_NAME}/${CURRENT_DIR}

require (
github.com/golang/mock v1.3.1
)
EOD
)

echo "$CONTENT" > go.mod
