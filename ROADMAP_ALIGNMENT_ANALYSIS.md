# ORCA Roadmap Alignment Analysis

## Summary

After reviewing the README.md roadmap and comparing it to the created milestones and issues, there are **significant misalignments** that need correction.

## README Roadmap vs Current Milestones

| README Phase | Current Milestone | Status | Notes |
|--------------|-------------------|--------|-------|
| Phase 1: MVP (Months 1-3) | v0.1-mvp | ✅ Aligned | Core features already implemented |
| Phase 2: Production (Months 4-6) | Split between v0.1, v0.2, v0.3 | ⚠️ Misaligned | Container runtime in wrong phase |
| Phase 3: NRP Integration (Months 7-9) | v0.4-nrp | ⚠️ Misaligned | Should be v0.3, missing features |
| Phase 4: Advanced (Months 9+) | None | ❌ Missing | No milestone or issues created |

## Detailed Analysis

### Phase 1: MVP (Months 1-3) → v0.1-mvp ✅
**README Says:**
- Core Virtual Kubelet provider ✅ **DONE** (already implemented)
- Basic pod-to-EC2 mapping ✅ **DONE** (already implemented)
- Explicit instance selection ✅ **DONE** (already implemented)
- Simple lifecycle management ✅ **DONE** (already implemented)

**Current Issues:**
- None needed - phase is complete!

**Assessment:** This phase is correctly represented and already implemented.

---

### Phase 2: Production Features (Months 4-6) → Should be v0.2 ⚠️
**README Says:**
- GPU support (all NVIDIA types)
- Container runtime integration
- kubectl logs and exec
- Spot instance handling

**Current State:**
- ❌ Container runtime → Issue #8 in **v0.1-mvp** (WRONG MILESTONE)
- ❌ kubectl logs → Issue #9 in **v0.1-mvp** (WRONG MILESTONE)
- ❌ kubectl exec → Issue #10 in **v0.1-mvp** (WRONG MILESTONE)
- ✅ GPU support → Issue #12 in **v0.2-gpu** (CORRECT)
- ❌ Spot instance handling → **MISSING ENTIRELY**

**Problems:**
1. Container runtime, logs, exec are critical and were put in v0.1-mvp based on TODO.md priority
2. However, README roadmap places them in Phase 2
3. Spot instance handling is mentioned in README but has no issue
4. Current v0.2-gpu milestone is too narrow (only GPU features)

**Recommendation:**
- Rename v0.2-gpu to v0.2-production
- Move container runtime, logs, exec to v0.2-production
- Add spot instance handling issue to v0.2-production
- Keep GPU capacity reservations in v0.2-production

---

### Phase 3: NRP Integration (Months 7-9) → Should be v0.3 ⚠️
**README Says:**
- Ceph storage auto-mounting
- NRP namespace awareness
- Multi-tenancy support
- Cost tracking and budgets

**Current State:**
- ❌ Ceph storage → **MISSING ENTIRELY**
- ❌ NRP namespace awareness → **MISSING ENTIRELY**
- ✅ Multi-tenancy → Issue #17 in **v0.4-nrp** (WRONG MILESTONE)
- ✅ Cost tracking/budgets → Issue #17 in **v0.4-nrp** (WRONG MILESTONE)

**Problems:**
1. Used v0.4 instead of v0.3 for NRP features
2. Missing critical Ceph storage integration
3. Missing NRP namespace awareness

**Recommendation:**
- Rename v0.4-nrp to v0.3-nrp
- Add Ceph storage integration issue
- Add NRP namespace awareness issue
- Move budget tracking issue to v0.3-nrp

---

### Phase 4: Advanced Features (Months 9+) → Missing ❌
**README Says:**
- Advanced scheduling algorithms
- Capacity planning
- Compliance features (HIPAA, FedRAMP)
- Multi-region support

**Current State:**
- ❌ Advanced scheduling → **MISSING**
- ❌ Capacity planning → **MISSING**
- ❌ Compliance features → **MISSING**
- ❌ Multi-region support → **MISSING**

**Problems:**
1. No milestone exists for Phase 4
2. No issues created for any Phase 4 features

**Recommendation:**
- Create v0.4-advanced milestone
- Create issues for scheduling, capacity planning, compliance, multi-region

---

### Current v0.3-production Milestone → Orphaned ⚠️
**Current Issues:**
- Issue #14: Provider unit tests
- Issue #15: CI/CD GitHub Actions
- Issue #16: IRSA security

**Problems:**
1. These are "production readiness" tasks from TODO.md
2. They don't correspond to any README roadmap phase
3. They should be distributed across other phases or moved to Phase 2

**Recommendation:**
- Move testing/CI/CD to v0.2-production (needed before production use)
- Move IRSA to v0.2-production (security before production)
- Delete v0.3-production milestone

---

## What TODO.md Says vs README

The misalignment happened because TODO.md and README.md have different organization:

**TODO.md Organization:**
- "Critical - Required for MVP" (container runtime, logs, exec)
- "Critical - GPU/ML Features" (capacity reservations)
- "High Priority" (testing, CI/CD, security)
- "Medium Priority" (monitoring, cost management)

**README.md Organization:**
- Phase 1: Core provider (DONE)
- Phase 2: Production features (container runtime, GPU, logs/exec, spot)
- Phase 3: NRP integration (Ceph, multi-tenancy, cost tracking)
- Phase 4: Advanced features (scheduling, compliance, multi-region)

## Recommended Actions

### 1. Rename Milestones
```bash
# v0.2-gpu → v0.2-production
gh api -X PATCH repos/scttfrdmn/orca/milestones/2 -f title="v0.2-production"

# v0.4-nrp → v0.3-nrp
gh api -X PATCH repos/scttfrdmn/orca/milestones/4 -f title="v0.3-nrp"
```

### 2. Create v0.4-advanced Milestone
```bash
gh api repos/scttfrdmn/orca/milestones \
  -f title="v0.4-advanced" \
  -f description="Phase 4: Advanced Features (Months 9+)" \
  -f due_on="2027-06-30T07:00:00Z"
```

### 3. Move Existing Issues
- Move #8, #9, #10 from v0.1-mvp to v0.2-production
- Move #14, #15, #16 from v0.3-production to v0.2-production
- Move #17 from v0.4-nrp to v0.3-nrp
- Delete v0.3-production milestone

### 4. Create Missing Issues
**Phase 2 Missing:**
- Spot instance handling

**Phase 3 Missing:**
- Ceph storage auto-mounting
- NRP namespace awareness

**Phase 4 Missing:**
- Advanced scheduling algorithms
- Capacity planning
- Compliance features (HIPAA, FedRAMP)
- Multi-region support

### 5. Update Issue Script
Create new `scripts/create-project-issues-from-readme.sh` that follows README roadmap exactly.

## Conclusion

The current issue structure follows TODO.md technical priorities but doesn't match the README's user-facing roadmap timeline. We need to:

1. ✅ Keep Phase 1 as-is (complete)
2. ⚠️ Consolidate Phase 2 with all production features
3. ⚠️ Rename and populate Phase 3 with NRP features
4. ❌ Create Phase 4 for advanced features
5. 🗑️ Remove orphaned v0.3-production milestone

This will align GitHub milestones with the public roadmap in README.md, making it clear to users and contributors what features are planned for each release.
