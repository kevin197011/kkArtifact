## 1. Web UI (Frontend) Cleanup
- [x] 1.1 Remove runtime `console.log` / `console.debug` statements from Web UI pages.
- [x] 1.2 Replace `console.error` used for “debugging” with:
  - user-facing feedback (`message.error`) where appropriate, and/or
  - a lightweight, optional debug logger that is disabled by default.
- [ ] 1.3 Ensure the Web UI build does not introduce new `console.*` logs (optional: lint rule to prevent `console.*` in production builds).

## 2. Server (Backend) Cleanup
- [x] 2.1 Ensure all debug-only logs are consistently gated behind a single debug switch (e.g. existing `DEBUG` env via `util.IsDebugMode()`).
- [x] 2.2 Replace any unconditional `fmt.Printf` debug outputs in long-running tasks with debug-gated logging.
- [x] 2.3 Confirm server error logs remain actionable but do not print secrets or full tokens.

## 3. Agent Cleanup
- [x] 3.1 Disable verbose request/response dumps (stderr “Request Details” blocks) by default.
- [x] 3.2 Add a single opt-in debug switch for the agent (env var and/or flag) to enable those dumps when needed.
- [x] 3.3 Ensure agent still reports errors clearly without relying on debug dumps.

## 4. Validation
- [x] 4.1 Verify Web UI `/audit-logs`, `/tokens`, `/inventory` have no `console.*` logs in normal runtime.
- [x] 4.2 Verify server normal mode logs contain no debug-only noise.
- [x] 4.3 Verify agent pull/push do not emit request/response dumps unless debug is enabled.

