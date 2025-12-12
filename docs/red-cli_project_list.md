# Documentation
[**Home**](../README.md) | [**Index**](index.md) | red-cli project list

## red-cli project list

List projects

### Synopsis

List all projects

```
red-cli project list [flags]
```

### Options

```
      --asc             Ascend order
  -h, --help            help for list
      --json            Output in JSON format
  -l, --limit int       Limit number of objects per page (default 25)
  -o, --offset int      skip this number of objects
  -p, --page int        List 25 objects per page (uses limit and offset)
  -q, --query string    Query for projects with name
  -s, --sort string     Sort field
      --status_id int   Filter on status ID
      --target_id int   Filter on target version ID
```

### Options inherited from parent commands

```
      --all     Ignore project-id
  -d, --debug   Show debug info and raw response
```

### SEE ALSO

* [red-cli project](./red-cli_project.md)	 - project
* [red-cli project list all](./red-cli_project_list_all.md)	 - List all projects

