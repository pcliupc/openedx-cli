---
name: openedx-skill-maker
version: 1.0.0
description: "创建 openedx CLI 的自定义 Skill"
metadata:
  requires:
    skills: ["openedx-shared"]
    bins: ["openedx"]
---

## Overview

This skill teaches you how to create custom skills for the OpenEdX CLI. Custom skills allow you to document extension API endpoints or new command patterns for AI Agents to use.

## Skill Structure

Every skill lives in its own directory under `skills/`:

```
skills/
  openedx-<domain>/
    SKILL.md                    # Required: skill definition
    references/                 # Optional: detailed command docs
      <command>.md
```

## Creating a SKILL.md

### YAML Frontmatter

```yaml
---
name: openedx-<domain>
version: 1.0.0
description: "<domain description for skill discovery>"
metadata:
  requires:
    skills: ["openedx-shared"]
    bins: ["openedx"]
  cliHelp: "openedx <domain> --help"
---
```

Rules:
- `name` must be unique, prefixed with `openedx-`
- `description` is used for skill discovery — write it so an AI Agent can match user intent to this skill
- `metadata.requires.skills` must include `"openedx-shared"` for all business skills
- `metadata.requires.bins` is always `["openedx"]`

### Markdown Body

Use this template for business-domain skills:

```markdown
## Prerequisites
- MUST read `../openedx-shared/SKILL.md` first

## Core Concepts
<!-- Define key resources and data models -->

## Resource Relationships
<!-- Show how resources relate to each other -->

## Commands
<!-- Table of commands with links to reference docs -->

## Common Workflows
<!-- Step-by-step operation sequences -->

## Error Handling
<!-- Domain-specific errors and resolutions -->
```

## Creating Reference Documents

Each reference file documents a single command in detail:

```markdown
# openedx <domain> <verb>

<one-line description>

## Usage

openedx <domain> <verb> [flags]

## Parameters

| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| --flag | Yes/No | value | Description |

## Examples

### Basic usage
openedx <domain> <verb> --required-flag value

## Output

Description of output JSON shape with example.

## Notes

- Implementation-specific caveats
```

Five required sections: Usage, Parameters, Examples, Output, Notes.

## Naming Conventions

- Skill directory: `openedx-<domain>` (lowercase, hyphenated)
- Reference files: `<domain>-<verb>.md` (e.g., `course-list.md`, `enrollment-add.md`)
- Command keys in registry: `<domain>.<verb>` (e.g., `course.list`, `enrollment.add`)

## Dependency Declaration

All business skills must declare `openedx-shared` as a dependency:

```yaml
metadata:
  requires:
    skills: ["openedx-shared"]
```

The first line of the Markdown body must be:

```markdown
## Prerequisites
- MUST read `../openedx-shared/SKILL.md` first
```