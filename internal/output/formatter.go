package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

// Format represents output format types
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
	FormatYAML  Format = "yaml"
)

// Formatter handles different output formats
type Formatter struct {
	format Format
	writer io.Writer
}

// NewFormatter creates a new formatter
func NewFormatter(format Format) *Formatter {
	return &Formatter{
		format: format,
		writer: os.Stdout,
	}
}

// SetWriter sets the output writer
func (f *Formatter) SetWriter(w io.Writer) {
	f.writer = w
}

// OutputTable outputs data in table format
func (f *Formatter) OutputTable(headers []string, rows [][]string) error {
	if f.format == FormatJSON {
		return f.outputJSON(f.tableToMap(headers, rows))
	}

	w := tabwriter.NewWriter(f.writer, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Print headers
	fmt.Fprintln(w, strings.Join(headers, "\t"))

	// Print separator
	separators := make([]string, len(headers))
	for i := range separators {
		separators[i] = strings.Repeat("-", len(headers[i]))
	}
	fmt.Fprintln(w, strings.Join(separators, "\t"))

	// Print rows
	for _, row := range rows {
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}

	return nil
}

// OutputJSON outputs data as JSON
func (f *Formatter) OutputJSON(data interface{}) error {
	return f.outputJSON(data)
}

// OutputSingle outputs a single item
func (f *Formatter) OutputSingle(item interface{}) error {
	switch f.format {
	case FormatJSON:
		return f.outputJSON(item)
	default:
		return f.outputJSON(item) // Default to JSON for single items
	}
}

func (f *Formatter) outputJSON(data interface{}) error {
	encoder := json.NewEncoder(f.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func (f *Formatter) tableToMap(headers []string, rows [][]string) []map[string]string {
	result := make([]map[string]string, len(rows))

	for i, row := range rows {
		item := make(map[string]string)
		for j, header := range headers {
			if j < len(row) {
				item[header] = row[j]
			}
		}
		result[i] = item
	}

	return result
}

// ParseFormat parses format string to Format type
func ParseFormat(s string) Format {
	switch strings.ToLower(s) {
	case "json":
		return FormatJSON
	case "yaml":
		return FormatYAML
	default:
		return FormatTable
	}
}
