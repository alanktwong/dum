---
name: Feature Request
about: Suggest a new feature or improvement
labels: enhancement
title: "[Feature] "
assignees: ''

body:
  - type: markdown
    attributes:
      value: "## User Story"
  - type: textarea
    id: user-story
    attributes:
      label: "User Story"
      description: "As a [persona], I want [goal] so that [benefit]."
      placeholder: "As a user, I want to..."
    validations:
      required: true

  - type: markdown
    attributes:
      value: "## Problem or Motivation"
  - type: textarea
    id: motivation
    attributes:
      label: "Problem or Motivation"
      description: "What problem does this solve? Why is this needed?"
      placeholder: "Describe the problem or motivation..."
    validations:
      required: true

  - type: markdown
    attributes:
      value: "## Proposed Solution"
  - type: textarea
    id: solution
    attributes:
      label: "Proposed Solution"
      description: "How should this work? What's your approach?"
      placeholder: "Describe your proposed solution..."
    validations:
      required: true

  - type: markdown
    attributes:
      value: "## Alternatives Considered"
  - type: textarea
    id: alternatives
    attributes:
      label: "Alternatives Considered"
      description: "Other approaches you considered?"
      placeholder: "Describe any alternative solutions..."
    validations:
      required: false

  - type: markdown
    attributes:
      value: "## Acceptance Criteria"
  - type: textarea
    id: acceptance
    attributes:
      label: "Acceptance Criteria"
      description: "How do we know this is done?"
      placeholder: |
        - [ ] Criterion 1
        - [ ] Criterion 2
        - [ ] Criterion 3
    validations:
      required: true

  - type: markdown
    attributes:
      value: "## Additional Context"
  - type: textarea
    id: additional
    attributes:
      label: "Additional Context"
      description: "Sketches, mockups, references, etc."
      placeholder: "Add any other context here..."
    validations:
      required: false
