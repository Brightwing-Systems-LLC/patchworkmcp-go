#!/bin/bash
# Usage: ./scripts/tag_version.sh <version>
# Example: ./scripts/tag_version.sh v0.2.0
#
# Go modules are published by pushing a version tag.
# No registry upload needed — `go get` fetches from the Git repo.

set -e

VERSION=${1}

if [[ -z "$VERSION" ]]; then
    echo "Error: Version is required"
    echo "Usage: $0 <version>"
    echo "Example: $0 v0.2.0"
    exit 1
fi

if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Error: Version must match vX.Y.Z format (e.g. v0.2.0)"
    exit 1
fi

echo "Tagging $VERSION..."
git tag "$VERSION"

echo "Tag created successfully!"
echo "Don't forget to:"
echo "1. Update CHANGELOG.md"
echo "2. git push origin main && git push origin --tags"
