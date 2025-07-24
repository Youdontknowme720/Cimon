package ui

import (
	"fmt"
	"log"

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

func StartView(workflows github.WorkflowRunsResponse, repo string, token string) {
	app := tview.NewApplication()
	pages := tview.NewPages()

	workflowTable := buildWorkflowTable(app, workflows, pages, repo, token)
	pages.AddPage("table", workflowTable, true, true)

	if err := app.SetRoot(pages, true).Run(); err != nil {
		log.Fatalf("Fehler beim Starten der App: %v", err)
	}
}

func buildWorkflowTable(app *tview.Application,
		workflows github.WorkflowRunsResponse,
		pages *tview.Pages,
		repo string,
		token string) *tview.Table {

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
			app.Stop()
			return nil
		}
		if event.Key() == tcell.KeyEnter {
			if row > 0 && row-1 < len(workflows.WorkflowRuns) {
				row_idx, _ := table.GetSelection()
				selectedWorkflow := workflows.WorkflowRuns[row_idx-1]
				jobs, err := selectedWorkflow.GetJobRuns(repo, token)
				if err != nil {
					modal := tview.NewModal().
						SetText(fmt.Sprintf("Fehler beim Laden der Jobs:\n\n%v selected workflow %+v", err, selectedWorkflow)).
						AddButtons([]string{"Zurück"}).
						SetDoneFunc(func(buttonIndex int, buttonLabel string) {
							pages.SwitchToPage("table")
						})
					pages.AddPage("error", modal, true, true)
					return nil
				}

				jobTable := BuildJobTable(jobs)
				pages.AddAndSwitchToPage("Jobs", jobTable, true)
				jobTable.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey{
					if ev.Rune() == 'b' {
						pages.SwitchToPage("table")
						return nil
					}
					return ev
				})
			}
			return nil
		}
		return event
	})
	table.SetTitle(" GitHub Workflows ").SetBorder(true)
	return table
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

func BuildJobTable(jobs []github.Job) *tview.Table{
	table := createTable([]string{"#", "Name", "Status"})
	for i, job := range jobs {
		color := statusColor(job.Conclusion)
		table.SetCell(i+1, 0, tview.NewTableCell(fmt.Sprintf("%-3d", i+1)))
		table.SetCell(i+1, 1, tview.NewTableCell(fmt.Sprintf("%-10s", job.Name)))
		table.SetCell(i+1, 2, tview.NewTableCell(fmt.Sprintf("[%-s]%-10s", color, job.Conclusion)))
	}

	if len(jobs) > 0 {
		table.Select(1, 0).SetFixed(1, 0).SetSelectable(true, false)
	}

	table.SetTitle(" GitHub Jobs ").SetBorder(true)
	return table
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
