# Verity CLI

Official command line interface for the [Verity API](https://verity.backworkai.com): Medicare coverage policies, medical code intelligence, prior authorization checks, claim validation, compliance review, and drug formulary evidence.

The CLI is designed for shell workflows, operational scripts, and quick ad hoc lookups.

## Installation

Install with Go:

```bash
go install github.com/backworkai/verity-cli/cmd/verity@v1.0.0
```

Make sure your Go bin directory is on `PATH`, commonly `$(go env GOPATH)/bin`.

Build from source:

```bash
git clone https://github.com/backworkai/verity-cli.git
cd verity-cli
go mod tidy
go build -o verity .
sudo mv verity /usr/local/bin/
```

Tagged releases publish pre-built binaries and checksums on [GitHub Releases](https://github.com/backworkai/verity-cli/releases).

## Quick Start

```bash
export VERITY_API_KEY=vrt_live_YOUR_API_KEY

verity health
verity check 76942 --include rvu,policies
verity prior-auth 76942 --state TX --diagnosis M54.5
verity policies list --query "ultrasound guidance" --type LCD
```

Get an API key from the [Verity dashboard](https://verity.backworkai.com/dashboard).

## Configuration

Configuration is resolved in this order:

1. Command line flags
2. Environment variables prefixed with `VERITY_`
3. `~/.verity.yaml`

```yaml
api_key: vrt_live_YOUR_API_KEY
base_url: https://verity.backworkai.com/api/v1
output: table
```

```bash
export VERITY_API_KEY=vrt_live_YOUR_API_KEY
export VERITY_BASE_URL=https://verity.backworkai.com/api/v1
export VERITY_OUTPUT=json
```

## Commands

### Code Lookup

```bash
verity check 76942
verity check 76942 --include rvu,policies --jurisdiction JM
verity batch 76942 99213 --include rvu,policies --output json
```

### Policies and Coverage

```bash
verity policies list --query "ultrasound guidance" --type LCD
verity policies get L33831 --include criteria,codes
verity policies compare 76942 --jurisdictions JM,JH,JK
verity policies changes --since 2026-01-01T00:00:00Z
verity coverage search "diabetes" --section indications --limit 10
verity evaluate L33831 --procedure 76942 --diagnosis M54.5
```

### Prior Authorization and Claims

```bash
verity prior-auth 76942 --state TX --diagnosis M54.5 --payer medicare
verity prior-auth research 27447 --payer "UnitedHealthcare" --state TX --sync
verity claims validate 99213 --diagnosis E11.9 --payer Medicare --state TX
```

### Spending, Compliance, and Drugs

```bash
verity spending T1019 T1020 --year 2023
verity compliance unreviewed --limit 10
verity compliance stats
verity compliance ack 123 --notes "Reviewed"
verity drugs formulary ozempic --payer all --limit 5
```

### Webhooks

```bash
verity webhooks list
verity webhooks create --url https://example.com/webhooks/verity --events policy.updated
verity webhooks test 123
```

## Output Formats

All commands support the global output flag:

```bash
verity check 76942 --output json
verity policies list --query diabetes --output yaml
```

Supported formats are `table`, `json`, and `yaml`.

## Global Flags

```text
--api-key string    Verity API key, or set VERITY_API_KEY
--base-url string   API base URL
--config string     Config file path
-o, --output        Output format: table, json, yaml
```

## Shell Completion

```bash
verity completion bash
verity completion zsh
verity completion fish
verity completion powershell
```

## Development

```bash
go mod tidy
go test ./...
go vet ./...
go build -o verity .
```

## Release

1. Push a tag such as `v1.0.0`.
2. The release workflow runs tests, cross-compiles macOS, Linux, and Windows binaries, writes SHA-256 checksums, and creates a GitHub Release.
3. `go install github.com/backworkai/verity-cli/cmd/verity@latest` resolves through the public Go module proxy after the tag is indexed.

## Support

- Documentation: https://verity.backworkai.com/docs
- Issues: https://github.com/backworkai/verity-cli/issues
- Email: support@verity.backworkai.com

## License

MIT
