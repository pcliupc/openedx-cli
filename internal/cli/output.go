package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/openedx/cli/internal/model"
)

// PrintJSON writes the given value as pretty-printed JSON to the writer.
func PrintJSON(w io.Writer, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("\n"))
	return err
}

// outputColumn defines a table column header and value extractor.
type outputColumn struct {
	Header string
	Get    func(v reflect.Value) string
}

// PrintTable writes the given value as an aligned table to the writer.
// It supports slice types (list results). Non-slice types fall back to JSON.
func PrintTable(w io.Writer, v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Slice {
		return PrintJSON(w, v)
	}

	if val.Len() == 0 {
		_, err := fmt.Fprintln(w, "(no results)")
		return err
	}

	elem := val.Index(0)
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}
	cols := outputColumnsForModel(elem.Type())
	if cols == nil {
		return PrintJSON(w, v)
	}

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

	var headerParts []string
	for i, h := range headers {
		headerParts = append(headerParts, outputPadRight(h, widths[i]))
	}
	fmt.Fprintln(w, strings.Join(headerParts, "  "))

	var sepParts []string
	for _, w := range widths {
		sepParts = append(sepParts, strings.Repeat("-", w))
	}
	fmt.Fprintln(w, strings.Join(sepParts, "  "))

	for _, row := range rows {
		var parts []string
		for i, cell := range row {
			parts = append(parts, outputPadRight(cell, widths[i]))
		}
		fmt.Fprintln(w, strings.Join(parts, "  "))
	}

	return nil
}

// outputColumnsForModel returns table column definitions for known model types.
func outputColumnsForModel(t reflect.Type) []outputColumn {
	switch t.Name() {
	case "Course":
		return []outputColumn{
			{"COURSE ID", outputFieldGetter("CourseID")},
			{"ORG", outputFieldGetter("Org")},
			{"TITLE", outputFieldGetter("Title")},
			{"PACING", outputFieldGetter("Pacing")},
			{"START", outputFieldGetter("Start")},
		}
	case "User":
		return []outputColumn{
			{"USERNAME", outputFieldGetter("Username")},
			{"EMAIL", outputFieldGetter("Email")},
			{"NAME", outputFieldGetter("Name")},
			{"ACTIVE", outputBoolGetter("IsActive")},
		}
	case "Enrollment":
		return []outputColumn{
			{"USERNAME", outputFieldGetter("Username")},
			{"COURSE ID", outputFieldGetter("CourseID")},
			{"MODE", outputFieldGetter("Mode")},
			{"ACTIVE", outputBoolGetter("IsActive")},
		}
	case "Grade":
		return []outputColumn{
			{"USERNAME", outputFieldGetter("Username")},
			{"COURSE ID", outputFieldGetter("CourseID")},
			{"PERCENT", outputFloatGetter("Percent")},
			{"GRADE", outputFieldGetter("LetterGrade")},
			{"PASSED", outputBoolGetter("Passed")},
			{"SECTION", outputFieldGetter("Section")},
		}
	case "Certificate":
		return []outputColumn{
			{"USERNAME", outputFieldGetter("Username")},
			{"COURSE ID", outputFieldGetter("CourseID")},
			{"TYPE", outputFieldGetter("CertificateType")},
			{"STATUS", outputFieldGetter("Status")},
			{"GRADE", outputFieldGetter("Grade")},
		}
	default:
		return nil
	}
}

func outputFieldGetter(name string) func(v reflect.Value) string {
	return func(v reflect.Value) string {
		f := v.FieldByName(name)
		if !f.IsValid() {
			return ""
		}
		return fmt.Sprintf("%v", f.Interface())
	}
}

func outputBoolGetter(name string) func(v reflect.Value) string {
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

func outputFloatGetter(name string) func(v reflect.Value) string {
	return func(v reflect.Value) string {
		f := v.FieldByName(name)
		if !f.IsValid() {
			return ""
		}
		return fmt.Sprintf("%.2f", f.Float())
	}
}

func outputPadRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

// PrintOutput dispatches to the appropriate output format handler.
func PrintOutput(w io.Writer, format string, v interface{}) error {
	switch format {
	case "json":
		return PrintJSON(w, v)
	case "table":
		return PrintTable(w, v)
	default:
		return PrintJSON(w, v)
	}
}

// ensure model types are referenced (avoids unused import).
var _ = model.Course{}
