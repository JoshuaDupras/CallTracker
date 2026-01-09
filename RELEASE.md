# Release Process

This document describes how to create and publish releases for the Fire Department Call Log application.

## Quick Release (Automated)

**The easiest way to create a release:**

```bash
go run scripts/release
```

The interactive release tool will:
1. ✅ Check for uncommitted changes
2. ✅ Prompt for version bump type (patch/minor/major/custom)
3. ✅ Update `wails.json` with the version
4. ✅ Prompt for release notes and update `CHANGELOG.md`
5. ✅ Run tests
6. ✅ Build the application
7. ✅ Create git commit and tag
8. ✅ Push to GitHub (with confirmation)
9. ✅ Trigger GitHub Actions to build and publish

**Note:** On first run, Go will automatically download required dependencies (Bubble Tea UI library).

---

## Manual Release Process

If you prefer to do it manually or the automated script fails:

## Versioning

This project follows [Semantic Versioning](https://semver.org/):
- **MAJOR** version for incompatible API changes or major feature overhauls
- **MINOR** version for new features in a backward-compatible manner
- **PATCH** version for backward-compatible bug fixes

Format: `vMAJOR.MINOR.PATCH` (e.g., `v1.0.0`, `v1.2.3`)

## Creating a Release

### 1. Prepare the Release

**Update Version Information:**
1. Update the version in `wails.json`:
   ```json
   {
     "version": "1.0.0"
   }
   ```

2. Update `README.md` if there are significant changes to document

3. Create/update `CHANGELOG.md` with release notes:
   ```markdown
   ## [1.0.0] - 2026-01-08
   
   ### Added
   - Initial public release
   - 12-step call entry wizard
   - Apparatus and responder tracking
   - Year-based call filtering
   - Statistics dashboard
   - PDF and CSV export
   
   ### Fixed
   - (List any bug fixes)
   
   ### Changed
   - (List any changes)
   ```

**Test Thoroughly:**
```powershell
# Clean build and test
wails build -clean

# Test the executable
.\build\bin\fd-call-log.exe

# Run Go tests
go test ./... -v
```

### 2. Commit and Tag

```powershell
# Stage all changes
git add .

# Commit with version message
git commit -m "Release v1.0.0"

# Create annotated tag
git tag -a v1.0.0 -m "Release version 1.0.0 - Initial public release"

# Push commits and tags
git push origin main
git push origin v1.0.0
```

### 3. GitHub Actions Automatic Build

When you push a tag (e.g., `v1.0.0`), GitHub Actions will automatically:
1. Build the application for Windows, Linux, and macOS
2. Create a GitHub Release
3. Upload the binaries as release assets

**Release artifacts:**
- `fd-call-log-windows-amd64.zip` - Windows executable
- `fd-call-log-linux-amd64.tar.gz` - Linux binary
- `fd-call-log-darwin-universal.tar.gz` - macOS universal app

### 4. Finalize GitHub Release

1. Go to https://github.com/YOUR_USERNAME/call-tracker-wails/releases
2. Find your new release (created automatically by GitHub Actions)
3. Edit the release notes if needed
4. Add detailed release notes from CHANGELOG.md
5. Mark as "Latest release" if appropriate
6. Uncheck "This is a pre-release" for stable versions

## Manual Release (Alternative)

If you need to create a release manually without GitHub Actions:

### Build for All Platforms

```powershell
# Windows
wails build -platform windows/amd64

# Linux (requires cross-compilation setup)
wails build -platform linux/amd64

# macOS (requires macOS machine or cross-compilation)
wails build -platform darwin/universal
```

### Create Release on GitHub

1. Go to https://github.com/YOUR_USERNAME/call-tracker-wails/releases/new
2. Click "Choose a tag" and type `v1.0.0` (or your version)
3. Click "Create new tag: v1.0.0 on publish"
4. Set release title: "v1.0.0 - Initial Release"
5. Add release notes from CHANGELOG.md
6. Attach the built binaries:
   - `fd-call-log.exe` (Windows)
   - `fd-call-log` (Linux)
   - `fd-call-log.app` (macOS, zipped)
7. Click "Publish release"

## Release Checklist

Before creating a release, ensure:

- [ ] All tests pass (`go test ./...`)
- [ ] Application builds successfully (`wails build`)
- [ ] Version number updated in `wails.json`
- [ ] CHANGELOG.md updated with release notes
- [ ] README.md is up to date
- [ ] No hardcoded passwords or secrets in code
- [ ] .gitignore excludes sensitive files
- [ ] All features documented
- [ ] Database migrations (if any) are documented
- [ ] Security warnings in README are clear

## Hotfix Process

For critical bug fixes that need immediate release:

1. Create a hotfix branch from main:
   ```powershell
   git checkout -b hotfix/v1.0.1 main
   ```

2. Make the fix and commit:
   ```powershell
   git commit -m "Fix critical bug in call saving"
   ```

3. Update version to `v1.0.1` in `wails.json`

4. Merge back to main:
   ```powershell
   git checkout main
   git merge hotfix/v1.0.1
   ```

5. Tag and push:
   ```powershell
   git tag -a v1.0.1 -m "Hotfix: Fix critical call saving bug"
   git push origin main
   git push origin v1.0.1
   ```

## Pre-releases

For beta or release candidate versions:

1. Use tag format: `v1.0.0-beta.1` or `v1.0.0-rc.1`
2. Mark as "This is a pre-release" on GitHub
3. Clearly label in release notes that it's not production-ready

## Distribution

After release, users can download:

**Windows:**
1. Download `fd-call-log-windows-amd64.zip`
2. Extract the ZIP file
3. Run `fd-call-log.exe`

**Linux:**
1. Download `fd-call-log-linux-amd64.tar.gz`
2. Extract: `tar -xzf fd-call-log-linux-amd64.tar.gz`
3. Run: `./fd-call-log`

**macOS:**
1. Download `fd-call-log-darwin-universal.tar.gz`
2. Extract: `tar -xzf fd-call-log-darwin-universal.tar.gz`
3. Drag `fd-call-log.app` to Applications folder
4. Right-click and select "Open" (first time only, to bypass Gatekeeper)

## Rollback

If a release has critical issues:

1. Mark the release as a pre-release on GitHub
2. Create a new hotfix release
3. Update README with known issues
4. Communicate to users via GitHub Discussions or Issues

## Support

After releasing:
- Monitor GitHub Issues for bug reports
- Respond to questions in GitHub Discussions
- Update documentation as needed
- Plan next release based on feedback
