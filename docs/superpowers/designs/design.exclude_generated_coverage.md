# Design: Excluding Generated Code from Test Coverage

**Date**: 2026-03-25
**Status**: Draft

## Background

The generated Go code (`*_enum.go`, `*_mocks.go`) currently lives in the same packages as source code (e.g., `pkg/types`, `pkg/plays`). The `go-test-coverage` tool can exclude these files at report time using file patterns, but this doesn't prevent them from being collected in the coverage profile.

This design explores three approaches to exclude generated code from test coverage more cleanly.

## Approaches

### B: Package Exclusion via `-coverpkg`

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

### C: Separate Packages for Generated Code

Move the enum type definitions AND generated code into dedicated subpackages like `pkg/types/gen/`, `pkg/plays/gen/`.

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

### D: Build Tags

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

| Aspect | B (coverpkg) | C (separate packages) | D (build tags) |
|--------|--------------|------------------------|-----------------|
| Implementation complexity | Low | High | Medium |
| Code restructuring | None | Significant | Medium |
| Maintenance burden | Medium | Low | Medium |
| Excludes at collection time | Yes | Yes | Yes |
| Import changes required | No | Yes | No |
| Tooling changes required | No | Yes | Yes |

## Recommendation

**Approach B (coverpkg)** is the simplest path forward given current codebase constraints.

However, if the long-term goal is cleaner package boundaries, **Approach C** provides the best separation but requires significant migration work.

## Open Questions

1. Should we maintain the exclusion list in the Makefile or in a config file?
2. How often do new generated packages get created? (affects maintenance burden of B)
3. Is the team open to restructuring imports for approach C?