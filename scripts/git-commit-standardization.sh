#!/bin/bash
# Script to commit all standardization changes

set -e

cd "$(dirname "$0")/.."

echo "🔍 Checking git status..."
git status --short

echo ""
echo "📝 Staging all changes..."
git add .

echo ""
echo "✅ Files staged. Creating commit..."

git commit -m "feat: Complete project standardization to match project-template-go

Major Changes:
============

## Documentation Site (MkDocs)
- Add mkdocs.yml with Material theme configuration
- Create comprehensive navigation structure (Getting Started, User Guide, Architecture, Development, API, Community)
- Add 30+ documentation pages (stubs ready for content)
- Custom CSS and JavaScript for ORCA branding
- Mermaid diagram support
- Light/dark mode toggle
- Full-text search

## GitHub Configuration
- Add 4 issue templates (Bug Report, Feature Request, Documentation, Technical Debt)
- Create labels.yml with 50+ organized labels (type, priority, area, use-case, status, milestones)
- Add comprehensive PR template
- Create GitHub Actions workflow for automatic docs deployment
- Enable GitHub Discussions integration

## Project Infrastructure
- Add .pre-commit-config.yaml with 11 hooks (Go, Markdown, YAML, Shell, Docker, Security)
- Create .goreleaser.yaml for multi-platform releases (Linux, macOS, Windows + Docker images)
- Enhance Makefile with docs-serve, docs-build, docs-deploy, pre-commit, release targets
- Add setup target for complete project initialization

## Assets
- Optimize logo images: 200px (35KB), 400px (154KB), 800px (710KB)
- Create icon suite: 32px (favicon), 64px, 128px (MkDocs), 256px, 512px (social)
- Update README with centered logo
- Add favicon and logo to MkDocs site

## Documentation
- Create GITHUB_SETUP_GUIDE.md with step-by-step GitHub configuration
- Add docs/SETUP_SUMMARY.md explaining all changes
- Create scripts/create-stub-docs.sh for generating doc pages

Technical Details:
=================
- MkDocs Material theme with custom ORCA blue branding
- GitHub Pages deployment via GitHub Actions
- Auto-generated changelog with semantic grouping
- Docker multi-arch images (amd64, arm64)
- Homebrew and AUR package support
- Cosign signing for releases
- Pre-commit hooks for code quality

Benefits:
========
- Professional documentation site at https://scttfrdmn.github.io/orca
- Consistent issue tracking and labeling
- Automated release process
- Code quality enforcement
- Community contribution guidelines
- Mobile-responsive docs with search

Next Steps:
==========
See GITHUB_SETUP_GUIDE.md for:
1. Enabling GitHub Pages
2. Applying labels with gh CLI
3. Setting up GitHub Discussions
4. Creating project board
5. Initial issue creation

Closes #N/A (initial standardization)"

echo ""
echo "✅ Commit created successfully!"
echo ""
echo "🚀 Ready to push with: git push origin main"
echo ""
echo "📖 After pushing, follow GITHUB_SETUP_GUIDE.md to complete GitHub configuration"
