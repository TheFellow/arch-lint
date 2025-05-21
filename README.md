# Go Arch Lint

`go-arch-lint` is a static analysis tool for Go projects that enforces architectural rules by analyzing import paths and package structures.
It helps maintain clean and consistent codebases by preventing unwanted dependencies and enforcing modular boundaries.

## Features

- Define custom rules to forbid specific imports.
- Use glob patterns to include or exclude files for analysis.
- Support for exceptions to allow specific imports in restricted contexts.
- Cross-platform compatibility with normalized file paths.

## Installation

TODO

## Usage

Run the linter with a configuration file:

```bash
./go-arch-lint -config=path/to/rules.yml
```

### Configuration

The linter uses a `rules.yml` file to define the rules for your project.
Below is an example configuration:

```yaml
specs:
  - name: no-experimental-imports
    files:
      include:
        - "example/alpha/{*.go,**/*.go}"
      exclude:
        - "example/alpha/internal/exception/*.go"
    rules:
      forbid:
        - "example/alpha/experimental"
      except:
        - "example/alpha/internal/excluded"
```

### Fields

- **name**: A descriptive name for the rule.
- **include**: Glob patterns specifying files to include in the analysis.
- **exclude**: Glob patterns specifying files to exclude from the analysis.
- **forbid**: Import paths that are forbidden.
- **except**: Import paths that are exceptions to the forbidden rules.

A `forbid` pattern supports a few special cases:
- `{variable}`: Matches a single path segment and captures it as a named variable.
- `*`: Matches a single path segment.
- `**`: Matches multiple path segments, including none.

An `except` patterns supports the same special cases as `forbid`, and one more
- `{variable}`: Matches a single path segment and captures it as a named variable.
- `{!variable}`: Does not match a path segment if it has the value captured in the `forbid` pattern.
- `*`: Matches a single path segment.
- `**`: Matches multiple path segments, including none.

## Example

Given the following project structure:

```
example/
├── alpha/
│   ├── experimental/
│   │   └── widget.go
│   ├── internal/
│   │   ├── exception/
│   │   │   └── test.go
│   │   └── excluded/
│   │       └── allowed.go
│   └── main.go
```

The linter will:

1. Analyze files matching the `include` patterns.
2. Exclude files matching the `exclude` patterns.
3. Report violations for imports matching the `forbid` patterns unless they match the `except` patterns.

## Output

On the happy path the linter will output
```
✔ go-arch-lint: no forbidden imports found.
```
and exit with code 0.

On the unhappy path the linter will output

```
go-arch-lint: [<rule name>] "path/to/file.go" imports "forbidden/package"
```

## Development

### Prerequisites

- Go 1.20 or later

### Running Tests

TODO

## Contributing

Contributions are welcome! Please open an issue or submit a pull request with your changes.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.