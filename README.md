# Cimon - GitLab Pipeline Monitor

**Cimon** is a terminal-based GitLab pipeline monitor written in Go using [tview](https://github.com/rivo/tview). It provides real-time monitoring of GitLab pipelines and jobs, multi-project management, and an intuitive terminal interface for DevOps workflows.

---

## Features

- **Multi-project monitoring** - Track pipelines across multiple GitLab projects simultaneously
- **Visual status indicators** - Clear emoji-based status display (success ✅, failed ❌, running 🔄, pending ⏳, canceled ⏹️)
- **Detailed job inspection** - View job stages, durations, and execution details
- **In-app configuration** - Add GitLab projects and tokens directly from the interface
- **Real-time updates** - On-demand refresh capabilities with loading indicators
- **Keyboard-driven navigation** - Efficient terminal-based workflow

---

## Installation

### Prerequisites
- Go 1.19 or higher
- GitLab personal access token with API read permissions

### Build from source
```bash
git clone https://github.com/Youdontknowme720/Cimon.git
cd Cimon
go build -o cimon
./cimon
```

---

## Quick Start

1. **Launch Cimon** - Run `./cimon` to open the Home screen
2. **Add your GitLab token** - Select **+ Add Token** and enter your personal access token
3. **Add a project** - Select **+ Add Project** and provide the GitLab project ID and display name
4. **Monitor pipelines** - Select your project to view real-time pipeline status

---

## Usage Guide

### Initial Setup

#### Adding a GitLab Token
1. Navigate to GitLab → Settings → Access Tokens
2. Create a token with `read_api` scope
3. In Cimon, select **+ Add Token**
4. Paste your token and save

#### Adding Projects
1. Find your GitLab project ID (visible in project settings or URL)
2. Select **+ Add Project** in Cimon
3. Enter the numeric **Project ID** 
4. Provide a descriptive **Project Name** for display
5. Save the configuration

### Monitoring Workflows

#### Pipeline Overview
- Select any configured project from the Home screen
- View pipelines with status indicators:
  - ✅ Success - Pipeline completed successfully
  - ❌ Failed - Pipeline failed
  - 🔄 Running - Pipeline currently executing
  - ⏳ Pending - Pipeline queued for execution
  - ⏹️ Canceled - Pipeline was canceled

#### Job Details
1. Select any pipeline to drill down into job details
2. Inspect individual jobs showing:
   - Job status and name
   - Execution stage
   - Runtime duration
3. Select specific jobs for detailed modal view

#### Data Management
- **Refresh**: Press `r` to update pipeline/job data
- **Auto-save**: All configuration changes are saved automatically
- **Loading indicators**: Visual feedback during data fetching

### Navigation Controls

| Key | Action |
|-----|--------|
| `r` | Refresh current data |
| `b` | Navigate back |
| `Esc` | Exit application |
| `Enter` | Select item |
| `Tab` | Navigate between elements |

---

## Configuration

Cimon automatically manages configuration in `config.yml` within the application directory.

### Configuration Structure
```yaml
token: "glpat-xxxxxxxxxxxxxxxxxxxx"
projects:
  - id: 12345678
    name: "Frontend Application"
  - id: 87654321
    name: "API Backend"
  - id: 11223344
    name: "DevOps Tools"
```

### Security Notes
- Store your `config.yml` securely and avoid committing it to version control
- Use GitLab tokens with minimal required permissions (`read_api`)
- Consider using environment-specific tokens for different GitLab instances

---

## Troubleshooting

### Common Issues

**"No pipelines found"**
- Verify project ID is correct (numeric identifier from GitLab)
- Ensure your token has `read_api` permissions
- Check if the project has any pipelines

**"Authentication failed"**
- Confirm your GitLab token is valid and not expired
- Verify token permissions include API access
- Test token manually with GitLab API

**"Connection timeout"**
- Check network connectivity to GitLab instance
- Verify GitLab server availability
- Consider firewall or proxy restrictions

### Getting Help
- Check the [Issues](https://github.com/Youdontknowme720/Cimon/issues) page for known problems
- Contribute bug reports with system information and error details

