# OpenEdX CLI Skills System Design

## Background

Lark CLI provides 19 skills that teach AI Agents how to use its CLI commands. Each skill is a structured Markdown document that an Agent reads before executing commands. The OpenEdX CLI should adopt the same pattern to make it usable by AI Agents (Claude, Copilot, etc.) in CI pipelines and development workflows.

## Goal

Create a skills system for the OpenEdX CLI that enables AI Agents to autonomously use the CLI for Open edX management operations. Skills are instruction documents, not executable code — they teach Agents when and how to invoke CLI commands.

## Design

### Skill Inventory (4+2 Structure)

Six skills total: 4 business-domain skills grouped by function, 2 auxiliary skills.

| Skill | Coverage | Description |
|---|---|---|
| `openedx-shared` | config, auth, diagnostics | Configuration, OAuth authentication, doctor checks, security rules |
| `openedx-course` | course, enrollment, role | Course lifecycle: create, import, export, outline, enrollments, role assignments |
| `openedx-user` | user | User management: create, list, get |
| `openedx-grade` | grade, certificate | Grades and certificates: grade listing, gradebook, certificate queries |
| `openedx-skill-maker` | — | Guide for creating custom skills |
| `openedx-openapi-explorer` | — | Guide for exploring raw Open edX API documentation |

### Directory Structure

```
skills/
  openedx-shared/
    SKILL.md
  openedx-course/
    SKILL.md
    references/
      course-list.md
      course-create.md
      course-import.md
      course-export.md
      course-outline.md
      enrollment-add.md
      enrollment-list.md
      enrollment-remove.md
      role-assign.md
  openedx-user/
    SKILL.md
    references/
      user-create.md
      user-list.md
      user-get.md
  openedx-grade/
    SKILL.md
    references/
      grade-list.md
      grade-gradebook.md
      certificate-list.md
  openedx-skill-maker/
    SKILL.md
  openedx-openapi-explorer/
    SKILL.md
```

Rules:
- Every skill must have a `SKILL.md` file (YAML frontmatter + Markdown body)
- Complex or high-frequency commands get dedicated files under `references/`
- Simple commands can be documented inline in SKILL.md

### SKILL.md Format

#### YAML Frontmatter

```yaml
---
name: openedx-course
version: 1.0.0
description: "课程管理：创建、导入、导出课程，管理注册和角色分配"
metadata:
  requires:
    skills: ["openedx-shared"]
    bins: ["openedx"]
  cliHelp: "openedx course --help"
---
```

Fields:
- `name` — unique skill identifier, used for dependency references
- `version` — semantic version, aligned with CLI release
- `description` — short description for skill discovery and trigger matching
- `metadata.requires.skills` — prerequisite skill dependencies (empty for `openedx-shared` since it is the base)
- `metadata.requires.bins` — required CLI binaries
- `metadata.cliHelp` — domain help command

#### Markdown Body (Business Skills)

```markdown
## Prerequisites
- MUST read `../openedx-shared/SKILL.md` first

## Core Concepts
<!-- Key concepts and data models for this domain -->

## Resource Relationships
<!-- Resource relationship tree (text) -->

## Commands
<!-- Command reference table with links to references/ -->

## Common Workflows
<!-- Common operation sequences within this domain -->

## Error Handling
<!-- Domain-specific error scenarios and handling advice -->
```

`openedx-shared` omits "Prerequisites" and "Resource Relationships", focusing on configuration, authentication, and diagnostics.

`openedx-skill-maker` and `openedx-openapi-explorer` use custom body structures tailored to their purpose.

### Skill Content Design

#### openedx-shared

- **Configuration** — YAML config file structure, profile concept, environment variable references for credentials
- **Authentication** — OAuth client credentials flow, in-memory token caching, credentials referenced by env var names (never stored as plaintext)
- **CLI Installation** — npm install and build-from-source methods
- **Diagnostics** — `openedx doctor` health check flow (base URL reachability → token acquisition → API availability), `openedx schema` for command-to-endpoint mapping inspection
- **Security Rules** — credentials never written to disk, non-interactive by default, stderr for logs, stdout for JSON data
- **Extension API** — configuring custom API endpoints, fallback trigger conditions (only 404/405/501)

#### openedx-course

- **Core Concepts** — Course, CourseOutline, Enrollment, RoleAssignment resource models
- **Resource Relationships**:

```
Course
├── outline (CourseOutline)
├── enrollments[] (Enrollment)
└── roles[] (RoleAssignment)
```

- **Commands** — 10 commands grouped by scenario (course management / enrollment management / role management), each linked to reference doc
- **Common Workflows** — create course and import content, batch enroll students, assign TA roles
- **Error Handling** — course not found (404), duplicate enrollment, insufficient permissions

#### openedx-user

- **Core Concepts** — User resource model
- **Commands** — 3 commands (create, list, get)
- **Common Workflows** — create user and immediately enroll, query by username
- **Error Handling** — username conflict, user not found

#### openedx-grade

- **Core Concepts** — Grade, Gradebook, Certificate resource models
- **Resource Relationships**:

```
Course
├── grades[] (Grade)
├── gradebook (Gradebook)
└── certificates[] (Certificate)
```

- **Commands** — 3 commands (grade list, gradebook get, certificate list)
- **Common Workflows** — query course grades, export gradebook, check certificate status
- **Error Handling** — no grade data, certificate not generated

#### openedx-skill-maker

- Teaches AI Agents how to create new custom skills for the OpenEdX CLI
- Includes SKILL.md template, reference document template, naming conventions, dependency declaration
- Follows the same pattern as lark-skill-maker

#### openedx-openapi-explorer

- Teaches AI Agents how to explore official Open edX API documentation
- Covers official API docs entry points (readthedocs), common endpoint index, how to map raw APIs to CLI extension configuration
- Helps Agents find and use raw APIs when built-in CLI commands are insufficient

### Reference Document Format

Each file under `references/` follows this template:

```markdown
# openedx course list

<one-line description>

## Usage

openedx course list [--profile <name>] [--limit <n>] [--output json|table]

## Parameters

| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| --profile | No | default | Config profile name |
| --limit | No | 50 | Max number of results |
| --output | No | json | Output format: json or table |

## Examples

### List all courses (JSON)
openedx course list

### View in table format
openedx course list --output table

### Use a specific profile
openedx course list --profile staging

## Output

JSON array with course_id, name, org, number, start, end fields.

## Notes

- Default output is JSON for pipeline processing
- Pagination handled automatically by CLI
- Use --limit for large course sets
```

Five fixed sections: Usage, Parameters, Examples, Output, Notes. Examples ordered from simple to complex.

### Installation and Distribution

**Install command:**
```bash
npx skills add openedx/cli -y -g
```

**Distribution:**
- Skills stored in `skills/` directory of the openedx-cli repository
- `npx skills add` pulls from GitHub and installs to the Agent's skill directory
- `-g` for global install, `-y` for auto-confirm

**Version alignment:**
- Skill versions track CLI versions
- New CLI commands trigger updates to corresponding skill docs
- Removed/renamed commands trigger skill updates with version bumps

**Discovery and triggering:**
- `description` field used for skill discovery
- Agent matches user intent to skill descriptions
- Unified `openedx-` prefix avoids conflicts with other tools
