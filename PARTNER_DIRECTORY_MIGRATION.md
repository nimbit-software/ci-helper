# Partner Directory Migration Summary

This document describes the Partner Directory functionality that was ported from the standalone CLI tool into the FlashPipe project.

## Overview

Partner Directory management capabilities have been integrated into FlashPipe, allowing you to:
- **Snapshot** (download) Partner Directory parameters from SAP CPI to local files
- **Deploy** (upload) Partner Directory parameters from local files to SAP CPI
- Support for both **String** and **Binary** parameters
- **Batch operations** for efficient processing of large parameter sets
- **Full sync mode** to keep local and remote in sync

## Architecture

The integration follows FlashPipe's existing patterns and consists of three main layers:

### 1. HTTP Client Layer (`internal/httpclnt/`)

**New Files:**
- `batch.go` - OData $batch request support for efficient bulk operations

**Key Features:**
- `BatchRequest` - Builds and executes OData multipart/mixed batch requests
- `BatchOperation` - Represents individual operations within a batch
- `BatchResponse` - Parses batch responses
- Helper functions for creating batch operations for Partner Directory parameters
- Default batch size of 90 operations per request

### 2. API Layer (`internal/api/`)

**New File:**
- `partnerdirectory.go` - Partner Directory API implementation

**Key Components:**
- `PartnerDirectory` - Main API client using FlashPipe's `HTTPExecuter`
- `StringParameter` - String parameter model
- `BinaryParameter` - Binary parameter model
- `BatchResult` - Results tracking for batch operations

**Methods:**
- `GetStringParameters()` - Fetch all string parameters
- `GetBinaryParameters()` - Fetch all binary parameters
- `GetStringParameter()` - Fetch single string parameter
- `GetBinaryParameter()` - Fetch single binary parameter
- `CreateStringParameter()` - Create new string parameter
- `UpdateStringParameter()` - Update existing string parameter
- `DeleteStringParameter()` - Delete string parameter
- `CreateBinaryParameter()` - Create new binary parameter
- `UpdateBinaryParameter()` - Update existing binary parameter
- `DeleteBinaryParameter()` - Delete binary parameter
- `BatchSyncStringParameters()` - Batch create/update string parameters
- `BatchSyncBinaryParameters()` - Batch create/update binary parameters
- `BatchDeleteStringParameters()` - Batch delete string parameters
- `BatchDeleteBinaryParameters()` - Batch delete binary parameters

### 3. Repository Layer (`internal/repo/`)

**New File:**
- `partnerdirectory.go` - File system operations for Partner Directory

**Key Components:**
- `PartnerDirectory` - Repository for managing local Partner Directory files

**Methods:**
- `GetLocalPIDs()` - Get all locally managed Partner IDs
- `WriteStringParameters()` - Write string parameters to properties files
- `WriteBinaryParameters()` - Write binary parameters to individual files
- `ReadStringParameters()` - Read string parameters from properties files
- `ReadBinaryParameters()` - Read binary parameters from files

**File Structure:**
```
partner-directory/
  {PID}/
    String.properties          # String parameters as key=value pairs
    Binary/
      {ParamId}.{ext}         # Binary parameter files
      _metadata.json          # Content type metadata
```

### 4. Command Layer (`internal/cmd/`)

**New Files:**
- `pd_snapshot.go` - Snapshot command implementation
- `pd_deploy.go` - Deploy command implementation

**Commands:**
- `flashpipe pd-snapshot` - Download parameters from SAP CPI
- `flashpipe pd-deploy` - Upload parameters to SAP CPI

## Usage

### Snapshot (Download) Parameters

Download all Partner Directory parameters from SAP CPI to local files:

```bash
# Using OAuth with environment variables
export FLASHPIPE_TMN_HOST="your-tenant.hana.ondemand.com"
export FLASHPIPE_OAUTH_HOST="your-tenant.authentication.eu10.hana.ondemand.com"
export FLASHPIPE_OAUTH_CLIENTID="your-client-id"
export FLASHPIPE_OAUTH_CLIENTSECRET="your-client-secret"

flashpipe pd-snapshot

# With explicit parameters
flashpipe pd-snapshot \
  --tmn-host "your-tenant.hana.ondemand.com" \
  --oauth-host "your-tenant.authentication.eu10.hana.ondemand.com" \
  --oauth-clientid "your-client-id" \
  --oauth-clientsecret "your-client-secret" \
  --resources-path "./partner-directory"

# Snapshot only specific PIDs
flashpipe pd-snapshot --pids "SAP_SYSTEM_001,CUSTOMER_API"

# Add-only mode (don't overwrite existing values)
flashpipe pd-snapshot --replace=false
```

### Deploy (Upload) Parameters

Upload Partner Directory parameters from local files to SAP CPI:

```bash
# Basic deploy
flashpipe pd-deploy

# Deploy with custom path
flashpipe pd-deploy --resources-path "./partner-directory"

# Deploy only specific PIDs
flashpipe pd-deploy --pids "SAP_SYSTEM_001,CUSTOMER_API"

# Add-only mode (don't update existing parameters)
flashpipe pd-deploy --replace=false

# Full sync (delete remote parameters not in local)
flashpipe pd-deploy --full-sync

# Dry run (see what would change)
flashpipe pd-deploy --dry-run
```

## Features

### Authentication
- **OAuth 2.0** client credentials flow (recommended)
- **Basic Authentication** (legacy support)
- Inherits authentication configuration from FlashPipe's global flags

### Modes

**Replace Mode (default):**
- `pd-snapshot`: Overwrites existing local files
- `pd-deploy`: Updates existing remote parameters

**Add-Only Mode (`--replace=false`):**
- `pd-snapshot`: Only adds new parameters, preserves existing local values
- `pd-deploy`: Only creates new parameters, skips existing ones

**Full Sync Mode (`--full-sync`, deploy only):**
- Deletes remote parameters not present in local files
- Only affects PIDs that have local directories
- Parameters in other PIDs are not touched
- Use with caution!

**Dry Run Mode (`--dry-run`, deploy only):**
- Shows what would be changed without making any changes
- Useful for testing and validation

### Filtering
- Use `--pids` to filter operations to specific Partner IDs
- Accepts comma-separated list: `--pids "PID1,PID2,PID3"`

### Binary Parameter Support
- Automatic content type detection from file extensions
- Metadata file (`_metadata.json`) stores content types
- Supported content types: xml, xsl, xsd, json, txt, zip, gz, zlib, crt
- Base64 encoding/decoding handled automatically

### Batch Processing
- Efficient OData $batch requests for bulk operations
- Default batch size: 90 operations per request
- Reduces API calls and improves performance
- Automatic retry and error handling

## Logging

Uses FlashPipe's `zerolog` for structured logging:
- `--debug` flag enables detailed logging
- Info-level logs for normal operations
- Debug-level logs for detailed operation tracking
- Warning-level logs for non-fatal errors

## Error Handling

- Continues processing on individual parameter errors
- Collects all errors and reports them at the end
- Returns non-zero exit code on errors
- Detailed error messages with PID/ID context

## Integration with FlashPipe

The Partner Directory functionality is fully integrated with FlashPipe:

1. **Uses FlashPipe's HTTP client** (`httpclnt.HTTPExecuter`)
2. **Follows FlashPipe's command structure** (cobra commands)
3. **Reuses authentication** (OAuth/Basic Auth from global flags)
4. **Uses FlashPipe's logging** (zerolog)
5. **Follows FlashPipe's patterns** (API, Repo, Command layers)
6. **Configuration file support** (via Viper)
7. **Analytics tracking** (command usage analytics)

## Configuration File Support

Partner Directory commands support FlashPipe's configuration file (`~/.flashpipe.yaml`):

```yaml
# Authentication
tmn-host: "your-tenant.hana.ondemand.com"
oauth-host: "your-tenant.authentication.eu10.hana.ondemand.com"
oauth-clientid: "your-client-id"
oauth-clientsecret: "your-client-secret"

# Command-specific (optional)
resources-path: "./partner-directory"
replace: true
```

## Migration from Standalone CLI

If you were using the standalone `ci-helper` CLI tool, the migration is straightforward:

### Command Changes
- `ci-helper pd-snapshot` → `flashpipe pd-snapshot`
- `ci-helper pd-deploy` → `flashpipe pd-deploy`

### Environment Variable Changes
- `CPI_URL` → `FLASHPIPE_TMN_HOST`
- `CPI_TOKEN_URL` → `FLASHPIPE_OAUTH_HOST` (just the host, not full URL)
- `CPI_CLIENT_ID` → `FLASHPIPE_OAUTH_CLIENTID`
- `CPI_CLIENT_SECRET` → `FLASHPIPE_OAUTH_CLIENTSECRET`

### Flag Changes
- `--cpi-url` → `--tmn-host`
- `--token-url` → `--oauth-host`
- `--client-id` → `--oauth-clientid`
- `--client-secret` → `--oauth-clientsecret`

### File Structure
The local file structure remains exactly the same, so existing Partner Directory folders can be used as-is.

## Performance Considerations

- Batch operations significantly reduce API calls (up to 90x improvement)
- Snapshot downloads all parameters in a few requests
- Deploy uses batch operations for creates, updates, and deletes
- Full sync mode queries all remote parameters once for comparison

## Best Practices

1. **Use OAuth** instead of Basic Auth for better security
2. **Test with `--dry-run`** before deploying changes
3. **Use `--pids` filter** for large tenants to process specific PIDs
4. **Enable `--debug`** for troubleshooting
5. **Store credentials in config file** (`~/.flashpipe.yaml`) instead of command line
6. **Use version control** for your partner-directory folder
7. **Be cautious with `--full-sync`** as it deletes remote parameters

## Examples

### Complete Workflow

```bash
# 1. Snapshot current state from production
flashpipe pd-snapshot \
  --resources-path "./partner-directory-prod" \
  --pids "PROD_SYSTEM"

# 2. Make local changes to String.properties or Binary files

# 3. Test deploy with dry run
flashpipe pd-deploy \
  --resources-path "./partner-directory-prod" \
  --pids "PROD_SYSTEM" \
  --dry-run

# 4. Deploy changes
flashpipe pd-deploy \
  --resources-path "./partner-directory-prod" \
  --pids "PROD_SYSTEM"
```

### CI/CD Pipeline Integration

```yaml
# Azure Pipelines example
- task: Bash@3
  displayName: 'Deploy Partner Directory Parameters'
  env:
    FLASHPIPE_TMN_HOST: $(CPI_HOST)
    FLASHPIPE_OAUTH_HOST: $(CPI_OAUTH_HOST)
    FLASHPIPE_OAUTH_CLIENTID: $(CPI_CLIENT_ID)
    FLASHPIPE_OAUTH_CLIENTSECRET: $(CPI_CLIENT_SECRET)
  inputs:
    targetType: 'inline'
    script: |
      flashpipe pd-deploy \
        --resources-path "./partner-directory" \
        --pids "$(PARTNER_IDS)" \
        --debug
```

## Technical Details

### Batch Request Format

The implementation uses OData V2 $batch format with multipart/mixed:
- Batch boundary: `batch_{counter}`
- Changeset boundary: `changeset_{counter}`
- Supports mixing query (GET) and changeset (POST/PUT/DELETE) operations
- Properly handles CSRF tokens for modifying operations

### File Format

**String Parameters (String.properties):**
```properties
PARAM_1=value1
PARAM_2=value2
PARAM_WITH_NEWLINE=line1\nline2
```

**Binary Parameters:**
- Individual files with appropriate extensions
- Base64 encoded in transit
- Metadata stored separately

**Metadata (_metadata.json):**
```json
{
  "certificate.crt": "application/x-x509-ca-cert",
  "config.xml": "application/xml"
}
```

## Troubleshooting

### Common Issues

**Authentication Errors:**
- Verify OAuth host is just the hostname (no `https://` prefix)
- Check that OAuth path is correct (default: `/oauth/token`)
- Ensure client credentials have appropriate permissions

**Batch Errors:**
- Check that batch size doesn't exceed server limits
- Review individual operation errors in the response
- Enable debug logging for detailed batch request/response

**File Encoding Issues:**
- Binary files are base64 encoded automatically
- String parameters escape special characters (\n, \r, \\)
- Ensure file extensions match content types in metadata

## Future Enhancements

Potential areas for improvement:
- Parallel batch execution
- Progress bars for large operations
- Diff view before deploy
- Import/export in different formats (JSON, YAML)
- Validation of parameter values
- Template support for multi-tenant deployments
