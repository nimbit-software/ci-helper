# Orchestrator YAML Configuration

## Overview

The Flashpipe orchestrator supports loading all configuration settings from a YAML file, making it easy to:
- Version control your deployment settings
- Share configurations across teams
- Use different configs for different environments (dev/qa/prod)
- Simplify CI/CD pipelines with consistent settings

## Quick Start

### Using Orchestrator Config File

```bash
# Load all settings from YAML
flashpipe orchestrator --orchestrator-config ./orchestrator-dev.yml

# Override specific settings via CLI
flashpipe orchestrator --orchestrator-config ./orchestrator-dev.yml \
  --deployment-prefix OVERRIDE
```

### Basic Configuration File

```yaml
# orchestrator-dev.yml
packagesDir: ./packages
deployConfig: ./dev-config.yml
deploymentPrefix: DEV
mode: update-and-deploy
parallelDeployments: 5
deployRetries: 5
deployDelaySeconds: 15
```

---

## Two-Phase Deployment Strategy

The orchestrator now uses a **two-phase approach** with **parallel deployment**:

### Phase 1: Update All Artifacts
1. Update all package metadata
2. Update all artifacts (MANIFEST.MF, parameters.prop, etc.)
3. Collect deployment tasks for Phase 2

### Phase 2: Deploy All Artifacts in Parallel
1. Group artifacts by package
2. Deploy artifacts in parallel (configurable concurrency)
3. Wait for all deployments to complete
4. Report results

**Benefits:**
- ✅ Faster deployments through parallelization
- ✅ All updates complete before any deployment starts
- ✅ Easier to track progress and failures
- ✅ Better error handling and reporting

---

## Configuration Reference

### Complete Configuration Schema

```yaml
# Required Settings
packagesDir: string          # Path to packages directory
deployConfig: string         # Path to deployment config (file/folder/URL)

# Optional: Filtering & Prefixing
deploymentPrefix: string     # Prefix for package/artifact IDs (e.g., "DEV", "PROD")
packageFilter: string        # Comma-separated package names to include
artifactFilter: string       # Comma-separated artifact names to include

# Optional: Config Loading
configPattern: string        # File pattern for folder scanning (default: "*.y*ml")
mergeConfigs: boolean        # Merge multiple configs (default: false)

# Optional: Execution Control
keepTemp: boolean            # Keep temporary files (default: false)
mode: string                 # Operation mode (see below)

# Optional: Deployment Settings
deployRetries: int           # Status check retries (default: 5)
deployDelaySeconds: int      # Delay between checks in seconds (default: 15)
parallelDeployments: int     # Max concurrent deployments (default: 3)
```

### Operation Modes

| Mode | Description | Updates | Deploys |
|------|-------------|---------|---------|
| `update-and-deploy` | Full lifecycle (default) | ✅ | ✅ |
| `update-only` | Only update artifacts | ✅ | ❌ |
| `deploy-only` | Only deploy artifacts | ❌ | ✅ |

---

## Deployment Settings Explained

### `parallelDeployments`

Controls how many artifacts are deployed concurrently **per package**.

```yaml
# Conservative (safe for rate limits)
parallelDeployments: 2

# Balanced (recommended)
parallelDeployments: 3

# Aggressive (faster, but may hit rate limits)
parallelDeployments: 10
```

**Recommendations:**
- **Development:** 5-10 (speed over safety)
- **Production:** 2-3 (safety over speed)
- **CI/CD:** 5-10 (optimize for pipeline speed)

### `deployRetries`

Number of times to check deployment status before giving up.

```yaml
# Quick fail (development)
deployRetries: 3

# Standard (recommended)
deployRetries: 5

# Patient (production)
deployRetries: 10
```

**Total wait time = `deployRetries` × `deployDelaySeconds`**

### `deployDelaySeconds`

Seconds to wait between deployment status checks.

```yaml
# Fast polling (may overload API)
deployDelaySeconds: 10

# Balanced (recommended)
deployDelaySeconds: 15

# Conservative (slower but safer)
deployDelaySeconds: 30
```

**Recommendations:**
- Small artifacts: 10-15 seconds
- Large artifacts: 20-30 seconds
- Complex flows: 30-60 seconds

---

## Configuration Examples

### Example 1: Development Environment

```yaml
# orchestrator-dev.yml
packagesDir: ./packages
deployConfig: ./configs/dev
deploymentPrefix: DEV
mode: update-and-deploy

# Fast deployment for quick iteration
parallelDeployments: 5
deployRetries: 5
deployDelaySeconds: 15

# Merge all configs in folder
mergeConfigs: true
configPattern: "*.yml"
```

**Usage:**
```bash
flashpipe orchestrator --orchestrator-config ./orchestrator-dev.yml
```

### Example 2: Production Environment

```yaml
# orchestrator-prod.yml
packagesDir: ./packages
deployConfig: ./configs/production.yml
deploymentPrefix: PROD
mode: update-and-deploy

# Conservative settings for production
parallelDeployments: 2
deployRetries: 10
deployDelaySeconds: 30

# Production safety
mergeConfigs: false
keepTemp: false
```

**Usage:**
```bash
flashpipe orchestrator --orchestrator-config ./orchestrator-prod.yml
```

### Example 3: CI/CD Pipeline

```yaml
# orchestrator-ci.yml
packagesDir: ./packages
deployConfig: https://raw.githubusercontent.com/myorg/configs/main/ci-config.yml
deploymentPrefix: CI
mode: update-and-deploy

# Optimize for speed
parallelDeployments: 10
deployRetries: 5
deployDelaySeconds: 10

# No filtering - deploy everything
packageFilter: ""
artifactFilter: ""
```

**Usage in CI/CD:**
```yaml
# .github/workflows/deploy.yml
- name: Deploy to CPI
  run: |
    flashpipe orchestrator --orchestrator-config ./orchestrator-ci.yml
  env:
    CPI_HOST: ${{ secrets.CPI_HOST }}
    CPI_USERNAME: ${{ secrets.CPI_USERNAME }}
    CPI_PASSWORD: ${{ secrets.CPI_PASSWORD }}
```

### Example 4: Testing Single Package

```yaml
# orchestrator-test.yml
packagesDir: ./packages
deployConfig: ./test-config.yml
deploymentPrefix: TEST
mode: update-only  # Don't deploy, just update

# Focus on single package
packageFilter: "MyTestPackage"

# Debug settings
keepTemp: true
parallelDeployments: 1
```

**Usage:**
```bash
flashpipe orchestrator --orchestrator-config ./orchestrator-test.yml
```

### Example 5: Selective Deployment

```yaml
# orchestrator-selective.yml
packagesDir: ./packages
deployConfig: ./configs
deploymentPrefix: QA

# Deploy only specific packages and artifacts
packageFilter: "CustomerIntegration,DeviceManagement"
artifactFilter: "CustomerSync,DeviceStatusUpdate"

mode: update-and-deploy
parallelDeployments: 3
```

---

## CLI Flag Override

CLI flags always **override** YAML configuration:

```yaml
# orchestrator.yml
deploymentPrefix: DEV
parallelDeployments: 3
```

```bash
# Override prefix to PROD
flashpipe orchestrator \
  --orchestrator-config ./orchestrator.yml \
  --deployment-prefix PROD

# Result: Uses PROD prefix (not DEV)
```

**Override Priority:**
1. CLI flags (highest)
2. YAML config
3. Defaults (lowest)

---

## Advanced Usage

### Multi-Environment Setup

```
configs/
├── orchestrator-dev.yml
├── orchestrator-qa.yml
├── orchestrator-prod.yml
└── deploy-configs/
    ├── dev/
    │   ├── packages-1.yml
    │   └── packages-2.yml
    ├── qa/
    │   └── packages.yml
    └── prod/
        └── packages.yml
```

**Deploy to different environments:**
```bash
# Development
flashpipe orchestrator --orchestrator-config configs/orchestrator-dev.yml

# QA
flashpipe orchestrator --orchestrator-config configs/orchestrator-qa.yml

# Production
flashpipe orchestrator --orchestrator-config configs/orchestrator-prod.yml
```

### Remote Configuration

Load config from GitHub/GitLab:

```yaml
# orchestrator-remote.yml
packagesDir: ./packages
deployConfig: https://raw.githubusercontent.com/myorg/configs/main/deploy.yml
deploymentPrefix: CICD
parallelDeployments: 5
```

**With authentication:**
```bash
flashpipe orchestrator \
  --orchestrator-config ./orchestrator-remote.yml \
  --auth-token $GITHUB_TOKEN \
  --auth-type bearer
```

### Debugging Failed Deployments

```yaml
# orchestrator-debug.yml
packagesDir: ./packages
deployConfig: ./configs
mode: update-only  # Stop before deployment

# Keep files for inspection
keepTemp: true

# Single-threaded for easier debugging
parallelDeployments: 1

# Verbose logging
# (use --debug flag)
```

**Usage:**
```bash
flashpipe orchestrator \
  --orchestrator-config ./orchestrator-debug.yml \
  --debug

# Inspect temporary files
ls -la /tmp/flashpipe-orchestrator-*/
```

---

## Performance Tuning

### Optimize for Speed

```yaml
# Maximum parallelism
parallelDeployments: 10

# Faster polling
deployRetries: 5
deployDelaySeconds: 10

# Merge configs for single deployment
mergeConfigs: true
```

**Expected speedup:** 3-5x faster than sequential

### Optimize for Reliability

```yaml
# Conservative parallelism
parallelDeployments: 2

# More retries, longer delays
deployRetries: 10
deployDelaySeconds: 30

# Process configs separately
mergeConfigs: false
```

**Trade-off:** Slower but more stable

### Optimize for API Rate Limits

```yaml
# Low parallelism
parallelDeployments: 1

# Standard retries with longer delays
deployRetries: 5
deployDelaySeconds: 20
```

---

## Monitoring & Logging

### Deployment Output

```
═══════════════════════════════════════════════════════════════════════
PHASE 1: UPDATING ALL PACKAGES AND ARTIFACTS
═══════════════════════════════════════════════════════════════════════

📦 Package: MyPackage
  Updating: MyArtifact1
    ✓ Updated successfully
  Updating: MyArtifact2
    ✓ Updated successfully

═══════════════════════════════════════════════════════════════════════
PHASE 2: DEPLOYING ALL ARTIFACTS IN PARALLEL
═══════════════════════════════════════════════════════════════════════
Total artifacts to deploy: 2
Max concurrent deployments: 3

📦 Deploying 2 artifacts for package: MyPackage
  → Deploying: MyArtifact1 (type: IntegrationFlow)
  → Deploying: MyArtifact2 (type: IntegrationFlow)
  ✓ Deployed: MyArtifact1
  ✓ Deployed: MyArtifact2
✓ All 2 artifacts deployed successfully for package MyPackage

═══════════════════════════════════════════════════════════════════════
📊 DEPLOYMENT SUMMARY
═══════════════════════════════════════════════════════════════════════
Packages Updated:   1
Packages Deployed:  1
Artifacts Updated:       2
Artifacts Deployed OK:   2
✓ All operations completed successfully!
```

---

## Troubleshooting

### Problem: Deployments are slow

**Solution 1:** Increase parallelism
```yaml
parallelDeployments: 10  # Up from 3
```

**Solution 2:** Reduce polling delay
```yaml
deployDelaySeconds: 10  # Down from 15
```

### Problem: Hitting API rate limits

**Solution:** Reduce parallelism
```yaml
parallelDeployments: 1  # Down from 3
deployDelaySeconds: 20  # Up from 15
```

### Problem: Deployments timing out

**Solution:** Increase retries and delay
```yaml
deployRetries: 10        # Up from 5
deployDelaySeconds: 30   # Up from 15
```

### Problem: Hard to debug which artifact failed

**Solution:** Use debug mode
```bash
flashpipe orchestrator \
  --orchestrator-config ./config.yml \
  --debug \
  --keep-temp
```

---

## Best Practices

### ✅ DO

- Version control your orchestrator config files
- Use different configs for different environments
- Set conservative values for production
- Use `keepTemp: true` when debugging
- Test with `update-only` mode first
- Monitor deployment logs for errors

### ❌ DON'T

- Don't set `parallelDeployments` too high (>10)
- Don't use same config for all environments
- Don't skip testing in non-prod first
- Don't ignore failed deployments in summary
- Don't commit sensitive credentials to YAML

---

## Migration from CLI Flags

### Before (CLI flags)

```bash
flashpipe orchestrator \
  --packages-dir ./packages \
  --deploy-config ./config.yml \
  --deployment-prefix DEV \
  --merge-configs \
  --update
```

### After (YAML config)

```yaml
# orchestrator.yml
packagesDir: ./packages
deployConfig: ./config.yml
deploymentPrefix: DEV
mergeConfigs: true
mode: update-and-deploy
```

```bash
flashpipe orchestrator --orchestrator-config ./orchestrator.yml
```

**Benefits:**
- Easier to read and maintain
- Version controlled settings
- Reusable across teams
- Consistent deployments

---

## See Also

- [Orchestrator Quick Start](./orchestrator-quickstart.md)
- [Deployment Config Examples](./examples/)
- [Partner Directory Configuration](./partner-directory-config-examples.md)