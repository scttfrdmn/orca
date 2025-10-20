#!/bin/bash
# Script to sync ALL labels from labels.yml to GitHub

set -e

REPO="scttfrdmn/orca"
LABELS_FILE=".github/labels.yml"

echo "🏷️  Syncing ALL GitHub labels from $LABELS_FILE"
echo "================================================"
echo ""

# Check if yq is installed
if ! command -v yq &> /dev/null; then
    echo "❌ yq not found. Installing via brew..."
    brew install yq
fi

echo "📋 Reading labels from $LABELS_FILE..."
label_count=$(yq eval 'length' "$LABELS_FILE")
echo "Found $label_count labels to sync"
echo ""

# Parse YAML and create/update each label
yq eval '.[] | .name' "$LABELS_FILE" | while read -r name; do
    color=$(yq eval ".[] | select(.name == \"$name\") | .color" "$LABELS_FILE")
    description=$(yq eval ".[] | select(.name == \"$name\") | .description" "$LABELS_FILE")

    # Try to create or update the label
    if gh label create "$name" --color "$color" --description "$description" --repo "$REPO" 2>/dev/null; then
        echo "✓ Created: $name"
    else
        # Label exists, try to edit it
        if gh label edit "$name" --color "$color" --description "$description" --repo "$REPO" 2>/dev/null; then
            echo "↻ Updated: $name"
        else
            echo "• Exists: $name"
        fi
    fi
done

echo ""
echo "✅ Label sync complete!"
echo ""
echo "📊 Current label count:"
gh label list --repo "$REPO" | wc -l | xargs echo "  Total labels:"
echo ""
echo "🔗 View labels: https://github.com/$REPO/labels"
