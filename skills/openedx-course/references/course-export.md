# openedx course export

Export a course as an OLX archive.

## Usage

```bash
openedx course export --course-id <course_id> [--output <directory>]
```

## Parameters

| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| --course-id | Yes | — | Course identifier to export |
| --output | No | current directory | Output directory for the exported archive |

## Examples

### Export a course
```bash
openedx course export --course-id "course-v1:MyOrg+CS101+2026_Spring"
```

### Export to a specific directory
```bash
openedx course export --course-id "course-v1:MyOrg+CS101+2026_Spring" --output ./exports
```

## Output

Job object (asynchronous operation):

```json
{
  "job_id": "xyz-456-ghi",
  "operation": "export",
  "status": "submitted",
  "submitted_at": "2026-04-15T10:00:00Z",
  "finished_at": "",
  "result": "",
  "artifacts": []
}
```

## Notes

- Export is asynchronous — the returned Job object tracks progress
- The exported file is an OLX tar.gz archive
- Use `--output` to control where the archive is saved