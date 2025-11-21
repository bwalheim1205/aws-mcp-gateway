# Contributing to AWS MCP Gateway

Thanks for your interest in contributing and welcome aboard! We want to maintain a readable and high quality code base to make it easy for people to contribute. Below are some of the requirements for the codebase. Happy coding!

---

## Prerequisites

Before contributing, please ensure you have:

* **Go installed** (minimum version `1.21`, but check `go.mod` for exact requirement)
* [`pre-commit`](https://pre-commit.com/) installed and configured
* Your editor/IDE set up for Go development (optional, but recommended)

---

## Setting Up Pre-commit

We use pre-commit hooks to enforce formatting, linting, and other checks.

1. Install pre-commit if you haven’t already:

```bash
pip install pre-commit
```

2. Install hooks for this project:

```bash
pre-commit install
```

3. Run hooks manually if needed:

```bash
pre-commit run --all-files
```

> **Note:** All pull requests (PRs) must pass pre-commit hooks before merging.

---

## Writing Code

When contributing:

* Write clear, readable Go code.
* Follow standard Go formatting (`gofmt`) and import conventions.
* Ensure your changes are isolated to the functionality you are adding or fixing.

---

## Unit Tests

* **All new functionality must include unit tests.**
* **Existing unit tests must pass.**
* You may **modify existing tests** if your changes require it, but do not break existing behavior without careful consideration.

Run tests locally:

```bash
go test ./...
```

---

## Pull Requests

When submitting a PR:

1. Ensure your branch is up-to-date with `main` (or the default branch).
2. Run all unit tests and pre-commit hooks locally.
3. Provide a clear description of your changes.
4. Include any relevant issue numbers (if applicable).

---

## Code Review

* PRs will be reviewed for functionality, code quality, and adherence to style.
* Reviewers may request changes; please address them before the PR is merged.

---

## Thank You

Thanks for helping improve this project! Following these guidelines ensures that contributions are high-quality, consistent, and easy to review.
