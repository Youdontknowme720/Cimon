package ui

import(
	"fmt"
	"github.com/rivo/tview"
	"github.com/Youdontknowme720/Cimon/gitlab"
)

type PipelineNode struct {
	Pipe utils.Pipeline
	ID int
	Status string
	Created string
	WebUrl string
}

func InitTree(app *tview.Application, projectID string, accessToken string) {
	root := tview.NewTreeNode("root")
	pipelines, _ := utils.GetPipelineStatus(projectID, accessToken)
	for _, pipeline := range pipelines{
		nodeName := fmt.Sprintf("Pipeline: %d", pipeline.ID)
		node := tview.NewTreeNode(nodeName).
			SetReference(PipelineNode{pipeline,
				pipeline.ID,
				pipeline.Status,
				pipeline.Created,
				pipeline.WebURL}).
			SetSelectable(true)
		root.AddChild(node)
	}
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)
	tree.SetSelectedFunc(func(node *tview.TreeNode){
		ref := node.GetReference()
		if len(node.GetChildren()) == 0{
			if treeNode, ok := ref.(PipelineNode); ok{
			failed := treeNode.Pipe.IsFailed()
			title := fmt.Sprintf("Has failed :%v", failed)
			newNode := tview.NewTreeNode(title).
				SetReference("Lel")
			node.AddChild(newNode)
		}
		}
	})
	app.SetRoot(tree, true)
	app.SetFocus(tree)
	app.Run()
}

