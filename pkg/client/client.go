package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	coolify "github.com/hongkongkiwi/coolifyme/internal/api"
	"github.com/hongkongkiwi/coolifyme/internal/config"
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

	// Create HTTP client with authentication
	httpClient := &http.Client{
		Transport: &authTransport{
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

// authTransport implements HTTP transport with Bearer token authentication
type authTransport struct {
	token string
	base  http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	return t.base.RoundTrip(req)
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
func (ac *ApplicationsClient) Start(ctx context.Context, uuidStr string, options *coolify.StartApplicationByUuidParams) error {
	appUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return fmt.Errorf("invalid UUID: %w", err)
	}

	resp, err := ac.client.API.StartApplicationByUuidWithResponse(ctx, appUUID, options)
	if err != nil {
		return fmt.Errorf("failed to start application: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("API error: %s", resp.Status())
	}

	return nil
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
func (ac *ApplicationsClient) Restart(ctx context.Context, uuidStr string) error {
	appUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return fmt.Errorf("invalid UUID: %w", err)
	}

	resp, err := ac.client.API.RestartApplicationByUuidWithResponse(ctx, appUUID)
	if err != nil {
		return fmt.Errorf("failed to restart application: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("API error: %s", resp.Status())
	}

	return nil
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
func (ac *ApplicationsClient) DeleteEnv(ctx context.Context, uuidStr string, envUuidStr string) (string, error) {
	appUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return "", fmt.Errorf("invalid UUID: %w", err)
	}

	envUUID, err := uuid.Parse(envUuidStr)
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
func (sc *ServicesClient) DeleteEnv(ctx context.Context, uuidStr string, envUuidStr string) (string, error) {
	serviceUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return "", fmt.Errorf("invalid UUID: %w", err)
	}

	envUUID, err := uuid.Parse(envUuidStr)
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

// DeployApplication deploys an application by UUID
func (dc *DeploymentsClient) DeployApplication(ctx context.Context, uuidStr string, force bool, branch string) error {
	return dc.DeployApplicationWithOptions(ctx, uuidStr, &DeployApplicationOptions{
		Force:  force,
		Branch: branch,
	})
}

// DeployApplicationWithOptions deploys an application with advanced options
func (dc *DeploymentsClient) DeployApplicationWithOptions(ctx context.Context, uuidStr string, options *DeployApplicationOptions) error {
	params := &coolify.DeployByTagOrUuidParams{
		Uuid:  &uuidStr,
		Force: &options.Force,
	}

	// Branch and PR are mutually exclusive
	if options.Branch != "" && options.PR != nil {
		return fmt.Errorf("cannot specify both branch and PR - they are mutually exclusive")
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
		return fmt.Errorf("failed to deploy application: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("API error: %s", resp.Status())
	}

	return nil
}

// DeployService deploys a service by starting it (services use start/restart for deployment)
func (dc *DeploymentsClient) DeployService(ctx context.Context, uuidStr string) error {
	return dc.client.Services().Start(ctx, uuidStr)
}

// List returns deployment history for an application
func (dc *DeploymentsClient) List(ctx context.Context, appUuidStr string) ([]coolify.Application, error) {
	appUUID, err := uuid.Parse(appUuidStr)
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

// GetByUuid returns a deployment by UUID
func (dc *DeploymentsClient) GetByUuid(ctx context.Context, uuidStr string) (*coolify.ApplicationDeploymentQueue, error) {
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
