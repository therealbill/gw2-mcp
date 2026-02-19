---
title: Contribute
---

# How to Contribute to GW2 MCP Server

**Goal**: Set up a local development environment, run quality checks, and submit a pull request to the GW2 MCP Server project.

**Time**: Approximately 20 minutes for initial setup.

## Prerequisites

Before starting, you should have:

- **Go 1.24+** installed -- verify with `go version`
- **Git** installed and configured with your GitHub account
- **Make** available on your system (GNU Make on Linux/macOS, or via Git Bash / MSYS2 on Windows)
- A GitHub account with permissions to fork repositories

## Development Setup

### 1. Fork and clone the repository

Fork the repository on GitHub, then clone your fork locally:

```bash
git clone https://github.com/<your-username>/gw2-mcp.git
cd gw2-mcp
```

### 2. Install development tools

Install `gofumpt`, `golangci-lint`, and `govulncheck`:

```bash
make tools
```

This runs `go install` for each tool, placing binaries in your `$GOPATH/bin`. Ensure that directory is on your `PATH`.

### 3. Build the project

Verify everything compiles:

```bash
make build
```

Expected result: a binary appears at `bin/gw2-mcp` (or `bin/gw2-mcp.exe` on Windows) with no errors.

### 4. Run the full pipeline

Confirm the project passes all checks before you make changes:

```bash
make all
```

This runs format, vet, lint, test, and build in sequence. If everything passes, your environment is ready.

## Running Tests

Run the test suite with race detection enabled:

```bash
make test
```

This executes `go test -v -race -coverprofile=coverage.out ./...`. All tests must pass before submitting a pull request.

To generate an HTML coverage report:

```bash
make coverage
```

Open `coverage.html` in your browser to inspect which lines are covered.

You can also run tests directly with Go if you prefer:

```bash
go test ./...
```

## Linting and Formatting

### Format your code

Run `gofumpt` (a stricter `gofmt`) and tidy module dependencies:

```bash
make format
```

This applies `gofumpt -w .` across the project and runs `go mod tidy`.

### Lint your code

Run `golangci-lint` to catch common issues:

```bash
make lint
```

Fix any reported issues before committing. The linter configuration lives in the project root and covers style, correctness, and performance checks.

### Run all quality checks at once

For a thorough check that includes vet, lint, tests, and a security vulnerability scan:

```bash
make check
```

## Commit Conventions

This project uses [Conventional Commits](https://www.conventionalcommits.org/). Every commit message must follow this format:

```
<type>: <short description>
```

Accepted types:

| Type | Use when |
|------|----------|
| `feat` | Adding a new feature or tool |
| `fix` | Fixing a bug |
| `docs` | Documentation-only changes |
| `refactor` | Code restructuring with no behavior change |
| `test` | Adding or updating tests |

Examples:

```
feat: add guild stash tool
fix: handle empty API response in wallet tool
docs: add configuration guide for Claude Desktop
refactor: extract shared HTTP client logic
test: add coverage for recipe search edge cases
```

Keep the description lowercase, imperative, and under 72 characters.

## Pull Request Workflow

### 1. Create a feature branch

Branch from `main` with a descriptive name:

```bash
git checkout -b feat/add-guild-treasury-tool
```

### 2. Make your changes

Edit the code. If you are adding a new tool, see [How to Add a New Tool](../how-to/add-a-new-tool/) for the step-by-step process.

### 3. Format and lint

```bash
make format
make lint
```

Fix any issues the linter reports.

### 4. Run tests

```bash
make test
```

All tests must pass with race detection enabled.

### 5. Commit with a conventional message

```bash
git add <changed-files>
git commit -m "feat: add guild treasury tool"
```

### 6. Push and open a pull request

```bash
git push origin feat/add-guild-treasury-tool
```

Open a pull request against the `main` branch on GitHub. In the PR description, include:

- What the change does and why
- Any new tools or endpoints added
- How you tested the change

### 7. Respond to review feedback

Address any review comments, push additional commits, and re-run `make check` to confirm everything still passes.

## Code Standards

Follow these standards for all contributions:

- **Format with gofumpt** -- never commit code that has not been run through `make format`
- **Lint with golangci-lint** -- all linter checks must pass cleanly with `make lint`
- **Write tests for core functionality** -- new tools and utility functions should include test coverage
- **Use race detection** -- tests run with `-race` by default; avoid data races in concurrent code
- **Keep commits atomic** -- each commit should represent one logical change

## Useful Make Targets

| Target | Description |
|--------|-------------|
| `make all` | Full pipeline: format, vet, lint, test, build |
| `make build` | Build the server binary |
| `make test` | Run tests with race detection |
| `make coverage` | Generate HTML coverage report |
| `make lint` | Run golangci-lint |
| `make format` | Run gofumpt and go mod tidy |
| `make tools` | Install development tools |
| `make check` | Vet, lint, test, and security scan |
| `make dev` | Run server with race detection (for local testing) |
| `make clean` | Remove build artifacts |
| `make help` | Print all available targets |

## Troubleshooting

### Problem: `gofumpt: command not found`
**Symptom**: `make format` fails because gofumpt is not installed.
**Cause**: Development tools have not been installed, or `$GOPATH/bin` is not on your `PATH`.
**Solution**: Run `make tools`, then confirm `$GOPATH/bin` is in your `PATH`:
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Problem: `golangci-lint` reports issues you did not introduce
**Symptom**: Lint failures appear in code you have not changed.
**Cause**: Your local `golangci-lint` version may differ from the project's expected version.
**Solution**: Re-run `make tools` to install the latest version, then run `make lint` again.

### Problem: Tests fail with race condition warnings
**Symptom**: `go test -race` reports data races.
**Cause**: Concurrent code accesses shared state without proper synchronization.
**Solution**: Use mutexes, channels, or other synchronization primitives. Run `make test` repeatedly to confirm the race is resolved -- race conditions can be intermittent.

### Problem: `make build` fails on Windows
**Symptom**: Build errors related to shell commands or paths.
**Cause**: GNU Make expects Unix-style commands.
**Solution**: Run Make from Git Bash or MSYS2 rather than PowerShell or Command Prompt.

## See also

- [How to Add a New Tool](../how-to/add-a-new-tool/) -- step-by-step guide for extending the server
- [Architecture](../explanation/architecture/) -- understand how the codebase is structured
