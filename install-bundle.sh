#!/bin/bash

# Constants
REPOSITORY=${REPOSITORY:-"agnosticeng/clickhouse-evm"}
TAG=${TAG:-"latest"}
OS_ARCH="$(uname | tr '[:upper:]' '[:lower:]')_$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')"

case "$TAG" in
 "latest") API_URL="https://api.github.com/repos/$REPOSITORY/releases/latest";;
 *) API_URL="https://api.github.com/repos/$REPOSITORY/releases/tags/$TAG";;
esac


# Fetch release info
RELEASE_JSON=$(wget -qO- "$API_URL")
RELEASE_TAG=$(echo "$RELEASE_JSON" | grep -m1 '"tag_name":' | cut -d '"' -f4)
BUNDLE_URL=$(echo "$RELEASE_JSON" | grep -oE "https.*$OS_ARCH.*\.tar\.gz")

if [ -z "$RELEASE_TAG" ] || [ -z "$BUNDLE_URL" ]; then
    echo "Failed to fetch release tag or download URL."
    exit 1
fi

echo "Installting bundle from $BUNDLE_URL"
su - clickhouse -c "wget -qO- $BUNDLE_URL | tar xvz -C /"

echo "Reloading binary UDF functions"
clickhouse client --query="system reload functions"

echo "Realoding SQL UDF functions"
for f in /var/lib/clickhouse/user_defined/*.sql; do clickhouse client --queries-file $f; done