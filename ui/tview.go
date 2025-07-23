package ui

import (
	"fmt"
	"log"

	"github.com/Youdontknowme720/Cimon/github"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func StartView(workflows github.WorkflowRunsResponse) {
	app := tview.NewApplication()
	pages := tview.NewPages()

	workflowTable := buildWorkflowTable(app, workflows)
	pages.AddPage("table", workflowTable, true, true)

	if err := app.SetRoot(pages, true).Run(); err != nil {
		log.Fatalf("Fehler beim Starten der App: %v", err)
	}
}

func buildWorkflowTable(app *tview.Application, workflows github.WorkflowRunsResponse) *tview.Table {
	table := tview.NewTable()
	table.SetBackgroundColor(tcell.ColorDefault)
	table.SetBorder(true)

	headers := []string{"#", "Name", "Status"}
	for i, h := range headers {
		cell := tview.NewTableCell(fmt.Sprintf("[::b]%s", h)).
			SetAlign(tview.AlignCenter).
			SetSelectable(false)
		table.SetCell(0, i, cell)
	}

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
