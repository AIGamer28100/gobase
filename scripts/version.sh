#!/bin/bash
# Version management script for GoBase
# Usage: ./version.sh [major|minor|patch|prerelease] [alpha|beta|rc]

set -e

VERSION_FILE="cmd/gobase/version.go"
CURRENT_VERSION=$(grep -o 'Version.*=.*"v[^"]*"' "$VERSION_FILE" | sed 's/.*"v\([^"]*\)".*/\1/')

echo "Current version: v$CURRENT_VERSION"

# Function to increment version numbers
increment_version() {
    local version=$1
    local type=$2
    local prerelease=$3
    
    # Remove prerelease suffix if exists
    base_version=$(echo "$version" | sed 's/-.*$//')
    
    # Split into major.minor.patch
    IFS='.' read -ra VERSION_PARTS <<< "$base_version"
    major=${VERSION_PARTS[0]}
    minor=${VERSION_PARTS[1]}
    patch=${VERSION_PARTS[2]}
    
    case $type in
        "major")
            major=$((major + 1))
            minor=0
            patch=0
            ;;
        "minor")
            minor=$((minor + 1))
            patch=0
            ;;
        "patch")
            patch=$((patch + 1))
            ;;
        "prerelease")
            # Keep current version, just change prerelease
            ;;
        *)
            echo "Invalid version type. Use: major, minor, patch, or prerelease"
            exit 1
            ;;
    esac
    
    new_version="$major.$minor.$patch"
    
    # Add prerelease suffix if specified
    if [ -n "$prerelease" ]; then
        case $prerelease in
            "alpha"|"beta"|"rc")
                new_version="$new_version-$prerelease"
                ;;
            *)
                echo "Invalid prerelease type. Use: alpha, beta, or rc"
                exit 1
                ;;
        esac
    fi
    
    echo "$new_version"
}

# Parse arguments
if [ $# -eq 0 ]; then
    echo "Usage: $0 [major|minor|patch|prerelease] [alpha|beta|rc]"
    echo "Current version: v$CURRENT_VERSION"
    exit 0
fi

TYPE=$1
PRERELEASE=$2

NEW_VERSION=$(increment_version "$CURRENT_VERSION" "$TYPE" "$PRERELEASE")

echo "New version will be: v$NEW_VERSION"
read -p "Continue? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Cancelled."
    exit 0
fi

# Update version in the file
sed -i "s/Version.*=.*\"v[^\"]*\"/Version     = \"v$NEW_VERSION\"/" "$VERSION_FILE"

echo "✅ Version updated to v$NEW_VERSION"
echo "📝 Updated file: $VERSION_FILE"
echo ""
echo "Next steps:"
echo "1. Commit the changes: git add . && git commit -m \"Bump version to v$NEW_VERSION\""
echo "2. Create a tag: git tag v$NEW_VERSION"
echo "3. Push: git push origin main --tags"
