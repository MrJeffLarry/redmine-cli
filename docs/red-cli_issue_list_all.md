# Documentation
[**Home**](../README.md) | [**Index**](index.md) | red-cli issue list all

## red-cli issue list all

List all issues

### Synopsis

List all issues and ignores project ID

```
red-cli issue list all [flags]
```

### Options

```
  -h, --help   help for all
```

### Options inherited from parent commands

```
      --all             Ignore project-id
      --asc             Ascend order
  -d, --debug           Show debug info and raw response
      --issue-urls      Show issue urls only
      --json            Output in JSON format
  -l, --limit int       Limit number of objects per page (default 25)
  -o, --offset int      skip this number of objects
  -p, --page int        List 25 objects per page (uses limit and offset)
      --project         Display project column
  -q, --query string    Query for issues with subject
      --rid string      Redmine instance ID (for multi-instance support)
  -s, --sort string     Sort field
      --status_id int   Filter on status ID
      --target_id int   Filter on target version ID
```

### SEE ALSO

* [red-cli issue list](./red-cli_issue_list.md)	 - List issues

