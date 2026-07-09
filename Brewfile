# Development toolchain for traveller-worldgen.
# Install everything with `task deps` (or `brew bundle`).

brew "go" # Go compiler and toolchain (version floor is in go.mod)
brew "go-task" # the `task` runner that drives this repo's Taskfile
brew "gofumpt" # stricter gofmt; formatting is CI-enforced (task fmt)
brew "golangci-lint" # meta-linter (task lint; config in .golangci.yml)
brew "gopls" # Go language server; also the .mcp.json MCP server
