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

## Prerequisites

MUST read `../openedx-shared/SKILL.md` first to understand configuration, authentication, and security rules.

## Core Concepts

- **Course** — a top-level learning unit identified by `course_id` (format: `course-v1:Org+Number+Run`). Contains org, number, run, title, pacing, start/end dates, and status.
- **CourseOutline** — hierarchical structure of a course's content: chapters → sequential → vertical → blocks. Each node has an id, title, type, and optional children.
- **Enrollment** — links a user to a course with a mode (audit, verified, professional, etc.). Tracks active/inactive state.
- **RoleAssignment** — grants a user a role (instructor, staff) within a course scope.

## Resource Relationships

```
Course (course_id)
├── outline → CourseOutline
│   └── chapters[] → SectionBlocks (recursive)
├── enrollments[] → Enrollment
│   └── (username, course_id, mode, is_active)
└── roles[] → RoleAssignment
    └── (username, course_id, role)
```

## Commands

### Course Management

| Command | Description | Reference |
|---------|-------------|-----------|
| `openedx course list` | List all courses | [course-list.md](references/course-list.md) |
| `openedx course get` | Get course details | inline below |
| `openedx course create` | Create a new course | [course-create.md](references/course-create.md) |
| `openedx course import` | Import course archive | [course-import.md](references/course-import.md) |
| `openedx course export` | Export course archive | [course-export.md](references/course-export.md) |
| `openedx course outline get` | Get course content outline | [course-outline.md](references/course-outline.md) |

### Enrollment Management

| Command | Description | Reference |
|---------|-------------|-----------|
| `openedx enrollment add` | Enroll a user in a course | [enrollment-add.md](references/enrollment-add.md) |
| `openedx enrollment list` | List enrollments | [enrollment-list.md](references/enrollment-list.md) |
| `openedx enrollment remove` | Remove a user's enrollment | [enrollment-remove.md](references/enrollment-remove.md) |

### Role Management

| Command | Description | Reference |
|---------|-------------|-----------|
| `openedx role assign` | Assign a role to a user | [role-assign.md](references/role-assign.md) |

`openedx course get` quick reference:

```bash
openedx course get --course-id <course_id>
```

Returns a single Course object.

## Common Workflows

### Create a course and import content

```bash
# 1. Create the course
openedx course create --org "MyOrg" --number "CS101" --run "2026_Spring" --title "Intro to CS"

# 2. Import course content from a tar.gz archive
openedx course import --course-id "course-v1:MyOrg+CS101+2026_Spring" --file ./content.tar.gz
```

### Batch enroll students

```bash
# Enroll each student
openedx enrollment add --course-id "course-v1:MyOrg+CS101+2026_Spring" --username alice --mode audit
openedx enrollment add --course-id "course-v1:MyOrg+CS101+2026_Spring" --username bob --mode verified
```

### Assign teaching staff

```bash
openedx role assign --course-id "course-v1:MyOrg+CS101+2026_Spring" --username prof_smith --role instructor
openedx role assign --course-id "course-v1:MyOrg+CS101+2026_Spring" --username ta_jones --role staff
```

## Error Handling

| Error | Cause | Resolution |
|-------|-------|------------|
| Course not found (404) | Invalid `course_id` | Verify course_id with `openedx course list` |
| Duplicate enrollment | User already enrolled | Check with `openedx enrollment list --course-id <id> --username <user>` |
| Permission denied (403) | Insufficient role | Ensure user has staff or instructor role on the course |
| Import job failed | Invalid archive format | Use Open edX compatible OLX tar.gz format |