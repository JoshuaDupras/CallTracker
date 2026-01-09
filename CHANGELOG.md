# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned
- Additional export formats
- Advanced search and filtering
- Call statistics and analytics dashboard
- Mobile companion app

## [0.0.1] - 2026-01-09

first release

## [1.0.0] - 2026-01-08

### Added
- Initial public release
- User authentication with 4-digit PIN system
- 12-step call entry wizard with validation
- Automatic call number generation (YYYY-NNN format)
- Apparatus tracking with checkbox selection
- Responder tracking with active user roster
- Year-based call filtering and viewing
- Statistics dashboard (total calls, mutual aid, call types)
- Call details modal with full incident information
- PDF export for individual calls and summary reports
- CSV export for data analysis
- Admin user management (add, edit, reset PIN, deactivate)
- Customizable picklists for dropdown values
- SQLite database with automatic schema initialization
- Audit logging for security and accountability
- Comprehensive README with installation guide
- GitHub Actions CI/CD for automated builds
- Cross-platform support (Windows, Linux, macOS)

### Security
- PIN hashing for user credentials
- Local-only data storage (no cloud/internet)
- Default admin PIN with mandatory change warning
- User deactivation instead of deletion

### Documentation
- Installation guide for non-technical users
- Backup and restore procedures
- Troubleshooting section
- Development guide for contributors
- Release process documentation
