# Changelog

## 0.1.0 (2026-02-24)

- Initial release
- MCP server integration via `RegisterFeedbackTool()`
- Manual registration via `NewFeedbackTool()` + `NewFeedbackHandler()`
- Core submission via `SendFeedback()`
- Retry logic with exponential backoff
- Context-aware cancellation during backoff
- Connection pooling via module-level http.Client
- Structured logging for failed submissions
