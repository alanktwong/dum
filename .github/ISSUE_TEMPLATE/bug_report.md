---
name: Bug Report
about: Report something that isn't working
labels: bug
title: "[Bug] "
assignees: ''

body:
  - type: markdown
    attributes:
      value: "## Description"
  - type: textarea
    id: description
    attributes:
      label: "Bug Description"
      description: "A clear and concise description of what the bug is."
      placeholder: "Describe the bug here..."
    validations:
      required: true

  - type: markdown
    attributes:
      value: "## Steps to Reproduce"
  - type: textarea
    id: steps
    attributes:
      label: "Steps to Reproduce"
      description: "How do you trigger this bug?"
      placeholder: |
        1. Go to '...'
        2. Click on '...'
        3. See error
    validations:
      required: true

  - type: markdown
    attributes:
      value: "## Expected Behavior"
  - type: textarea
    id: expected
    attributes:
      label: "Expected Behavior"
      description: "What should happen?"
      placeholder: "Describe what you expected to happen..."
    validations:
      required: true

  - type: markdown
    attributes:
      value: "## Actual Behavior"
  - type: textarea
    id: actual
    attributes:
      label: "Actual Behavior"
      description: "What actually happens?"
      placeholder: "Describe what actually happened..."
    validations:
      required: true

  - type: markdown
    attributes:
      value: "## Environment"
  - type: textarea
    id: environment
    attributes:
      label: "Environment"
      description: "OS, Go version, dum version, etc."
      placeholder: |
        - OS: macOS 14.0
        - Go: 1.26.0
        - dum: latest
    validations:
      required: false

  - type: markdown
    attributes:
      value: "## Additional Context"
  - type: textarea
    id: additional
    attributes:
      label: "Additional Context"
      description: "Screenshots, error messages, logs, etc."
      placeholder: "Add any other context about the problem here..."
    validations:
      required: false

  - type: dropdown
    id: severity
    attributes:
      label: "Severity"
      options:
        - Critical
        - Major
        - Minor
    validations:
      required: false
