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
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		a.handlePipelineSelected(node, proj.ID)
	})
	return tree
}

func (a *App) handlePipelineSelected(node *tview.TreeNode, projectID int) {
	ref := node.GetReference()

	if ref == nil {
		return
	}

	switch v := ref.(type) {
	case gitlab.Pipeline:
		page := a.createJobPage(projectID, v.ID)
		a.pages.AddPage("JobPage", page, true, true)
		a.pages.SwitchToPage("JobPage")
	}
}

func (a *App) handlePipelineClick(projectID string) *tview.TreeView {
	pipelines, err := gitlab.GetAllPipelines(projectID, a.token, 3)
	if err != nil {
		panic(err)
	}
	root := tview.NewTreeNode("Pipelines").
		SetColor(tcell.ColorDarkOrange)

	for _, p := range pipelines {
		commitMessage, err := gitlab.GetCommit(projectID, p.Sha, a.token)
		if err != nil {
			panic(err)
		}

		statusText := gitlab.StatusEmoji(p.Status)

		pipelineNode := tview.NewTreeNode(statusText + " " + commitMessage.Message).
			SetReference(p).
			SetSelectable(true)
		root.AddChild(pipelineNode)
	}

	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	return tree
}
