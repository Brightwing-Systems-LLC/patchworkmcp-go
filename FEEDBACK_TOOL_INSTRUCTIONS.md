# Feedback Tool: Agent Instructions Guide

## Why Server Instructions Are Critical

Registering the feedback tool on your MCP server is **not enough**. You **must** include explicit instructions that tell agents when and how to call it. Without instructions, agents will see the tool but will almost never call it on their own.

**The `server.WithInstructions(...)` call is what makes the difference between a feedback tool that collects real signal and one that sits unused.**

MCP agents do not automatically infer when to use tools. They rely on server-provided instructions to guide behavior. If you skip this step, you will receive zero feedback — even from capable agents that hit real limitations during their work.

## How to Add Instructions

When creating your MCP server, pass an instruction string via `server.WithInstructions(...)`:

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

## Writing Effective Instructions

Good instructions should tell the agent:

1. **When to call the tool** — on encountering any limitation, missing capability, or gap.
2. **How urgently** — before responding to the user, not as an afterthought.
3. **What to include** — be specific about what was needed and what would have helped.

### Recommended Template

```
If you encounter a limitation — a missing tool, incomplete data, wrong format,
or any gap that prevents you from fully completing the user's request — call
the feedback tool BEFORE responding to the user. Be specific about what you
needed and what would have helped.
```

### Domain-Specific Examples

Tailor the instructions to the tools your server provides. Here are examples for common domains:

**Database server:**
```
If a query returns unexpected results, or you need a filter, join, aggregation,
or table that isn't available, report it via the feedback tool before responding
to the user. Include the query you tried and what result you expected.
```

**File management server:**
```
If you cannot read, write, move, or transform a file in the way the user needs —
due to a missing tool, unsupported format, or permission issue — call the feedback
tool before responding. Describe the operation you attempted and what blocked you.
```

**API integration server:**
```
If an API call fails, returns incomplete data, or you need an endpoint or parameter
that doesn't exist, call the feedback tool before responding to the user. Include
the endpoint you tried, the parameters you used, and what you expected to receive.
```

## What Happens Without Instructions

| Scenario | With Instructions | Without Instructions |
|---|---|---|
| Agent hits a missing tool | Reports the gap via feedback tool | Silently works around it or tells the user it can't help |
| Agent gets incomplete data | Files specific feedback about what was missing | Returns partial results without reporting the limitation |
| Agent needs a different format | Reports the format mismatch | Attempts manual conversion or fails quietly |
| You want to improve your server | You get actionable feedback to prioritize | You have no visibility into what agents actually need |

## Common Mistakes

1. **Registering the tool without instructions.** The tool exists but agents don't know when to use it. This is the most common mistake.

2. **Instructions that are too vague.** Saying "use the feedback tool when appropriate" gives the agent no real guidance. Be specific about the situations that should trigger feedback.

3. **Instructions that are too narrow.** If you only mention one scenario (e.g., "missing tools"), agents won't report other types of gaps like incomplete data or wrong formats.

4. **Placing instructions only in documentation.** Agents don't read your docs — they read the `instructions` field on the MCP server. The instructions must be passed programmatically via `server.WithInstructions(...)`.

## Verifying It Works

After adding instructions, verify that the feedback tool is being called:

1. **Test manually** — Connect an agent (Claude Desktop, Cursor, Claude Code, etc.) and give it a task that you know will hit a limitation. Confirm that it calls the feedback tool.
2. **Check your PatchworkMCP dashboard** — Log in at [patchworkmcp.com](https://patchworkmcp.com) and verify that feedback submissions are appearing.
3. **Check server logs** — If submissions fail, look for lines prefixed with `PATCHWORKMCP_UNSENT_FEEDBACK` in your server logs. These contain the full payload for later replay.
