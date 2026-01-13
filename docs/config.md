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

### Multi-Instance Support

Red-cli supports managing multiple Redmine instances simultaneously using the `--rid` flag. This allows you to work with different Redmine servers without having to logout and login repeatedly.

#### Login to Multiple Instances

To login to multiple Redmine instances, use the `red-cli auth login` command and you will be promted to enter the name of the instance, this name will be used as the Redmine Instance ID (RID) for future commands.

#### Use Multiple Instances

To use multiple Redmine instances, you need to specify a unique Redmine instance ID (RID) for each instance when using the cli. This ID can be the id of the instance or name you choose when authenticated.

```bash
# Use to instance 1
red-cli issue list --rid 1

# Login to instance 2
red-cli issue list --rid server2
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

#### Logout from Specific Instance

To logout from a specific instance:

```bash
# will be promted to choose server to logout
red-cli auth logout
```

This will remove only the credentials for choosed server, leaving other instances intact.

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