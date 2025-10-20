#!/bin/bash
# Script to set all project items to Todo status

set -e

PROJECT_NUMBER=3
OWNER="scttfrdmn"
PROJECT_ID="PVT_kwHOAC31Us4BF9CP"

echo "📋 Setting Project Items to Todo Status"
echo "========================================"
echo ""

# Get Status field ID and Todo option ID
FIELD_ID=$(gh project field-list "$PROJECT_NUMBER" --owner "$OWNER" --format json | jq -r '.fields[] | select(.name == "Status") | .id')
TODO_OPTION_ID=$(gh project field-list "$PROJECT_NUMBER" --owner "$OWNER" --format json | jq -r '.fields[] | select(.name == "Status") | .options[] | select(.name == "Todo") | .id')

echo "Project: #$PROJECT_NUMBER (ORCA Development)"
echo "Status Field ID: $FIELD_ID"
echo "Todo Option ID: $TODO_OPTION_ID"
echo ""

# Get all items and set them to Todo
updated=0
skipped=0

gh project item-list "$PROJECT_NUMBER" --owner "$OWNER" --format json --limit 100 | jq -r '.items[] | "\(.id)|\(.content.number)|\(.status // "none")"' | while IFS='|' read -r item_id issue_num current_status; do
    if [ "$current_status" = "Todo" ]; then
        echo "• Issue #$issue_num already in Todo"
        ((skipped++)) || true
    else
        echo "→ Setting issue #$issue_num to Todo (was: $current_status)"
        if gh project item-edit --id "$item_id" --project-id "$PROJECT_ID" --field-id "$FIELD_ID" --single-select-option-id "$TODO_OPTION_ID" 2>&1 >/dev/null; then
            echo "  ✓ Updated #$issue_num"
            ((updated++)) || true
        else
            echo "  ✗ Failed to update #$issue_num"
        fi
    fi
done

echo ""
echo "✅ Status update complete!"
echo ""
echo "🔗 View project: https://github.com/users/$OWNER/projects/$PROJECT_NUMBER"
