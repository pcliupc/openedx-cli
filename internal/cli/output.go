package cli

import (
	"encoding/json"
	"errors"
	"io"
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

// PrintTable writes the given value as a table to the writer.
// Table format is not yet supported.
func PrintTable(w io.Writer, v interface{}) error {
	return errors.New("table format not yet supported")
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
