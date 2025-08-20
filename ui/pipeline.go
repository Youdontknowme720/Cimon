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
	// Hauptcontainer für die Pipeline-Seite
	container := tview.NewFlex().SetDirection(tview.FlexRow)

	// Header für die Pipeline-Seite
	header := a.createPipelineHeader(proj)

	// Pipeline-Tree
	tree := a.handlePipelineClick(fmt.Sprint(proj.ID))

	// Tree-Styling anwenden
	a.stylePipelineTree(tree, proj)

	// Footer für Pipeline-Seite
	footer := a.createPipelineFooter()

	// Input-Handling
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

	// Selection-Handler
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		a.handlePipelineSelected(node, proj.ID)
	})

	// Container zusammenbauen
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

	// Header-Text mit Projekt-Info
	headerText := fmt.Sprintf(
		"🔧 [::bu]%s[::-] - Pipeline Overview\n[::d]Projekt-ID: %d | Letzte Aktualisierung: %s[::-]",
		proj.Name,
		proj.ID,
		time.Now().Format("15:04:05"),
	)

	header.SetText(headerText)

	// Styling
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
	// Tree-Styling - jede Methode einzeln, da nicht alle chainable sind
	tree.SetBorder(true)
	tree.SetBorderColor(ColorBorder)
	tree.SetTitle(fmt.Sprintf(" 📋 Pipelines für %s ", proj.Name))
	tree.SetTitleAlign(tview.AlignLeft)
	tree.SetTitleColor(ColorPrimary)

	// TreeView-spezifische Einstellungen
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

	// Pipeline-Daten laden
	pipelines, err := gitlab.GetAllPipelines(projectID, a.token, 3)
	if err != nil {
		errorNode := tview.NewTreeNode("❌ Fehler beim Laden der Pipelines: " + err.Error()).
			SetColor(ColorDanger).
			SetSelectable(false)
		tree.SetRoot(errorNode)
		return tree
	}

	// Haupt-Root-Node erstellen
	root := tview.NewTreeNode("🔧 Pipelines (" + fmt.Sprint(len(pipelines)) + ")").
		SetColor(ColorPrimary).
		SetExpanded(true).
		SetSelectable(false)

	// Pipeline-Nodes hinzufügen
	for i, p := range pipelines {
		pipelineNode := a.createPipelineNode(projectID, p, i+1)
		root.AddChild(pipelineNode)
	}

	tree.SetRoot(root)
	tree.SetCurrentNode(root)

	// Erstes Kind automatisch auswählen, falls vorhanden
	if len(root.GetChildren()) > 0 {
		tree.SetCurrentNode(root.GetChildren()[0])
	}

	return tree
}

func (a *App) createPipelineNode(projectID string, pipeline gitlab.Pipeline, index int) *tview.TreeNode {
	// Commit-Message holen
	commitMessage, err := gitlab.GetCommit(projectID, pipeline.Sha, a.token)
	if err != nil {
		commitMessage = &gitlab.Commit{Message: "Unbekannte Commit-Nachricht"}
	}

	// Status-Emoji und Farbe
	statusEmoji := gitlab.StatusEmoji(pipeline.Status)
	nodeColor := a.getStatusColor(pipeline.Status)

	// Commit-Message kürzen falls zu lang
	message := commitMessage.Message
	if len(message) > 60 {
		message = message[:57] + "..."
	}

	// Node-Text formatieren
	nodeText := fmt.Sprintf("%s Pipeline #%d: %s",
		statusEmoji,
		pipeline.ID,
		strings.TrimSpace(message))

	// Zusätzliche Pipeline-Info (SHA anzeigen)
	if len(pipeline.Sha) >= 8 {
		shortSha := pipeline.Sha[:8]
		nodeText += fmt.Sprintf(" [gray](%s)[white]", shortSha)
	}

	// Node erstellen
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
	// Loading-State setzen
	loadingNode := tview.NewTreeNode("⏳ Aktualisiere Pipelines...").
		SetColor(ColorPrimary).
		SetSelectable(false)

	tree.SetRoot(loadingNode)

	// Async refresh (in einer echten Implementierung mit Goroutine)
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
