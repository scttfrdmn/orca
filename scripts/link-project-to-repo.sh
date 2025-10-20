#!/bin/bash
# Script to link GitHub Project to repository and add all issues

set -e

REPO="scttfrdmn/orca"
PROJECT_NUMBER=3
OWNER="scttfrdmn"

echo "🔗 Linking GitHub Project to Repository"
echo "========================================"
echo ""

echo "📋 Project: #$PROJECT_NUMBER (ORCA Development)"
echo "📦 Repository: $REPO"
echo ""

echo "1️⃣  Linking project to repository..."
gh project link "$PROJECT_NUMBER" --repo "$REPO" --owner "$OWNER" 2>/dev/null || echo "   (Already linked)"

echo ""
echo "2️⃣  Adding all open issues to project..."

# Get all open issues
issue_numbers=$(gh issue list --repo "$REPO" --state open --limit 100 --json number --jq '.[].number')

added=0
skipped=0

for issue in $issue_numbers; do
    if gh project item-add "$PROJECT_NUMBER" --owner "$OWNER" --url "https://github.com/$REPO/issues/$issue" 2>/dev/null; then
        echo "   ✓ Added issue #$issue"
        ((added++))
    else
        echo "   • Issue #$issue already in project"
        ((skipped++))
    fi
done

echo ""
echo "✅ Project linking complete!"
echo ""
echo "📊 Summary:"
echo "   Issues added: $added"
echo "   Already in project: $skipped"
echo ""
echo "🔗 View project: https://github.com/users/$OWNER/projects/$PROJECT_NUMBER"
echo "🔗 View in repo: https://github.com/$REPO/projects"
echo ""
echo "💡 Next steps:"
echo "   1. Configure custom fields (Priority, Area, Use Case)"
echo "   2. Enable workflow automations"
echo "   3. Create additional views (by priority, by area)"
