# Contributing to POS API

Thank you for your interest in contributing to POS API! This document provides guidelines and steps for contributing.

## 📋 Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Focus on the best outcome for the project

## 🚀 Getting Started

1. Fork the repository
2. Clone your fork: `git clone <your-fork-url>`
3. Create a branch: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Submit a pull request

## 💻 Development Setup

```bash
# Install dependencies
go mod tidy

# Run tests
make test

# Run linter (requires golangci-lint)
make lint

# Format code
make fmt
```

## 📝 Coding Standards

### Go Style

- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use `gofmt` for formatting
- Use meaningful variable and function names
- Write comments for exported functions

### Project Structure

- `handler/` - HTTP handlers (thin layer, validation only)
- `service/` - Business logic
- `repository/` - Database access
- `dto/` - Request/Response structures
- `models/` - Domain models

### Error Handling

```go
// Good
if err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}

// Bad
if err != nil {
    return err
}
```

## 🔀 Git Workflow

### Branch Naming

- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation
- `refactor/` - Code refactoring

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add customer loyalty points
fix: resolve stock calculation bug
docs: update API documentation
refactor: simplify transaction service
```

## 🧪 Testing

- Write tests for new features
- Ensure all tests pass before submitting PR
- Aim for meaningful test coverage

```bash
# Run tests
make test

# Run tests with coverage
make test-coverage
```

## 📬 Pull Request Process

1. Update documentation if needed
2. Add/update tests
3. Ensure CI passes
4. Request review from maintainers
5. Address feedback promptly

## ❓ Questions?

Open an issue for any questions or concerns.

Thank you for contributing! 🎉
