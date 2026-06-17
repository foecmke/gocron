# Contributing Guide

We warmly welcome community contributions! Whether it's bug fixes, new feature development, documentation improvements, or issue reports, all are valuable support for the project.

## Git Commit Convention

To maintain clear and consistent project history, please follow this commit message format:

### Basic Format

```text
type(scope): subject

body

footer
```

### Detailed Component Description

- **type**: Commit type (required)
- **scope**: Affected scope (optional) - such as web, api, auth, etc., in parentheses
- **subject**: Brief description (required) - one line explaining what this commit does
- **body**: Detailed description (optional) - explain why this change was made and how it was implemented
- **footer**: Footer information (optional) - such as closed issue numbers, breaking change notes, etc.

### Type Descriptions

| Type | Description |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation update |
| `style` | Code formatting (no functional changes, e.g., whitespace, semi-colons) |
| `refactor` | Code refactoring (neither a fix nor a feature) |
| `test` | Test related |
| `chore` | Build process or auxiliary tool changes |
| `perf` | Performance optimization |
| `ci` | CI/CD related |

### Scope Examples

| Scope | Description |
|-------|-------------|
| `web` | Frontend related |
| `api` | Backend API related |
| `node` | Task node related |
| `auth` | Authentication related |
| `i18n` | Internationalization related |
| `log` | Logging related |
| `db` | Database related |

### Real Examples

**Example 1: New Feature**

```text
feat(auth): add two-factor authentication

To improve system security, implemented TOTP-based two-factor authentication:
- Support for Google Authenticator and other auth apps
- Provide both QR code and manual entry setup methods
- Complete enable/disable workflow

Closes #123
```

**Example 2: Bug Fix**

```bash
fix(web): fix task list pagination display issue
```

**Example 3: Documentation**

```bash
docs(readme): update installation instructions
```

## How to Contribute

1. **Fork the project** - Click the Fork button in the top right of the GitHub repository
2. **Create a branch** - `git checkout -b feature/your-feature`
3. **Commit code** - Follow the commit convention above
4. **Push branch** - `git push origin feature/your-feature`
5. **Create PR** - Create a Pull Request on GitHub

## Development Guidelines

- Run `make test` before committing to ensure tests pass
- Add corresponding test cases for new features
- Update relevant documentation for important features
- Maintain consistent code style
- Keep commit messages concise and clear in Chinese or English

We look forward to your participation in making gocron even better! 🚀
