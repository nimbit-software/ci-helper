# Orchestrator Quick Start Guide

## What's New? 🎉

The Flashpipe orchestrator has been upgraded with three major enhancements:

1. **📝 YAML Configuration** - Load all settings from a config file
2. **⚡ Parallel Deployment** - Deploy multiple artifacts simultaneously (3-5x faster!)
3. **🔄 Two-Phase Strategy** - Update all artifacts first, then deploy in parallel

---

## Quick Start in 30 Seconds

### 1. Create Orchestrator Config

```yaml
# orchestrator.yml
packagesDir: ./packages
deployConfig: ./deploy-config.yml
deploymentPrefix: DEV
mode: update-and-deploy
parallelDeployments: 5
```

### 2. Run the Orchestrator

```bash
flashpipe orchestrator --orchestrator-config ./orchestrator.yml
```

That's it! 🚀

---

## Before vs After

### ❌ Old Way (Many CLI Flags)

```bash
flashpipe orchestrator \
  --packages-dir ./packages \
  --deploy-config ./config.yml \
  --deployment-prefix DEV \
  --merge-configs \
  --parallel-deployments 5 \
  --deploy-retries 10 \
  --deploy-delay 20 \
  --update
```

### ✅ New Way (Simple!)

```bash
flashpipe orchestrator --orchestrator-config ./orchestrator.yml
```

---

## How It Works Now

### Two-Phase Deployment

```
╔═══════════════════════════════════════════════════════════╗
║              PHASE 1: UPDATE EVERYTHING                   ║
╚═══════════════════════════════════════════════════════════╝
  📦 Package 1
     ✓ Update Artifact A
     ✓ Update Artifact B
     ✓ Update Artifact C
  
  📦 Package 2
     ✓ Update Artifact D
     ✓ Update Artifact E

╔═══════════════════════════════════════════════════════════╗
║          PHASE 2: DEPLOY IN PARALLEL (5x faster!)         ║
╚═══════════════════════════════════════════════════════════╝
  📦 Package 1: Deploying 3 artifacts...
     → Deploy A, B, C (all at once!)
     ✓ All deployed in ~2 minutes
  
  📦 Package 2: Deploying 2 artifacts...
     → Deploy D, E (all at once!)
     ✓ All deployed in ~2 minutes

Total Time: ~4 minutes (instead of ~20 minutes!)
```

---

## Essential Configurations

### Development (Fast & Loose)

```yaml
packagesDir: ./packages
deployConfig: ./dev-config.yml
deploymentPrefix: DEV
parallelDeployments: 10    # Maximum speed
deployRetries: 5
deployDelaySeconds: 10
```

### Production (Safe & Reliable)

```yaml
packagesDir: ./packages
deployConfig: ./prod-config.yml
deploymentPrefix: PROD
parallelDeployments: 2     # Conservative
deployRetries: 10          # More retries
deployDelaySeconds: 30     # Longer waits
```

### CI/CD Pipeline

```yaml
packagesDir: ./packages
deployConfig: https://raw.githubusercontent.com/org/repo/main/config.yml
deploymentPrefix: CI
parallelDeployments: 5
mode: update-and-deploy
```

---

## Common Use Cases

### Deploy Everything

```bash
flashpipe orchestrator --orchestrator-config ./orchestrator.yml
```

### Deploy Specific Packages

```bash
flashpipe orchestrator \
  --orchestrator-config ./orchestrator.yml \
  --package-filter "Package1,Package2"
```

### Update Only (No Deployment)

```bash
flashpipe orchestrator \
  --orchestrator-config ./orchestrator.yml \
  --update-only
```

### Deploy Only (Skip Updates)

```bash
flashpipe orchestrator \
  --orchestrator-config ./orchestrator.yml \
  --deploy-only
```

### Debug Mode

```bash
flashpipe orchestrator \
  --orchestrator-config ./orchestrator.yml \
  --keep-temp \
  --debug
```

---

## Configuration Options

### All Available Settings

```yaml
# Required
packagesDir: ./packages              # Where your packages are
deployConfig: ./config.yml           # Deployment configuration

# Optional: Filtering
deploymentPrefix: DEV                # Prefix for all IDs
packageFilter: "Pkg1,Pkg2"          # Only these packages
artifactFilter: "Art1,Art2"         # Only these artifacts

# Optional: Config Loading
configPattern: "*.yml"               # File pattern for folders
mergeConfigs: true                   # Merge multiple configs

# Optional: Behavior
keepTemp: false                      # Keep temp files for debugging
mode: update-and-deploy              # See modes below

# Optional: Performance Tuning
parallelDeployments: 3               # Concurrent deployments
deployRetries: 5                     # Status check retries
deployDelaySeconds: 15               # Delay between checks
```

### Operation Modes

| Mode | What It Does |
|------|--------------|
| `update-and-deploy` | Update artifacts, then deploy (default) |
| `update-only` | Only update, skip deployment |
| `deploy-only` | Only deploy, skip updates |

---

## Performance Tuning

### Speed vs Safety Trade-offs

```yaml
# Maximum Speed (Development)
parallelDeployments: 10
deployRetries: 3
deployDelaySeconds: 10

# Balanced (Recommended)
parallelDeployments: 5
deployRetries: 5
deployDelaySeconds: 15

# Maximum Safety (Production)
parallelDeployments: 2
deployRetries: 10
deployDelaySeconds: 30
```

### Expected Deployment Times

| Artifacts | Sequential | Parallel (5x) | Speedup |
|-----------|-----------|---------------|---------|
| 5 | 10 min | 2 min | **5x faster** |
| 10 | 20 min | 4 min | **5x faster** |
| 20 | 40 min | 8 min | **5x faster** |

---

## Sample Output

```
Starting flashpipe orchestrator
Deployment Strategy: Two-phase with parallel deployment
  Phase 1: Update all artifacts
  Phase 2: Deploy all artifacts in parallel (max 5 concurrent)

═══════════════════════════════════════════════════════════════
PHASE 1: UPDATING ALL PACKAGES AND ARTIFACTS
═══════════════════════════════════════════════════════════════

📦 Package: CustomerIntegration
  Updating: CustomerSync
    ✓ Updated successfully
  Updating: CustomerDataTransform
    ✓ Updated successfully

═══════════════════════════════════════════════════════════════
PHASE 2: DEPLOYING ALL ARTIFACTS IN PARALLEL
═══════════════════════════════════════════════════════════════
Total artifacts to deploy: 2
Max concurrent deployments: 5

📦 Deploying 2 artifacts for package: CustomerIntegration
  → Deploying: CustomerSync (type: IntegrationFlow)
  → Deploying: CustomerDataTransform (type: IntegrationFlow)
  ✓ Deployed: CustomerSync
  ✓ Deployed: CustomerDataTransform
✓ All 2 artifacts deployed successfully

═══════════════════════════════════════════════════════════════
📊 DEPLOYMENT SUMMARY
═══════════════════════════════════════════════════════════════
Packages Updated:   1
Packages Deployed:  1
Artifacts Updated:       2
Artifacts Deployed OK:   2
✓ All operations completed successfully!
═══════════════════════════════════════════════════════════════
```

---

## Troubleshooting

### Deployments Are Slow

**Increase parallelism:**
```yaml
parallelDeployments: 10
deployDelaySeconds: 10
```

### Hitting API Rate Limits

**Reduce parallelism:**
```yaml
parallelDeployments: 1
deployDelaySeconds: 20
```

### Deployments Timing Out

**Increase retries:**
```yaml
deployRetries: 10
deployDelaySeconds: 30
```

---

## Migration from Old Orchestrator

Your old commands still work! But here's how to upgrade:

### Old Command
```bash
flashpipe orchestrator --packages-dir ./packages --deploy-config ./config.yml --update
```

### New Command
```yaml
# orchestrator.yml
packagesDir: ./packages
deployConfig: ./config.yml
mode: update-and-deploy
```

```bash
flashpipe orchestrator --orchestrator-config ./orchestrator.yml
```

**Bonus: Automatic parallel deployment! 🚀**

---

## CI/CD Integration

### GitHub Actions

```yaml
- name: Deploy to CPI
  run: flashpipe orchestrator --orchestrator-config ./orchestrator-ci.yml
  env:
    CPI_HOST: ${{ secrets.CPI_HOST }}
    CPI_USERNAME: ${{ secrets.CPI_USERNAME }}
    CPI_PASSWORD: ${{ secrets.CPI_PASSWORD }}
```

### GitLab CI

```yaml
deploy:
  script:
    - flashpipe orchestrator --orchestrator-config ./orchestrator-ci.yml
  environment: production
```

---

## Next Steps

1. **Try it out** - Create an orchestrator config file
2. **Test with update-only** - Verify updates work correctly
3. **Deploy to dev** - Use parallel deployment in development
4. **Tune performance** - Adjust parallelism for your environment
5. **Deploy to prod** - Use conservative settings for production

---

## Complete Example

```yaml
# orchestrator-dev.yml
packagesDir: ./packages
deployConfig: ./configs/dev
deploymentPrefix: DEV
packageFilter: ""
artifactFilter: ""
configPattern: "*.yml"
mergeConfigs: true
keepTemp: false
mode: update-and-deploy
deployRetries: 5
deployDelaySeconds: 15
parallelDeployments: 5
```

```bash
# Deploy to development
flashpipe orchestrator --orchestrator-config ./orchestrator-dev.yml

# Override for testing single package
flashpipe orchestrator \
  --orchestrator-config ./orchestrator-dev.yml \
  --package-filter "TestPackage" \
  --update-only
```

---

## Documentation

- 📘 [Full YAML Configuration Guide](./docs/orchestrator-yaml-config.md)
- 📊 [Detailed Enhancements](./ORCHESTRATOR_ENHANCEMENTS.md)
- 📝 [Example Configurations](./docs/examples/orchestrator-config-example.yml)

---

## Summary

**What You Get:**
- ⚡ **3-5x faster deployments** through parallelization
- 📝 **Simpler configuration** via YAML files
- 🔍 **Better visibility** with two-phase approach
- 🎯 **Tunable performance** for different environments
- ✅ **Backward compatible** - old commands still work!

**Get Started:**
```bash
flashpipe orchestrator --orchestrator-config ./orchestrator.yml
```

Happy deploying! 🚀