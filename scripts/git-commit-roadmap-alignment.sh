#!/bin/bash
# Script to commit roadmap alignment changes

set -e

cd "$(dirname "$0")/.."

echo "🎯 Committing Roadmap Alignment Changes"
echo "========================================"
echo ""

echo "📋 Changes summary:"
echo "  - ROADMAP_ALIGNMENT_ANALYSIS.md: Analysis of misalignment"
echo "  - ROADMAP_ALIGNMENT_COMPLETE.md: Completion summary"
echo "  - docs/index.md: Added roadmap section with milestone links"
echo ""

echo "📝 Staging alignment documentation..."
git add ROADMAP_ALIGNMENT_ANALYSIS.md
git add ROADMAP_ALIGNMENT_COMPLETE.md
git add docs/index.md

echo ""
echo "✅ Files staged. Creating commit..."

git commit -m "docs: Align project roadmap with GitHub milestones

Comprehensive roadmap alignment between README.md phases and GitHub milestones.

## What Changed

### Milestones
- Renamed v0.2-gpu → v0.2-production (Phase 2)
- Renamed v0.4-nrp → v0.3-nrp (Phase 3)
- Created v0.4-advanced (Phase 4)
- Deleted v0.3-production (orphaned)
- Updated all milestone descriptions to match README phases

### Issues
- Moved #8, #9, #10 from v0.1-mvp to v0.2-production
- Created #14: [Phase 2] Spot Instance Support
- Created #15: [Phase 3] Ceph Storage Auto-Mounting
- Created #16: [Phase 3] NRP Namespace Awareness
- Created #17: [Phase 4] Advanced Scheduling
- Created #18: [Phase 4] Capacity Planning
- Created #19: [Phase 4] Compliance Features (HIPAA/FedRAMP)
- Created #20: [Phase 4] Multi-Region Support
- Closed #1-7: Duplicate issues

### Documentation
- Updated docs/index.md with roadmap section
- Added links to GitHub milestones in documentation
- Created ROADMAP_ALIGNMENT_ANALYSIS.md
- Created ROADMAP_ALIGNMENT_COMPLETE.md

## Alignment Verified

✅ Phase 1: MVP - All README items complete or tracked
✅ Phase 2: Production - All 4 README items have issues
✅ Phase 3: NRP Integration - All 4 README items have issues
✅ Phase 4: Advanced - All 4 README items have issues

## Result

GitHub project structure now perfectly aligns with public-facing README roadmap.
Users and contributors can easily track progress by phase.

## Related

- README.md lines 207-232: Roadmap section
- GitHub Milestones: https://github.com/scttfrdmn/orca/milestones
- Project Board: https://github.com/scttfrdmn/orca/projects/3"

echo ""
echo "✅ Commit created successfully!"
echo ""
echo "🚀 Ready to push with: git push origin main"
echo ""
echo "📊 Milestone Summary:"
echo "  v0.1-mvp:        Phase 1 ✅ Complete"
echo "  v0.2-production: Phase 2 🚧 6 issues"
echo "  v0.3-nrp:        Phase 3 ⏳ 2 issues"
echo "  v0.4-advanced:   Phase 4 ⏳ 4 issues"
echo "  v1.0:            Release ⏳ TBD"
