package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/Youdontknowme720/Cimonv2/gitlab"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *App) createJobPage(projectID int, pipelineID int) tview.Primitive {
	container := tview.NewFlex().SetDirection(tview.FlexRow)

	header := a.createJobHeader(projectID, pipelineID)

	tree := a.handleJobClick(fmt.Sprint(projectID), pipelineID)

	a.styleJobTree(tree, projectID, pipelineID)

	footer := a.createJobFooter()

	tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'b', 'B':
				a.showNotification("Zurück zur Pipeline-Ansicht...", ColorSuccess)
				a.pages.SwitchToPage(PagePipeline)
				return nil
			case 'r', 'R':
				a.showNotification("Aktualisiere Jobs...", ColorPrimary)
				a.refreshJobs(tree, projectID, pipelineID)
				return nil
			}
		case tcell.KeyEsc:
			a.pages.SwitchToPage(PagePipeline)
			return nil
		}
		return event
	})

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		a.handleJobSelected(node, projectID, pipelineID)
	})

	container.
		AddItem(header, 4, 0, false).
		AddItem(tree, 0, 1, true).
		AddItem(footer, 2, 0, false)

	return container
}

func (a *App) createJobHeader(projectID int, pipelineID int) *tview.TextView {
	header := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignCenter)

	headerText := fmt.Sprintf(
		"⚙️ [::bu]Jobs für Pipeline #%d[::-]\n[::d]Projekt-ID: %d | Letzte Aktualisierung: %s[::-]",
		pipelineID,
		projectID,
		time.Now().Format("15:04:05"),
	)

	header.SetText(headerText)

	header.SetBorder(true)
	header.SetBorderColor(ColorPrimary)
	header.SetTitle(" 🔨 Job Übersicht ")
	header.SetTitleAlign(tview.AlignCenter)
	header.SetTitleColor(ColorAccent)

	return header
}

func (a *App) createJobFooter() *tview.TextView {
	footer := tview.NewTextView().
		SetText("[::b]Navigation:[::-] [yellow]↑/↓[::-] Auswählen | [yellow]Enter[::-] Job-Details | [yellow]L[::-] Logs | [yellow]B[::-] Zurück | [yellow]R[::-] Aktualisieren | [yellow]Esc[::-] Pipeline").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetRegions(true)

	footer.SetBorder(true)
	footer.SetBorderColor(ColorSecondary)
	footer.SetTitle(" Job-Steuerung ")
	footer.SetTitleAlign(tview.AlignCenter)
	footer.SetTitleColor(ColorAccent)

	return footer
}

func (a *App) styleJobTree(tree *tview.TreeView, projectID int, pipelineID int) {
	tree.SetBorder(true)
	tree.SetBorderColor(ColorBorder)
	tree.SetTitle(fmt.Sprintf(" 📋 Jobs für Pipeline #%d ", pipelineID))
	tree.SetTitleAlign(tview.AlignLeft)
	tree.SetTitleColor(ColorPrimary)

	tree.SetGraphics(true)
	tree.SetTopLevel(1)
}

func (a *App) handleJobClick(projectID string, pipelineID int) *tview.TreeView {
	loadingNode := tview.NewTreeNode("⏳ Lade Jobs...").
		SetColor(ColorPrimary).
		SetSelectable(false)

	tree := tview.NewTreeView().
		SetRoot(loadingNode).
		SetCurrentNode(loadingNode)

	jobs, err := gitlab.GetJobDetails(projectID, pipelineID, a.token)
	if err != nil {
		errorNode := tview.NewTreeNode("❌ Fehler beim Laden der Jobs: " + err.Error()).
			SetColor(ColorDanger).
			SetSelectable(false)
		tree.SetRoot(errorNode)
		return tree
	}

	root := tview.NewTreeNode(fmt.Sprintf("⚙️ Jobs (%d)", len(jobs))).
		SetColor(ColorPrimary).
		SetExpanded(true).
		SetSelectable(false)

	stats := a.calculateJobStats(jobs)

	if len(jobs) > 0 {
		statsText := fmt.Sprintf("📊 Status: %s", a.formatJobStats(stats))
		statsNode := tview.NewTreeNode(statsText).
			SetColor(ColorSecondary).
			SetSelectable(false)
		root.AddChild(statsNode)

		separatorNode := tview.NewTreeNode("─────────────────").
			SetColor(ColorBorder).
			SetSelectable(false)
		root.AddChild(separatorNode)
	}

	for i, job := range jobs {
		jobNode := a.createJobNode(job, i+1)
		root.AddChild(jobNode)
	}

	tree.SetRoot(root)
	tree.SetCurrentNode(root)

	children := root.GetChildren()
	if len(children) > 2 {
		tree.SetCurrentNode(children[2]) // Skip stats und separator
	}

	return tree
}

func (a *App) createJobNode(job gitlab.Job, index int) *tview.TreeNode {
	statusEmoji := gitlab.StatusEmoji(job.Status)
	nodeColor := a.getJobStatusColor(job.Status)

	jobName := job.Name
	if len(jobName) > 50 {
		jobName = jobName[:47] + "..."
	}

	durationText := ""
	if job.Duration > 0 {
		duration := time.Duration(job.Duration) * time.Second
		durationText = fmt.Sprintf(" [gray](%v)[white]", duration.Round(time.Second))
	}

	nodeText := fmt.Sprintf("%s %s%s",
		statusEmoji,
		jobName,
		durationText)

	if job.Stage != "" {
		nodeText += fmt.Sprintf(" [darkgray][%s][white]", job.Stage)
	}

	jobNode := tview.NewTreeNode(nodeText).
		SetReference(job).
		SetColor(nodeColor).
		SetSelectable(true)

	return jobNode
}

func (a *App) getJobStatusColor(status string) tcell.Color {
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
	case "skipped":
		return tcell.ColorDarkGray
	case "manual":
		return tcell.ColorPurple
	default:
		return ColorText
	}
}

func (a *App) calculateJobStats(jobs []gitlab.Job) map[string]int {
	stats := make(map[string]int)
	for _, job := range jobs {
		stats[strings.ToLower(job.Status)]++
	}
	return stats
}

func (a *App) formatJobStats(stats map[string]int) string {
	parts := []string{}

	if count, ok := stats["success"]; ok && count > 0 {
		parts = append(parts, fmt.Sprintf("✅%d", count))
	}
	if count, ok := stats["failed"]; ok && count > 0 {
		parts = append(parts, fmt.Sprintf("❌%d", count))
	}
	if count, ok := stats["running"]; ok && count > 0 {
		parts = append(parts, fmt.Sprintf("🔄%d", count))
	}
	if count, ok := stats["pending"]; ok && count > 0 {
		parts = append(parts, fmt.Sprintf("⏳%d", count))
	}
	if count, ok := stats["canceled"]; ok && count > 0 {
		parts = append(parts, fmt.Sprintf("🚫%d", count))
	}

	if len(parts) == 0 {
		return "Keine Jobs"
	}

	return strings.Join(parts, " | ")
}

func (a *App) handleJobSelected(node *tview.TreeNode, projectID int, pipelineID int) {
	ref := node.GetReference()
	if ref == nil {
		a.showNotification("Keine Job-Daten verfügbar", ColorWarning)
		return
	}

	switch job := ref.(type) {
	case gitlab.Job:
		// Job-Details Modal anzeigen
		a.showJobDetailsModal(job, projectID, pipelineID)
	default:
		a.showNotification("Unbekannter Job-Typ", ColorDanger)
	}
}

func (a *App) showJobDetailsModal(job gitlab.Job, projectID int, pipelineID int) {
	details := fmt.Sprintf(
		"Job: %s\n"+
			"Status: %s %s\n"+
			"Stage: %s\n"+
			"ID: %d\n",
		job.Name,
		gitlab.StatusEmoji(job.Status), job.Status,
		job.Stage,
		job.ID)

	if job.Duration > 0 {
		duration := time.Duration(job.Duration) * time.Second
		details += fmt.Sprintf("Dauer: %v\n", duration.Round(time.Second))
	}
}

func (a *App) refreshJobs(tree *tview.TreeView, projectID int, pipelineID int) {
	loadingNode := tview.NewTreeNode("⏳ Aktualisiere Jobs...").
		SetColor(ColorPrimary).
		SetSelectable(false)

	tree.SetRoot(loadingNode)

	go func() {
		time.Sleep(500 * time.Millisecond)

		a.app.QueueUpdateDraw(func() {
			newTree := a.handleJobClick(fmt.Sprint(projectID), pipelineID)
			tree.SetRoot(newTree.GetRoot())

			a.showNotification("Jobs aktualisiert!", ColorSuccess)
		})
	}()
}
