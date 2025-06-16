// Package client provides HTTP client functionality for interacting with the Coolify API.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	coolify "github.com/hongkongkiwi/coolifyme/internal/api"
	"github.com/hongkongkiwi/coolifyme/internal/config"
	"github.com/hongkongkiwi/coolifyme/internal/logger"
)

// Client wraps the generated Coolify API client
type Client struct {
	API    *coolify.ClientWithResponses
	config *config.Config
}

// New creates a new Coolify client
func New(cfg *config.Config) (*Client, error) {
	if cfg.APIToken == "" {
		return nil, fmt.Errorf("API token is required")
	}

	// Create HTTP client with authentication and logging
	httpClient := &http.Client{
		Transport: &loggingTransport{
			token: cfg.APIToken,
			base:  http.DefaultTransport,
		},
	}

	// Create the API client
	apiClient, err := coolify.NewClientWithResponses(cfg.BaseURL, coolify.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	return &Client{
		API:    apiClient,
		config: cfg,
	}, nil
}

// loggingTransport implements HTTP transport with Bearer token authentication and request/response logging
type loggingTransport struct {
	token string
	base  http.RoundTripper
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	// Set authentication headers
	req.Header.Set("Authorization", "Bearer "+t.token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Log request details if debug logging is enabled
	logger.Debug("API Request",
		"method", req.Method,
		"url", req.URL.String(),
		"headers", formatHeaders(req.Header),
	)

	// Log request body if present
	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err == nil {
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			if len(bodyBytes) > 0 {
				logger.Debug("API Request Body", "body", string(bodyBytes))
			}
		}
	}

	// Make the request
	resp, err := t.base.RoundTrip(req)
	duration := time.Since(start)

	if err != nil {
		logger.Debug("API Request Failed",
			"method", req.Method,
			"url", req.URL.String(),
			"duration", duration.String(),
			"error", err.Error(),
		)
		return resp, err
	}

	// Log response details
	logger.Debug("API Response",
		"method", req.Method,
		"url", req.URL.String(),
		"status", resp.Status,
		"duration", duration.String(),
		"headers", formatHeaders(resp.Header),
	)

	// Log response body if debug logging and it's a small response
	if resp.Body != nil && resp.ContentLength < 10000 { // Only log small responses
		bodyBytes, err := io.ReadAll(resp.Body)
		if err == nil {
			resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			if len(bodyBytes) > 0 {
				logger.Debug("API Response Body", "body", string(bodyBytes))
			}
		}
	}

	return resp, nil
}

// formatHeaders formats HTTP headers for logging (excluding sensitive ones)
func formatHeaders(headers http.Header) string {
	var formatted []string
	for key, values := range headers {
		if strings.ToLower(key) == "authorization" {
			formatted = append(formatted, fmt.Sprintf("%s: [REDACTED]", key))
		} else {
			formatted = append(formatted, fmt.Sprintf("%s: %s", key, strings.Join(values, ", ")))
		}
	}
	return strings.Join(formatted, "; ")
}

// Applications returns an applications client
func (c *Client) Applications() *ApplicationsClient {
	return &ApplicationsClient{client: c}
}

// Projects returns a projects client
func (c *Client) Projects() *ProjectsClient {
	return &ProjectsClient{client: c}
}

// Servers returns a servers client
func (c *Client) Servers() *ServersClient {
	return &ServersClient{client: c}
}

// Services returns a services client
func (c *Client) Services() *ServicesClient {
	return &ServicesClient{client: c}
}

// Deployments returns a deployments client
func (c *Client) Deployments() *DeploymentsClient {
	return &DeploymentsClient{client: c}
}

// Databases returns a databases client
func (c *Client) Databases() *DatabasesClient {
	return &DatabasesClient{client: c}
}

// PrivateKeys returns a private keys client
func (c *Client) PrivateKeys() *PrivateKeysClient {
	return &PrivateKeysClient{client: c}
}

// Resources returns a resources client
func (c *Client) Resources() *ResourcesClient {
	return &ResourcesClient{client: c}
}

// Teams returns a teams client
func (c *Client) Teams() *TeamsClient {
	return &TeamsClient{client: c}
}

// System returns a system client
func (c *Client) System() *SystemClient {
	return &SystemClient{client: c}
}

// ApplicationsClient handles application-related operations
type ApplicationsClient struct {
	client *Client
}

// List returns all applications
func (ac *ApplicationsClient) List(ctx context.Context) ([]coolify.Application, error) {
	resp, err := ac.client.API.ListApplicationsWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list applications: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response body")
	}

	return *resp.JSON200, nil
}

// CreatePublic creates a new application from a public repository
func (ac *ApplicationsClient) CreatePublic(ctx context.Context, req coolify.CreatePublicApplicationJSONRequestBody) (*coolify.Application, error) {
	resp, err := ac.client.API.CreatePublicApplicationWithResponse(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON201 == nil {
		return nil, fmt.Errorf("empty response body")
	}

	// Note: The API returns just a UUID, we'd need to fetch the full application
	// This is a simplified implementation
	return nil, nil
}

// Get returns an application by UUID
func (ac *ApplicationsClient) Get(ctx context.Context, uuidStr string) (*coolify.Application, error) {
	appUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID: %w", err)
	}

	resp, err := ac.client.API.GetApplicationByUuidWithResponse(ctx, appUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get application: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response body")
	}

	return resp.JSON200, nil
}

// Delete deletes an application by UUID
func (ac *ApplicationsClient) Delete(ctx context.Context, uuidStr string, options *coolify.DeleteApplicationByUuidParams) error {
	appUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return fmt.Errorf("invalid UUID: %w", err)
	}

	resp, err := ac.client.API.DeleteApplicationByUuidWithResponse(ctx, appUUID, options)
	if err != nil {
		return fmt.Errorf("failed to delete application: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("API error: %s", resp.Status())
	}

	return nil
}

// Update updates an application by UUID
func (ac *ApplicationsClient) Update(ctx context.Context, uuidStr string, req coolify.UpdateApplicationByUuidJSONRequestBody) (string, error) {
	appUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return "", fmt.Errorf("invalid UUID: %w", err)
	}

	resp, err := ac.client.API.UpdateApplicationByUuidWithResponse(ctx, appUUID, req)
	if err != nil {
		return "", fmt.Errorf("failed to update application: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil || resp.JSON200.Uuid == nil {
		return "", fmt.Errorf("empty response body")
	}

	return *resp.JSON200.Uuid, nil
}

// CreatePrivateGithubApp creates a new application from a private GitHub app repository
func (ac *ApplicationsClient) CreatePrivateGithubApp(ctx context.Context, req coolify.CreatePrivateGithubAppApplicationJSONRequestBody) (*coolify.Application, error) {
	resp, err := ac.client.API.CreatePrivateGithubAppApplicationWithResponse(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	return nil, nil // API returns UUID, would need to fetch full application
}

// CreatePrivateDeployKey creates a new application from a private repository with deploy key
func (ac *ApplicationsClient) CreatePrivateDeployKey(ctx context.Context, req coolify.CreatePrivateDeployKeyApplicationJSONRequestBody) (*coolify.Application, error) {
	resp, err := ac.client.API.CreatePrivateDeployKeyApplicationWithResponse(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	return nil, nil // API returns UUID, would need to fetch full application
}

// CreateDockerfile creates a new application from a Dockerfile
func (ac *ApplicationsClient) CreateDockerfile(ctx context.Context, req coolify.CreateDockerfileApplicationJSONRequestBody) (*coolify.Application, error) {
	resp, err := ac.client.API.CreateDockerfileApplicationWithResponse(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	return nil, nil // API returns UUID, would need to fetch full application
}

// CreateDockerImage creates a new application from a Docker image
func (ac *ApplicationsClient) CreateDockerImage(ctx context.Context, req coolify.CreateDockerimageApplicationJSONRequestBody) (*coolify.Application, error) {
	resp, err := ac.client.API.CreateDockerimageApplicationWithResponse(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	return nil, nil // API returns UUID, would need to fetch full application
}

// CreateDockerCompose creates a new application from a Docker Compose file
func (ac *ApplicationsClient) CreateDockerCompose(ctx context.Context, req coolify.CreateDockercomposeApplicationJSONRequestBody) (*coolify.Application, error) {
	resp, err := ac.client.API.CreateDockercomposeApplicationWithResponse(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	return nil, nil // API returns UUID, would need to fetch full application
}

// Start starts an application
func (ac *ApplicationsClient) Start(ctx context.Context, uuidStr string, options *coolify.StartApplicationByUuidParams) (*StartResponse, error) {
	appUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID: %w", err)
	}

	resp, err := ac.client.API.StartApplicationByUuidWithResponse(ctx, appUUID, options)
	if err != nil {
		return nil, fmt.Errorf("failed to start application: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response body")
	}

	startResponse := &StartResponse{}
	if resp.JSON200.Message != nil {
		startResponse.Message = *resp.JSON200.Message
	}
	if resp.JSON200.DeploymentUuid != nil {
		startResponse.DeploymentUUID = *resp.JSON200.DeploymentUuid
	}

	return startResponse, nil
}

// Stop stops an application
func (ac *ApplicationsClient) Stop(ctx context.Context, uuidStr string) error {
	appUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return fmt.Errorf("invalid UUID: %w", err)
	}

	resp, err := ac.client.API.StopApplicationByUuidWithResponse(ctx, appUUID)
	if err != nil {
		return fmt.Errorf("failed to stop application: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("API error: %s", resp.Status())
	}

	return nil
}

// Restart restarts an application
func (ac *ApplicationsClient) Restart(ctx context.Context, uuidStr string) (*RestartResponse, error) {
	appUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID: %w", err)
	}

	resp, err := ac.client.API.RestartApplicationByUuidWithResponse(ctx, appUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to restart application: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response body")
	}

	restartResponse := &RestartResponse{}
	if resp.JSON200.Message != nil {
		restartResponse.Message = *resp.JSON200.Message
	}
	if resp.JSON200.DeploymentUuid != nil {
		restartResponse.DeploymentUUID = *resp.JSON200.DeploymentUuid
	}

	return restartResponse, nil
}

// GetLogs gets application logs
func (ac *ApplicationsClient) GetLogs(ctx context.Context, uuidStr string, params *coolify.GetApplicationLogsByUuidParams) (string, error) {
	appUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return "", fmt.Errorf("invalid UUID: %w", err)
	}

	resp, err := ac.client.API.GetApplicationLogsByUuidWithResponse(ctx, appUUID, params)
	if err != nil {
		return "", fmt.Errorf("failed to get application logs: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil || resp.JSON200.Logs == nil {
		return "", fmt.Errorf("empty response body")
	}

	return *resp.JSON200.Logs, nil
}

// ListEnvs lists environment variables for an application
func (ac *ApplicationsClient) ListEnvs(ctx context.Context, uuidStr string) ([]coolify.EnvironmentVariable, error) {
	appUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID: %w", err)
	}

	resp, err := ac.client.API.ListEnvsByApplicationUuidWithResponse(ctx, appUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to list environment variables: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response body")
	}

	return *resp.JSON200, nil
}

// CreateEnv creates an environment variable for an application
func (ac *ApplicationsClient) CreateEnv(ctx context.Context, uuidStr string, req coolify.CreateEnvByApplicationUuidJSONRequestBody) (string, error) {
	appUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return "", fmt.Errorf("invalid UUID: %w", err)
	}

	resp, err := ac.client.API.CreateEnvByApplicationUuidWithResponse(ctx, appUUID, req)
	if err != nil {
		return "", fmt.Errorf("failed to create environment variable: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return "", fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON201 == nil || resp.JSON201.Uuid == nil {
		return "", fmt.Errorf("empty response body")
	}

	return *resp.JSON201.Uuid, nil
}

// UpdateEnv updates an environment variable for an application
func (ac *ApplicationsClient) UpdateEnv(ctx context.Context, uuidStr string, req coolify.UpdateEnvByApplicationUuidJSONRequestBody) (string, error) {
	appUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return "", fmt.Errorf("invalid UUID: %w", err)
	}

	resp, err := ac.client.API.UpdateEnvByApplicationUuidWithResponse(ctx, appUUID, req)
	if err != nil {
		return "", fmt.Errorf("failed to update environment variable: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return "", fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON201 == nil || resp.JSON201.Message == nil {
		return "", fmt.Errorf("empty response body")
	}

	return *resp.JSON201.Message, nil
}

// UpdateEnvs updates multiple environment variables for an application
func (ac *ApplicationsClient) UpdateEnvs(ctx context.Context, uuidStr string, req coolify.UpdateEnvsByApplicationUuidJSONRequestBody) (string, error) {
	appUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return "", fmt.Errorf("invalid UUID: %w", err)
	}

	resp, err := ac.client.API.UpdateEnvsByApplicationUuidWithResponse(ctx, appUUID, req)
	if err != nil {
		return "", fmt.Errorf("failed to update environment variables: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return "", fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON201 == nil || resp.JSON201.Message == nil {
		return "", fmt.Errorf("empty response body")
	}

	return *resp.JSON201.Message, nil
}

// DeleteEnv deletes an environment variable for an application
func (ac *ApplicationsClient) DeleteEnv(ctx context.Context, uuidStr string, envUUIDStr string) (string, error) {
	appUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return "", fmt.Errorf("invalid UUID: %w", err)
	}

	envUUID, err := uuid.Parse(envUUIDStr)
	if err != nil {
		return "", fmt.Errorf("invalid env UUID: %w", err)
	}

	resp, err := ac.client.API.DeleteEnvByApplicationUuidWithResponse(ctx, appUUID, envUUID)
	if err != nil {
		return "", fmt.Errorf("failed to delete environment variable: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil || resp.JSON200.Message == nil {
		return "", fmt.Errorf("empty response body")
	}

	return *resp.JSON200.Message, nil
}

// ProjectsClient handles project-related operations
type ProjectsClient struct {
	client *Client
}

// List returns all projects
func (pc *ProjectsClient) List(ctx context.Context) ([]coolify.Project, error) {
	resp, err := pc.client.API.ListProjectsWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response body")
	}

	return *resp.JSON200, nil
}

// Create creates a new project
func (pc *ProjectsClient) Create(ctx context.Context, req coolify.CreateProjectJSONRequestBody) (string, error) {
	resp, err := pc.client.API.CreateProjectWithResponse(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create project: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return "", fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON201 == nil || resp.JSON201.Uuid == nil {
		return "", fmt.Errorf("empty response body")
	}

	return *resp.JSON201.Uuid, nil
}

// Get returns a project by UUID
func (pc *ProjectsClient) Get(ctx context.Context, uuidStr string) (*coolify.Project, error) {
	resp, err := pc.client.API.GetProjectByUuidWithResponse(ctx, uuidStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response body")
	}

	return resp.JSON200, nil
}

// Delete deletes a project by UUID
func (pc *ProjectsClient) Delete(ctx context.Context, uuidStr string) error {
	projectUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return fmt.Errorf("invalid UUID: %w", err)
	}

	resp, err := pc.client.API.DeleteProjectByUuidWithResponse(ctx, projectUUID)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("API error: %s", resp.Status())
	}

	return nil
}

// Update updates a project by UUID
func (pc *ProjectsClient) Update(ctx context.Context, uuidStr string, req coolify.UpdateProjectByUuidJSONRequestBody) (*coolify.Project, error) {
	projectUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID: %w", err)
	}

	resp, err := pc.client.API.UpdateProjectByUuidWithResponse(ctx, projectUUID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON201 == nil {
		return nil, fmt.Errorf("empty response body")
	}

	// Convert the response to a full Project object
	project := &coolify.Project{}
	if resp.JSON201.Uuid != nil {
		project.Uuid = resp.JSON201.Uuid
	}
	if resp.JSON201.Name != nil {
		project.Name = resp.JSON201.Name
	}
	if resp.JSON201.Description != nil {
		project.Description = resp.JSON201.Description
	}

	return project, nil
}

// GetEnvironment returns an environment by name or UUID within a project
func (pc *ProjectsClient) GetEnvironment(ctx context.Context, projectUUID, environmentNameOrUUID string) (*coolify.Environment, error) {
	resp, err := pc.client.API.GetEnvironmentByNameOrUuidWithResponse(ctx, projectUUID, environmentNameOrUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response body")
	}

	return resp.JSON200, nil
}

// ServersClient handles server-related operations
type ServersClient struct {
	client *Client
}

// List returns all servers
func (sc *ServersClient) List(ctx context.Context) ([]coolify.Server, error) {
	resp, err := sc.client.API.ListServersWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list servers: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response body")
	}

	return *resp.JSON200, nil
}

// Create creates a new server
func (sc *ServersClient) Create(ctx context.Context, req coolify.CreateServerJSONRequestBody) (string, error) {
	resp, err := sc.client.API.CreateServerWithResponse(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create server: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return "", fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON201 == nil || resp.JSON201.Uuid == nil {
		return "", fmt.Errorf("empty response body")
	}

	return *resp.JSON201.Uuid, nil
}

// Get returns a server by UUID
func (sc *ServersClient) Get(ctx context.Context, uuidStr string) (*coolify.Server, error) {
	resp, err := sc.client.API.GetServerByUuidWithResponse(ctx, uuidStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get server: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response body")
	}

	return resp.JSON200, nil
}

// Delete deletes a server by UUID
func (sc *ServersClient) Delete(ctx context.Context, uuidStr string) error {
	serverUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return fmt.Errorf("invalid UUID: %w", err)
	}

	resp, err := sc.client.API.DeleteServerByUuidWithResponse(ctx, serverUUID)
	if err != nil {
		return fmt.Errorf("failed to delete server: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("API error: %s", resp.Status())
	}

	return nil
}

// Update updates a server by UUID
func (sc *ServersClient) Update(ctx context.Context, uuidStr string, req coolify.UpdateServerByUuidJSONRequestBody) (*coolify.Server, error) {
	resp, err := sc.client.API.UpdateServerByUuidWithResponse(ctx, uuidStr, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update server: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON201 == nil {
		return nil, fmt.Errorf("empty response body")
	}

	return resp.JSON201, nil
}

// GetResources returns resources for a server by UUID (returns as JSON string per API spec)
func (sc *ServersClient) GetResources(ctx context.Context, uuidStr string) (string, error) {
	resp, err := sc.client.API.GetResourcesByServerUuidWithResponse(ctx, uuidStr)
	if err != nil {
		return "", fmt.Errorf("failed to get server resources: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil {
		return "", fmt.Errorf("empty response body")
	}

	// Convert to JSON string for consistent API interface
	jsonBytes, err := json.Marshal(*resp.JSON200)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}

	return string(jsonBytes), nil
}

// GetDomains returns domains for a server by UUID (returns as JSON string per API spec)
func (sc *ServersClient) GetDomains(ctx context.Context, uuidStr string) (string, error) {
	resp, err := sc.client.API.GetDomainsByServerUuidWithResponse(ctx, uuidStr)
	if err != nil {
		return "", fmt.Errorf("failed to get server domains: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil {
		return "", fmt.Errorf("empty response body")
	}

	// Convert to JSON string for consistent API interface
	jsonBytes, err := json.Marshal(*resp.JSON200)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}

	return string(jsonBytes), nil
}

// Validate validates a server by UUID
func (sc *ServersClient) Validate(ctx context.Context, uuidStr string) (string, error) {
	resp, err := sc.client.API.ValidateServerByUuidWithResponse(ctx, uuidStr)
	if err != nil {
		return "", fmt.Errorf("failed to validate server: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return "", fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON201 == nil || resp.JSON201.Message == nil {
		return "", fmt.Errorf("empty response body")
	}

	return *resp.JSON201.Message, nil
}

// ServicesClient handles service-related operations
type ServicesClient struct {
	client *Client
}

// List returns all services
func (sc *ServicesClient) List(ctx context.Context) ([]coolify.Service, error) {
	resp, err := sc.client.API.ListServicesWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response body")
	}

	return *resp.JSON200, nil
}

// Get returns a service by UUID
func (sc *ServicesClient) Get(ctx context.Context, uuidStr string) (*coolify.Service, error) {
	resp, err := sc.client.API.GetServiceByUuidWithResponse(ctx, uuidStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response body")
	}

	return resp.JSON200, nil
}

// Start starts a service
func (sc *ServicesClient) Start(ctx context.Context, uuidStr string) error {
	serviceUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return fmt.Errorf("invalid UUID: %w", err)
	}

	resp, err := sc.client.API.StartServiceByUuidWithResponse(ctx, serviceUUID)
	if err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("API error: %s", resp.Status())
	}

	return nil
}

// Stop stops a service
func (sc *ServicesClient) Stop(ctx context.Context, uuidStr string) error {
	serviceUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return fmt.Errorf("invalid UUID: %w", err)
	}

	resp, err := sc.client.API.StopServiceByUuidWithResponse(ctx, serviceUUID)
	if err != nil {
		return fmt.Errorf("failed to stop service: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("API error: %s", resp.Status())
	}

	return nil
}

// Restart restarts a service
func (sc *ServicesClient) Restart(ctx context.Context, uuidStr string) error {
	serviceUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return fmt.Errorf("invalid UUID: %w", err)
	}

	resp, err := sc.client.API.RestartServiceByUuidWithResponse(ctx, serviceUUID, nil)
	if err != nil {
		return fmt.Errorf("failed to restart service: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("API error: %s", resp.Status())
	}

	return nil
}

// Create creates a new service
func (sc *ServicesClient) Create(ctx context.Context, req coolify.CreateServiceJSONRequestBody) (string, error) {
	resp, err := sc.client.API.CreateServiceWithResponse(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create service: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return "", fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON201 == nil || resp.JSON201.Uuid == nil {
		return "", fmt.Errorf("empty response body")
	}

	return *resp.JSON201.Uuid, nil
}

// Delete deletes a service by UUID
func (sc *ServicesClient) Delete(ctx context.Context, uuidStr string, options *coolify.DeleteServiceByUuidParams) error {
	resp, err := sc.client.API.DeleteServiceByUuidWithResponse(ctx, uuidStr, options)
	if err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("API error: %s", resp.Status())
	}

	return nil
}

// Update updates a service by UUID
func (sc *ServicesClient) Update(ctx context.Context, uuidStr string, req coolify.UpdateServiceByUuidJSONRequestBody) (string, error) {
	serviceUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return "", fmt.Errorf("invalid UUID: %w", err)
	}

	resp, err := sc.client.API.UpdateServiceByUuidWithResponse(ctx, serviceUUID, req)
	if err != nil {
		return "", fmt.Errorf("failed to update service: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil || resp.JSON200.Uuid == nil {
		return "", fmt.Errorf("empty response body")
	}

	return *resp.JSON200.Uuid, nil
}

// ListEnvs lists environment variables for a service
func (sc *ServicesClient) ListEnvs(ctx context.Context, uuidStr string) ([]coolify.EnvironmentVariable, error) {
	serviceUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID: %w", err)
	}

	resp, err := sc.client.API.ListEnvsByServiceUuidWithResponse(ctx, serviceUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to list environment variables: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response body")
	}

	return *resp.JSON200, nil
}

// CreateEnv creates an environment variable for a service
func (sc *ServicesClient) CreateEnv(ctx context.Context, uuidStr string, req coolify.CreateEnvByServiceUuidJSONRequestBody) (string, error) {
	serviceUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return "", fmt.Errorf("invalid UUID: %w", err)
	}

	resp, err := sc.client.API.CreateEnvByServiceUuidWithResponse(ctx, serviceUUID, req)
	if err != nil {
		return "", fmt.Errorf("failed to create environment variable: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return "", fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON201 == nil || resp.JSON201.Uuid == nil {
		return "", fmt.Errorf("empty response body")
	}

	return *resp.JSON201.Uuid, nil
}

// UpdateEnv updates an environment variable for a service
func (sc *ServicesClient) UpdateEnv(ctx context.Context, uuidStr string, req coolify.UpdateEnvByServiceUuidJSONRequestBody) (string, error) {
	serviceUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return "", fmt.Errorf("invalid UUID: %w", err)
	}

	resp, err := sc.client.API.UpdateEnvByServiceUuidWithResponse(ctx, serviceUUID, req)
	if err != nil {
		return "", fmt.Errorf("failed to update environment variable: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return "", fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON201 == nil || resp.JSON201.Message == nil {
		return "", fmt.Errorf("empty response body")
	}

	return *resp.JSON201.Message, nil
}

// UpdateEnvs updates multiple environment variables for a service
func (sc *ServicesClient) UpdateEnvs(ctx context.Context, uuidStr string, req coolify.UpdateEnvsByServiceUuidJSONRequestBody) (string, error) {
	serviceUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return "", fmt.Errorf("invalid UUID: %w", err)
	}

	resp, err := sc.client.API.UpdateEnvsByServiceUuidWithResponse(ctx, serviceUUID, req)
	if err != nil {
		return "", fmt.Errorf("failed to update environment variables: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return "", fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON201 == nil || resp.JSON201.Message == nil {
		return "", fmt.Errorf("empty response body")
	}

	return *resp.JSON201.Message, nil
}

// DeleteEnv deletes an environment variable for a service
func (sc *ServicesClient) DeleteEnv(ctx context.Context, uuidStr string, envUUIDStr string) (string, error) {
	serviceUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return "", fmt.Errorf("invalid UUID: %w", err)
	}

	envUUID, err := uuid.Parse(envUUIDStr)
	if err != nil {
		return "", fmt.Errorf("invalid env UUID: %w", err)
	}

	resp, err := sc.client.API.DeleteEnvByServiceUuidWithResponse(ctx, serviceUUID, envUUID)
	if err != nil {
		return "", fmt.Errorf("failed to delete environment variable: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil || resp.JSON200.Message == nil {
		return "", fmt.Errorf("empty response body")
	}

	return *resp.JSON200.Message, nil
}

// DeploymentsClient handles deployment-related operations
type DeploymentsClient struct {
	client *Client
}

// DeployApplicationOptions contains options for deploying an application
type DeployApplicationOptions struct {
	Force  bool
	Branch string
	PR     *int
}

// DeploymentResult contains information about a triggered deployment
type DeploymentResult struct {
	Message        string `json:"message"`
	ResourceUUID   string `json:"resource_uuid"`
	DeploymentUUID string `json:"deployment_uuid"`
}

// DeployResponse contains the response from a deployment request
type DeployResponse struct {
	Deployments []DeploymentResult `json:"deployments"`
}

// StartResponse represents the response from starting an application
type StartResponse struct {
	Message        string `json:"message"`
	DeploymentUUID string `json:"deployment_uuid"`
}

// RestartResponse represents the response from restarting an application
type RestartResponse struct {
	Message        string `json:"message"`
	DeploymentUUID string `json:"deployment_uuid"`
}

// StopResponse represents the response from stopping an application
type StopResponse struct {
	Message string `json:"message"`
}

// DeployApplication deploys an application by UUID
func (dc *DeploymentsClient) DeployApplication(ctx context.Context, uuidStr string, force bool, branch string) (*DeployResponse, error) {
	return dc.DeployApplicationWithOptions(ctx, uuidStr, &DeployApplicationOptions{
		Force:  force,
		Branch: branch,
	})
}

// DeployApplicationWithOptions deploys an application with advanced options
func (dc *DeploymentsClient) DeployApplicationWithOptions(ctx context.Context, uuidStr string, options *DeployApplicationOptions) (*DeployResponse, error) {
	params := &coolify.DeployByTagOrUuidParams{
		Uuid:  &uuidStr,
		Force: &options.Force,
	}

	// Branch and PR are mutually exclusive
	if options.Branch != "" && options.PR != nil {
		return nil, fmt.Errorf("cannot specify both branch and PR - they are mutually exclusive")
	}

	// If branch is specified, we need to deploy from a specific tag/branch
	if options.Branch != "" {
		params.Tag = &options.Branch
	}

	// If PR is specified, deploy from a specific pull request
	if options.PR != nil {
		params.Pr = options.PR
	}

	resp, err := dc.client.API.DeployByTagOrUuidWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy application: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil || resp.JSON200.Deployments == nil {
		return nil, fmt.Errorf("empty response body")
	}

	// Convert the response to our struct
	result := &DeployResponse{
		Deployments: make([]DeploymentResult, 0, len(*resp.JSON200.Deployments)),
	}

	for _, deployment := range *resp.JSON200.Deployments {
		deploymentResult := DeploymentResult{}
		if deployment.Message != nil {
			deploymentResult.Message = *deployment.Message
		}
		if deployment.ResourceUuid != nil {
			deploymentResult.ResourceUUID = *deployment.ResourceUuid
		}
		if deployment.DeploymentUuid != nil {
			deploymentResult.DeploymentUUID = *deployment.DeploymentUuid
		}
		result.Deployments = append(result.Deployments, deploymentResult)
	}

	return result, nil
}

// DeployService deploys a service by starting it (services use start/restart for deployment)
func (dc *DeploymentsClient) DeployService(ctx context.Context, uuidStr string) error {
	return dc.client.Services().Start(ctx, uuidStr)
}

// List returns deployment history for an application
func (dc *DeploymentsClient) List(ctx context.Context, appUUIDStr string) ([]coolify.Application, error) {
	appUUID, err := uuid.Parse(appUUIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID: %w", err)
	}

	resp, err := dc.client.API.ListDeploymentsByAppUuidWithResponse(ctx, appUUID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response body")
	}

	return *resp.JSON200, nil
}

// ListAll returns all deployments
func (dc *DeploymentsClient) ListAll(ctx context.Context) ([]coolify.ApplicationDeploymentQueue, error) {
	resp, err := dc.client.API.ListDeploymentsWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response body")
	}

	return *resp.JSON200, nil
}

// GetByUUID returns a deployment by UUID
func (dc *DeploymentsClient) GetByUUID(ctx context.Context, uuidStr string) (*coolify.ApplicationDeploymentQueue, error) {
	resp, err := dc.client.API.GetDeploymentByUuidWithResponse(ctx, uuidStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response body")
	}

	return resp.JSON200, nil
}

// Watch monitors a deployment until it completes or fails
func (dc *DeploymentsClient) Watch(ctx context.Context, uuidStr string) error {
	fmt.Printf("🔄 Monitoring deployment %s...\n", uuidStr)

	for {
		deployment, err := dc.GetByUUID(ctx, uuidStr)
		if err != nil {
			return fmt.Errorf("failed to get deployment status: %w", err)
		}

		if deployment.Status == nil {
			return fmt.Errorf("deployment status is unknown")
		}

		status := *deployment.Status
		fmt.Printf("📊 Status: %s\n", status)

		// Check if deployment is finished (success or failure)
		switch status {
		case "finished", "success", "completed":
			fmt.Printf("✅ Deployment completed successfully!\n")
			return nil
		case "failed", "error", "cancelled":
			fmt.Printf("❌ Deployment failed with status: %s\n", status)
			if deployment.Logs != nil && *deployment.Logs != "" {
				fmt.Printf("📝 Recent logs:\n%s\n", *deployment.Logs)
			}
			return fmt.Errorf("deployment failed")
		case "running", "in_progress", "building", "deploying":
			// Continue monitoring
			fmt.Printf("⏳ Deployment in progress...\n")
		default:
			fmt.Printf("ℹ️  Unknown status: %s\n", status)
		}

		// Wait before next check
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(5 * time.Second):
			// Continue loop
		}
	}
}

// DeployMultiple deploys multiple applications by their UUIDs
func (dc *DeploymentsClient) DeployMultiple(ctx context.Context, uuids []string, options *DeployApplicationOptions) (*DeployResponse, error) {
	if len(uuids) == 0 {
		return nil, fmt.Errorf("no UUIDs provided")
	}

	// Join UUIDs with commas as the API supports comma-separated lists
	uuidList := strings.Join(uuids, ",")

	params := &coolify.DeployByTagOrUuidParams{
		Uuid:  &uuidList,
		Force: &options.Force,
	}

	// If branch is specified, we need to deploy from a specific tag/branch
	if options.Branch != "" {
		params.Tag = &options.Branch
	}

	// If PR is specified, deploy from a specific pull request
	if options.PR != nil {
		params.Pr = options.PR
	}

	resp, err := dc.client.API.DeployByTagOrUuidWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy applications: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil || resp.JSON200.Deployments == nil {
		return nil, fmt.Errorf("empty response body")
	}

	// Convert the response to our struct
	result := &DeployResponse{
		Deployments: make([]DeploymentResult, 0, len(*resp.JSON200.Deployments)),
	}

	for _, deployment := range *resp.JSON200.Deployments {
		deploymentResult := DeploymentResult{}
		if deployment.Message != nil {
			deploymentResult.Message = *deployment.Message
		}
		if deployment.ResourceUuid != nil {
			deploymentResult.ResourceUUID = *deployment.ResourceUuid
		}
		if deployment.DeploymentUuid != nil {
			deploymentResult.DeploymentUUID = *deployment.DeploymentUuid
		}
		result.Deployments = append(result.Deployments, deploymentResult)
	}

	return result, nil
}

// ListWithPagination returns deployment history for an application with pagination support
func (dc *DeploymentsClient) ListWithPagination(ctx context.Context, appUUIDStr string, skip, take int) ([]coolify.Application, error) {
	appUUID, err := uuid.Parse(appUUIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID: %w", err)
	}

	params := &coolify.ListDeploymentsByAppUuidParams{}
	if skip > 0 {
		params.Skip = &skip
	}
	if take > 0 {
		params.Take = &take
	}

	resp, err := dc.client.API.ListDeploymentsByAppUuidWithResponse(ctx, appUUID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response body")
	}

	return *resp.JSON200, nil
}

// DatabasesClient handles database-related operations
type DatabasesClient struct {
	client *Client
}

// List returns all databases (currently returns raw string as API is not fully implemented)
func (dc *DatabasesClient) List(ctx context.Context) (string, error) {
	resp, err := dc.client.API.ListDatabasesWithResponse(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list databases: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil {
		return "", fmt.Errorf("empty response body")
	}

	return *resp.JSON200, nil
}

// Get returns a database by UUID (currently returns raw string as API is not fully implemented)
func (dc *DatabasesClient) Get(ctx context.Context, uuidStr string) (string, error) {
	dbUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return "", fmt.Errorf("invalid UUID: %w", err)
	}

	resp, err := dc.client.API.GetDatabaseByUuidWithResponse(ctx, dbUUID)
	if err != nil {
		return "", fmt.Errorf("failed to get database: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil {
		return "", fmt.Errorf("empty response body")
	}

	return *resp.JSON200, nil
}

// Start starts a database
func (dc *DatabasesClient) Start(ctx context.Context, uuidStr string) error {
	dbUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return fmt.Errorf("invalid UUID: %w", err)
	}

	resp, err := dc.client.API.StartDatabaseByUuidWithResponse(ctx, dbUUID)
	if err != nil {
		return fmt.Errorf("failed to start database: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("API error: %s", resp.Status())
	}

	return nil
}

// Stop stops a database
func (dc *DatabasesClient) Stop(ctx context.Context, uuidStr string) error {
	dbUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return fmt.Errorf("invalid UUID: %w", err)
	}

	resp, err := dc.client.API.StopDatabaseByUuidWithResponse(ctx, dbUUID)
	if err != nil {
		return fmt.Errorf("failed to stop database: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("API error: %s", resp.Status())
	}

	return nil
}

// Restart restarts a database
func (dc *DatabasesClient) Restart(ctx context.Context, uuidStr string) error {
	dbUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return fmt.Errorf("invalid UUID: %w", err)
	}

	resp, err := dc.client.API.RestartDatabaseByUuidWithResponse(ctx, dbUUID, nil)
	if err != nil {
		return fmt.Errorf("failed to restart database: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("API error: %s", resp.Status())
	}

	return nil
}

// Delete deletes a database by UUID
func (dc *DatabasesClient) Delete(ctx context.Context, uuidStr string, options *coolify.DeleteDatabaseByUuidParams) error {
	dbUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return fmt.Errorf("invalid UUID: %w", err)
	}

	resp, err := dc.client.API.DeleteDatabaseByUuidWithResponse(ctx, dbUUID, options)
	if err != nil {
		return fmt.Errorf("failed to delete database: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("API error: %s", resp.Status())
	}

	return nil
}

// Update updates a database by UUID
func (dc *DatabasesClient) Update(ctx context.Context, uuidStr string, req coolify.UpdateDatabaseByUuidJSONRequestBody) error {
	dbUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return fmt.Errorf("invalid UUID: %w", err)
	}

	resp, err := dc.client.API.UpdateDatabaseByUuidWithResponse(ctx, dbUUID, req)
	if err != nil {
		return fmt.Errorf("failed to update database: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("API error: %s", resp.Status())
	}

	return nil
}

// CreatePostgreSQL creates a new PostgreSQL database
func (dc *DatabasesClient) CreatePostgreSQL(ctx context.Context, req coolify.CreateDatabasePostgresqlJSONRequestBody) error {
	resp, err := dc.client.API.CreateDatabasePostgresqlWithResponse(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create PostgreSQL database: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("API error: %s", resp.Status())
	}

	return nil
}

// CreateMySQL creates a new MySQL database
func (dc *DatabasesClient) CreateMySQL(ctx context.Context, req coolify.CreateDatabaseMysqlJSONRequestBody) error {
	resp, err := dc.client.API.CreateDatabaseMysqlWithResponse(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create MySQL database: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("API error: %s", resp.Status())
	}

	return nil
}

// CreateRedis creates a new Redis database
func (dc *DatabasesClient) CreateRedis(ctx context.Context, req coolify.CreateDatabaseRedisJSONRequestBody) error {
	resp, err := dc.client.API.CreateDatabaseRedisWithResponse(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create Redis database: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("API error: %s", resp.Status())
	}

	return nil
}

// CreateMongoDB creates a new MongoDB database
func (dc *DatabasesClient) CreateMongoDB(ctx context.Context, req coolify.CreateDatabaseMongodbJSONRequestBody) error {
	resp, err := dc.client.API.CreateDatabaseMongodbWithResponse(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create MongoDB database: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("API error: %s", resp.Status())
	}

	return nil
}

// CreateClickHouse creates a new ClickHouse database
func (dc *DatabasesClient) CreateClickHouse(ctx context.Context, req coolify.CreateDatabaseClickhouseJSONRequestBody) error {
	resp, err := dc.client.API.CreateDatabaseClickhouseWithResponse(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create ClickHouse database: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("API error: %s", resp.Status())
	}

	return nil
}

// CreateDragonfly creates a new Dragonfly database
func (dc *DatabasesClient) CreateDragonfly(ctx context.Context, req coolify.CreateDatabaseDragonflyJSONRequestBody) error {
	resp, err := dc.client.API.CreateDatabaseDragonflyWithResponse(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create Dragonfly database: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("API error: %s", resp.Status())
	}

	return nil
}

// CreateKeyDB creates a new KeyDB database
func (dc *DatabasesClient) CreateKeyDB(ctx context.Context, req coolify.CreateDatabaseKeydbJSONRequestBody) error {
	resp, err := dc.client.API.CreateDatabaseKeydbWithResponse(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create KeyDB database: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("API error: %s", resp.Status())
	}

	return nil
}

// CreateMariaDB creates a new MariaDB database
func (dc *DatabasesClient) CreateMariaDB(ctx context.Context, req coolify.CreateDatabaseMariadbJSONRequestBody) error {
	resp, err := dc.client.API.CreateDatabaseMariadbWithResponse(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create MariaDB database: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("API error: %s", resp.Status())
	}

	return nil
}

// TeamsClient handles team-related operations
type TeamsClient struct {
	client *Client
}

// List returns all teams
func (tc *TeamsClient) List(ctx context.Context) ([]coolify.Team, error) {
	resp, err := tc.client.API.ListTeamsWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list teams: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response body")
	}

	return *resp.JSON200, nil
}

// Get returns a team by ID
func (tc *TeamsClient) Get(ctx context.Context, teamID int) (*coolify.Team, error) {
	resp, err := tc.client.API.GetTeamByIdWithResponse(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response body")
	}

	return resp.JSON200, nil
}

// GetMembers returns members of a team by team ID
func (tc *TeamsClient) GetMembers(ctx context.Context, teamID int) ([]coolify.User, error) {
	resp, err := tc.client.API.GetMembersByTeamIdWithResponse(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get team members: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response body")
	}

	return *resp.JSON200, nil
}

// GetCurrent returns the current team
func (tc *TeamsClient) GetCurrent(ctx context.Context) (*coolify.Team, error) {
	resp, err := tc.client.API.GetCurrentTeamWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current team: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response body")
	}

	return resp.JSON200, nil
}

// GetCurrentMembers returns members of the current team
func (tc *TeamsClient) GetCurrentMembers(ctx context.Context) ([]coolify.User, error) {
	resp, err := tc.client.API.GetCurrentTeamMembersWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current team members: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response body")
	}

	return *resp.JSON200, nil
}

// SystemClient handles system-related operations
type SystemClient struct {
	client *Client
}

// Version returns the system version
func (sc *SystemClient) Version(ctx context.Context) (string, error) {
	resp, err := sc.client.API.VersionWithResponse(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get version: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil {
		return "", fmt.Errorf("empty response body")
	}

	return *resp.JSON200, nil
}

// Healthcheck performs a health check
func (sc *SystemClient) Healthcheck(ctx context.Context) (string, error) {
	resp, err := sc.client.API.HealthcheckWithResponse(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to perform healthcheck: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil {
		return "", fmt.Errorf("empty response body")
	}

	return *resp.JSON200, nil
}

// EnableAPI enables the API
func (sc *SystemClient) EnableAPI(ctx context.Context) (string, error) {
	resp, err := sc.client.API.EnableApiWithResponse(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to enable API: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil || resp.JSON200.Message == nil {
		return "", fmt.Errorf("empty response body")
	}

	return *resp.JSON200.Message, nil
}

// DisableAPI disables the API
func (sc *SystemClient) DisableAPI(ctx context.Context) (string, error) {
	resp, err := sc.client.API.DisableApiWithResponse(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to disable API: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil || resp.JSON200.Message == nil {
		return "", fmt.Errorf("empty response body")
	}

	return *resp.JSON200.Message, nil
}

// PrivateKeysClient handles private key-related operations
type PrivateKeysClient struct {
	client *Client
}

// List returns all private keys
func (pkc *PrivateKeysClient) List(ctx context.Context) ([]coolify.PrivateKey, error) {
	resp, err := pkc.client.API.ListPrivateKeysWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list private keys: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response body")
	}

	return *resp.JSON200, nil
}

// Create creates a new private key
func (pkc *PrivateKeysClient) Create(ctx context.Context, req coolify.CreatePrivateKeyJSONRequestBody) (string, error) {
	resp, err := pkc.client.API.CreatePrivateKeyWithResponse(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create private key: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return "", fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON201 == nil || resp.JSON201.Uuid == nil {
		return "", fmt.Errorf("empty response body")
	}

	return *resp.JSON201.Uuid, nil
}

// Get returns a private key by UUID
func (pkc *PrivateKeysClient) Get(ctx context.Context, uuidStr string) (*coolify.PrivateKey, error) {
	resp, err := pkc.client.API.GetPrivateKeyByUuidWithResponse(ctx, uuidStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get private key: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response body")
	}

	return resp.JSON200, nil
}

// Update updates a private key
func (pkc *PrivateKeysClient) Update(ctx context.Context, req coolify.UpdatePrivateKeyJSONRequestBody) (string, error) {
	resp, err := pkc.client.API.UpdatePrivateKeyWithResponse(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to update private key: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return "", fmt.Errorf("API error: %s", resp.Status())
	}

	if resp.JSON201 == nil || resp.JSON201.Uuid == nil {
		return "", fmt.Errorf("empty response body")
	}

	return *resp.JSON201.Uuid, nil
}

// Delete deletes a private key by UUID
func (pkc *PrivateKeysClient) Delete(ctx context.Context, uuidStr string) error {
	resp, err := pkc.client.API.DeletePrivateKeyByUuidWithResponse(ctx, uuidStr)
	if err != nil {
		return fmt.Errorf("failed to delete private key: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("API error: %s", resp.Status())
	}

	return nil
}

// ResourcesClient handles resource-related operations
type ResourcesClient struct {
	client *Client
}

// List returns all resources
func (rc *ResourcesClient) List(ctx context.Context) (string, error) {
	resp, err := rc.client.API.ListResourcesWithResponse(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list resources: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("API error: %s", resp.Status())
	}

	// Note: API returns string according to OpenAPI spec
	if resp.JSON200 == nil {
		return "", fmt.Errorf("empty response body")
	}

	return *resp.JSON200, nil
}
