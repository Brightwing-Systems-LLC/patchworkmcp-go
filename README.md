# PatchworkMCP - Go SDK

PatchworkMCP helps MCP server authors understand how agents actually use their servers — what works, what's missing, and what to build next. This SDK adds three things to your server:

1. **A feedback tool** that agents call when they hit a limitation (missing tool, incomplete data, wrong format). This creates a structured stream of real-world gap reports visible on your PatchworkMCP dashboard.
2. **Server instructions** that tell agents when and how to use the feedback tool. Without these, agents see the tool but don't know when to call it.
3. **A heartbeat monitor** that pings PatchworkMCP every 60 seconds so you can track uptime and see which tools your server exposes.

None of these change your server's existing behavior. They're additive — your existing tools, resources, and prompts stay exactly as they are.

## Quick Start

### 1. Create a PatchworkMCP account

Go to [patchworkmcp.com](https://patchworkmcp.com), create a team, register your server, and generate an API key. You'll need:
- Your **API key**
- Your **server slug** (the identifier for your server on PatchworkMCP)

### 2. Install the SDK

```bash
go get github.com/Brightwing-Systems-LLC/patchworkmcp-go
```

No extra dependencies beyond [mcp-go](https://github.com/mark3labs/mcp-go) and the Go standard library.

### 3. Set environment variables

| Variable | Description | Required |
|---|---|---|
| `PATCHWORKMCP_API_KEY` | Your API key from patchworkmcp.com | Yes |
| `PATCHWORKMCP_SERVER_SLUG` | Your server's slug from patchworkmcp.com | Yes |
| `PATCHWORKMCP_URL` | API endpoint (default: `https://patchworkmcp.com`) | No |

### 4. Add to your server

```go
import (
    feedback "github.com/Brightwing-Systems-LLC/patchworkmcp-go"
    "github.com/mark3labs/mcp-go/server"
)

// -- Server instructions --
// This is what tells agents to use the feedback tool. Without it,
// agents may see the tool but won't know when to call it.
const instructions = `If you encounter a limitation — a missing tool, incomplete data, wrong format,
or any gap that prevents you from fully completing the user's request — call
the feedback tool BEFORE responding to the user. Be specific about what you
needed and what would have helped.`

s := server.NewMCPServer("my-server", "1.0.0",
    server.WithInstructions(instructions),
)

// -- Feedback tool --
// Adds a "feedback" tool that agents call when they hit a gap.
feedback.RegisterFeedbackTool(s, nil)

// -- Heartbeat monitor --
// Sends a ping every 60s with your server slug and tool list.
mw := feedback.StartMiddleware(&feedback.MiddlewareOptions{
    ToolNames: []string{"my_tool_1", "my_tool_2"},
})
defer mw.Stop()
```

That's it. Three additions, each one line.

## What Each Piece Does

### Feedback Tool

When an agent can't find the right tool, gets incomplete results, or has to work around a limitation, it calls the `feedback` tool with structured data:

- **what_i_needed** — the capability or data it was looking for
- **what_i_tried** — which tools it tried and what happened
- **gap_type** — category: `missing_tool`, `incomplete_results`, `missing_parameter`, `wrong_format`, `other`
- **suggestion** — the agent's idea for what would help

These reports appear on your PatchworkMCP dashboard, giving you a prioritized list of what to build next based on real agent usage.

### Server Instructions (Critical)

**You must add instructions to your MCP server telling the agent to use the feedback tool.** This is the single most important step in the integration. Without explicit instructions, agents will silently ignore the feedback tool — even though it appears in their tool list.

The `instructions` parameter on your MCP server is what makes the feedback tool useful. It tells agents: "if you hit a wall, report it before responding." Key principles:

1. **Tell the agent it is required.** Agents treat server instructions as authoritative. If you don't say "you must call the feedback tool," they won't.
2. **Specify when to call it.** List the concrete scenarios: missing tool, incomplete results, wrong format, about to say "not possible."
3. **Say to call it BEFORE responding.** If the agent responds first, it rarely circles back to submit feedback.
4. **Ask for specifics.** Generic feedback like "something was missing" is not actionable.

Without instructions, PatchworkMCP receives zero signal about what your server is missing. **No instructions = no feedback.**

For the full guide on writing effective agent instructions, see [FEEDBACK_TOOL_INSTRUCTIONS.md](FEEDBACK_TOOL_INSTRUCTIONS.md).

### Heartbeat Monitor

The middleware sends a heartbeat to PatchworkMCP every 60 seconds containing:
- Your server slug
- How many tools your server exposes
- The list of tool names

This powers uptime monitoring on your dashboard and lets PatchworkMCP track which tools are live.

## Alternative Integration Patterns

### Manual tool registration

```go
tool := feedback.NewFeedbackTool()
handler := feedback.NewFeedbackHandler(nil)
s.AddTool(tool, handler)
```

### Override configuration

```go
feedback.RegisterFeedbackTool(s, &feedback.Options{
    PatchworkURL: "https://custom.example.com",
    ServerSlug:   "my-server",
    APIKey:       "my-key",
})

mw := feedback.StartMiddleware(&feedback.MiddlewareOptions{
    PatchworkURL: "https://custom.example.com",
    ServerSlug:   "my-server",
    APIKey:       "my-key",
    ToolNames:    []string{"my_tool_1"},
})
```

## Reliability

- Feedback submissions retry up to 2 times with exponential backoff (500ms, 1000ms)
- Retries on 429 (rate limit) and 5xx (server error) status codes
- Context-aware: respects cancellation during backoff
- On failure, logs the full payload with `PATCHWORKMCP_UNSENT_FEEDBACK` prefix for later replay
- Never panics — always returns a user-facing message
- Heartbeats are fire-and-forget; failures are logged but don't affect your server

## License

MIT
