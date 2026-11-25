package openlist

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"aqua/pkg/transport/client"
	"aqua/services/admin/types"
)

// -----------------------------------------------------------------------------
// Constants
// -----------------------------------------------------------------------------

const (
	apiPathLogin  = "/api/auth/login/hash"
	apiPathPut    = "/api/fs/put"
	apiPathList   = "/api/fs/list"
	apiPathRemove = "/api/fs/remove"
	apiPathMkdir  = "/api/fs/mkdir"
)

// -----------------------------------------------------------------------------
// OpenList Domain Types & Errors
// -----------------------------------------------------------------------------

// Response is the generic payload wrapper returned by OpenList API.
type Response struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// APIError represents a domain-level error returned by the OpenList Server.
type APIError struct {
	HTTPStatus int
	Code       int
	Message    string
	RequestURL string
	Headers    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("openlist API error: [http %d] code=%d msg=%q url=%s", e.HTTPStatus, e.Code, e.Message, e.RequestURL)
}

// isTokenInvalidErr checks whether an error indicates a stale or invalid JWT token.
func isTokenInvalidErr(err error) bool {
	if err == nil {
		return false
	}
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.HTTPStatus == 401 || apiErr.Code == 401 || strings.Contains(apiErr.Message, "token is invalidated")
	}
	return strings.Contains(err.Error(), "401") || strings.Contains(err.Error(), "token is invalidated")
}

// isStoreNotFoundErr checks whether an error indicates the target store/directory doesn't exist.
func isStoreNotFoundErr(err error) bool {
	if err == nil {
		return false
	}
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		msg := strings.ToLower(apiErr.Message)
		return strings.Contains(msg, "store not found")
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "store not found")
}

// -----------------------------------------------------------------------------
// Client Setup
// -----------------------------------------------------------------------------

// Client handles communication with an OpenList server.
type Client struct {
	cfg    types.OpenListConfig
	client *http.Client
	mu     sync.Mutex
	token  string
	until  time.Time
}

// NewClient initializes a new OpenList API client.
func NewClient(cfg types.OpenListConfig) *Client {
	return &Client{
		cfg:    cfg,
		client: client.NewHTTPClient(30 * time.Second),
	}
}

// -----------------------------------------------------------------------------
// Authentication Flow
// -----------------------------------------------------------------------------

func (c *Client) login(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Use cached token if valid
	if c.token != "" && time.Now().Before(c.until) {
		return c.token, nil
	}

	reqBody := map[string]string{
		"username": c.cfg.Username,
		"password": c.cfg.Password,
	}

	var data struct {
		Token string `json:"token"`
	}

	if err := c.doHTTPRequest(ctx, http.MethodPost, apiPathLogin, "", reqBody, &data, nil); err != nil {
		return "", fmt.Errorf("openlist login: %w", err)
	}

	if data.Token == "" {
		return "", errors.New("openlist login: token empty in response")
	}

	c.token = data.Token
	// Cache the token proactively (usually OpenList tokens last a few hours, caching for 50m is a safe heuristic)
	c.until = time.Now().Add(50 * time.Minute)
	return c.token, nil
}

func (c *Client) invalidateToken() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.token = ""
	c.until = time.Time{}
}

// -----------------------------------------------------------------------------
// HTTP Utilities
// -----------------------------------------------------------------------------

// doAuthRequest extends doHTTPRequest with automatic token injection and transparent 401 retry.
func (c *Client) doAuthRequest(ctx context.Context, method, apiPath string, reqBody any, outData any, headers map[string]string) error {
	execute := func() error {
		token, err := c.login(ctx)
		if err != nil {
			return err
		}
		return c.doHTTPRequest(ctx, method, apiPath, token, reqBody, outData, headers)
	}

	err := execute()
	// Auto-retry once unconditionally if the token just expired
	if err != nil && isTokenInvalidErr(err) {
		c.invalidateToken()
		return execute()
	}
	return err
}

// doHTTPRequest wraps all boilerplate for net/http requests, header formatting, and API error surfacing.
func (c *Client) doHTTPRequest(ctx context.Context, method, apiPath, token string, reqBody any, outData any, headers map[string]string) error {
	base := strings.TrimRight(c.cfg.Endpoint, "/")
	reqURL := base + apiPath

	// 1. Prepare payload
	var r io.Reader
	if reqBody != nil {
		if reader, ok := reqBody.(io.Reader); ok {
			r = reader
		} else {
			b, err := json.Marshal(reqBody)
			if err != nil {
				return fmt.Errorf("marshal request body: %w", err)
			}
			r = bytes.NewReader(b)
		}
	}

	// 2. Prepare request
	req, err := http.NewRequestWithContext(ctx, method, reqURL, r)
	if err != nil {
		return fmt.Errorf("create http request: %w", err)
	}

	// OpenList expects the raw token in Authorization (no "Bearer " scheme)
	if token != "" {
		req.Header.Set("Authorization", token)
	}
	if reqBody != nil && r != reqBody {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// 3. Dispatch HTTP Call
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("execute http request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// 4. Capture raw body and construct error traces if anything fails
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response bytes: %w", err)
	}

	apiErrorTrace := func(code int, msg string) *APIError {
		return &APIError{
			HTTPStatus: resp.StatusCode,
			Code:       code,
			Message:    msg,
			RequestURL: reqURL,
			Headers:    formatHeaders(req.Header, token), // Save safely-masked headers into the error for debugging
		}
	}

	if resp.StatusCode != http.StatusOK {
		return apiErrorTrace(resp.StatusCode, string(raw))
	}

	var olResp Response
	if err := json.Unmarshal(raw, &olResp); err != nil {
		return fmt.Errorf("decode openlist wrapper: %w (raw payload: %s)", err, string(raw))
	}

	if olResp.Code != 200 {
		return apiErrorTrace(olResp.Code, olResp.Message)
	}

	// 5. Decode OpenList successful `data` field
	if outData != nil && len(olResp.Data) > 0 && string(olResp.Data) != "null" {
		if err := json.Unmarshal(olResp.Data, outData); err != nil {
			return fmt.Errorf("decode openlist payload struct: %w", err)
		}
	}

	return nil
}

// -----------------------------------------------------------------------------
// Logging Redaction Internals
// -----------------------------------------------------------------------------

func maskToken(token string) string {
	if len(token) <= 10 {
		return "***"
	}
	return token[:5] + "***" + token[len(token)-5:]
}

func formatHeaders(header http.Header, token string) string {
	b, _ := json.MarshalIndent(header, "", "  ")
	s := string(b)
	if token != "" {
		s = strings.ReplaceAll(s, token, maskToken(token))
	}
	return s
}

// -----------------------------------------------------------------------------
// Public Client Operations
// -----------------------------------------------------------------------------

// Mkdir recursively creates directories at the specific virtual target path.
func (c *Client) Mkdir(ctx context.Context, path string) error {
	reqBody := map[string]string{
		"path": path,
	}
	if err := c.doAuthRequest(ctx, http.MethodPost, apiPathMkdir, reqBody, nil, nil); err != nil {
		return fmt.Errorf("openlist call Mkdir: %w", err)
	}
	return nil
}

// Upload buffers a file to the specified target path. Retries folder creation on store gap.
// Returns the expected full Web URL of the file.
func (c *Client) Upload(ctx context.Context, destPath string, content io.Reader) (string, error) {
	headers := map[string]string{
		"Content-Type": "application/octet-stream",
		"File-Path":    url.PathEscape(destPath),
		"As-Task":      "false",
	}

	err := c.doAuthRequest(ctx, http.MethodPut, apiPathPut, content, nil, headers)

	// Recoverable error: Auto-mkdir if the parent directory does not exist yet (like rolling date folders)
	if err != nil && isStoreNotFoundErr(err) {
		dirPath := destPath[:strings.LastIndex(destPath, "/")]
		if errMkdir := c.Mkdir(ctx, dirPath); errMkdir == nil {
			// One time retry upload if Mkdir succeeded
			err = c.doAuthRequest(ctx, http.MethodPut, apiPathPut, content, nil, headers)
		}
	}

	if err != nil {
		// Embed comprehensive context to help the developer quickly diagnose it in the web service logs
		var apiErr *APIError
		if errors.As(err, &apiErr) {
			return "", fmt.Errorf("upload rejected by OpenList (%s): code=%d msg=%q\nheaders=%s", destPath, apiErr.Code, apiErr.Message, apiErr.Headers)
		}
		return "", fmt.Errorf("upload failed (%s): %w", destPath, err)
	}

	base := strings.TrimRight(c.cfg.Endpoint, "/")
	return base + "/" + strings.TrimLeft(destPath, "/"), nil
}

// List returns files and directories present at dirPath.
func (c *Client) List(ctx context.Context, dirPath string, page, perPage int) ([]FsObject, int, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 100
	}

	reqBody := ListRequest{
		Path:    dirPath,
		Page:    page,
		PerPage: perPage,
		Refresh: true,
	}

	var data struct {
		Content []FsObject `json:"content"`
		Total   int        `json:"total"`
	}

	if err := c.doAuthRequest(ctx, http.MethodPost, apiPathList, reqBody, &data, nil); err != nil {
		return nil, 0, fmt.Errorf("openlist call List (%s): %w", dirPath, err)
	}

	return data.Content, data.Total, nil
}

// Remove permanently deletes the files specified by names under a parent directory.
func (c *Client) Remove(ctx context.Context, dir string, names []string) error {
	if len(names) == 0 {
		return nil
	}

	reqBody := RemoveRequest{
		Dir:   dir,
		Names: names,
	}

	if err := c.doAuthRequest(ctx, http.MethodPost, apiPathRemove, reqBody, nil, nil); err != nil {
		return fmt.Errorf("openlist call Remove (%s): %w", dir, err)
	}

	return nil
}
