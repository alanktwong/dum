# Design: Excluding Generated Code from Test Coverage

**Date**: 2026-03-25
**Status**: Approved - Approach B (separate packages)

## Background

The generated Go code (`*_enum.go`, `*_mocks.go`) currently lives in the same packages as source 
code (e.g., `pkg/types`, `pkg/plays`). The `go-test-coverage` tool can exclude these files
at report time using file patterns, but this doesn't prevent them from being collected in
the coverage profile.

This design explores approaches to exclude generated code from test coverage more cleanly.

---

## Approaches

### A: Package Exclusion via `-coverpkg`

Modify the Makefile test target to exclude packages containing generated code at collection time.

```makefile
# Exclude packages with generated code from coverage collection
PACKAGES_TO_COVER := $(shell go list ./... | grep -v -E 'pkg/types|pkg/plays|pkg/tasks')

test:
    go test -coverpkg=./... ...
```

**Pros:**
- Excludes at collection time, not report time
- No code restructuring required
- Simple change to Makefile

**Cons:**
- Need to maintain list of excluded packages
- May exclude too much if package names change
- Doesn't work well if generated code is interspersed with source in same package
- Requires knowing exactly which packages contain generated code

---

### B: Separate Packages for Generated Code

Move the enum type definitions AND generated code into dedicated subpackages
like `pkg/types/gen/`, `pkg/plays/gen/`.

Structure:
```
pkg/
  types/
    gen/
      task_type.go       # source: type TaskType string
      task_type_enum.go  # generated
    task_type.go         # removed (moved to gen/)
```

Update all imports from `awong/dotfiles/pkg/types` to `awong/dotfiles/pkg/types/gen`.

**Pros:**
- True package-level exclusion via `-coverpkg`
- Clean separation of concerns
- Can exclude entire packages from coverage collection
- Generated code is never part of source package metrics

**Cons:**
- Widespread import changes (all files importing `pkg/types` need updating)
- Requires updating mockery and go-enum generation configs
- Moving type definitions changes the public API of the package
- Risk of breaking existing imports during transition

---

### C: Build Tags

Add build tags to generated files to exclude them at compilation time.

In `*_enum.go` and `*_mocks.go`:
```go
//go:build ignore_coverage
// +build ignore_coverage

package types
...
```

Then in Makefile:
```makefile
go test -cover -tags ignore_coverage ./...
```

**Pros:**
- Excludes at compilation time, not collection time
- Works with go-enum and mockery configurations

**Cons:**
- Requires modifying generation templates (go-enum, mockery)
- Build tags may have unintended side effects
- More complex to maintain
- May not work well with all tooling

---

## Comparison Matrix

| Aspect                      | A (coverpkg) | B (separate packages) | C (build tags) |
|-----------------------------|--------------|-----------------------|----------------|
| Implementation complexity   | Low          | High                  | Medium         |
| Code restructuring          | None         | Significant           | Medium         |
| Maintenance burden          | Medium       | Low                   | Medium         |
| Excludes at collection time | Yes          | Yes                   | Yes            |
| Import changes required     | No           | Yes                   | No             |
| Tooling changes required    | No           | Yes                   | Yes            |

## Recommendation

**Approach A (coverpkg)** is the simplest path forward given current codebase constraints.

However, if the long-term goal is cleaner package boundaries, **Approach B** provides the best
separation but requires significant migration work.

**Approach C (build tags)** is a middle ground but requires tooling changes.

## Open Questions

1. **Where to maintain exclusion list?** → Config file (not Makefile)
2. **How often are new generated packages created?** → Unpredictable; could stabilize or keep growing with new CLIs
3. **Open to restructuring imports for approach B?** → Yes
4. **Are tooling changes for approach C acceptable?** → No

---

## Decision

Given answers to (1) and (3), **Approach B (separate packages)** is recommended:
 
- Uses config file for exclusion list
- Clean package boundaries with generated code in dedicated subpackages
- Eliminates need to maintain exclusion list long-term
- Follows Go best practices for generated code organization