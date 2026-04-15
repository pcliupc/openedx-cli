# openedx course create

Create a new course on the platform.

## Usage

```bash
openedx course create --org <org> --number <number> --run <run> --title <title>
```

## Parameters

| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| --org | Yes | — | Organization identifier (e.g. "MyOrg") |
| --number | Yes | — | Course number (e.g. "CS101") |
| --run | Yes | — | Course run identifier (e.g. "2026_Spring") |
| --title | Yes | — | Human-readable course title |

## Examples

### Create a new course
```bash
openedx course create --org "MyOrg" --number "CS101" --run "2026_Spring" --title "Introduction to Computer Science"
```

### Create with a specific profile
```bash
openedx course create --org "MyOrg" --number "CS101" --run "2026_Spring" --title "Intro to CS" --profile staging
```

## Output

Single Course object:

```json
{
  "course_id": "course-v1:MyOrg+CS101+2026_Spring",
  "org": "MyOrg",
  "number": "CS101",
  "run": "2026_Spring",
  "title": "Introduction to Computer Science",
  "pacing": "",
  "start": "",
  "end": "",
  "status": ""
}
```

## Notes

- The returned `course_id` is auto-generated from org+number+run in the format `course-v1:{Org}+{Number}+{Run}`
- Save the returned `course_id` for subsequent operations (import, enrollment, etc.)
- Newly created courses have empty start/end dates and pacing — configure these in the Open edX Studio