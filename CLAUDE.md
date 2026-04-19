# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a minimal Go module named `hello` using Go 1.26.2. The project is in its initial state with only a `go.mod` file.

## Development Commands

```bash
# Build the module
go build ./...

# Run all tests
go test ./...

# Run a specific test (using -run with test name pattern)
go test -run TestName ./...

# Format code
go fmt ./...

# Vet code for common mistakes
go vet ./...

# Tidy dependencies
go mod tidy

# Add a dependency
go get <package-path>
```

## Architecture

This is a new Go module with no source code yet. Standard Go project structure should be followed as the codebase grows:

- Place main packages in `main.go` or `cmd/` directory
- Place library code in `internal/` (private) or `pkg/` (public) directories
- Place tests alongside source files with `_test.go` suffix
