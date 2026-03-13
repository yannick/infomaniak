# infomaniak

CLI tool to manage domains and nameservers via the [Infomaniak API](https://developer.infomaniak.com).

## Requirements

- Go 1.23+
- An Infomaniak API token with the `domain` scope ([create one here](https://manager.infomaniak.com/v3/ng/accounts/token/list))

## Installation

### From source

```sh
git clone https://github.com/yannick/infomaniak.git
cd infomaniak
make deps
make build
# binary is at ./bin/infomaniak
```

### Install to GOPATH

```sh
make install
```

## Configuration

infomaniak reads configuration from three sources in this order of precedence:

1. **CLI flags** (`--token`, `--account-id`)
2. **Environment variables** (`INFOMANIAK_TOKEN`, `INFOMANIAK_ACCOUNT_ID`)
3. **Config file** (`~/.infomaniak.yaml` or `./.infomaniak.yaml`)

### Config file

Create `~/.infomaniak.yaml`:

```yaml
token: "your-api-token-here"
account_id: "12345"
```

Or specify a custom path:

```sh
infomaniak --config /path/to/config.yaml domains list
```

### Environment variables

```sh
export INFOMANIAK_TOKEN="your-api-token-here"
export INFOMANIAK_ACCOUNT_ID="12345"
```

## Usage

### List domains

```sh
infomaniak domains list
```

```
NAME                   TLD      EXPIRES
foo.ch               ch       2026-04-30
bar.ch              ch       2027-01-29
me1337.net             net      2027-01-30
```

### Show domain details

```sh
infomaniak domains show example.ch
```

```
Name:            example.ch
TLD:             ch
Premium:         false
Created:         2024-01-15
Expires:         2025-01-15
DNS Anycast:     false
DNSSEC:          true
Domain Privacy:  false
```

### Update nameservers

```sh
infomaniak domains update-ns example.ch --nameservers ns1.example.ch,ns2.example.ch
```

```
Nameservers for example.ch updated successfully.
```

Use `--verify` to check nameserver availability before applying:

```sh
infomaniak domains update-ns example.ch --nameservers ns1.example.ch,ns2.example.ch --verify
```

## Output formats

All commands support three output modes:

### Default (table)

Human-readable tabular output shown in the examples above.

### JSON (`--json`)

Full API response as pretty-printed JSON, useful for scripting with `jq`:

```sh
infomaniak domains list --json
```

```json
[
  {
    "id": 0,
    "name": "foo.ch",
    "tld": "ch",
    "is_premium": false,
    "created_at": 1460380100,
    "expires_at": 1777500000,
    "options": {
      "dns_anycast": false,
      "renewal_warranty": false,
      "domain_privacy": false,
      "dnssec": false
    },
    "contacts": { ... }
  }
]
```

Extract just domain names with `jq`:

```sh
infomaniak domains list --json | jq -r '.[].name'
```

### Simple (`--simple`)

Plain domain names, one per line. Useful for piping to other tools:

```sh
infomaniak domains list --simple
```

```
foo.ch
bar.ch
```

Example — check DNS for all domains:

```sh
infomaniak domains list --simple | xargs -I{} dig +short {} NS
```

`--json` and `--simple` are mutually exclusive.

## Development

```sh
make deps       # download dependencies
make build      # build for current platform
make test       # run tests with race detector
make lint       # run go vet, gofmt, and golangci-lint
make clean      # remove build artifacts
```

## CI

GitHub Actions runs on every push and PR to `main`:

- **test**: `go vet`, `gofmt`, `go test -race`
- **lint**: `golangci-lint`
- **build**: cross-compiles for linux, darwin, and windows (amd64 + arm64)

## License

MIT
