# Contributing to gosqlpp MCP Proxy

Thank you for your interest in contributing to gosqlpp MCP Proxy! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Code Style](#code-style)
- [Reporting Issues](#reporting-issues)

## Code of Conduct

This project adheres to a code of conduct that we expect all contributors to follow. Please be respectful and constructive in all interactions.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/go-sqlpp-mcp-proxy.git
   cd go-sqlpp-mcp-proxy
   ```
3. **Add the upstream remote**:
   ```bash
   git remote add upstream https://github.com/your-org/go-sqlpp-mcp-proxy.git
   ```

## Development Setup

### Prerequisites
- Go 1.24.5 or later
- Git
- A text editor or IDE

### Initial Setup
```bash
# Install dependencies
go mod download

# Build the project
go build -o mcp_sqlpp_proxy main.go

# Run tests
go test ./...
```

## Making Changes

### Before You Start
1. **Check existing issues** to see if your change is already being worked on
2. **Create an issue** for new features or significant changes to discuss the approach
3. **Keep changes focused** - one feature or fix per pull request

### Development Workflow
1. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following the code style guidelines

3. **Test your changes**:
   ```bash
   # Run all tests
   go test ./...
   
   # Test build
   go build -o test_binary main.go && rm test_binary
   
   # Manual testing
   ./mcp_sqlpp_proxy --help
   ```

4. **Commit your changes**:
   ```bash
   git add .
   git commit -m "Add feature: brief description"
   ```

## Testing

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### Writing Tests
- Add tests for new functionality
- Ensure existing tests still pass
- Follow Go testing conventions
- Use table-driven tests where appropriate

### Manual Testing
Test both transport modes:
```bash
# Test stdio mode
./mcp_sqlpp_proxy --transport stdio

# Test HTTP mode
./mcp_sqlpp_proxy --transport http --port 8080 --xfer-port 8891
```

## Submitting Changes

### Pull Request Process
1. **Update your branch** with the latest upstream changes:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Push your changes**:
   ```bash
   git push origin feature/your-feature-name
   ```

3. **Create a pull request** on GitHub with:
   - Clear title and description
   - Reference to related issues
   - Screenshots or examples if applicable
   - Test results

### Pull Request Guidelines
- **Keep PRs focused** on a single feature or fix
- **Write clear commit messages** following conventional commit format
- **Update documentation** if your changes affect usage
- **Add tests** for new functionality
- **Ensure CI passes** before requesting review

## Code Style

### Go Style Guidelines
- Follow standard Go formatting: `go fmt ./...`
- Use `go vet` to check for common issues
- Follow effective Go practices
- Use meaningful variable and function names
- Add comments for exported functions and complex logic

### Code Organization
- Keep functions focused and small
- Use appropriate error handling
- Follow Go naming conventions
- Group related functionality together

### Example Code Style
```go
// Good: Clear function name and documentation
func runStdioProxy(logger *log.Logger) error {
    cmd := exec.Command(sqlppPath, "-t", "stdio")
    // ... implementation
    return nil
}

// Good: Proper error handling
if err := cmd.Start(); err != nil {
    return fmt.Errorf("failed to start mcp_sqlpp: %w", err)
}
```

## Reporting Issues

### Bug Reports
When reporting bugs, please include:
- **Go version**: `go version`
- **Operating system** and version
- **Steps to reproduce** the issue
- **Expected behavior**
- **Actual behavior**
- **Log output** if available
- **Configuration** used

### Feature Requests
For feature requests, please include:
- **Use case**: Why is this feature needed?
- **Proposed solution**: How should it work?
- **Alternatives considered**: Other approaches you've thought about
- **Additional context**: Any other relevant information

### Issue Templates
Use the provided issue templates when available, or follow the structure above.

## Development Tips

### Debugging
- Use the generated log files for debugging
- Add temporary log statements for complex issues
- Test both stdio and HTTP modes
- Use `go run -race main.go` to check for race conditions

### Performance
- Profile memory usage for high-traffic scenarios
- Consider the impact of logging on performance
- Test with realistic data volumes

### Documentation
- Update README.md for user-facing changes
- Update code comments for internal changes
- Add examples for new features

## Questions?

If you have questions about contributing:
- Check existing issues and discussions
- Create a new issue with the "question" label
- Reach out to maintainers

Thank you for contributing to gosqlpp MCP Proxy!
