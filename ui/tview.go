package ui

import (
	"fmt"
	"github.com/Youdontknowme720/Cimon/github"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
)
var app *tview.Application
var pages *tview.Pages
func StartView(workflows github.WorkflowRunsResponse) {
	app := tview.NewApplication()
	pages := tview.NewPages()
	workflowList := buildWorkflowList(workflows)
	pages.AddPage("list", workflowList, true, true)

	if err := app.SetRoot(pages, true).Run(); err != nil {
		log.Fatalf("Run failed: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("Run failed: %v", err)
	}
}


func buildWorkflowList(workflows github.WorkflowRunsResponse) *tview.List {
	list := tview.NewList().
		SetSelectedTextColor(tcell.ColorDarkOrange)
	for _, wf := range workflows.WorkflowRuns {
		var conclusionColor = "red"
		label := fmt.Sprintf("%s", wf.DisplayTitle)
		if wf.Conclusion == "success" {
			conclusionColor = "darkgreen"
		}
		desc := fmt.Sprintf("\t [%s]%s[-]",conclusionColor, wf.Conclusion)
		list.AddItem(label, desc, 0,nil)
	}

	list.AddItem("Beenden", "Programm schließen", 'q', func() {
		app.Stop()
	})

	list.SetBorder(true).SetTitle(" GitHub Workflows ").SetTitleAlign(tview.AlignLeft)

	return list
}
