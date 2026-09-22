# arch-lint

`arch-lint` is a static analysis tool for Go projects that enforces architectural rules by analyzing import paths and package structures.
It helps maintain clean and consistent codebases by preventing unwanted dependencies and enforcing modular boundaries.

## Features

- Use glob patterns to include or exclude packages for analysis.
- Define custom rules to forbid specific imports.
- Support exceptions to allow forbidden imports in restricted contexts.

## Installation

```
go install github.com/TheFellow/arch-lint@latest
```

## Usage

Run the linter with a configuration file:

```bash
./arch-lint -config=path/to/rules.yml
```

### Configuration

The linter reads the file selected by `-config` (default: `-config=config.yaml`) to define the rules for your project.
Below is an example configuration:

```yaml
specs:
  - name: no-experimental-imports
    packages:
      include:
        - "example/alpha/**"
      exclude:
        - "example/alpha/internal/exception/**"
    rules:
      forbid:
        - "example/alpha/experimental"
      except:
        - "example/alpha/internal/excluded"
      exempt:
        - "example/alpha/common"
```

Set `include_tests: true` to enable test loading in the original CLI. In the analyzer,
`include_tests: false` skips packages whose import path or package name ends in
`_test`; it does not filter individual in-package `_test.go` files from a pass.

Configuration files are validated against a built-in YAML schema before the linter runs.
Invalid files will cause arch-lint to exit with an error.

### Fields

- **name**: A descriptive name for the rule.
- **include**: Glob patterns specifying packages to include in the analysis.
- **exclude**: Glob patterns specifying packages to exclude from the analysis.
- **forbid**: Import paths that are forbidden.
- **except**: Importer package paths allowed to use forbidden imports.
- **exempt**: Import paths that are exempt from `forbid` rules.

`packages.include` and `packages.exclude` use doublestar globs, including brace
alternatives such as `example/{beta,delta}/**`. Rule patterns (`forbid`, `except`,
and `exempt`) use a separate, capture-aware syntax; they are not regular expressions
or the full doublestar glob language. Dots and other literal characters match exactly.

A `forbid` pattern supports a few special cases:
- `*`: Matches a single path segment.
- `**`: Matches multiple path segments, including none.
- `{variable}`: Matches a single path segment and captures it as a named variable.

All three rule fields also support whole-segment literal alternatives such as
`{beta,delta}`. `*`, `**`, captures, and alternatives must occupy complete segments.
`foo/**` matches `foo` and its descendants, but never `foobar` or `foo-v2`.
`foo/**/bar` also matches `foo/bar`. Variable names use letters, digits, and
underscores, cannot start with a digit, and must be unique within a pattern.
Malformed rule patterns and package globs are rejected when loading configuration.

An `except` pattern supports the same special cases as `forbid`, and one more
- `*`: Matches a single path segment.
- `**`: Matches multiple path segments, including none.
- `{variable}`: Matches this path segment when its value matches the one captured in the `forbid` pattern.
- `{!variable}`: Matches this path segment when its value **does not** match the one captured in the `forbid` pattern.

An `exempt` pattern supports the same special cases as `forbid`, and one more
- `*`: Matches a single path segment.
- `**`: Matches multiple path segments, including none.
- `{variable}`: Matches the value captured in the `forbid` pattern.
- `{!variable}`: Matches this path segment when its value **does not** match the one captured in the `forbid` pattern.

### How it works

First all packages in scope for analysis are collected.
That is, all packages that match the `include` glob patterns and do not match the `exclude` glob patterns.

Then each package in scope is analyzed.
During analysis there are two packages under consideration:
- The package being analyzed (the `current` package).
- The package being imported (the `imported` package).

The `current` package is forbidden from importing the `imported` package
if the `imported` package matches a `forbid` pattern.

Once forbidden, the `imported` package will be allowed if:
- The `current` package matches an `except` pattern.
- The `imported` package matches an `exempt` pattern.

Exceptions are spec-wide: an exception may use captures from any matching
`forbid` pattern, so reordering `forbid` entries does not change the result.
Each exception is evaluated against one forbid match at a time; captures from
separate patterns are not merged. An exception with an unbound variable does
not match, including negated references such as `{!domain}`.

`forbid: ["**"]` includes standard-library imports such as `fmt`. Use `exempt`
entries to allow particular imports, or scope `forbid` to a narrower path.

## Output

On the happy path the linter will output
```
✔ arch-lint: no forbidden imports found.
```
and exit with code 0.

On the unhappy path the linter will output

```
arch-lint: [<rule name>] package "path/to"  imports "forbidden/package"
```

and exit with code 1.

## go/analysis Integration

Since v0.0.13 (commit `7476878`), arch-lint also ships as a `go/analysis` Analyzer, which means it can run as a standalone singlechecker binary or integrate directly into golangci-lint.

Version v0.0.12 (`f31600a`) predates this integration and has neither the
singlechecker nor the golangci-lint module plugin.

### Standalone Singlechecker

Build and run the singlechecker binary:

```bash
go build -o arch-lint-checker ./cmd/arch-lint
./arch-lint-checker -config=.arch-lint.yml ./...
```

The singlechecker supports all standard `go/analysis` flags (`-json`, `-c=N`, `-test=false`, etc.). Config is resolved by walking up the directory tree for `.arch-lint.yml`, or you can pass `-config` explicitly.

### golangci-lint Module Plugin

arch-lint can be integrated into golangci-lint v2 as a [module plugin](https://golangci-lint.run/plugins/module-plugins/). This compiles arch-lint into a custom golangci-lint build (no fragile `.so` plugins needed).

1. Create a `.custom-gcl.yml` in your project:

```yaml
version: v2.0.0
plugins:
  - module: "github.com/TheFellow/arch-lint"
    import: "github.com/TheFellow/arch-lint/plugin"
    version: latest
```

2. Build your custom golangci-lint:

```bash
golangci-lint custom
```

3. Configure `.golangci.yml`:

```yaml
linters:
  enable:
    - archlint

linters-settings:
  custom:
    archlint:
      type: module
      description: Enforces architectural import boundaries
      settings:
        config: ".arch-lint.yml"
```

4. Run:

```bash
./custom-gcl run ./...
```

The plugin exposes the same Analyzer that the singlechecker uses, so behavior is identical.

## Development

### Prerequisites

- Go 1.23 or later

### Running Tests

```
go test ./...
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request with your changes.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.