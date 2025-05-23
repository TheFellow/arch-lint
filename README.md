# Go Arch Lint

`go-arch-lint` is a static analysis tool for Go projects that enforces architectural rules by analyzing import paths and package structures.
It helps maintain clean and consistent codebases by preventing unwanted dependencies and enforcing modular boundaries.

## Features

- Define custom rules to forbid specific imports.
- Use glob patterns to include or exclude files for analysis.
- Support for exceptions to allow specific imports in restricted contexts.
- Cross-platform compatibility with normalized file paths.

## Installation

```
go install github.com/TheFellow/go-arch-lint@latest
```

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
- `*`: Matches a single path segment.
- `**`: Matches multiple path segments, including none.
- `{variable}`: Matches a single path segment and captures it as a named variable.

An `except` patterns supports the same special cases as `forbid`, and one more
- `*`: Matches a single path segment.
- `**`: Matches multiple path segments, including none.
- `{variable}`: Matches this path segment when its value matches the one captured in the `forbid` pattern.
- `{!variable}`: Matches this path segment when its value **does not** match the one captured in the `forbid` pattern.


The linter will:

- For all files in scope, which is
   - Files matching the `include` pattern(s)
   - Files not matching the `exclude` pattern(s)
- For each file matching a `forbid` pattern:
   - Report a linting error, unless
   - The import matches an `except` pattern

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

and exit with code 1.

## Development

### Prerequisites

- Go 1.20 or later

### Running Tests

TODO: Write some tests lol

## Contributing

Contributions are welcome! Please open an issue or submit a pull request with your changes.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.