# Backwork CLI

Official command line interface for the [Backwork API](https://backworkhealth.com): Medicare coverage policies, medical code intelligence, prior authorization checks, claim validation, compliance review, and drug formulary evidence.

The CLI is designed for shell workflows, operational scripts, and quick ad hoc lookups.

This tool was called `verity` until recently. See [Migrating from the Verity CLI](#migrating-from-the-verity-cli) if you already have it installed.

## Installation

Build from source. This works today:

```bash
git clone https://github.com/tylergibbs1/verity-cli.git
cd verity-cli
go mod tidy
go build -o backwork ./cmd/backwork
sudo mv backwork /usr/local/bin/
```

`go install` does **not** work yet:

```bash
# Fails right now: the module path was renamed but the GitHub repository
# is still published as tylergibbs1/verity-cli, so the proxy cannot fetch it.
go install github.com/tylergibbs1/backwork-cli/cmd/backwork@latest
```

That command starts working once the repository itself is renamed to `backwork-cli` and a tag is pushed. Until then, use the source build above.

Make sure your Go bin directory is on `PATH`, commonly `$(go env GOPATH)/bin`.

Tagged releases publish pre-built binaries and checksums on [GitHub Releases](https://github.com/tylergibbs1/verity-cli/releases).

## Quick Start

```bash
export BACKWORK_API_KEY=bwk_live_YOUR_API_KEY

backwork health
backwork check 76942 --include rvu,policies
backwork prior-auth 76942 --state TX --diagnosis M54.5
backwork policies list --query "ultrasound guidance" --type LCD
```

Get an API key from the [Backwork dashboard](https://backworkhealth.com/dashboard).

## Migrating from the Verity CLI

Nothing you have set up stops working. The CLI reads the new names first and falls back to the old ones.

| | Old | New |
|---|---|---|
| Binary | `verity` | `backwork` |
| Config file | `~/.verity.yaml` | `~/.backwork.yaml` |
| Env vars | `VERITY_API_KEY`, `VERITY_BASE_URL`, `VERITY_OUTPUT` | `BACKWORK_API_KEY`, `BACKWORK_BASE_URL`, `BACKWORK_OUTPUT` |
| API key prefix | `vrt_live_...` | `bwk_live_...` |

- **Config file.** If `~/.backwork.yaml` does not exist, the CLI reads `~/.verity.yaml` instead. Rename the file when convenient.
- **Environment variables.** `BACKWORK_*` is read first. If it is unset, the matching `VERITY_*` variable is used.
- **API keys.** Your existing `vrt_live_...` keys remain valid. The API accepts both prefixes, so there is no need to rotate them. Newly issued keys use `bwk_live_...`.
- **Binary name.** The old `verity` binary is not removed by the source build above. Delete it yourself once you have switched: `sudo rm /usr/local/bin/verity`.

## Configuration

Configuration is resolved in this order:

1. Command line flags
2. Environment variables prefixed with `BACKWORK_` (falling back to `VERITY_`)
3. `~/.backwork.yaml` (falling back to `~/.verity.yaml`)

```yaml
api_key: bwk_live_YOUR_API_KEY
base_url: https://backworkhealth.com/api/v1
output: table
```

```bash
export BACKWORK_API_KEY=bwk_live_YOUR_API_KEY
export BACKWORK_BASE_URL=https://backworkhealth.com/api/v1
export BACKWORK_OUTPUT=json
```

## Commands

### Code Lookup

```bash
backwork check 76942
backwork check 76942 --include rvu,policies --jurisdiction JM
backwork batch 76942 99213 --include rvu,policies --output json
```

### Policies and Coverage

```bash
backwork policies list --query "ultrasound guidance" --type LCD
backwork policies get L33831 --include criteria,codes
backwork policies compare 76942 --jurisdictions JM,JH,JK
backwork policies changes --since 2026-01-01T00:00:00Z
backwork coverage search "diabetes" --section indications --limit 10
backwork evaluate L33831 --procedure 76942 --diagnosis M54.5
```

### Prior Authorization and Claims

```bash
backwork prior-auth 76942 --state TX --diagnosis M54.5 --payer medicare
backwork prior-auth research 27447 --payer "UnitedHealthcare" --state TX --sync
backwork claims validate 99213 --diagnosis E11.9 --payer Medicare --state TX
```

### Spending, Compliance, and Drugs

```bash
backwork spending T1019 T1020 --year 2023
backwork compliance unreviewed --limit 10
backwork compliance stats
backwork compliance ack 123 --notes "Reviewed"
backwork drugs formulary ozempic --payer all --limit 5
```

### Webhooks

```bash
backwork webhooks list
backwork webhooks create --url https://example.com/webhooks/backwork --events policy.updated
backwork webhooks test 123
```

## Output Formats

All commands support the global output flag:

```bash
backwork check 76942 --output json
backwork policies list --query diabetes --output yaml
```

Supported formats are `table`, `json`, and `yaml`.

## Global Flags

```text
--api-key string    Backwork API key, or set BACKWORK_API_KEY
--base-url string   API base URL
--config string     Config file path
-o, --output        Output format: table, json, yaml
```

## Shell Completion

```bash
backwork completion bash
backwork completion zsh
backwork completion fish
backwork completion powershell
```

## Development

```bash
go mod tidy
go test ./...
go vet ./...
go build -o backwork .
```

## Release

1. Push a tag such as `v1.0.0`.
2. The release workflow runs tests, cross-compiles macOS, Linux, and Windows binaries, writes SHA-256 checksums, and creates a GitHub Release.
3. `go install github.com/tylergibbs1/backwork-cli/cmd/backwork@latest` resolves through the public Go module proxy once the GitHub repository has been renamed to `backwork-cli` and the tag is indexed. It fails before that rename.

## Support

- Documentation: https://backworkhealth.com/docs
- Issues: https://github.com/tylergibbs1/verity-cli/issues
- Email: support@backworkhealth.com

## License

MIT
