package ui

import (
	"fmt"

	"github.com/Youdontknowme720/Cimonv2/config"
	"github.com/Youdontknowme720/Cimonv2/gitlab"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *App) createPipelinePage(proj config.GitLabProject) tview.Primitive {
	tree := a.handlePipelineClick(fmt.Sprint(proj.ID))
	tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == 'b' {
			a.pages.SwitchToPage(PageHome)
			return nil
		}
		return event
	})
	return tree
}

func (a *App) handlePipelineClick(projectID string) *tview.TreeView {
	pipelines, err := gitlab.GetAllPipelines(projectID, a.token, 3)
	if err != nil {
		panic(err)
	}
	root := tview.NewTreeNode("Pipelines").
		SetColor(tcell.ColorGreen)

	for _, p := range pipelines {
		pipelineNode := tview.NewTreeNode(fmt.Sprint(p.ID)).
			SetReference(p).
			SetSelectable(true)
		root.AddChild(pipelineNode)
	}

	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	return tree
}
