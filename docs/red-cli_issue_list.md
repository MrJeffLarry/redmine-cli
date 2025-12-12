# Documentation
[**Home**](../README.md) | [**Index**](index.md) | red-cli issue list

## red-cli issue list

List issues

### Synopsis

List all issues

```
red-cli issue list [flags]
```

### Options

```
      --asc             Ascend order
  -h, --help            help for list
      --issue-urls      Show issue urls only
      --json            Output in JSON format
  -l, --limit int       Limit number of objects per page (default 25)
  -o, --offset int      skip this number of objects
  -p, --page int        List 25 objects per page (uses limit and offset)
      --project         Display project column
  -q, --query string    Query for issues with subject
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

* [red-cli issue](./red-cli_issue.md)	 - issue
* [red-cli issue list all](./red-cli_issue_list_all.md)	 - List all issues
* [red-cli issue list me](./red-cli_issue_list_me.md)	 - List all my issues

