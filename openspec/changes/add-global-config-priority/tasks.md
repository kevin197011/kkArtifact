## 1. Configuration Loading

- [x] 1.1 Update `GetGlobalConfigPath()` to support `/etc/kkArtifact/config.yml` (with capital A) as primary path, fallback to `/etc/kkartifact/kkartifact.yml` for backward compatibility
- [x] 1.2 Modify `config.Load()` to accept command-line overrides parameter
- [x] 1.3 Implement `mergeConfigsWithOverrides()` function that merges: global config → local config → command-line overrides
- [x] 1.4 Update `mergeConfigs()` to handle `ignore` field merging (combine arrays, remove duplicates, command-line takes precedence)

## 2. Command-Line Interface

- [x] 2.1 Add `--ignore` flag to `push` command (accepts comma-separated or multiple flags)
- [x] 2.2 Add `--ignore` flag to `pull` command (accepts comma-separated or multiple flags)
- [x] 2.3 Add `--server-url` flag to push/pull commands (optional, overrides config)
- [x] 2.4 Add `--token` flag to push/pull commands (optional, overrides config)
- [x] 2.5 Add `--concurrency` flag to push/pull commands (optional, overrides config)
- [x] 2.6 Update `runPush()` to pass command-line overrides to `config.Load()`
- [x] 2.7 Update `runPull()` to pass command-line overrides to `config.Load()`

## 3. Ignore Pattern Merging

- [x] 3.1 Implement `mergeIgnorePatterns()` function that:
  - Combines ignore patterns from global config, local config, and command-line
  - Removes duplicate patterns (command-line patterns take precedence)
  - Preserves order: global → local → command-line
- [x] 3.2 Update `mergeConfigsWithOverrides()` to use the new ignore merging logic
- [ ] 3.3 Add unit tests for ignore pattern merging scenarios

## 4. Testing

- [ ] 4.1 Add unit tests for configuration loading with global config priority
- [ ] 4.2 Add unit tests for command-line parameter overrides
- [ ] 4.3 Add unit tests for ignore pattern merging (global + local + command-line)
- [ ] 4.4 Add integration tests for push/pull with command-line overrides
- [ ] 4.5 Test backward compatibility with existing config files

## 5. Documentation

- [ ] 5.1 Update agent help text to document new `--ignore` flag
- [ ] 5.2 Update agent help text to document other override flags (`--server-url`, `--token`, `--concurrency`)
- [ ] 5.3 Document configuration priority order in README or agent documentation
