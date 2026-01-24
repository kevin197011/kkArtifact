# Change: Remove Debug Logs Across Web UI, Server, and Agent

## Why

The project currently contains multiple **debug-style logs** that leak into user environments:

- **Web UI**: `console.log/console.error` in runtime code (e.g. audit logs page, inventory page, tokens page)
- **Agent**: verbose request/response dumps to stderr (request headers, response bodies) intended for debugging
- **Server**: a mix of always-on logs and debug-gated logs (`DEBUG`), plus some `fmt.Printf` outputs in scheduled tasks

These logs have downsides:
- Noisy production output and harder troubleshooting of real incidents
- Potentially exposes sensitive operational details (URLs, headers, payload fragments)
- Inconsistent behavior across components (some are gated, some always-on)

This change standardizes logging so that **debug logs are removed or disabled by default**, and **only emitted when explicitly enabled**.

## What Changes

- **MODIFIED**: Web UI removes runtime `console.*` debug logs; user-facing errors remain via UI messaging.
- **MODIFIED**: Agent request/response “dump” logs are disabled by default and only enabled via an explicit debug flag/environment variable.
- **MODIFIED**: Server debug logs are consistently gated behind a single debug switch; scheduled tasks avoid `fmt.Printf` in normal mode.

**BREAKING**: None.

## Impact

- **Affected specs**: `artifact-web-ui`, `artifact-api`, `artifact-agent`
- **Affected code**:
  - Web UI pages containing `console.*` logs (e.g. `web-ui/src/pages/AuditLogs.tsx`, `InventoryPage.tsx`, `Tokens.tsx`)
  - Agent client verbose debug dumps (`agent/internal/client/*`)
  - Server debug gating / scheduled tasks (`server/internal/util/debug.go`, scheduler tasks, handlers)
- **User impact**: Cleaner production output; debugging remains possible via explicit opt-in.

