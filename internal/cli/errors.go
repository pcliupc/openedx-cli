package cli

import (
	"encoding/json"
	"fmt"
	"io"
)

// CLIError represents a structured error that can be printed as JSON.
type CLIError struct {
	Code        string `json:"error"`
	Message     string `json:"message"`
	Resource    string `json:"resource,omitempty"`
	Suggestion  string `json:"suggestion,omitempty"`
}

// NewCLIError creates a new CLIError with the given error code and message.
func NewCLIError(errorCode, message string) *CLIError {
	return &CLIError{
		Code:    errorCode,
		Message: message,
	}
}

// Print writes the error as pretty-printed JSON to the writer.
func (e *CLIError) Print(w io.Writer) error {
	data, err := json.MarshalIndent(e, "", "  ")
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

// Error satisfies the error interface.
func (e *CLIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
