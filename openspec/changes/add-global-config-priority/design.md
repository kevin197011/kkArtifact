# Design: Global Config Priority and Command-Line Override

## Context

The agent currently supports:
- Global config: `/etc/kkartifact/kkartifact.yml`
- Local config: `.kkartifact.yml` (in current directory)
- Priority: Local config overrides global config

Users need:
1. Support for `/etc/kkArtifact/config.yml` path (with capital A)
2. Command-line flags to override config values
3. Merged `ignore` patterns (not replaced)

## Goals

- Support `/etc/kkArtifact/config.yml` as primary global config path
- Allow all config fields to be overridden via command-line flags
- Merge `ignore` patterns from all sources (global + local + command-line)
- Maintain backward compatibility with existing config files

## Non-Goals

- Changing the local config file name (remains `.kkartifact.yml`)
- Supporting multiple global config paths simultaneously
- Environment variable overrides (out of scope)

## Decisions

### Decision 1: Global Config Path Priority

**What**: Support both `/etc/kkArtifact/config.yml` and `/etc/kkartifact/kkartifact.yml` with priority order.

**Why**: 
- User expects `/etc/kkArtifact/config.yml` (with capital A)
- Existing code uses `/etc/kkartifact/kkartifact.yml` (lowercase)
- Need backward compatibility

**Implementation**:
- Try `/etc/kkArtifact/config.yml` first
- Fallback to `/etc/kkartifact/kkartifact.yml` if not found
- If neither exists, continue with local config only

### Decision 2: Configuration Priority Order

**What**: Priority order: Global config → Local config → Command-line flags

**Why**: 
- Global config provides system-wide defaults
- Local config allows project-specific overrides
- Command-line flags allow one-time overrides

**Implementation**:
- Load global config first (if exists)
- Merge local config (if exists) on top of global
- Apply command-line overrides last

### Decision 3: Ignore Pattern Merging

**What**: Merge `ignore` patterns from all sources, with command-line taking precedence for duplicates.

**Why**:
- Users want to combine global ignore patterns (e.g., `logs/`, `tmp/`) with command-line specific patterns (e.g., `test.rb`)
- Command-line should override global/local patterns if duplicate

**Implementation**:
- Combine all ignore patterns into a single array
- Remove duplicates, keeping the last occurrence (command-line patterns appear last)
- Preserve order: global → local → command-line

### Decision 4: Command-Line Flag Format

**What**: `--ignore` flag accepts comma-separated values or multiple flags.

**Why**:
- Flexible: `--ignore "logs/, tmp/, *.log"` or `--ignore logs/ --ignore tmp/`
- Consistent with common CLI patterns

**Implementation**:
- Support both `--ignore "pattern1, pattern2"` and `--ignore pattern1 --ignore pattern2`
- Parse comma-separated values and combine with multiple flags

## Risks / Trade-offs

### Risk: Path Case Sensitivity

**Risk**: Linux filesystems are case-sensitive, so `/etc/kkArtifact/config.yml` and `/etc/kkartifact/kkartifact.yml` are different files.

**Mitigation**: Try both paths in order, use whichever exists.

### Risk: Ignore Pattern Conflicts

**Risk**: If global config has `ignore: ["logs/"]` and command-line has `--ignore "logs/"`, we might have duplicates.

**Mitigation**: Remove duplicates during merge, keeping command-line version (appears last).

### Trade-off: Configuration Complexity

**Trade-off**: More configuration sources increase complexity.

**Mitigation**: Clear priority order and documentation. Most users will use config files, advanced users can use flags.

## Migration Plan

1. **Phase 1**: Add support for `/etc/kkArtifact/config.yml` (backward compatible)
2. **Phase 2**: Add command-line override flags (optional, doesn't break existing usage)
3. **Phase 3**: Implement ignore pattern merging (backward compatible - existing behavior preserved if no command-line flags used)

**Rollback**: If issues arise, revert to previous config loading logic. No data migration needed.

## Open Questions

- Should we support environment variable overrides? (Deferred - not in scope)
- Should we support multiple global config paths? (No - one global config is sufficient)
