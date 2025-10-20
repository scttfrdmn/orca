# GitHub Setup Guide for ORCA

This guide walks you through completing the ORCA project setup on GitHub, including enabling GitHub Pages, applying labels, and initializing the project board.

## ✅ Completed

The following has already been configured locally:

### 1. GitHub Pages Configuration
- ✅ `mkdocs.yml` - Complete MkDocs Material configuration
- ✅ `.github/workflows/docs.yml` - Auto-deploy workflow
- ✅ `docs/` structure with 30+ pages
- ✅ Custom styling and JavaScript
- ✅ Optimized logo and icon assets

### 2. Issue Templates
- ✅ Bug Report template (`.github/ISSUE_TEMPLATE/bug_report.yml`)
- ✅ Feature Request template (`.github/ISSUE_TEMPLATE/feature_request.yml`)
- ✅ Documentation template (`.github/ISSUE_TEMPLATE/documentation.yml`)
- ✅ Technical Debt template (`.github/ISSUE_TEMPLATE/technical_debt.yml`)
- ✅ Config with discussion links (`.github/ISSUE_TEMPLATE/config.yml`)

### 3. Labels Configuration
- ✅ `.github/labels.yml` with 50+ labels organized by:
  - Type (bug, enhancement, documentation, technical-debt)
  - Priority (critical, high, medium, low)
  - Area (provider, aws, instances, config, etc.)
  - Use Case (gpu-training, cost-optimization, etc.)
  - Status (triage, ready, in-progress, etc.)
  - Milestones (v0.1-mvp, v0.2-gpu, v0.3-production, etc.)

### 4. Pull Request Template
- ✅ `.github/pull_request_template.md` - Comprehensive PR checklist

### 5. Project Infrastructure
- ✅ `.pre-commit-config.yaml` - 11 pre-commit hooks
- ✅ `.goreleaser.yaml` - Multi-platform release automation
- ✅ Enhanced Makefile with docs, release, and setup targets

## 📋 Next Steps

Complete these steps in order:

### Step 1: Commit and Push Changes

```bash
cd /Users/scttfrdmn/src/orca

# Check status
git status

# Add all new files
git add .

# Commit
git commit -m "feat: Complete project standardization

- Add MkDocs documentation site with Material theme
- Create GitHub issue templates (bug, feature, docs, tech debt)
- Add comprehensive label configuration (50+ labels)
- Create PR template
- Add pre-commit hooks configuration
- Configure GoReleaser for multi-platform releases
- Enhance Makefile with docs and release targets
- Optimize logo and icon assets (multiple sizes)
- Create 30+ documentation stub pages
- Add GitHub Actions workflow for docs deployment

Closes #X (if applicable)"

# Push to GitHub
git push origin main
```

### Step 2: Enable GitHub Pages

1. Go to your repository on GitHub: https://github.com/scttfrdmn/orca
2. Click **Settings** → **Pages**
3. Under "Build and deployment":
   - Source: **GitHub Actions**
4. Click **Save**
5. Wait 2-3 minutes for the first deployment
6. Visit: https://scttfrdmn.github.io/orca

### Step 3: Apply GitHub Labels

```bash
# Install GitHub CLI if not already installed
# brew install gh

# Authenticate (if needed)
gh auth login

# Apply labels from configuration file
gh label sync --labels .github/labels.yml --force

# Verify labels were created
gh label list
```

This will create all 50+ labels with proper colors and descriptions.

### Step 4: Enable GitHub Discussions

1. Go to repository **Settings** → **General**
2. Scroll to **Features**
3. Check **✓ Discussions**
4. Click **Set up discussions**
5. Create categories:
   - **Q&A** - Ask questions
   - **Ideas** - Feature discussions
   - **Show and Tell** - Share projects
   - **General** - General discussions

### Step 5: Create GitHub Project (Optional but Recommended)

#### Option A: Using GitHub Web UI

1. Go to repository **Projects** tab
2. Click **New project**
3. Select **Board** template
4. Name: **ORCA Development**
5. Add custom columns:
   - 📋 Backlog
   - 🔖 Ready
   - 🚀 In Progress
   - 👀 In Review
   - ✅ Done

#### Option B: Using GitHub CLI

```bash
# Create project
gh project create --owner scttfrdmn --title "ORCA Development"

# Add fields (customize as needed)
gh project field-create --owner scttfrdmn --project "ORCA Development" \
  --name "Priority" --data-type SINGLE_SELECT \
  --single-select-options "Critical,High,Medium,Low"

gh project field-create --owner scttfrdmn --project "ORCA Development" \
  --name "Area" --data-type SINGLE_SELECT \
  --single-select-options "Provider,AWS,Instances,Config,Docs,Tests"
```

### Step 6: Configure Project Automation

Add automation rules:

1. Go to project → **⋯** → **Workflows**
2. Enable:
   - **Item added to project** → Move to "Backlog"
   - **Item closed** → Move to "Done"
   - **Pull request merged** → Move to "Done"
   - **Item reopened** → Move to "Ready"

### Step 7: Create Initial Issues

Create issues for the MVP work:

```bash
# Example: Create an issue for MVP work
gh issue create \
  --title "Implement container runtime integration" \
  --body "Need to add container runtime (containerd) to EC2 instances for running actual workloads." \
  --label "enhancement,priority: critical,area: provider,milestone: v0.1-mvp"

# You can also use the web UI with the new templates
```

### Step 8: Setup Branch Protection (Recommended)

```bash
# Protect main branch
gh api repos/scttfrdmn/orca/branches/main/protection \
  --method PUT \
  --field required_status_checks[strict]=true \
  --field required_status_checks[contexts][]=build \
  --field required_status_checks[contexts][]=test \
  --field enforce_admins=false \
  --field required_pull_request_reviews[required_approving_review_count]=1 \
  --field required_pull_request_reviews[dismiss_stale_reviews]=true

# Or use web UI: Settings → Branches → Add rule
```

## 📊 Verify Setup

After completing all steps, verify:

### Documentation Site
- [ ] Visit https://scttfrdmn.github.io/orca
- [ ] Check navigation works
- [ ] Verify logo displays correctly
- [ ] Test search functionality
- [ ] Check light/dark mode toggle

### Labels
```bash
gh label list | wc -l  # Should show 50+ labels
```

### Issue Templates
1. Go to **Issues** → **New issue**
2. Verify you see 4 templates:
   - 🐛 Bug Report
   - ✨ Feature Request
   - 📚 Documentation
   - 🔧 Technical Debt

### Discussions
- [ ] Visit Discussions tab
- [ ] Verify categories are present
- [ ] Test creating a discussion

### Projects
- [ ] Visit Projects tab
- [ ] Verify project board exists
- [ ] Test adding an issue to the project

## 🎯 Recommended First Issues

Create these issues to populate your project board:

1. **Critical MVP Issues:**
   - Container runtime integration
   - kubectl logs support (via CloudWatch)
   - kubectl exec support (via Systems Manager)
   - GPU capacity reservations support

2. **Documentation Issues:**
   - Complete user guide pages
   - Add architecture diagrams
   - Create example manifests
   - Write troubleshooting guide

3. **Good First Issues:**
   - Add more unit tests
   - Improve error messages
   - Add godoc comments
   - Fix broken documentation links

## 📈 Project Management Workflow

Suggested workflow:

1. **Issues** → Discussion → Triage → Label → Add to Project
2. **Project Board** → Backlog → Ready → In Progress → Review → Done
3. **Pull Requests** → Linked to issues → Review → Merge → Auto-close issue

## 🔄 Ongoing Maintenance

Weekly tasks:

```bash
# Sync labels if you update .github/labels.yml
gh label sync --labels .github/labels.yml --force

# Triage new issues
gh issue list --label triage

# Review PRs awaiting review
gh pr list --label "in-review"

# Check documentation build
make docs-build
```

## 🚀 Quick Command Reference

```bash
# Documentation
make docs-serve          # Serve locally at http://localhost:8000
make docs-build          # Build documentation
make docs-deploy         # Manual deploy to GitHub Pages

# Labels
gh label sync --labels .github/labels.yml --force

# Issues
gh issue list
gh issue create --template feature_request.yml
gh issue view 123

# Projects
gh project list --owner scttfrdmn
gh project item-list --owner scttfrdmn --project "ORCA Development"

# PRs
gh pr list
gh pr create
gh pr view 123
```

## 📚 Additional Resources

- [GitHub Projects Documentation](https://docs.github.com/en/issues/planning-and-tracking-with-projects)
- [GitHub Actions for Pages](https://docs.github.com/en/pages/getting-started-with-github-pages/using-custom-workflows-with-github-pages)
- [MkDocs Material](https://squidfunk.github.io/mkdocs-material/)
- [GitHub CLI Manual](https://cli.github.com/manual/)

---

**Questions?** Open a discussion or create an issue!
