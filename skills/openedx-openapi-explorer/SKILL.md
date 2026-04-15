---
name: openedx-openapi-explorer
version: 1.0.0
description: "探索 Open edX 原始 API 文档，查找 CLI 未覆盖的端点"
metadata:
  requires:
    skills: ["openedx-shared"]
    bins: ["openedx"]
---

## Overview

This skill helps you find and use raw Open edX APIs that are not covered by the CLI's built-in commands. Use this when a user's request cannot be fulfilled by any existing CLI command.

## When to Use

Use this skill when:
- The user needs an operation not available through any `openedx` command
- You need to find the correct API endpoint for an extension mapping
- You want to understand what APIs exist for a specific Open edX feature

## Official API Documentation

Open edX API documentation is available at:
- **Open edX API docs (readthedocs):** https://docs.openedx.org/en/latest/developers/references/api.html
- **Course Structure API:** `/api/course_structure/v0/`
- **Grades API:** `/api/grades/v1/`
- **Enrollment API:** `/api/enrollment/v1/`
- **User API:** `/api/user/v1/`
- **Certificate API:** `/api/certificates/v0/`

## Finding Endpoints

### Step 1: Check what the CLI already covers

```bash
openedx schema all
```

This shows all built-in command-to-endpoint mappings. If the command exists, use the corresponding CLI command instead of raw API calls.

### Step 2: Search the API docs

Browse the readthedocs API reference for the relevant domain. Key sections:
- **Core APIs** — courses, enrollments, users, grades
- **Platform APIs** — analytics, discovery, ecommerce
- **LMS APIs** — certificates, schedules, teams

### Step 3: Configure as extension

Once you find the endpoint, add it to the config as an extension:

```yaml
extensions:
  <domain>.<verb>:
    method: GET
    url: https://openedx.example.com/api/<path>
```

### Step 4: Test with doctor

```bash
openedx doctor verify <domain>.<verb>
```

## Common API Patterns

Open edX APIs follow consistent patterns:

| Pattern | Example |
|---------|---------|
| List resources | `GET /api/<resource>/v1/` |
| Get a resource | `GET /api/<resource>/v1/{id}/` |
| Create a resource | `POST /api/<resource>/v1/` |
| Update a resource | `PATCH /api/<resource>/v1/{id}/` |
| Delete a resource | `DELETE /api/<resource>/v1/{id}/` |

Authentication: all APIs require `Authorization: Bearer <token>` header.

## Extension Configuration Tips

- Use the same command key format as built-in commands: `<domain>.<verb>`
- Extension URLs must be full URLs including the base path
- The CLI sends the same auth token to extension endpoints as to public APIs
- Extension responses should match the shape of the corresponding public API response for best normalizer compatibility
