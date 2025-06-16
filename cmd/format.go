package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// OutputFormat represents different output formats
type OutputFormat string

// Output format constants
const (
	// FormatJSON represents JSON output format
	FormatJSON OutputFormat = "json"
	// FormatYAML represents YAML output format
	FormatYAML OutputFormat = "yaml"
	// FormatTable represents tabular output format
	FormatTable OutputFormat = "table"
	// FormatCSV represents CSV output format
	FormatCSV OutputFormat = "csv"
	// FormatCustom represents custom template-based output format
	FormatCustom OutputFormat = "custom"
	// FormatWide represents wide table output format showing more columns
	FormatWide OutputFormat = "wide"
	// FormatName represents name-only output format
	FormatName OutputFormat = "name"
	// FormatTemplate represents Go template-based output format
	FormatTemplate OutputFormat = "template"
)

// FormatOptions holds formatting configuration
type FormatOptions struct {
	Format       OutputFormat
	Columns      []string
	NoHeaders    bool
	Wide         bool
	Template     string
	SortBy       string
	SortReverse  bool
	ShowKind     bool
	CustomFormat string
}

// ParseFormatOptions parses format options from command flags
func ParseFormatOptions(cmd *cobra.Command) *FormatOptions {
	options := &FormatOptions{}

	if format, err := cmd.Flags().GetString("output"); err == nil && format != "" {
		// Parse custom format strings like "table(name,status,url)"
		if strings.HasPrefix(format, "table(") && strings.HasSuffix(format, ")") {
			options.Format = FormatTable
			columnsStr := strings.TrimSuffix(strings.TrimPrefix(format, "table("), ")")
			if columnsStr != "" {
				options.Columns = strings.Split(columnsStr, ",")
				for i, col := range options.Columns {
					options.Columns[i] = strings.TrimSpace(col)
				}
			}
		} else if strings.HasPrefix(format, "custom(") && strings.HasSuffix(format, ")") {
			options.Format = FormatCustom
			options.CustomFormat = strings.TrimSuffix(strings.TrimPrefix(format, "custom("), ")")
		} else {
			options.Format = OutputFormat(format)
		}
	} else {
		options.Format = FormatTable
	}

	if columns, err := cmd.Flags().GetString("columns"); err == nil && columns != "" {
		options.Columns = strings.Split(columns, ",")
		for i, col := range options.Columns {
			options.Columns[i] = strings.TrimSpace(col)
		}
	}

	if noHeaders, err := cmd.Flags().GetBool("no-headers"); err == nil {
		options.NoHeaders = noHeaders
	}

	if sortBy, err := cmd.Flags().GetString("sort-by"); err == nil && sortBy != "" {
		options.SortBy = sortBy
	}

	if sortReverse, err := cmd.Flags().GetBool("sort-reverse"); err == nil {
		options.SortReverse = sortReverse
	}

	if showKind, err := cmd.Flags().GetBool("show-kind"); err == nil {
		options.ShowKind = showKind
	}

	return options
}

// FormatOutput formats and outputs data according to the specified options
func FormatOutput(data interface{}, options *FormatOptions) error {
	switch options.Format {
	case FormatJSON:
		return outputJSON(data)
	case FormatYAML:
		return outputYAML(data)
	case FormatTable:
		return outputTable(data, options)
	case FormatCSV:
		return outputCSV(data, options)
	case FormatWide:
		return outputWide(data, options)
	case FormatName:
		return outputNameOnly(data)
	case FormatCustom:
		return outputCustom(data, options)
	default:
		return outputTable(data, options)
	}
}

func outputJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func outputYAML(data interface{}) error {
	encoder := yaml.NewEncoder(os.Stdout)
	defer func() {
		if err := encoder.Close(); err != nil {
			// Log error but don't fail the operation since output may already be displayed
			fmt.Fprintf(os.Stderr, "Warning: failed to close YAML encoder: %v\n", err)
		}
	}()
	return encoder.Encode(data)
}

func outputTable(data interface{}, options *FormatOptions) error {
	rows, headers := extractTableData(data, options)

	if len(rows) == 0 {
		fmt.Println("No data to display")
		return nil
	}

	// Sort data if requested
	if options.SortBy != "" {
		sortTableData(rows, headers, options.SortBy, options.SortReverse)
	}

	// Create table writer
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer func() {
		if err := w.Flush(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to flush table writer: %v\n", err)
		}
	}()

	// Print headers if not disabled
	if !options.NoHeaders {
		if _, err := fmt.Fprintln(w, strings.Join(headers, "\t")); err != nil {
			return fmt.Errorf("failed to write table headers: %w", err)
		}

		// Print separator line
		separators := make([]string, len(headers))
		for i, header := range headers {
			separators[i] = strings.Repeat("-", len(header))
		}
		if _, err := fmt.Fprintln(w, strings.Join(separators, "\t")); err != nil {
			return fmt.Errorf("failed to write table separators: %w", err)
		}
	}

	// Print data rows
	for _, row := range rows {
		if _, err := fmt.Fprintln(w, strings.Join(row, "\t")); err != nil {
			return fmt.Errorf("failed to write table row: %w", err)
		}
	}

	return nil
}

func outputCSV(data interface{}, options *FormatOptions) error {
	rows, headers := extractTableData(data, options)

	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	// Write headers if not disabled
	if !options.NoHeaders {
		if err := writer.Write(headers); err != nil {
			return fmt.Errorf("failed to write CSV headers: %w", err)
		}
	}

	// Write data rows
	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return nil
}

func outputWide(data interface{}, options *FormatOptions) error {
	// Wide format shows more columns and details
	wideOptions := *options
	wideOptions.Columns = nil // Show all available columns
	return outputTable(data, &wideOptions)
}

func outputNameOnly(data interface{}) error {
	// Extract just the name field from each item
	items := reflectToSlice(data)
	for _, item := range items {
		if name := extractField(item, "name"); name != "" {
			fmt.Println(name)
		} else if name := extractField(item, "Name"); name != "" {
			fmt.Println(name)
		}
	}
	return nil
}

func outputCustom(data interface{}, options *FormatOptions) error {
	// Custom format using Go template-like syntax
	format := options.CustomFormat
	if format == "" {
		return fmt.Errorf("custom format string is required")
	}

	items := reflectToSlice(data)
	for _, item := range items {
		output := format

		// Replace placeholders with actual values
		v := reflect.ValueOf(item)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		if v.Kind() == reflect.Struct {
			t := v.Type()
			for i := 0; i < v.NumField(); i++ {
				field := t.Field(i)
				value := v.Field(i)

				// Convert value to string
				var strValue string
				if value.Kind() == reflect.Ptr {
					if value.IsNil() {
						strValue = ""
					} else {
						strValue = fmt.Sprintf("%v", value.Elem().Interface())
					}
				} else {
					strValue = fmt.Sprintf("%v", value.Interface())
				}

				// Replace placeholders (case-insensitive)
				placeholder := "{" + strings.ToLower(field.Name) + "}"
				output = strings.ReplaceAll(strings.ToLower(output), placeholder, strValue)

				// Also try with original case
				placeholder = "{" + field.Name + "}"
				output = strings.ReplaceAll(output, placeholder, strValue)
			}
		}

		fmt.Println(output)
	}

	return nil
}

func extractTableData(data interface{}, options *FormatOptions) ([][]string, []string) {
	items := reflectToSlice(data)
	if len(items) == 0 {
		return nil, nil
	}

	// Get headers from the first item
	headers := getHeaders(items[0], options)

	// Extract data rows
	var rows [][]string
	for _, item := range items {
		row := extractRow(item, headers, options)
		rows = append(rows, row)
	}

	return rows, headers
}

func reflectToSlice(data interface{}) []interface{} {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	var items []interface{}

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			items = append(items, v.Index(i).Interface())
		}
	default:
		// Single item
		items = append(items, data)
	}

	return items
}

func getHeaders(item interface{}, options *FormatOptions) []string {
	if len(options.Columns) > 0 {
		// Use specified columns
		headers := make([]string, len(options.Columns))
		copy(headers, options.Columns)
		return headers
	}

	// Auto-detect headers from struct fields
	var headers []string
	v := reflect.ValueOf(item)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() == reflect.Struct {
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			// Skip unexported fields
			if field.PkgPath != "" {
				continue
			}

			// Use json tag name if available, otherwise field name
			tag := field.Tag.Get("json")
			if tag != "" && tag != "-" {
				// Remove omitempty and other options
				if idx := strings.Index(tag, ","); idx != -1 {
					tag = tag[:idx]
				}
				headers = append(headers, strings.ToUpper(tag))
			} else {
				headers = append(headers, strings.ToUpper(field.Name))
			}
		}
	}

	return headers
}

func extractRow(item interface{}, headers []string, _ *FormatOptions) []string {
	row := make([]string, len(headers))

	for i, header := range headers {
		row[i] = extractField(item, header)
	}

	return row
}

func extractField(item interface{}, fieldName string) string {
	v := reflect.ValueOf(item)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Sprintf("%v", item)
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)

		// Check field name (case-insensitive)
		if strings.EqualFold(field.Name, fieldName) {
			return formatFieldValue(v.Field(i))
		}

		// Check json tag
		tag := field.Tag.Get("json")
		if tag != "" && tag != "-" {
			if idx := strings.Index(tag, ","); idx != -1 {
				tag = tag[:idx]
			}
			if strings.EqualFold(tag, fieldName) {
				return formatFieldValue(v.Field(i))
			}
		}
	}

	return ""
}

func formatFieldValue(value reflect.Value) string {
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return ""
		}
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.String:
		return value.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'f', -1, 64)
	case reflect.Bool:
		return strconv.FormatBool(value.Bool())
	default:
		// Try to format as time
		if value.Type() == reflect.TypeOf(time.Time{}) {
			t := value.Interface().(time.Time)
			return t.Format("2006-01-02 15:04:05")
		}
		return fmt.Sprintf("%v", value.Interface())
	}
}

func sortTableData(rows [][]string, headers []string, sortBy string, reverse bool) {
	// Find the column index to sort by
	sortIndex := -1
	for i, header := range headers {
		if strings.EqualFold(header, sortBy) {
			sortIndex = i
			break
		}
	}

	if sortIndex == -1 {
		return // Column not found
	}

	sort.Slice(rows, func(i, j int) bool {
		if sortIndex >= len(rows[i]) || sortIndex >= len(rows[j]) {
			return false
		}

		result := strings.Compare(rows[i][sortIndex], rows[j][sortIndex])
		if reverse {
			return result > 0
		}
		return result < 0
	})
}

// AddFormatFlags adds formatting flags to a command
func AddFormatFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("output", "o", "", "Output format (json, yaml, table, csv, wide, name, custom(format))")
	cmd.Flags().String("columns", "", "Comma-separated list of columns to display")
	cmd.Flags().Bool("no-headers", false, "Don't print headers")
	cmd.Flags().String("sort-by", "", "Sort by column name")
	cmd.Flags().Bool("sort-reverse", false, "Reverse sort order")
	cmd.Flags().Bool("show-kind", false, "Show resource kind/type")
}

// formatCmd demonstrates format options
var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "Format command examples and testing",
	Long:  "Examples and testing for various output formats",
}

var formatExamplesCmd = &cobra.Command{
	Use:   "examples",
	Short: "Show format examples",
	Long:  "Display examples of different output formats",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🎨 Output Format Examples")
		fmt.Println("========================")
		fmt.Println()

		fmt.Println("Basic formats:")
		fmt.Println("  --output json                 # JSON format")
		fmt.Println("  --output yaml                 # YAML format")
		fmt.Println("  --output table                # Table format (default)")
		fmt.Println("  --output csv                  # CSV format")
		fmt.Println("  --output wide                 # Wide table with all columns")
		fmt.Println("  --output name                 # Names only")
		fmt.Println()

		fmt.Println("Table with specific columns:")
		fmt.Println("  --output \"table(name,status,url)\"")
		fmt.Println("  --columns name,status,url")
		fmt.Println()

		fmt.Println("Custom format:")
		fmt.Println("  --output \"custom({name} is {status})\"")
		fmt.Println()

		fmt.Println("Sorting:")
		fmt.Println("  --sort-by name                # Sort by name column")
		fmt.Println("  --sort-reverse                # Reverse sort order")
		fmt.Println()

		fmt.Println("Other options:")
		fmt.Println("  --no-headers                  # Don't show column headers")
		fmt.Println("  --show-kind                   # Include resource type")

		return nil
	},
}

func init() {
	formatCmd.AddCommand(formatExamplesCmd)
}
