#!/bin/bash

# Script to help test the Go module publishing workflow
# This script creates a test tag and pushes it to trigger the publishing workflow

set -e

echo "🚀 Go Module Publishing Test Script"
echo "=================================="

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "❌ Error: Not in a git repository"
    exit 1
fi

# Check if we have uncommitted changes
if ! git diff-index --quiet HEAD --; then
    echo "❌ Error: You have uncommitted changes. Please commit or stash them first."
    echo "Uncommitted files:"
    git diff --name-only
    exit 1
fi

# Get current branch
CURRENT_BRANCH=$(git branch --show-current)
echo "📍 Current branch: $CURRENT_BRANCH"

# Check if we're on main branch
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo "⚠️  Warning: You're not on the main branch. Publishing should typically be done from main."
    read -p "Do you want to continue? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "❌ Aborted"
        exit 1
    fi
fi

# Get the latest version
echo "🔍 Checking existing tags..."
LATEST_TAG=$(git tag -l 'v*' | sort -V | tail -n 1)
if [ -z "$LATEST_TAG" ]; then
    echo "📝 No existing tags found. Starting with v0.1.0"
    NEW_VERSION="v0.1.0"
else
    echo "📝 Latest tag: $LATEST_TAG"
    
    # Extract version number and increment patch
    VERSION_NUMBER=${LATEST_TAG#v}
    IFS='.' read -ra VERSION_PARTS <<< "$VERSION_NUMBER"
    MAJOR=${VERSION_PARTS[0]}
    MINOR=${VERSION_PARTS[1]}
    PATCH=${VERSION_PARTS[2]}
    
    # Increment patch version
    NEW_PATCH=$((PATCH + 1))
    NEW_VERSION="v${MAJOR}.${MINOR}.${NEW_PATCH}"
fi

echo "🏷️  New version: $NEW_VERSION"

# Confirm before proceeding
read -p "Do you want to create and push tag $NEW_VERSION? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ Aborted"
    exit 1
fi

# Run tests before tagging
echo "🧪 Running tests..."
if ! go test ./...; then
    echo "❌ Tests failed. Please fix them before publishing."
    exit 1
fi

# Create and push the tag
echo "🏷️  Creating tag $NEW_VERSION..."
git tag -a "$NEW_VERSION" -m "Release $NEW_VERSION"

echo "📤 Pushing tag to remote..."
git push origin "$NEW_VERSION"

echo "✅ Tag $NEW_VERSION created and pushed successfully!"
echo ""
echo "🎉 The publishing workflow should now be triggered automatically."
echo "📊 Check the GitHub Actions tab to monitor the progress:"
echo "   https://github.com/control-theory/vercelreceiver/actions"
echo ""
echo "📦 Once published, your module will be available at:"
echo "   https://pkg.go.dev/github.com/control-theory/vercelreceiver@$NEW_VERSION"
echo ""
echo "🔍 You can also manually trigger the publishing workflow at:"
echo "   https://github.com/control-theory/vercelreceiver/actions/workflows/publish-module.yml"
