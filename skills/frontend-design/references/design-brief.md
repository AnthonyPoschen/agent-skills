# `DESIGN.md` Template

Use `DESIGN.md` as the project's durable record of design decisions. Read it
before a material design consultation and update it after new evidence changes
the direction. Record decisions and reasons, not raw conversation notes.

For a new user-facing flow, substantial redesign, or work handed to another
implementer, add a brief using the template below. Remove headings that do not
affect the decision. Write unknowns as questions and label assumptions so they
can be revised cheaply.

## Start with one document

Use this structure while the design record is still small:

```markdown
# Design

## Shared decisions

- [Decision]: [reason and scope]

## Feature briefs

### Account settings

[Brief]

### Invoice dashboard

[Brief]
```

## Split when the document no longer scans

When `DESIGN.md` exceeds about 1,000 lines or its areas have independent
owners, questions, or change histories, keep `DESIGN.md` as a router:

```markdown
# Design

This file indexes the project's design knowledge. Read the relevant brief
before changing a user-facing flow.

## Shared decisions

- [Decision]: [reason and scope]

## Design briefs

- [Account settings](design/account-settings.md)
- [Invoice dashboard](design/invoice-dashboard.md)
```

Use lowercase kebab-case names in `design/`. Each child document begins with a
specific title and contains one brief. Keep cross-cutting decisions, the index,
and a short explanation of the structure in the root `DESIGN.md`.

## Brief template

## Outcome

- **User outcome:**
- **Primary action or decision:**
- **Success signal:**
- **Why now:**

## Audience And Context

- **Primary users and their familiarity:**
- **Where and when they use it:**
- **Devices, viewports, and input methods:**
- **Accessibility needs:**

## Workflow

1. **Entry point:**
2. **Happy path:**
3. **Alternatives and edge cases:**
4. **Feedback after success, failure, or interruption:**

## Content, Data, And States

- **Information that must be visible first:**
- **Supporting information:**
- **Realistic data shape and density:**
- **Required states:** normal, loading, empty, validation, error, success,
  permission, destructive, selected, and any domain-specific state.

## Constraints

- **Existing product and design-system conventions:**
- **Technical limits or component constraints:**
- **Accessibility, privacy, legal, or performance requirements:**
- **Available and licensed fonts, imagery, icons, and brand assets:**

## Direction And Scope

- **Visual direction:**
- **References and the specific lesson to take from each:**
- **Anti-references or patterns to avoid:**
- **In scope:**
- **Out of scope:**
- **Defaults and open questions:**

## Handoff

- **Decisions ready for implementation:**
- **Assumptions to validate:**
- **Prototype or usability question, if needed:**
- **Verification plan:**
