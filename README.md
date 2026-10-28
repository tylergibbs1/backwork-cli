# Backwork CLI

[![CI](https://github.com/tylergibbs1/backwork-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/tylergibbs1/backwork-cli/actions/workflows/ci.yml)

Official command line interface for the [Backwork API](https://backworkhealth.com): Medicare coverage policies, medical code intelligence, prior authorization checks, claim validation, compliance review, and drug formulary evidence. It's built for shell workflows, operational scripts, and quick ad hoc lookups.

## Quick start

```bash
go install github.com/tylergibbs1/backwork-cli/cmd/backwork@latest
export BACKWORK_API_KEY=bwk_live_YOUR_API_KEY

backwork health
backwork check 76942 --include rvu,policies
backwork prior-auth 76942 --state TX --diagnosis M54.5
backwork policies list --query "ultrasound guidance" --type LCD
```

Get an API key from the [Backwork dashboard](https://backworkhealth.com/dashboard). Full API docs are at [backworkhealth.com/docs](https://backworkhealth.com/docs).

## Features

- **Code lookup**: `check` and `batch` return code details, RVUs, and the policies that list a code.
- **Policies and coverage**: search, fetch, and compare LCDs, Articles, and NCDs across MAC jurisdictions, track policy changes, and evaluate a case against a policy.
- **Prior authorization and claims**: check prior auth from Backwork's policies, run background research on a payer's public site, and validate claim denial risk.
- **Spending, compliance, and drugs**: spending data by code and year, unreviewed policy changes with acknowledgement, and commercial PBM formulary evidence.
- **Webhooks**: list, create, update, delete, and test webhook endpoints.
- **Script-friendly output**: every command prints `table`, `json`, or `yaml`.

## Installation

Install the latest release with Go:

```bash
go install github.com/tylergibbs1/backwork-cli/cmd/backwork@latest
```

Make sure your Go bin directory is on `PATH`, commonly `$(go env GOPATH)/bin`.

Tagged releases publish pre-built binaries and checksums on [GitHub Releases](https://github.com/tylergibbs1/backwork-cli/releases).

Or build from source:

```bash
git clone https://github.com/tylergibbs1/backwork-cli.git
cd backwork-cli
go build -o backwork ./cmd/backwork
sudo mv backwork /usr/local/bin/
```

## Configuration

Configuration is resolved in this order:

1. Command line flags
2. Environment variables prefixed with `BACKWORK_`
3. `~/.backwork.yaml`

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

Run `backwork <command> --help` for every flag.

### Code lookup

```bash
backwork check 76942
backwork check 76942 --include rvu,policies --jurisdiction JM
backwork batch 76942 99213 --include rvu,policies --output json
```

### Policies and coverage

```bash
backwork policies list --query "ultrasound guidance" --type LCD
backwork policies get L33831 --include criteria,codes
backwork policies compare 76942 --jurisdictions JM,JH,JK
backwork policies changes --since 2026-01-01T00:00:00Z
backwork coverage search "diabetes" --section indications --limit 10
backwork evaluate L33831 --procedure 76942 --diagnosis M54.5
backwork jurisdictions
```

### Prior authorization and claims

```bash
backwork prior-auth 76942 --state TX --diagnosis M54.5 --payer medicare
backwork prior-auth research 27447 --payer "UnitedHealthcare" --state TX --sync
backwork prior-auth research-status <research-id>
backwork claims validate 99213 --diagnosis E11.9 --payer Medicare --state TX
```

### Spending, compliance, and drugs

```bash
backwork spending T1019 T1020 --year 2023
backwork compliance unreviewed --limit 10
backwork compliance stats
backwork compliance ack 123 --notes "Reviewed"
backwork compliance bulk-ack 123 124 --notes "Reviewed"
backwork drugs formulary ozempic --payer all --limit 5
```

### Webhooks

```bash
backwork webhooks list
backwork webhooks create --url https://example.com/webhooks/backwork --events compliance.acknowledged
backwork webhooks test 123
backwork webhooks delete 123
```

### Output formats

All commands support the global output flag. Supported formats are `table`, `json`, and `yaml`.

```bash
backwork check 76942 --output json
backwork policies list --query diabetes --output yaml
```

### Global flags

```text
--api-key string    Backwork API key, or set BACKWORK_API_KEY
--base-url string   API base URL
--config string     Config file path
-o, --output        Output format: table, json, yaml
```

### Shell completion

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

The `Makefile` wraps the same steps (`make build`, `make test`) and `make build-all` cross-compiles macOS, Linux, and Windows binaries into `dist/`. CI runs vet and tests on every pull request.

### Release

1. Push a tag such as `v1.1.0`.
2. The release workflow runs tests, cross-compiles macOS, Linux, and Windows binaries, writes SHA-256 checksums, and creates a GitHub Release.
3. `go install github.com/tylergibbs1/backwork-cli/cmd/backwork@latest` resolves through the public Go module proxy after the tag is indexed.

## Support

- Documentation: https://backworkhealth.com/docs
- Issues: https://github.com/tylergibbs1/backwork-cli/issues
- Email: support@backworkhealth.com

## License

MIT
