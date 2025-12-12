# Multi-Instance Support Guide

[**Home**](../README.md) | [**Index**](index.md) | [**Config**](config.md)

## Overview

Red-cli supports managing multiple Redmine instances simultaneously. This feature is particularly useful when you work with multiple Redmine servers (e.g., production, staging, different organizations) and need to switch between them without repeatedly logging in and out.

## Getting Started

### Login to Multiple Instances

To login to a specific Redmine instance, use the `--rid` flag followed by a unique identifier (e.g., "1", "2", "prod", "staging"):

```bash
# Login to instance 1 (e.g., your production server)
red-cli auth login --rid 1
# Follow the prompts to enter server URL and credentials

# Login to instance 2 (e.g., your staging server)
red-cli auth login --rid 2
# Follow the prompts to enter server URL and credentials
```

You can use any identifier that makes sense to you:

```bash
# Using descriptive names
red-cli auth login --rid prod
red-cli auth login --rid staging
red-cli auth login --rid client-a
```

### Using Commands with Different Instances

Once logged in, you can target specific instances using the `--rid` flag:

```bash
# List issues from production
red-cli issue list --rid prod

# Create an issue in staging
red-cli issue create --rid staging

# View projects in client-a instance
red-cli project list --rid client-a

# View user info from instance 1
red-cli user me --rid 1
```

### Default Instance

The first instance you login to becomes your default instance. When you run commands without the `--rid` flag, they will target the default instance:

```bash
# These commands use the default instance
red-cli issue list
red-cli project list
```

### Managing Instances

#### Logout from a Specific Instance

To logout from a specific instance without affecting others:

```bash
red-cli auth logout --rid staging
```

This removes only the credentials for the "staging" instance.

#### Logout from All Instances

To logout from all instances (legacy behavior):

```bash
red-cli auth logout
```

## Configuration Storage

Multi-instance configurations are stored in `~/.red/config.json` with the following structure:

```json
{
    "instances": {
        "prod": {
            "server": "https://redmine-prod.example.com",
            "api-key": "your-prod-api-key",
            "user-id": 1,
            "project-id": 23,
            "editor": "",
            "pager": ""
        },
        "staging": {
            "server": "https://redmine-staging.example.com",
            "api-key": "your-staging-api-key",
            "user-id": 5,
            "project-id": 42,
            "editor": "",
            "pager": ""
        }
    },
    "default-instance": "prod"
}
```

## Use Cases

### Scenario 1: Production and Staging Environments

```bash
# Setup
red-cli auth login --rid prod
# Enter: https://redmine.company.com

red-cli auth login --rid staging
# Enter: https://staging.redmine.company.com

# Work with production
red-cli issue list --rid prod
red-cli issue view 123 --rid prod

# Test in staging
red-cli issue create --rid staging
red-cli issue edit 456 --rid staging
```

### Scenario 2: Multiple Clients

```bash
# Setup for different clients
red-cli auth login --rid client-a
red-cli auth login --rid client-b
red-cli auth login --rid client-c

# Switch between clients easily
red-cli issue list --rid client-a
red-cli project list --rid client-b
red-cli user me --rid client-c
```

### Scenario 3: Personal and Work Redmine

```bash
# Setup
red-cli auth login --rid work
red-cli auth login --rid personal

# Use work Redmine (default)
red-cli issue list

# Use personal Redmine
red-cli issue create --rid personal
```

## Backward Compatibility

Red-cli maintains full backward compatibility with single-instance usage:

- If you never use the `--rid` flag, red-cli works exactly as before
- Existing configurations are automatically recognized
- You can migrate to multi-instance mode at any time by simply using `--rid`
- Single-instance configurations continue to work without modification

## Tips

1. **Use meaningful identifiers**: Choose `--rid` values that make sense to you (e.g., "prod", "staging", "client-name")
2. **Set a sensible default**: Login to your most frequently used instance first to make it the default
3. **Create aliases**: Add shell aliases for frequently used instances:
   ```bash
   alias red-prod='red-cli --rid prod'
   alias red-staging='red-cli --rid staging'
   ```

## Troubleshooting

### How do I know which instance is the default?

Check your config file at `~/.red/config.json` - the `default-instance` field shows your default.

### Can I change the default instance?

Yes, simply login to the desired instance again:
```bash
red-cli auth login --rid new-default
```

### What happens to my existing configuration?

Your existing single-instance configuration continues to work without changes. When you first use `--rid`, your config is automatically migrated to support multiple instances while preserving your existing credentials.

## See Also

- [Configuration Guide](config.md)
- [Authentication Commands](red-cli_auth.md)
- [Login Command](red-cli_auth_login.md)
- [Logout Command](red-cli_auth_logout.md)
