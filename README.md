# Cimon - GitLab Pipeline Monitor

**Cimon** is a terminal-based GitLab pipeline monitor written in Go using [tview](https://github.com/rivo/tview). It allows you to quickly view the status of your GitLab pipelines and jobs, manage multiple projects, and receive real-time updates directly from your terminal.

---

## Features

- Monitor pipelines for multiple GitLab projects
- View pipeline statuses with visual indicators (success, failed, running, pending, canceled)
- Inspect job details including stage, duration, and status
- Add multiple GitLab projects and tokens directly from the interface
- Refresh pipelines and jobs on-demand

---

## Installation

1. Clone the repository:

```bash

git clone https://github.com/Youdontknowme720/Cimon.git
cd Cimon

```

2. Build the project with go:
```bash

go build -o cimon

```
3. Run the application
```bash

./cimon

```
```

## Usage

When you first start **Cimon**, you will see the **Home screen**, which lists your configured GitLab projects. From here, you can:

### Add a GitLab Token

1. Select the **+ Add Token** button.
2. Enter your GitLab personal access token (PAT) in the form.
   - This token is required to fetch pipelines and job details from GitLab.
3. Save the token to continue.

### Add a GitLab Project

1. Select the **+ Add Project** button.
2. Enter the **Project ID** from GitLab (numerical identifier of the repository).
3. Enter a **Project Name** of your choice (this is only for display in Cimon).
4. Save the project.

> After adding a project and token, your configuration is stored automatically.

### Viewing Pipelines

1. After adding a project, select it from the Home screen.
2. You will see a table listing the latest pipelines with **status icons**.

Each pipeline displays:

- Status emoji
- Commit message (truncated if too long)
- Short SHA of the commit

### Viewing Jobs

1. Select a pipeline to view its jobs.
2. Each job shows:
   - Status emoji
   - Job name (truncated if too long)
   - Duration
   - Stage
3. Select a job to view details in a modal window.

### Refreshing Data

- Press `r` while viewing pipelines or jobs to refresh the data.
- The interface will show a loading indicator during updates.

### Navigation

- `b` → Go back to the previous screen  
- `Esc` → Exit Cimon

## Configuration

Cimon stores your configuration in a **YAML file** (e.g., `config.yml`) located in the application folder. The configuration contains:

- Your GitLab token
- The list of configured projects (Project ID + Project Name)

> When adding a token or project via the Home screen, this file is updated automatically.

### Example `config.yml`:

```yaml
token: "YOUR_GITLAB_PERSONAL_ACCESS_TOKEN"
projects:
  - id: 123456
    name: "My Project"
  - id: 987654
    name: "Another Project"
```
```
Replace `YOUR_GITLAB_PERSONAL_ACCESS_TOKEN` with the token you added in the app. The project IDs are numerical GitLab repository identifiers, and the names can be any descriptive string you like.
