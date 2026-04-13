package cmd

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

// tableColumn defines a column header and a function to extract the cell value.
type tableColumn struct {
	Header string
	Get    func(v reflect.Value) string
}

// printTable renders list results as an aligned table. Non-slice types and
// unknown element types fall back to JSON output.
func printTable(w io.Writer, v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Slice {
		return printJSON(w, v)
	}

	if val.Len() == 0 {
		_, err := fmt.Fprintln(w, "(no results)")
		return err
	}

	// Get columns from the first element's concrete type.
	elem := val.Index(0)
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}
	cols := columnsForModel(elem.Type())
	if cols == nil {
		return printJSON(w, v)
	}

	// Build header and rows.
	headers := make([]string, len(cols))
	for i, c := range cols {
		headers[i] = c.Header
	}

	rows := make([][]string, val.Len())
	for i := 0; i < val.Len(); i++ {
		item := val.Index(i)
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}
		rows[i] = make([]string, len(cols))
		for j, c := range cols {
			rows[i][j] = c.Get(item)
		}
	}

	// Calculate column widths.
	widths := make([]int, len(cols))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Print header.
	var headerParts []string
	for i, h := range headers {
		headerParts = append(headerParts, padRight(h, widths[i]))
	}
	fmt.Fprintln(w, strings.Join(headerParts, "  "))

	// Print separator.
	var sepParts []string
	for _, w := range widths {
		sepParts = append(sepParts, strings.Repeat("-", w))
	}
	fmt.Fprintln(w, strings.Join(sepParts, "  "))

	// Print data rows.
	for _, row := range rows {
		var parts []string
		for i, cell := range row {
			parts = append(parts, padRight(cell, widths[i]))
		}
		fmt.Fprintln(w, strings.Join(parts, "  "))
	}

	return nil
}

// columnsForModel returns table column definitions for known model types.
func columnsForModel(t reflect.Type) []tableColumn {
	switch t.Name() {
	case "Course":
		return []tableColumn{
			{"COURSE ID", fieldGetter("CourseID")},
			{"ORG", fieldGetter("Org")},
			{"TITLE", fieldGetter("Title")},
			{"PACING", fieldGetter("Pacing")},
			{"START", fieldGetter("Start")},
		}
	case "User":
		return []tableColumn{
			{"USERNAME", fieldGetter("Username")},
			{"EMAIL", fieldGetter("Email")},
			{"NAME", fieldGetter("Name")},
			{"ACTIVE", boolGetter("IsActive")},
		}
	case "Enrollment":
		return []tableColumn{
			{"USERNAME", fieldGetter("Username")},
			{"COURSE ID", fieldGetter("CourseID")},
			{"MODE", fieldGetter("Mode")},
			{"ACTIVE", boolGetter("IsActive")},
		}
	case "Grade":
		return []tableColumn{
			{"USERNAME", fieldGetter("Username")},
			{"COURSE ID", fieldGetter("CourseID")},
			{"PERCENT", floatGetter("Percent")},
			{"GRADE", fieldGetter("LetterGrade")},
			{"PASSED", boolGetter("Passed")},
			{"SECTION", fieldGetter("Section")},
		}
	case "Certificate":
		return []tableColumn{
			{"USERNAME", fieldGetter("Username")},
			{"COURSE ID", fieldGetter("CourseID")},
			{"TYPE", fieldGetter("CertificateType")},
			{"STATUS", fieldGetter("Status")},
			{"GRADE", fieldGetter("Grade")},
		}
	default:
		return nil
	}
}

// fieldGetter returns a function that reads a string field by name.
func fieldGetter(name string) func(v reflect.Value) string {
	return func(v reflect.Value) string {
		f := v.FieldByName(name)
		if !f.IsValid() {
			return ""
		}
		return fmt.Sprintf("%v", f.Interface())
	}
}

// boolGetter returns a function that reads a bool field by name.
func boolGetter(name string) func(v reflect.Value) string {
	return func(v reflect.Value) string {
		f := v.FieldByName(name)
		if !f.IsValid() {
			return ""
		}
		if f.Bool() {
			return "true"
		}
		return "false"
	}
}

// floatGetter returns a function that reads a float64 field by name.
func floatGetter(name string) func(v reflect.Value) string {
	return func(v reflect.Value) string {
		f := v.FieldByName(name)
		if !f.IsValid() {
			return ""
		}
		return fmt.Sprintf("%.2f", f.Float())
	}
}

// padRight pads a string with spaces to the given width.
func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}
