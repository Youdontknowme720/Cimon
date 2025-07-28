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

type JobNode struct{
	Job utils.Job
	Name string
	Status string
}

func InitTree(app *tview.Application, projectID string, accessToken string) {
	root := tview.NewTreeNode("root")
	pipelines, _ := utils.GetPipelineStatus(projectID, accessToken)
	for _, pipeline := range pipelines{
		node := buildPipelineNode(pipeline)
		root.AddChild(node)
	}
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)
	tree.SetSelectedFunc(func(node *tview.TreeNode){
		ref := node.GetReference()
		if len(node.GetChildren()) == 0{
			if treeNode, ok := ref.(PipelineNode); ok{
				jobs, _ := utils.GetJobDetails(projectID, treeNode.ID, accessToken)
				for _, job := range jobs{
					nodeName := fmt.Sprintf("Job %s", job.Name)
					newNode := tview.NewTreeNode(nodeName).
						SetReference(JobNode{job,
							job.Name,
							job.Status})
					node.AddChild(newNode)
				}
		}
		}
	})
	app.SetRoot(tree, true)
	app.SetFocus(tree)
	app.Run()
}

func buildPipelineNode(pipe utils.Pipeline) *tview.TreeNode{
	nodeName := fmt.Sprintf("Pipeline: %d", pipe.ID)
	node := tview.NewTreeNode(nodeName).
		SetReference(PipelineNode{pipe,
			pipe.ID,
			pipe.Status,
			pipe.Created,
			pipe.WebURL}).
		SetSelectable(true)
	return node
}
