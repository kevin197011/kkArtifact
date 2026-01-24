## Context

The codebase has mixed logging patterns:
- UI uses `console.*`
- Server uses `log.Printf`, plus a partial debug-gate (`DEBUG`)
- Agent prints detailed request/response dumps intended for troubleshooting

## Goals

- Debug logs are **not emitted by default** in production or normal usage.
- Debug logs can be **explicitly enabled** for troubleshooting.
- Logs must not leak sensitive information (tokens, credentials).

## Decisions

### Decision 1: “Debug logs” definition

Debug logs include:
- `console.log`, `console.debug` in Web UI
- Agent stderr “Request Details” / “Response Body” dump blocks
- Server logs that are only helpful for tracing internals and not for normal operation

### Decision 2: Default-off + explicit opt-in

- Server continues to use existing `DEBUG` env gate (via `util.IsDebugMode()`), expanded consistently where needed.
- Agent introduces a single debug gate (env/flag) and wraps verbose dumps behind it.
- Web UI avoids `console.*` at runtime; any optional debug output should be behind an explicit enablement (e.g. build-time flag or storage key).

### Decision 3: Keep operational logs

Operational logs (startup info, critical warnings/errors) remain, but should:
- avoid full token output
- avoid printing request/response bodies unless explicitly enabled

