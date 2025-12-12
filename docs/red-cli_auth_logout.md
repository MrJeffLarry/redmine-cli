# Documentation
[**Home**](../README.md) | [**Index**](index.md) | red-cli auth logout

## red-cli auth logout

logout from Redmine

### Synopsis

Reset and logout from Redmine server. When using multi-instance mode with `--rid`, only the specified instance is logged out. Without `--rid`, the default instance is logged out.

```
red-cli auth logout [flags]
```

### Examples

```bash
# Logout from default instance
red-cli auth logout

# Logout from a specific instance
red-cli auth logout --rid prod
red-cli auth logout --rid staging
red-cli auth logout --rid 2
```

### Options

```
  -h, --help   help for logout
```

### Options inherited from parent commands

```
      --all          Ignore project-id
  -d, --debug        Show debug info and raw response
      --rid string   Redmine instance ID (for multi-instance support)
```

### SEE ALSO

* [red-cli auth](./red-cli_auth.md)	 - auth to Redmine
* [Multi-Instance Guide](./multi-instance.md) - Guide for managing multiple Redmine instances

