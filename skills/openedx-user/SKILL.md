---
name: openedx-user
version: 1.0.0
description: "用户管理：创建、查询、列表"
metadata:
  requires:
    skills: ["openedx-shared"]
    bins: ["openedx"]
  cliHelp: "openedx user --help"
---

## Prerequisites

MUST read `../openedx-shared/SKILL.md` first to understand configuration, authentication, and security rules.

## Core Concepts

- **User** — a platform user identified by `username`. Has email, full name, active status, and creation timestamp. Users can be enrolled in courses and assigned roles.

## Commands

| Command | Description | Reference |
|---------|-------------|-----------|
| `openedx user create` | Create a new user | [user-create.md](references/user-create.md) |
| `openedx user list` | List users | [user-list.md](references/user-list.md) |
| `openedx user get` | Get user details | [user-get.md](references/user-get.md) |

## Common Workflows

### Create a user and immediately enroll

```bash
# 1. Create the user
openedx user create --username alice --email alice@example.com --name "Alice Smith"

# 2. Enroll in a course
openedx enrollment add --course-id "course-v1:MyOrg+CS101+2026_Spring" --username alice
```

### Query a user and assign a role

```bash
# 1. Verify the user exists
openedx user get --username prof_smith

# 2. Assign a course role
openedx role assign --course-id "course-v1:MyOrg+CS101+2026_Spring" --username prof_smith --role instructor
```

## Error Handling

| Error | Cause | Resolution |
|-------|-------|------------|
| Username conflict (400/409) | Username already exists | Choose a different username or use `user get` to check |
| User not found (404) | Invalid username | Verify username with `openedx user list` |
| Permission denied (403) | Insufficient role | Ensure your profile has user management permissions |
