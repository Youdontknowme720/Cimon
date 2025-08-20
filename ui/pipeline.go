package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/Youdontknowme720/Cimonv2/config"
	"github.com/Youdontknowme720/Cimonv2/gitlab"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *App) createPipelinePage(proj config.GitLabProject) tview.Primitive {
	container := tview.NewFlex().SetDirection(tview.FlexRow)

	header := a.createPipelineHeader(proj)

	tree := a.handlePipelineClick(fmt.Sprint(proj.ID))

	a.stylePipelineTree(tree, proj)

	footer := a.createPipelineFooter()

	tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'b', 'B':
				a.showNotification("Zurück zur Startseite...", ColorSuccess)
				a.pages.SwitchToPage(PageHome)
				return nil
			case 'r', 'R':
				a.showNotification("Aktualisiere Pipelines...", ColorPrimary)
				a.refreshPipelines(tree, proj)
				return nil
			}
		case tcell.KeyEsc:
			a.pages.SwitchToPage(PageHome)
			return nil
		}
		return event
	})

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		a.handlePipelineSelected(node, proj.ID)
	})

	container.
		AddItem(header, 4, 0, false).
		AddItem(tree, 0, 1, true).
		AddItem(footer, 2, 0, false)

	return container
}

func (a *App) createPipelineHeader(proj config.GitLabProject) *tview.TextView {
	header := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignCenter)

	headerText := fmt.Sprintf(
		"🔧 [::bu]%s[::-] - Pipeline Overview\n[::d]Projekt-ID: %d | Letzte Aktualisierung: %s[::-]",
		proj.Name,
		proj.ID,
		time.Now().Format("15:04:05"),
	)

	header.SetText(headerText)

	header.SetBorder(true)
	header.SetBorderColor(ColorPrimary)
	header.SetTitle(" 🚀 Pipeline Status ")
	header.SetTitleAlign(tview.AlignCenter)
	header.SetTitleColor(ColorAccent)

	return header
}

func (a *App) createPipelineFooter() *tview.TextView {
	footer := tview.NewTextView().
		SetText("[::b]Navigation:[::-] [yellow]↑/↓[::-] Auswählen | [yellow]Enter[::-] Jobs anzeigen | [yellow]B[::-] Zurück | [yellow]R[::-] Aktualisieren | [yellow]Esc[::-] Home").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetRegions(true)

	footer.SetBorder(true)
	footer.SetBorderColor(ColorSecondary)
	footer.SetTitle(" Steuerung ")
	footer.SetTitleAlign(tview.AlignCenter)
	footer.SetTitleColor(ColorAccent)

	return footer
}

func (a *App) stylePipelineTree(tree *tview.TreeView, proj config.GitLabProject) {
	tree.SetBorder(true)
	tree.SetBorderColor(ColorBorder)
	tree.SetTitle(fmt.Sprintf(" 📋 Pipelines für %s ", proj.Name))
	tree.SetTitleAlign(tview.AlignLeft)
	tree.SetTitleColor(ColorPrimary)

	tree.SetGraphics(true)
	tree.SetTopLevel(1)
}

func (a *App) handlePipelineSelected(node *tview.TreeNode, projectID int) {
	ref := node.GetReference()
	if ref == nil {
		a.showNotification("Keine Pipeline-Daten verfügbar", ColorWarning)
		return
	}

	switch v := ref.(type) {
	case gitlab.Pipeline:
		a.showNotification(fmt.Sprintf("Lade Jobs für Pipeline #%d...", v.ID), ColorSuccess)
		page := a.createJobPage(projectID, v.ID)
		a.pages.AddPage("JobPage", page, true, true)
		a.pages.SwitchToPage("JobPage")
	default:
		a.showNotification("Unbekannter Pipeline-Typ", ColorDanger)
	}
}

func (a *App) handlePipelineClick(projectID string) *tview.TreeView {
	// Loading-Anzeige während des Ladens
	loadingNode := tview.NewTreeNode("⏳ Lade Pipelines...").
		SetColor(ColorPrimary).
		SetSelectable(false)

	tree := tview.NewTreeView().
		SetRoot(loadingNode).
		SetCurrentNode(loadingNode)

	pipelines, err := gitlab.GetAllPipelines(projectID, a.token, 3)
	if err != nil {
		errorNode := tview.NewTreeNode("❌ Fehler beim Laden der Pipelines: " + err.Error()).
			SetColor(ColorDanger).
			SetSelectable(false)
		tree.SetRoot(errorNode)
		return tree
	}

	root := tview.NewTreeNode("🔧 Pipelines (" + fmt.Sprint(len(pipelines)) + ")").
		SetColor(ColorPrimary).
		SetExpanded(true).
		SetSelectable(false)

	for i, p := range pipelines {
		pipelineNode := a.createPipelineNode(projectID, p, i+1)
		root.AddChild(pipelineNode)
	}

	tree.SetRoot(root)
	tree.SetCurrentNode(root)

	if len(root.GetChildren()) > 0 {
		tree.SetCurrentNode(root.GetChildren()[0])
	}

	return tree
}

func (a *App) createPipelineNode(projectID string, pipeline gitlab.Pipeline, index int) *tview.TreeNode {
	commitMessage, err := gitlab.GetCommit(projectID, pipeline.Sha, a.token)
	if err != nil {
		commitMessage = &gitlab.Commit{Message: "Unbekannte Commit-Nachricht"}
	}

	statusEmoji := gitlab.StatusEmoji(pipeline.Status)
	nodeColor := a.getStatusColor(pipeline.Status)

	message := commitMessage.Message
	if len(message) > 60 {
		message = message[:57] + "..."
	}

	nodeText := fmt.Sprintf("%s Pipeline #%d: %s",
		statusEmoji,
		pipeline.ID,
		strings.TrimSpace(message))

	if len(pipeline.Sha) >= 8 {
		shortSha := pipeline.Sha[:8]
		nodeText += fmt.Sprintf(" [gray](%s)[white]", shortSha)
	}

	pipelineNode := tview.NewTreeNode(nodeText).
		SetReference(pipeline).
		SetColor(nodeColor).
		SetSelectable(true)

	return pipelineNode
}

func (a *App) getStatusColor(status string) tcell.Color {
	switch strings.ToLower(status) {
	case "success":
		return ColorSuccess
	case "failed":
		return ColorDanger
	case "running":
		return ColorPrimary
	case "pending":
		return ColorWarning
	case "canceled", "cancelled":
		return tcell.ColorGray
	default:
		return ColorText
	}
}

func (a *App) refreshPipelines(tree *tview.TreeView, proj config.GitLabProject) {
	loadingNode := tview.NewTreeNode("⏳ Aktualisiere Pipelines...").
		SetColor(ColorPrimary).
		SetSelectable(false)

	tree.SetRoot(loadingNode)

	go func() {
		time.Sleep(500 * time.Millisecond) // Simulate loading

		a.app.QueueUpdateDraw(func() {
			newTree := a.handlePipelineClick(fmt.Sprint(proj.ID))
			// Tree-Inhalt aktualisieren
			tree.SetRoot(newTree.GetRoot())

			a.showNotification("Pipelines aktualisiert!", ColorSuccess)
		})
	}()
}
