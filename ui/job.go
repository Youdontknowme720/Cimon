package ui

import (
	"fmt"

	"github.com/Youdontknowme720/Cimonv2/gitlab"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *App) createJobPage(projectID int, pipelineID int) tview.Primitive {
	tree := a.handleJobClick(fmt.Sprint(projectID), pipelineID)
	tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == 'b' {
			a.pages.SwitchToPage(PagePipeline)
			return nil
		}
		return event
	})
	return tree
}

func (a *App) handleJobClick(projectID string, pipelineID int) *tview.TreeView {
	jobs, err := gitlab.GetJobDetails(projectID, pipelineID, a.token)
	root := tview.NewTreeNode("Jobs").SetColor(tcell.ColorDarkOrange)

	if err != nil {
		root.SetText(fmt.Sprintf("Fehler: %v", err))
		return tview.NewTreeView().SetRoot(root).SetCurrentNode(root)
	}

	for _, j := range jobs {
		statusText := gitlab.StatusEmoji(j.Status)
		jobNode := tview.NewTreeNode(statusText + " " + j.Name).
			SetReference(j).
			SetSelectable(true)
		root.AddChild(jobNode)
	}

	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	return tree
}
