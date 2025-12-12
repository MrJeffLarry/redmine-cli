[Home](index.md)
----

## Config

The red-cli tool supports both local and global configuration using the config.json file. You can override the global configuration with a local configuration.

For example, if you use one repository for one project and another repository for a different project, you can create a folder called .red in the first repository and inside it, create a file called config.json. This file can contain and override one or more configuration options, as shown below:

```bash
repo1/.red/config.json
```

This file can contain the following JSON code:

```json
{
    "project-id": 12
}
```

The value of project-id will then override the corresponding value in the global configuration.


### Complete config list options

#### Single Instance Mode (Legacy)

```json
{
    "server": "https://redmine.example.com",
    "api-key": "randomkeyfromredmine",
    "user-id": 1,
    "project-id": 23,
    "editor": "",
    "pager": "",
}
```

#### Multi-Instance Mode

For managing multiple Redmine instances, use the following structure:

```json
{
    "instances": {
        "1": {
            "server": "https://redmine1.example.com",
            "api-key": "apikey1",
            "user-id": 1,
            "project-id": 23,
            "editor": "",
            "pager": ""
        },
        "2": {
            "server": "https://redmine2.example.com",
            "api-key": "apikey2",
            "user-id": 5,
            "project-id": 42,
            "editor": "",
            "pager": ""
        }
    },
    "default-instance": "1"
}
```

### Multi-Instance Support

Red-cli supports managing multiple Redmine instances simultaneously using the `--rid` flag. This allows you to work with different Redmine servers without having to logout and login repeatedly.

#### Login to Multiple Instances

To login to a specific instance, use the `--rid` flag:

```bash
# Login to instance 1
red-cli auth login --rid 1

# Login to instance 2
red-cli auth login --rid 2
```

#### Using Commands with Different Instances

Once logged in to multiple instances, you can use any command with the `--rid` flag to target a specific instance:

```bash
# View issues from instance 1
red-cli issue list --rid 1

# Create an issue in instance 2
red-cli issue create --rid 2

# List projects in instance 1
red-cli project list --rid 1
```

#### Default Instance

When you login with `--rid` for the first time, that instance becomes the default. If you don't specify `--rid` in subsequent commands, the default instance will be used.

To use a different instance without specifying `--rid` every time, login to that instance again to make it the default.

#### Logout from Specific Instance

To logout from a specific instance:

```bash
# Logout from instance 2
red-cli auth logout --rid 2
```

This will remove only the credentials for instance 2, leaving other instances intact.

#### Backward Compatibility

Red-cli maintains full backward compatibility. If you don't use the `--rid` flag, it works exactly as before, storing a single set of credentials in the traditional format. This means:

- Existing users don't need to change anything
- Single-instance usage works without any changes
- You can start using multi-instance support whenever needed by simply adding the `--rid` flag

### Editor

By default, red-cli looks for the environment variables `EDITOR`, `GIT_EDITOR`, or `VISUAL` to determine which text editor to use. If none of these environment variables are set, red-cli will use Notepad on Windows and Nano on other operating systems.

However, you can override the default text editor by using the following command:

```bash
red-cli config editor set vim
```

This command sets the text editor to Vim. You can replace "vim" with the name of any other text editor you prefer. Once you have set the text editor, red-cli will use that editor whenever you need to edit text, such as when creating or updating an issue.

### Pager

By default, red-cli prints output directly to the screen. However, you can change this behavior by modifying the configuration with the following command:

```bash
red-cli config pager set less
```

This command sets the pager to "less". You can replace "less" with the name of any other pager you prefer. Once you have set the pager, red-cli will use that pager whenever it needs to display a large amount of output. This can make it easier to read and navigate through long lists or issues.