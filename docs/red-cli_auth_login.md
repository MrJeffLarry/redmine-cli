# Documentation
[**Home**](../README.md) | [**Index**](index.md) | red-cli auth login

## red-cli auth login

login to Redmine

### Synopsis

Authenticate to Redmine server. Optionally use the `--rid` flag to login to a specific instance when managing multiple Redmine servers.

```
red-cli auth login [flags]
```

### Examples

```bash
# Login to default instance
red-cli auth login

# Login to a specific instance
red-cli auth login --rid prod
red-cli auth login --rid staging
red-cli auth login --rid 2
```

### Options

```
  -h, --help   help for login
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

