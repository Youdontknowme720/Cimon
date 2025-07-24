package ui

import (
	"fmt"
	"github.com/Youdontknowme720/Cimon/github"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type AppState struct{
	App *tview.Application
	Pages *tview.Pages
	Token string
	Repo string
}

type WorkflowSelectCallback func(workflow github.Workflow)

func StartView(app *tview.Application, repo string, token string) error {
	workflows, err := github.GetWorkflowStatus(repo, 5, token)
	if err != nil {
		return fmt.Errorf("Fehler beim Abrufen der Workflows: %w", err)
	}

	pages := tview.NewPages()
	state := &AppState{
		App:   app,
		Pages: pages,
		Repo:  repo,
		Token: token,
	}

	ShowWorkflowTable(state, workflows)

	// Start der App
	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		return fmt.Errorf("Fehler beim Starten der App: %w", err)
	}

	return nil
}

func ShowWorkflowTable(state *AppState,
		workflows github.WorkflowRunsResponse,
		) {

	table := createTable([]string{"#", "Name", "Status"})
	for i, wf := range workflows.WorkflowRuns {
		color := statusColor(wf.Conclusion)
		table.SetCell(i+1, 0, tview.NewTableCell(fmt.Sprintf("%d  ", i+1)))
		table.SetCell(i+1, 1, tview.NewTableCell(wf.DisplayTitle))
		table.SetCell(i+1, 2, tview.NewTableCell(fmt.Sprintf("[%s]%s", color, wf.Conclusion)))
	}
	table.Select(1, 0).SetFixed(1, 0).SetSelectable(true, false)
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		row, _ := table.GetSelection()
		switch event.Rune() {
		case 'j':
			if row < table.GetRowCount()-1 {
				table.Select(row+1, 0)
			}
			return nil
		case 'k':
			if row > 1 {
				table.Select(row-1, 0)
			}
			return nil
		case 'q':
			state.App.Stop()
			return nil
		}
		if event.Key() == tcell.KeyEnter && row > 0 {
			selected := workflows.WorkflowRuns[row-1]
			onWorkflowEnter(state, selected)
			return nil
		}
		return event
	})
	table.SetTitle(" GitHub Workflows ").SetBorder(true)
	state.Pages.AddAndSwitchToPage("Workflows", table, true)
}

func statusColor(conclusion string) string {
	switch conclusion {
	case "success":
		return "darkgreen"
	case "failure":
		return "red"
	case "cancelled":
		return "gray"
	default:
		return "yellow"
	}
}

func onWorkflowEnter(state *AppState, workflow github.Workflow) {
	jobs, err := workflow.GetJobRuns(state.Repo, state.Token)
	if err != nil {
		fmt.Sprintf("Fehler beim Laden der Jobs: %v", err)
		return
	}
	ShowJobTable(state, jobs)
}

func ShowJobTable(state *AppState, jobs []github.Job) {
	table := createTable([]string{"#", "Name", "Status"})

	for i, job := range jobs {
		color := statusColor(job.Conclusion)
		table.SetCell(i+1, 0, tview.NewTableCell(fmt.Sprintf("%d", i+1)))
		table.SetCell(i+1, 1, tview.NewTableCell(job.Name))
		table.SetCell(i+1, 2, tview.NewTableCell(fmt.Sprintf("[%s]%s", color, job.Conclusion)))
	}

	table.Select(1, 0).SetFixed(1, 0).SetSelectable(true, false)
	table.SetTitle(" GitHub Jobs ").SetBorder(true)

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		row, _ := table.GetSelection()
		switch event.Rune() {
		case 'j':
			if row < table.GetRowCount()-1 {
				table.Select(row+1, 0)
			}
			return nil
			case 'k':
				if row > 1 {
					table.Select(row-1, 0)
				}
				return nil
				case 'b':
					state.Pages.SwitchToPage("Workflows")
					return nil
					case 'q':
						state.App.Stop()
						return nil
		}
		return event
	})

	state.Pages.AddAndSwitchToPage("Jobs", table, true)
}

func createTable(headers []string) *tview.Table {
	table := tview.NewTable()
	table.SetSelectedStyle(tcell.StyleDefault.
		Background(tcell.ColorBlue).
		Foreground(tcell.ColorWhite)).
		SetSelectable(true, false)
	table.SetBackgroundColor(tcell.ColorDefault)
	table.SetBorder(true)
	for i, h := range headers {
		cell := tview.NewTableCell(fmt.Sprintf("[::b]%s", h)).
			SetAlign(tview.AlignCenter).
			SetSelectable(false)
		table.SetCell(0, i, cell)
	}
	return table
}
