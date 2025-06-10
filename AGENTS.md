# AGENTS instructions

To ensure this repository uses the correct Go toolchain during CI and for Codex,
install Go 1.24 before running any builds or tests. Run the helper script
`./scripts/setup-go-1.24.sh` which downloads Go 1.24.0 and installs it to
`/usr/local/go`. The script also updates the PATH for the current session.
