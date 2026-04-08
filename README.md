# cloudping-go

A simple CLI to measure real HTTP round-trip latency from your machine to cloud regions.

`cloudping-go` sends direct requests to region-specific endpoints for AWS, GCP, and Azure, then sorts results from fastest to slowest.

## Features

- Ping many cloud regions in parallel
- Filter by provider and/or region code
- Average latency over multiple requests per region
- Colored output table for quick scanning
- Region list command to discover supported codes

## Supported Providers

The embedded dataset currently includes:

- AWS: 34 regions
- Azure: 48 regions
- GCP: 42 regions

Total: 124 regions.

## Requirements

- Go 1.26+
- Network access to cloud endpoints

## Quick Start

Run directly:

```bash
go run . list
go run . ping
```

Build binary:

```bash
go build -o cloudping-go .
./cloudping-go list
./cloudping-go ping --providers aws,gcp --count 5
```

## Usage

```bash
cloudping-go [command] [flags]
```

### Commands

- `list` - show supported providers, region names/codes, and endpoint URLs
- `ping` - run latency tests and print sorted results

### Global Flags

- `-p, --providers` provider filter (comma-separated), example: `aws,gcp`
- `-r, --regions` region code filter (comma-separated), example: `us-east-1,europe-west9`

### Ping Flags

- `-n, --count` number of requests per region used for averaging (default: `3`)

## Examples

```bash
# List all regions
cloudping-go list

# List only Azure regions
cloudping-go list --providers azure

# Ping everything
cloudping-go ping

# Ping only AWS + GCP
cloudping-go ping --providers aws,gcp

# Ping selected regions, averaging 5 requests each
cloudping-go ping --regions us-east-1,europe-west9 -n 5
```

## How It Works

- Region metadata is embedded from `internal/pinger/regions.json`
- A warm-up request is sent first to reduce DNS/TLS setup noise
- Each region is measured `n` times and averaged
- Requests use a 3-second HTTP timeout
- Pings run concurrently (up to 5 workers)

## Tech Stack

- [Cobra](https://github.com/spf13/cobra) for CLI structure
- [progressbar](https://github.com/schollz/progressbar) for progress display
- [table](https://github.com/rodaine/table) and [fatih/color](https://github.com/fatih/color) for terminal output
- Go standard library packages: `net/http`, `embed`, `time`, `encoding/json`

## References

- https://gitlab.com/leonhard-llc/cloudping.info/
- https://pkg.go.dev/net/http
- https://pkg.go.dev/embed
- https://pkg.go.dev/time
- https://pkg.go.dev/encoding/json
- https://github.com/GoogleCloudPlatform/gcping
- https://global.gcping.com/api/endpoints
- https://github.com/schollz/progressbar
- https://github.com/spf13/cobra
- https://www.azurespeed.com/
- https://docs.aws.amazon.com/general/latest/gr/ddb.html

## License

MIT.
