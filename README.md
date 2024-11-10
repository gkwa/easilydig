# Easilydig - DynamoDB Usage Metrics Manager

Command-line tool for managing usage metrics in DynamoDB.

## Prerequisites

- Go 1.23+
- AWS credentials configured
- DynamoDB table `usage-metrics` in us-east-1

## Installation

```bash
go install github.com/gkwa/easilydig@latest
```

## Usage

### Append Metrics

Add new usage metric from JSON file:

```bash
easilydig append /path/to/data.json
```

Expected JSON format:
```json
{
  "date": "2024-11-10",
  "amount": 61.25,
  "amountUnits": "GB", 
  "total": 400,
  "totalUnits": "GB",
  "overage": 0,
  "overageUnits": "GB",
  "scrapedAt": "2024-11-10T05:56:21.800Z"
}
```

### Fetch Metrics

Get all metrics from DynamoDB:

```bash
easilydig fetch
```

### Verbosity

Increase logging verbosity:

```bash
easilydig -v append data.json
easilydig -vv append data.json
```

### Logging Format

Output logs in JSON format:

```bash
easilydig --log-format=json append data.json
```

## Development

Build from source:

```bash
git clone https://github.com/gkwa/easilydig
cd easilydig
go build
```

Run Tests:

```bash
go test ./... -v
```

Run Linter:

```bash 
golangci-lint run
```

## Quick Commands

```bash
# Build
go build

# Run with local file
./easilydig append testdata/sample.json

# Run tests
go test -v ./...

# Enable debug logging
./easilydig -v append data.json

# Show version
./easilydig version

# Output help
./easilydig --help
```