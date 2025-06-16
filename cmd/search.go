package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"

	coolify "github.com/hongkongkiwi/coolifyme/internal/api"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search across all resources",
	Long:  "Search applications, services, servers, and databases with powerful filtering",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		query := args[0]
		resourceType, _ := cmd.Flags().GetString("type")
		status, _ := cmd.Flags().GetString("status")
		tag, _ := cmd.Flags().GetString("tag")
		caseSensitive, _ := cmd.Flags().GetBool("case-sensitive")
		limit, _ := cmd.Flags().GetInt("limit")
		jsonOutput, _ := cmd.Flags().GetBool("json")

		ctx := context.Background()
		results := &SearchResults{}

		// Search based on resource type filter
		if resourceType == "" || resourceType == "applications" || resourceType == "apps" {
			if err := searchApplications(ctx, client, query, status, tag, caseSensitive, results); err != nil {
				fmt.Printf("⚠️  Failed to search applications: %v\n", err)
			}
		}

		if resourceType == "" || resourceType == "services" || resourceType == "svc" {
			if err := searchServices(ctx, client, query, status, tag, caseSensitive, results); err != nil {
				fmt.Printf("⚠️  Failed to search services: %v\n", err)
			}
		}

		if resourceType == "" || resourceType == "servers" || resourceType == "srv" {
			if err := searchServers(ctx, client, query, status, tag, caseSensitive, results); err != nil {
				fmt.Printf("⚠️  Failed to search servers: %v\n", err)
			}
		}

		if resourceType == "" || resourceType == "databases" || resourceType == "db" {
			if err := searchDatabases(ctx, client, query, status, tag, caseSensitive, results); err != nil {
				fmt.Printf("⚠️  Failed to search databases: %v\n", err)
			}
		}

		// Apply limit
		if limit > 0 {
			results.ApplyLimit(limit)
		}

		// Output results
		if jsonOutput {
			output, err := json.MarshalIndent(results, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(output))
			return nil
		}

		displaySearchResults(results, query)
		return nil
	},
}

// findCmd represents the find command for more specific searches
var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Find resources with advanced filters",
	Long:  "Find resources using advanced filtering options",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		name, _ := cmd.Flags().GetString("name")
		status, _ := cmd.Flags().GetString("status")
		tag, _ := cmd.Flags().GetString("tag")
		resourceType, _ := cmd.Flags().GetString("type")
		jsonOutput, _ := cmd.Flags().GetBool("json")

		if name == "" && status == "" && tag == "" {
			return fmt.Errorf("at least one filter must be specified (--name, --status, or --tag)")
		}

		ctx := context.Background()
		results := &SearchResults{}

		// Search based on resource type filter
		if resourceType == "" || resourceType == "applications" || resourceType == "apps" {
			if err := findApplications(ctx, client, name, status, tag, results); err != nil {
				fmt.Printf("⚠️  Failed to find applications: %v\n", err)
			}
		}

		if resourceType == "" || resourceType == "services" || resourceType == "svc" {
			if err := findServices(ctx, client, name, status, tag, results); err != nil {
				fmt.Printf("⚠️  Failed to find services: %v\n", err)
			}
		}

		if resourceType == "" || resourceType == "servers" || resourceType == "srv" {
			if err := findServers(ctx, client, name, status, tag, results); err != nil {
				fmt.Printf("⚠️  Failed to find servers: %v\n", err)
			}
		}

		// Output results
		if jsonOutput {
			output, err := json.MarshalIndent(results, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(output))
			return nil
		}

		displaySearchResults(results, fmt.Sprintf("filters: name=%s status=%s tag=%s", name, status, tag))
		return nil
	},
}

// SearchResults holds the results from a search operation across different resource types
type SearchResults struct {
	Applications []SearchResultApp    `json:"applications"`
	Services     []SearchResultSvc    `json:"services"`
	Servers      []SearchResultServer `json:"servers"`
	Databases    []SearchResultDB     `json:"databases"`
	TotalCount   int                  `json:"total_count"`
}

// SearchResultApp represents an application in search results
type SearchResultApp struct {
	UUID   string `json:"uuid"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Type   string `json:"type"`
	URL    string `json:"url,omitempty"`
}

// SearchResultSvc represents a service in search results
type SearchResultSvc struct {
	UUID   string `json:"uuid"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Type   string `json:"type"`
}

// SearchResultServer represents a server in search results
type SearchResultServer struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	IP          string `json:"ip"`
	Status      string `json:"status"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}

// SearchResultDB represents a database in search results
type SearchResultDB struct {
	UUID   string `json:"uuid"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Type   string `json:"type"`
}

func searchApplications(ctx context.Context, client interface{}, query, status, tag string, caseSensitive bool, results *SearchResults) error {
	// Type assertion to get the actual client
	c, ok := client.(interface {
		Applications() interface {
			List(context.Context) ([]coolify.Application, error)
		}
	})
	if !ok {
		return fmt.Errorf("invalid client type")
	}

	apps, err := c.Applications().List(ctx)
	if err != nil {
		return err
	}

	for _, app := range apps {
		if matchesSearch(app, query, status, tag, caseSensitive) {
			result := SearchResultApp{
				Type: "application",
			}
			if app.Uuid != nil {
				result.UUID = *app.Uuid
			}
			if app.Name != nil {
				result.Name = *app.Name
			}
			if app.Status != nil {
				result.Status = *app.Status
			}
			if app.Fqdn != nil {
				result.URL = *app.Fqdn
			}
			results.Applications = append(results.Applications, result)
		}
	}
	return nil
}

func searchServices(ctx context.Context, client interface{}, query, status, tag string, caseSensitive bool, results *SearchResults) error {
	c, ok := client.(interface {
		Services() interface {
			List(context.Context) ([]coolify.Service, error)
		}
	})
	if !ok {
		return fmt.Errorf("invalid client type")
	}

	services, err := c.Services().List(ctx)
	if err != nil {
		return err
	}

	for _, svc := range services {
		if matchesSearchService(svc, query, status, tag, caseSensitive) {
			result := SearchResultSvc{
				Type: "service",
			}
			if svc.Uuid != nil {
				result.UUID = *svc.Uuid
			}
			if svc.Name != nil {
				result.Name = *svc.Name
			}
			// Services don't have a status field in the API model
			result.Status = StatusUnknown
			results.Services = append(results.Services, result)
		}
	}
	return nil
}

func searchServers(ctx context.Context, client interface{}, query, status, tag string, caseSensitive bool, results *SearchResults) error {
	c, ok := client.(interface {
		Servers() interface {
			List(context.Context) ([]coolify.Server, error)
		}
	})
	if !ok {
		return fmt.Errorf("invalid client type")
	}

	servers, err := c.Servers().List(ctx)
	if err != nil {
		return err
	}

	for _, srv := range servers {
		if matchesSearchServer(srv, query, status, tag, caseSensitive) {
			result := SearchResultServer{
				Type: "server",
			}
			if srv.Uuid != nil {
				result.UUID = *srv.Uuid
			}
			if srv.Name != nil {
				result.Name = *srv.Name
			}
			if srv.Ip != nil {
				result.IP = *srv.Ip
			}
			if srv.Description != nil {
				result.Description = *srv.Description
			}
			// Determine status from validation
			if srv.ValidationLogs != nil {
				result.Status = StatusValidated
			} else {
				result.Status = StatusUnknown
			}
			results.Servers = append(results.Servers, result)
		}
	}
	return nil
}

func searchDatabases(ctx context.Context, client interface{}, query, status, tag string, caseSensitive bool, results *SearchResults) error {
	return fmt.Errorf("database search not yet implemented")
}

func findApplications(ctx context.Context, client interface{}, name, status, tag string, results *SearchResults) error {
	return searchApplications(ctx, client, name, status, tag, false, results)
}

func findServices(ctx context.Context, client interface{}, name, status, tag string, results *SearchResults) error {
	return searchServices(ctx, client, name, status, tag, false, results)
}

func findServers(ctx context.Context, client interface{}, name, status, tag string, results *SearchResults) error {
	return searchServers(ctx, client, name, status, tag, false, results)
}

func matchesSearch(app coolify.Application, query, status, tag string, caseSensitive bool) bool {
	// Search in name, description, and other fields
	searchFields := []string{}

	if app.Name != nil {
		searchFields = append(searchFields, *app.Name)
	}
	if app.Description != nil {
		searchFields = append(searchFields, *app.Description)
	}
	if app.Fqdn != nil {
		searchFields = append(searchFields, *app.Fqdn)
	}
	if app.GitRepository != nil {
		searchFields = append(searchFields, *app.GitRepository)
	}

	// Check query match
	queryMatches := query == "" || containsText(strings.Join(searchFields, " "), query, caseSensitive)

	// Check status filter
	statusMatches := status == "" || (app.Status != nil && *app.Status == status)

	// Note: Tag filtering would require additional API support
	tagMatches := tag == ""

	return queryMatches && statusMatches && tagMatches
}

func matchesSearchService(svc coolify.Service, query, status, tag string, caseSensitive bool) bool {
	searchFields := []string{}

	if svc.Name != nil {
		searchFields = append(searchFields, *svc.Name)
	}
	if svc.Description != nil {
		searchFields = append(searchFields, *svc.Description)
	}

	queryMatches := query == "" || containsText(strings.Join(searchFields, " "), query, caseSensitive)
	// Services don't have a status field, so status filtering is not supported
	statusMatches := status == ""
	tagMatches := tag == ""

	return queryMatches && statusMatches && tagMatches
}

func matchesSearchServer(srv coolify.Server, query, status, tag string, caseSensitive bool) bool {
	searchFields := []string{}

	if srv.Name != nil {
		searchFields = append(searchFields, *srv.Name)
	}
	if srv.Description != nil {
		searchFields = append(searchFields, *srv.Description)
	}
	if srv.Ip != nil {
		searchFields = append(searchFields, *srv.Ip)
	}

	queryMatches := query == "" || containsText(strings.Join(searchFields, " "), query, caseSensitive)

	// For servers, we check validation status
	serverStatus := StatusUnknown
	if srv.ValidationLogs != nil {
		serverStatus = StatusValidated
	}
	statusMatches := status == "" || serverStatus == status
	tagMatches := tag == ""

	return queryMatches && statusMatches && tagMatches
}

func containsText(text, query string, caseSensitive bool) bool {
	if !caseSensitive {
		text = strings.ToLower(text)
		query = strings.ToLower(query)
	}

	// Support wildcard patterns
	if strings.Contains(query, "*") {
		pattern := strings.ReplaceAll(regexp.QuoteMeta(query), `\*`, `.*`)
		if !caseSensitive {
			pattern = "(?i)" + pattern
		}
		matched, _ := regexp.MatchString(pattern, text)
		return matched
	}

	return strings.Contains(text, query)
}

// ApplyLimit applies a limit to the search results across all resource types
func (sr *SearchResults) ApplyLimit(limit int) {
	count := 0

	// Limit applications
	if count < limit && len(sr.Applications) > 0 {
		remaining := limit - count
		if len(sr.Applications) > remaining {
			sr.Applications = sr.Applications[:remaining]
		}
		count += len(sr.Applications)
	}

	// Limit services
	if count < limit && len(sr.Services) > 0 {
		remaining := limit - count
		if len(sr.Services) > remaining {
			sr.Services = sr.Services[:remaining]
		}
		count += len(sr.Services)
	}

	// Limit servers
	if count < limit && len(sr.Servers) > 0 {
		remaining := limit - count
		if len(sr.Servers) > remaining {
			sr.Servers = sr.Servers[:remaining]
		}
		count += len(sr.Servers)
	}

	// Limit databases
	if count < limit && len(sr.Databases) > 0 {
		remaining := limit - count
		if len(sr.Databases) > remaining {
			sr.Databases = sr.Databases[:remaining]
		}
	}
}

func displaySearchResults(results *SearchResults, query string) {
	totalResults := len(results.Applications) + len(results.Services) + len(results.Servers) + len(results.Databases)

	fmt.Printf("🔍 Search Results for: %s\n", query)
	fmt.Printf("====================================\n\n")

	if totalResults == 0 {
		fmt.Println("📭 No results found")
		return
	}

	// Display Applications
	if len(results.Applications) > 0 {
		fmt.Printf("📱 Applications (%d)\n", len(results.Applications))
		fmt.Println("-------------------")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		if _, err := fmt.Fprintln(w, "UUID\tNAME\tSTATUS\tURL"); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to write application headers: %v\n", err)
		}
		if _, err := fmt.Fprintln(w, "----\t----\t------\t---"); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to write application separators: %v\n", err)
		}

		for _, app := range results.Applications {
			if _, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", app.UUID, app.Name, app.Status, app.URL); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to write application row: %v\n", err)
			}
		}
		if err := w.Flush(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to flush application table: %v\n", err)
		}
		fmt.Println()
	}

	// Display Services
	if len(results.Services) > 0 {
		fmt.Printf("🔧 Services (%d)\n", len(results.Services))
		fmt.Println("---------------")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		if _, err := fmt.Fprintln(w, "UUID\tNAME\tSTATUS"); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to write service headers: %v\n", err)
		}
		if _, err := fmt.Fprintln(w, "----\t----\t------"); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to write service separators: %v\n", err)
		}

		for _, svc := range results.Services {
			if _, err := fmt.Fprintf(w, "%s\t%s\t%s\n", svc.UUID, svc.Name, svc.Status); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to write service row: %v\n", err)
			}
		}
		if err := w.Flush(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to flush service table: %v\n", err)
		}
		fmt.Println()
	}

	// Display Servers
	if len(results.Servers) > 0 {
		fmt.Printf("🖥️  Servers (%d)\n", len(results.Servers))
		fmt.Println("-------------")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		if _, err := fmt.Fprintln(w, "UUID\tNAME\tIP\tSTATUS\tDESCRIPTION"); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to write server headers: %v\n", err)
		}
		if _, err := fmt.Fprintln(w, "----\t----\t--\t------\t-----------"); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to write server separators: %v\n", err)
		}

		for _, srv := range results.Servers {
			if _, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", srv.UUID, srv.Name, srv.IP, srv.Status, srv.Description); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to write server row: %v\n", err)
			}
		}
		if err := w.Flush(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to flush server table: %v\n", err)
		}
		fmt.Println()
	}

	fmt.Printf("📊 Total: %d results\n", totalResults)
}

func init() {
	// Search command flags
	searchCmd.Flags().StringP("type", "T", "", "Resource type to search (applications, services, servers, databases)")
	searchCmd.Flags().String("status", "", "Filter by status")
	searchCmd.Flags().String("tag", "", "Filter by tag")
	searchCmd.Flags().BoolP("case-sensitive", "c", false, "Case sensitive search")
	searchCmd.Flags().IntP("limit", "L", 0, "Limit number of results (0 = no limit)")
	searchCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// Find command flags
	findCmd.Flags().StringP("name", "n", "", "Filter by name pattern (supports wildcards)")
	findCmd.Flags().String("status", "", "Filter by status")
	findCmd.Flags().String("tag", "", "Filter by tag")
	findCmd.Flags().StringP("type", "T", "", "Resource type to search (applications, services, servers)")
	findCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
}
