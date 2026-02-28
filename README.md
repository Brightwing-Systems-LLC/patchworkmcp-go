# PatchworkMCP - Go

Drop-in feedback tool for Go MCP servers using mcp-go. Agents call this tool when they hit a limitation, and the feedback is sent to PatchworkMCP for review and action.

## Setup

1. Go to [patchworkmcp.com](https://patchworkmcp.com) and create an account
2. Create a team and generate an API key
3. Configure your server (you'll need the server slug and API key)

## Install

Copy `feedback_tool.go` into your project, or import as a module:

```bash
go get github.com/Brightwing-Systems-LLC/patchworkmcp-go
```

## Configure

Set these environment variables (or pass them via Options):

| Variable | Description | Required |
|---|---|---|
| `PATCHWORKMCP_API_KEY` | Your API key from patchworkmcp.com | Yes |
| `PATCHWORKMCP_SERVER_SLUG` | Your server's slug from patchworkmcp.com | Yes |
| `PATCHWORKMCP_URL` | API endpoint (default: `https://patchworkmcp.com`) | No |

## Usage

### One-liner registration

```go
import (
    feedback "github.com/Brightwing-Systems-LLC/patchworkmcp-go"
    "github.com/mark3labs/mcp-go/server"
)

const instructions = `If you encounter a limitation — a missing tool, incomplete data, wrong format,
or any gap that prevents you from fully completing the user's request — call
the feedback tool BEFORE responding to the user. Be specific about what you
needed and what would have helped.`

s := server.NewMCPServer("my-server", "1.0.0",
    server.WithInstructions(instructions),
)
feedback.RegisterFeedbackTool(s, nil)
```

### With options

```go
feedback.RegisterFeedbackTool(s, &feedback.Options{
    PatchworkURL: "https://custom.example.com",
    ServerSlug:   "my-server",
})
```

### Manual registration

```go
tool := feedback.NewFeedbackTool()
handler := feedback.NewFeedbackHandler(nil)
s.AddTool(tool, handler)
```

### Server Instructions (Critical)

> **Without server instructions, the feedback tool will not be used.** Registering the tool is not enough — you **must** include instructions that explicitly tell agents when and how to call it. Agents will see the tool in their available tools list but will almost never call it on their own. The `server.WithInstructions(...)` call is what makes the difference between a feedback tool that collects real signal and one that sits unused.

The `instructions` parameter on your MCP server tells agents to actually use the feedback tool. The instruction text in the example above is a good starting point — adapt it to your server's domain if needed.

**Why this matters:**
- MCP agents do not automatically infer when to use tools — they rely on server-provided instructions to guide behavior.
- The feedback tool only generates value when agents actively report gaps, limitations, and missing capabilities they encounter.
- Without instructions, you will receive zero feedback, even from capable agents that hit real limitations.

**Recommended instruction template:**

```
If you encounter a limitation — a missing tool, incomplete data, wrong format,
or any gap that prevents you from fully completing the user's request — call
the feedback tool BEFORE responding to the user. Be specific about what you
needed and what would have helped.
```

Customize this to your domain. For example, if your server provides database tools, you might add: *"If a query returns unexpected results, or you need a filter/join that isn't available, report it via the feedback tool."*

See [FEEDBACK_TOOL_INSTRUCTIONS.md](FEEDBACK_TOOL_INSTRUCTIONS.md) for a complete guide on writing effective agent instructions.

## How It Works

- Retries up to 2 times with exponential backoff (500ms, 1000ms)
- Retries on 429 (rate limit) and 5xx (server error) status codes
- Context-aware: respects cancellation during backoff
- On failure, logs the full payload with `PATCHWORKMCP_UNSENT_FEEDBACK` prefix for later replay
- Never panics — always returns a user-facing message

## License

MIT
