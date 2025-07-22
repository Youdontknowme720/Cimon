package ui

import (
	"fmt"
	"github.com/Youdontknowme720/Cimon/github"
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
	list := tview.NewList()

	for _, wf := range workflows.WorkflowRuns {
		label := fmt.Sprintf("%s [%s]", wf.Name, wf.Status)
		desc := fmt.Sprintf("Conclusion: %s", wf.Conclusion)
		list.AddItem(label, desc, 0, func(name string, status string, conclusion string) func() {
			return func() {
			}
		}(wf.Name, wf.Status, wf.Conclusion))
	}

	list.AddItem("Beenden", "Programm schließen", 'q', func() {
		app.Stop()
	})

	list.SetBorder(true).SetTitle(" GitHub Workflows ").SetTitleAlign(tview.AlignLeft)

	return list
}
