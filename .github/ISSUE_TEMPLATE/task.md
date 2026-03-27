---
name: Task / Chore
about: Maintenance, refactoring, or cleanup tasks
labels: task
title: "[Task] "
assignees: ''

body:
  - type: markdown
    attributes:
      value: "## Description"
  - type: textarea
    id: description
    attributes:
      label: "Task Description"
      description: "What needs to be done?"
      placeholder: "Describe the task..."
    validations:
      required: true

  - type: markdown
    attributes:
      value: "## Motivation / Background"
  - type: textarea
    id: motivation
    attributes:
      label: "Motivation / Background"
      description: "Why does this need to be done?"
      placeholder: "Why is this task needed?"
    validations:
      required: false

  - type: markdown
    attributes:
      value: "## Scope"
  - type: textarea
    id: scope
    attributes:
      label: "Scope"
      description: "What's in scope? What's out of scope?"
      placeholder: |
        **In scope:**
        - Item 1
        - Item 2

        **Out of scope:**
        - Item 3
        - Item 4
    validations:
      required: false

  - type: markdown
    attributes:
      value: "## Technical Notes"
  - type: textarea
    id: technical
    attributes:
      label: "Technical Notes"
      description: "Any technical considerations, files affected, etc."
      placeholder: "Files, dependencies, considerations..."
    validations:
      required: false

  - type: markdown
    attributes:
      value: "## Success Criteria"
  - type: textarea
    id: success
    attributes:
      label: "Success Criteria"
      description: "How do we know this is done?"
      placeholder: |
        - [ ] Criterion 1
        - [ ] Criterion 2
    validations:
      required: false

  - type: dropdown
    id: priority
    attributes:
      label: "Priority"
      options:
        - High
        - Medium
        - Low
    validations:
      required: false
