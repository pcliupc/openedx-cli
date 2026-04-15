# openedx course outline get

Get the hierarchical content outline of a course.

## Usage

```bash
openedx course outline get --course-id <course_id> [--username <user>] [--depth <n>] [--block-types <types>]
```

## Parameters

| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| --course-id | Yes | — | Course identifier |
| --username | No | — | Filter blocks visible to a specific user |
| --depth | No | 0 (all) | Depth of block tree to return (0 = unlimited) |
| --block-types | No | all | Comma-separated block types to include (e.g. "chapter,sequential,vertical") |

## Examples

### Get full course outline
```bash
openedx course outline get --course-id "course-v1:MyOrg+CS101+2026_Spring"
```

### Get outline visible to a specific student
```bash
openedx course outline get --course-id "course-v1:MyOrg+CS101+2026_Spring" --username alice
```

### Limit depth to chapters only
```bash
openedx course outline get --course-id "course-v1:MyOrg+CS101+2026_Spring" --depth 1
```

### Filter by block type
```bash
openedx course outline get --course-id "course-v1:MyOrg+CS101+2026_Spring" --block-types "chapter,sequential"
```

## Output

CourseOutline object with nested block structure:

```json
{
  "course_id": "course-v1:MyOrg+CS101+2026_Spring",
  "chapters": [
    {
      "id": "block-v1:MyOrg+CS101+2026_Spring+type@chapter+block@ch1",
      "title": "Week 1: Introduction",
      "type": "chapter",
      "children": [
        {
          "id": "block-v1:MyOrg+CS101+2026_Spring+type@sequential+block@seq1",
          "title": "Lesson 1",
          "type": "sequential",
          "children": []
        }
      ]
    }
  ]
}
```

## Notes

- The outline is a recursive tree: chapters → sequentials → verticals → components
- Use `--depth 1` for a quick overview of chapter titles only
- Use `--username` to see what a specific learner can access (respects content gating)