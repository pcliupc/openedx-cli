# openedx course import

Import course content from an OLX archive into an existing course.

## Usage

```bash
openedx course import --course-id <course_id> --file <path>
```

## Parameters

| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| --course-id | Yes | — | Target course identifier |
| --file | Yes | — | Path to the course archive file (tar.gz) |

## Examples

### Import course content
```bash
openedx course import --course-id "course-v1:MyOrg+CS101+2026_Spring" --file ./my-course.tar.gz
```

### Import on a staging instance
```bash
openedx course import --course-id "course-v1:MyOrg+CS101+2026_Spring" --file ./my-course.tar.gz --profile staging
```

## Output

Job object (asynchronous operation):

```json
{
  "job_id": "abc-123-def",
  "operation": "import",
  "status": "submitted",
  "submitted_at": "2026-04-15T10:00:00Z",
  "finished_at": "",
  "result": "",
  "artifacts": []
}
```

## Notes

- Import is asynchronous — the returned Job object tracks progress
- The archive must be in Open edX OLX format (tar.gz)
- The target course must already exist (use `openedx course create` first)
- Check job status via the platform's task API if needed