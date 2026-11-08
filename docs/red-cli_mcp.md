# Documentation
[**Home**](../README.md) | [**Index**](index.md) | red-cli mcp

## red-cli mcp

Start an MCP server

### Synopsis

Run red-cli as a Model Context Protocol (MCP) server over stdio.

This command exposes Redmine operations as MCP tools so AI assistants can
retrieve and update Redmine data through your existing red-cli authentication.

Available tools:
	- list_issues
	- get_issue
	- create_issue
	- update_issue
	- add_note
	- list_projects
	- get_current_user

The server communicates using JSON-RPC on stdin/stdout and is intended to be
started by an MCP client configuration (for example Claude Desktop).

```
red-cli mcp [flags]
```

### Examples

```
# Run MCP server directly (for local testing)
red-cli mcp

# Example MCP client configuration
{
  "mcpServers": {
    "redmine": {
      "command": "red-cli",
      "args": ["mcp"]
    }
  }
}
```

### Options

```
  -h, --help   help for mcp
```

### Options inherited from parent commands

```
      --all          Ignore project-id
  -d, --debug        Show debug info and raw response
      --rid string   Redmine server name or ID
```

### SEE ALSO

* [red-cli](./red-cli.md)	 - Redmine CLI

