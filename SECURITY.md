# Security Policy

## Reporting Vulnerabilities

If you discover a security vulnerability, please report it responsibly:

1. **Do NOT** open a public GitHub issue
2. Email the maintainer directly (see GitHub profile)
3. Include a description of the vulnerability
4. Include steps to reproduce if possible

## Security Considerations

### File System Access
- CodeMetrics reads source files from the file system
- It does NOT write to files (except via `--file` output option)
- It does NOT execute any code being analyzed
- It does NOT send data to external services

### Input Validation
- File paths are validated before reading
- Unsupported file extensions are handled gracefully
- Malformed source code is handled without crashes

### Dependencies
- tree-sitter is a well-maintained parsing library
- No network dependencies at runtime
- All dependencies are vendored or checksummed

## Scope

This security policy applies to:
- The `codemetrics` CLI tool
- The `internal/` packages
- The `cmd/codemetrics` binary

Out of scope:
- Test files and fixtures
- Documentation
- CI/CD configuration

## Updates

Security updates will be released as patch versions (e.g., 1.0.0 -> 1.0.1).
