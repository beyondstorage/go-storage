---
author: Xuanwo <github@xuanwo.io>
status: finished
updated_at: 2020-03-06
deprecates:
  - design/3-support-service-init-via-config-string.md
---

# Proposal: Remove config string

## Background

[storage][] added config string support in proposal [3-support-service-init-via-config-string][], and updated in proposal [9-remove-storager-init][]. These proposals allow users init services like following:

```go
srv, store, err := coreutils.Open("fs:///?work_dir=/path/to/dir")
if err != nil {
    log.Fatalf("service init failed: %v", err)
}
```

With time goes on, and deeper understanding of [storage][]'s configuration, I found `config string` is not a good solution. `config string` does have some benefits: simple string, easy to construct, easy to understand(?).

However, after some experience on demo project [bard][], which is a paste bin service built upon [storage][], I found the `config string` deeply influences end user side configuration. [bard][]'s config looks like following:

```yaml
public_url: http://127.0.0.1:8080
listen: 127.0.0.1:8080

key: xxxxxxx
max_file_size: 104857600

database:
  type: sqlite3
  connection: "/tmp/bard/db"

storage: "fs:///?work_dir=/tmp/bard/data"
```

Every application built upon [storage][] either exposes config string to end user directly or writes a format config function to convert their own config to [storage][] config string. This is unexpected.

Not only that, config string also makes it hard to construct type safe pairs. We need to parse them from string and can't have effective use of existing configuration formats.

## Proposal

So I propose following changes:

- Remove the idea of config string
- Refactor `Open(cfg string)` to `Open(t string, opt []*types.Pair)`

Add a config type `Config` to help developer parse pairs:

```go
type Config struct {
    Type    string
    Options map[string]string
}

func (c *Config) Parse() (t string, []*types.Pair, error) {}
```

## Rationale

None

## Compatibility

Following packages will be affected:

- `coreutils`
- `pkg/config`

## Implementation

Most of the work would be done by the author of this proposal.

[storage]: https://github.com/Xuanwo/storage
[3-support-service-init-via-config-string]: https://github.com/Xuanwo/storage/blob/master/docs/design/3-support-service-init-via-config-string.md
[9-remove-storager-init]: https://github.com/Xuanwo/storage/blob/master/docs/design/9-remove-storager-init.md
[bard]: https://github.com/Xuanwo/bard