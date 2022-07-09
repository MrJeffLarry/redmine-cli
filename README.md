# redmine-cli

### Usage for override global config

You can override the global config with a local config, for example if you use one repo for one project and other for other project you can create a folder inside current working directory **.red** inside that create a file called **config.json** this can then contain and override one or more options below

```bash
.red/config.json
```

contains 

```json
{
    "project-id": 12
}
```

this will then override the project

## Command Tree

red
- issue
- - create
- - list
- - - all
- - - me
- project
- - list
- - - all
- - - me
- user
- - me
- config
- - global
- - - project
- - - - set
- - local
- - - project
- - - - set

## Config

**Complete config list options**

```json
{
    "server": "https://redmine.example.com",
    "api-key": "randomkeyfromredmine",
    "user-id": 1,
    "project-id": 23,
}
```