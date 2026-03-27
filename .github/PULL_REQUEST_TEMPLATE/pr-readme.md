# PR Template Guide

Use this guide to select the appropriate PR template for your change.

## Templates

| Template | Use When |
|----------|----------|
| **Brief** | Small changes (under 50 lines), cosmetic fixes, obvious bug fixes, documentation updates |
| **Build / CI** | Changes to build system, CI/CD, or tooling |
| **Challenge-Solution-Effect** | Architectural changes, significant trade-off decisions, changes that unlock new capabilities |
| **Situation-Action-Result** | Bug fixes where problem context matters, before/after state changes |
| **User Story** | Feature development with clear requirements and acceptance criteria |
| **What-Why-How** | Technical changes where context isn't obvious from the diff |

## Examples

### Brief
- Fixing a typo
- Adding an aria-label for accessibility
- Updating a dependency version

### Build / CI
- Adding a new CI pipeline
- Changing the build system
- Updating linting configuration

### Challenge-Solution-Effect
- Rewriting a core system
- Implementing a new caching strategy
- Changing the database architecture

### Situation-Action-Result
- Fixing a bug that users reported
- Addressing a security vulnerability
- Fixing data corruption

### User Story
- Adding a new feature
- Implementing a user request
- Building a new API endpoint

### What-Why-How
- Adding rate limiting
- Implementing a new algorithm
- Performance optimizations
