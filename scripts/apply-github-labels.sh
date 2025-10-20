#!/bin/bash
# Script to apply GitHub labels from labels.yml using gh CLI

set -e

REPO="scttfrdmn/orca"
LABELS_FILE=".github/labels.yml"

echo "🏷️  Applying GitHub labels to $REPO"
echo ""

# Check if gh is authenticated
if ! gh auth status >/dev/null 2>&1; then
    echo "❌ GitHub CLI not authenticated. Run: gh auth login"
    exit 1
fi

# Parse YAML and create labels (requires yq or manual parsing)
# For now, we'll create the most important labels manually

echo "Creating type labels..."
gh label create "bug" --color "d73a4a" --description "Something isn't working correctly" --repo "$REPO" 2>/dev/null || echo "  - bug exists"
gh label create "enhancement" --color "a2eeef" --description "New feature or request" --repo "$REPO" 2>/dev/null || echo "  - enhancement exists"
gh label create "documentation" --color "0075ca" --description "Improvements or additions to documentation" --repo "$REPO" 2>/dev/null || echo "  - documentation exists"
gh label create "technical-debt" --color "fbca04" --description "Code refactoring or improvement needed" --repo "$REPO" 2>/dev/null || echo "  - technical-debt exists"

echo ""
echo "Creating priority labels..."
gh label create "priority: critical" --color "b60205" --description "Highest priority - blocking research work" --repo "$REPO" 2>/dev/null || echo "  - priority: critical exists"
gh label create "priority: high" --color "d93f0b" --description "High priority - should be addressed soon" --repo "$REPO" 2>/dev/null || echo "  - priority: high exists"
gh label create "priority: medium" --color "fbca04" --description "Medium priority - important but not urgent" --repo "$REPO" 2>/dev/null || echo "  - priority: medium exists"
gh label create "priority: low" --color "0e8a16" --description "Low priority - nice to have" --repo "$REPO" 2>/dev/null || echo "  - priority: low exists"

echo ""
echo "Creating area labels..."
gh label create "area: provider" --color "c5def5" --description "Virtual Kubelet provider" --repo "$REPO" 2>/dev/null || echo "  - area: provider exists"
gh label create "area: aws" --color "c5def5" --description "AWS integration" --repo "$REPO" 2>/dev/null || echo "  - area: aws exists"
gh label create "area: instances" --color "c5def5" --description "Instance selection" --repo "$REPO" 2>/dev/null || echo "  - area: instances exists"
gh label create "area: docs" --color "c5def5" --description "Documentation" --repo "$REPO" 2>/dev/null || echo "  - area: docs exists"

echo ""
echo "Creating status labels..."
gh label create "triage" --color "ededed" --description "Needs initial review and prioritization" --repo "$REPO" 2>/dev/null || echo "  - triage exists"
gh label create "ready" --color "0e8a16" --description "Ready to be worked on" --repo "$REPO" 2>/dev/null || echo "  - ready exists"
gh label create "in-progress" --color "fbca04" --description "Currently being worked on" --repo "$REPO" 2>/dev/null || echo "  - in-progress exists"

echo ""
echo "Creating special labels..."
gh label create "good first issue" --color "7057ff" --description "Good for newcomers" --repo "$REPO" 2>/dev/null || echo "  - good first issue exists"
gh label create "help wanted" --color "008672" --description "Extra attention needed from community" --repo "$REPO" 2>/dev/null || echo "  - help wanted exists"

echo ""
echo "Creating milestone labels..."
gh label create "milestone: v0.1-mvp" --color "bfd4f2" --description "Part of v0.1 MVP release" --repo "$REPO" 2>/dev/null || echo "  - milestone: v0.1-mvp exists"
gh label create "milestone: v0.2-gpu" --color "bfd4f2" --description "Part of v0.2 GPU features" --repo "$REPO" 2>/dev/null || echo "  - milestone: v0.2-gpu exists"

echo ""
echo "✅ Core labels created!"
echo ""
echo "📋 To see all labels: gh label list --repo $REPO"
echo ""
echo "💡 TIP: For complete label setup, install yq and use:"
echo "    brew install yq"
echo "    yq e '.[] | .name, .color, .description' .github/labels.yml"
