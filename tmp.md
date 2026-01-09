version on config?
version: 1.0

overrides global config:
repo1/
  .red/
    config.json
  other-files...

~/.red/
  config.json


```json
{
    "server": "https://redmine.example.com",
    "api-key": "randomkeyfromredmine",
    "project": "my-project",
    "project-id": 23,
    "user-id": 1,
    "editor": "nano",
    "pager": "less"
}
```

version 2.0

```json
{
    "version": "2.0",
    "servers": [
        {
            "id": 1,
            "name": "Redmine 1",
            "server": "https://redmine1.example.com",
            "api-key": "apikey1",
            "user-id": 1,
            "project-id": 23
        },
        {
            "id": 2,
            "name": "Redmine 2",
            "server": "https://redmine2.example.com",
            "api-key": "apikey2",
            "user-id": 5,
            "project-id": 42
        }
    ],
    "default-server": 1,
    "editor": "nano",
    "pager": "less"
}
```
