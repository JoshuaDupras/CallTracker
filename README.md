# Fire Department Call Log

A desktop application for tracking fire department emergency calls. Simple to use with no internet connection required!

## What This Application Does

This program helps fire departments:
- Track emergency calls with automatic call numbers
- Record who responded and what equipment was used
- Keep notes about each incident
- View call history by year with statistics
- Export reports to PDF and CSV formats

Everything is stored on your computer - no internet needed, no monthly fees, no complicated setup.

## Installation Guide

### First Time Setup (One-Time Only)

**⚠️ IMPORTANT**: You need to install some free software before using this application. These steps require some technical knowledge.

#### Step 1: Install Go Programming Language
1. Visit https://go.dev/dl/
2. Download the Windows installer (the big blue button)
3. Run the installer and click "Next" through all the screens
4. When done, **restart your computer**

#### Step 2: Install Node.js
1. Visit https://nodejs.org/
2. Download the "LTS" version (recommended for most users)
3. Run the installer and click "Next" through all the screens
4. When done, **restart your computer** again

#### Step 3: Install Wails (Command Line Required)
1. Press the Windows key and type `PowerShell`
2. Right-click on "Windows PowerShell" and select "Run as administrator"
3. Copy and paste this line, then press Enter:
   ```powershell
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```
4. Wait for it to finish (you'll see the cursor blinking when done)
5. Close PowerShell

**✅ You're done with setup! You only need to do this once.**

---

### Building the Application

Now that you have the required software installed, you can build the program:

**For Development (includes version info):**
1. Press the Windows key and type `PowerShell`
2. Right-click on "Windows PowerShell" and select "Run as administrator"
3. Navigate to the application folder:
   ```powershell
   cd path\to\call-tracker-wails
   ```
   (Replace `path\to\call-tracker-wails` with wherever you downloaded this project)
4. Build with version info:
   ```powershell
   go run scripts/build
   ```
   This embeds git commit info and build timestamp into the executable.

**Standard Build (without version info):**
1. Follow steps 1-3 above
2. Build the application:
   ```powershell
   wails build
   ```
3. Wait 1-2 minutes while it builds (you'll see progress messages)

**✅ The program is now ready!** You'll find it at: `build\bin\fd-call-log.exe`

---

### Running the Application

**Easy Way (After Building):**
1. Open File Explorer
2. Navigate to the project folder
3. Go into the `build\bin` folder
4. Double-click `fd-call-log.exe`

**Alternative (Build and Run Together):**
1. Open PowerShell as administrator
2. Navigate to the project folder
3. Run:
   ```powershell
   wails build; .\build\bin\fd-call-log.exe
   ```

---

### First Login

When you first start the application:
- **Username**: Admin User
- **PIN**: 1234

**⚠️ IMPORTANT**: Change this PIN immediately after your first login!
1. Click the settings/profile icon
2. Select "Change PIN"
3. Choose a new 4-digit PIN you'll remember

---

## Features

- ✅ **Simple PIN login** - No complex passwords
- ✅ **Automatic call numbering** - Calls are numbered like 2026-001, 2026-002, etc.
- ✅ **Call logging wizard** - Step-by-step questions to record each call
- ✅ **Track apparatus** - Select which trucks/equipment responded
- ✅ **Track responders** - Mark which firefighters were on scene
- ✅ **Time tracking** - Record when dispatched, en route, on scene, and cleared
- ✅ **Call history** - View all calls by year with statistics
- ✅ **Clickable details** - Click any call to see full information
- ✅ **Export reports** - Generate PDF or CSV reports
- ✅ **Works offline** - No internet connection needed
- ✅ **Admin tools** - Add/edit users and dropdown options

---

## Common Questions

### "Where is my data stored?"
In a file called `fd-calls.db` in the same folder as the program. **Back this file up regularly!**

### "Can I use this on multiple computers?"
Yes, but you'll need to copy the `fd-calls.db` file between computers. Each computer runs independently.

### "How do I update the program?"
1. Download the new version
2. **Copy your `fd-calls.db` file to a safe place first!**
3. Build the new version
4. Copy your `fd-calls.db` file back

### "It won't start!"
- Make sure you closed all other copies of the program
- Try restarting your computer
- Delete the `fd-calls.db` file (⚠️ this erases all data!) and start fresh

### "I forgot my PIN!"
If you're an admin, another admin can reset it. If there are no admins available, you'll need to delete the database file and start over.

---

## Technical Information (For Developers)

### Technology Stack
- **Backend**: Go 1.21+ with Wails v2.11.0
- **Frontend**: Vanilla JavaScript, HTML, CSS
- **Database**: SQLite (embedded, no server required)
- **Platform**: Windows (can be built for macOS/Linux)

### Prerequisites
- Go 1.21 or later
- Node.js 16+ and npm
- Wails CLI v2
- C compiler (MinGW-w64 or TDM-GCC for Windows)

### Development Commands

```powershell
# Install dependencies
go mod tidy
cd frontend && npm install && cd ..

# Development mode with hot reload
wails dev

# Production build
wails build

# Run without Wails CLI
go run .

# Clean build (if having issues)
wails build -clean
```

### Project Structure

```
call-tracker-wails/
├── fd-call-log.exe         # Run this file! (after building)
├── fd-calls.db             # Your data (BACK THIS UP!)
├── app.go                  # Backend application methods
├── main.go                 # Entry point
├── wails.json             # Configuration
├── go.mod                 # Go dependencies
├── build/
│   └── bin/
│       └── fd-call-log.exe # Built application
├── internal/
│   ├── db/                # Database code
│   └── export/            # PDF/CSV export code
└── frontend/
    └── src/               # User interface files
```

### Database Schema

The application uses SQLite with these tables:
- **users** - Fire department members with PIN authentication
- **calls** - Emergency call records with all incident details
- **picklists** - Dropdown values (call types, towns, apparatus, etc.)
- **call_apparatus** - Which trucks/equipment responded to each call
- **call_responders** - Which firefighters responded to each call
- **audit_log** - Activity tracking for security

### Call Data Model

Each call includes:
- **incident_number**: Auto-generated (e.g., 2026-001)
- **call_type**: Type of emergency (fire, EMS, MVA, etc.)
- **mutual_aid**: Whether giving or receiving assistance
- **address**: Location of incident
- **town**: Jurisdiction
- **location_notes**: Additional location details
- **dispatched**: When the call came in
- **enroute**: When units left the station
- **on_scene**: When units arrived
- **clear**: When units became available again
- **narrative**: Detailed description of incident
- **apparatus**: List of equipment used
- **responders**: List of personnel who responded

### Default Users

The application starts with one admin account:

| Name       | PIN  | Role  |
|------------|------|-------|
| Admin User | 1234 | Admin |

**⚠️ Change this PIN immediately after first login!**

To add more users:
1. Log in as admin
2. Click "Manage Users"
3. Click "Add User"
4. Fill in details and assign a 4-digit PIN

### Customizing Dropdown Options

As an admin, you can customize the dropdown lists:
1. Log in as admin
2. Click "Manage Settings" or "Picklists"
3. Choose a category (Call Types, Towns, Apparatus, etc.)
4. Add, edit, or remove options

Categories include:
- **Call Types**: Fire, EMS, MVA, Hazmat, etc.
- **Towns**: Jurisdictions you serve
- **Mutual Aid Agencies**: Neighboring departments
- **Apparatus**: Engine 1, Ladder 1, Rescue 1, etc.

---

## Troubleshooting

### For Everyone

**"The program won't start"**
1. Make sure you closed all other windows of this program
2. Try restarting your computer
3. Check if your antivirus is blocking it

**"I can't find my calls"**
- Use the year dropdown at the top to switch between years
- Calls are organized by the year they were dispatched

**"How do I backup my data?"**
1. Close the program completely
2. Copy the file `fd-calls.db` to a USB drive or cloud storage
3. Keep multiple backups in different places!

### For Technical Users

**"go: command not found"**
- Ensure Go is installed and in your PATH
- Restart your terminal after installation
- Check with: `go version`

**"wails: command not found"**
- Run: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- Ensure `%USERPROFILE%\go\bin` is in your PATH
- On Windows: Add `C:\Users\[YourUsername]\go\bin` to PATH

**"gcc: executable file not found"**
- Install MinGW-w64 or TDM-GCC
- Add the `bin` directory to your PATH
- Verify with: `gcc --version`

**"Database is locked"**
- Close all instances of the application
- Kill any hung processes: `taskkill /F /IM fd-call-log.exe`
- Try again

**Build fails with frontend errors**
```powershell
cd frontend
npm install
cd ..
wails build -clean
```

---

## Development Guide

### Adding Backend Features
1. Add methods to `app.go`
2. Methods are auto-exposed to frontend
3. Rebuild to generate JavaScript bindings
4. Call from frontend: `window.go.main.App.MethodName()`

### Adding Frontend Features  
1. Edit files in `frontend/src/`
   - `index.html` - Structure
   - `style.css` - Appearance
   - `app.js` - Logic
2. Use `wails dev` for live reload during development
3. Build production version: `wails build`

### Database Queries
- Add to `internal/db/queries_*.go`
- Follow existing patterns for consistency
- Use prepared statements to prevent SQL injection

### Export Features
- PDF generation: `internal/export/pdf.go`
- CSV generation: `internal/export/csv.go`
- Both use the Call model from database

---

## Building for Distribution

### Windows
```powershell
wails build -clean
```
Output: `build\bin\fd-call-log.exe`

### macOS
```bash
wails build -clean
```
Output: `build/bin/fd-call-log.app`

### Linux
```bash
wails build -clean
```
Output: `build/bin/fd-call-log`

### Creating an Installer
Wails can create installers:
```powershell
wails build -clean -nsis
```

---

## Backup and Restore

### Backing Up Data
**Method 1: Manual Copy**
1. Close the application completely
2. Copy `fd-calls.db` to a safe location
3. Label it with the date (e.g., `fd-calls-2026-01-08.db`)

**Method 2: Export Reports**
1. Open the application
2. Go to Call List
3. Use Export to CSV to save call data
4. These can be opened in Excel

### Restoring Data
1. Close the application
2. Copy your backup `fd-calls.db` file
3. Paste it into the application folder (replacing the current one)
4. Start the application

---

## Security Notes

- **PINs are hashed** - Not stored in plain text
- **Database is local** - No data sent over internet
- **Audit log** - Tracks who did what and when
- **No default remote access** - Application only runs locally

### Best Practices
1. Change default admin PIN immediately
2. Use unique PINs for each user
3. Back up your database regularly
4. Keep the computer physically secure
5. Don't share user PINs
6. Remove users who leave the department

---

## Version History

### Current Version
- 12-step call entry wizard
- Apparatus and responder tracking
- Year-based call filtering
- Statistics dashboard
- Export to PDF and CSV
- User and picklist management

---

## Contributing and Development

### Running Tests
```powershell
go test ./... -v
```

### Making Changes
1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make your changes and commit: `git commit -m "Add my feature"`
4. Push to your fork: `git push origin feature/my-feature`
5. Open a Pull Request

### Release Process

See [RELEASE.md](RELEASE.md) for detailed instructions on creating releases.

**Quick automated release:**
```bash
go run scripts/release
```

The interactive tool will guide you through the entire release process.

**Manual release:**
```powershell
# Update version in wails.json
# Update CHANGELOG.md
git add .
git commit -m "Release v1.0.0"
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin main
git push origin v1.0.0
```

GitHub Actions will automatically build binaries for Windows, Linux, and macOS and create a release.

---
