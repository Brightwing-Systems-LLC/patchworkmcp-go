// PatchworkMCP middleware for Go MCP servers.
//
// Provides heartbeat monitoring. Call StartMiddleware() after your server
// starts. It sends periodic heartbeat pings to PatchworkMCP so you can
// monitor server uptime and tool availability.
//
// Configuration via environment (or Options):
//   PATCHWORKMCP_URL          - default: https://patchworkmcp.com
//   PATCHWORKMCP_API_KEY      - required API key
//   PATCHWORKMCP_SERVER_SLUG  - required server identifier

package feedback

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const heartbeatInterval = 60 * time.Second

type heartbeatPayload struct {
	ServerSlug string   `json:"server_slug"`
	ToolCount  int      `json:"tool_count"`
	ToolNames  []string `json:"tool_names"`
}

// PatchworkMiddleware sends periodic heartbeats to PatchworkMCP.
type PatchworkMiddleware struct {
	apiURL     string
	apiKey     string
	serverSlug string
	toolNames  []string
	cancel     context.CancelFunc
}

// MiddlewareOptions configures the PatchworkMCP middleware.
type MiddlewareOptions struct {
	// PatchworkURL overrides PATCHWORKMCP_URL.
	PatchworkURL string
	// APIKey overrides PATCHWORKMCP_API_KEY.
	APIKey string
	// ServerSlug overrides PATCHWORKMCP_SERVER_SLUG.
	ServerSlug string
	// ToolNames is the list of tool names your server exposes.
	ToolNames []string
}

func (o *MiddlewareOptions) resolveURL() string {
	if o != nil && o.PatchworkURL != "" {
		return o.PatchworkURL
	}
	return getEnv("PATCHWORKMCP_URL", "https://patchworkmcp.com")
}

func (o *MiddlewareOptions) resolveKey() string {
	if o != nil && o.APIKey != "" {
		return o.APIKey
	}
	return os.Getenv("PATCHWORKMCP_API_KEY")
}

func (o *MiddlewareOptions) resolveSlug() string {
	if o != nil && o.ServerSlug != "" {
		return o.ServerSlug
	}
	if v := os.Getenv("PATCHWORKMCP_SERVER_SLUG"); v != "" {
		return v
	}
	return "unknown"
}

// StartMiddleware creates and starts the PatchworkMCP heartbeat middleware.
// Call this after your MCP server starts. Pass the list of tool names your
// server exposes for accurate heartbeat reporting.
//
//	mw := feedback.StartMiddleware(&feedback.MiddlewareOptions{
//	    ToolNames: []string{"my_tool_1", "my_tool_2"},
//	})
//	defer mw.Stop()
func StartMiddleware(opts *MiddlewareOptions) *PatchworkMiddleware {
	apiKey := opts.resolveKey()
	slug := opts.resolveSlug()

	if apiKey == "" || slug == "" || slug == "unknown" {
		fmt.Fprintln(os.Stderr, "PatchworkMCP middleware not started: missing API_KEY or SERVER_SLUG")
		return &PatchworkMiddleware{}
	}

	var toolNames []string
	if opts != nil {
		toolNames = opts.ToolNames
	}

	ctx, cancel := context.WithCancel(context.Background())
	mw := &PatchworkMiddleware{
		apiURL:     opts.resolveURL(),
		apiKey:     apiKey,
		serverSlug: slug,
		toolNames:  toolNames,
		cancel:     cancel,
	}

	go mw.heartbeatLoop(ctx)
	fmt.Fprintf(os.Stderr, "PatchworkMCP middleware started for %s\n", slug)
	return mw
}

// Stop cancels the heartbeat loop.
func (mw *PatchworkMiddleware) Stop() {
	if mw.cancel != nil {
		mw.cancel()
	}
}

func (mw *PatchworkMiddleware) heartbeatLoop(ctx context.Context) {
	// Send an initial heartbeat immediately.
	mw.sendHeartbeat()

	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			mw.sendHeartbeat()
		}
	}
}

func (mw *PatchworkMiddleware) sendHeartbeat() {
	payload := heartbeatPayload{
		ServerSlug: mw.serverSlug,
		ToolCount:  len(mw.toolNames),
		ToolNames:  mw.toolNames,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return
	}

	endpoint := mw.apiURL + "/api/v1/heartbeat/"
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	if mw.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+mw.apiKey)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "PatchworkMCP heartbeat failed: %v\n", err)
		return
	}
	resp.Body.Close()

	if resp.StatusCode >= 400 {
		fmt.Fprintf(os.Stderr, "PatchworkMCP heartbeat returned %d\n", resp.StatusCode)
	}
}
