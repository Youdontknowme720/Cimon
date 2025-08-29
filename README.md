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

2. Build the project with go

```bash
go build -o cimon

3. Run the application

```bash
./cimon
