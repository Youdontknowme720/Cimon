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

	table := a.handlePipelineClick(fmt.Sprint(proj.ID))

	a.stylePipelineTable(table, proj)

	footer := a.createPipelineFooter()

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'b', 'B':
				a.showNotification("Zurück zur Startseite...", ColorSuccess)
				a.pages.SwitchToPage(PageHome)
				return nil
			case 'r', 'R':
				a.showNotification("Aktualisiere Pipelines...", ColorPrimary)
				a.refreshPipelines(table, proj)
				return nil
			}
		case tcell.KeyEsc:
			a.pages.SwitchToPage(PageHome)
			return nil
		case tcell.KeyEnter:
			a.handlePipelineSelected(table, proj.ID)
			return nil
		}
		return event
	})

	container.
		AddItem(header, 5, 0, false).
		AddItem(table, 0, 1, true).
		AddItem(footer, 2, 0, false)

	return container
}

func (a *App) stylePipelineTable(table *tview.Table, proj config.GitLabProject) {
	table.SetBorder(true)
	table.SetBorderColor(ColorOrange)
	table.SetTitle(fmt.Sprintf(" 📋 Pipelines für %s ", proj.Name))
	table.SetTitleAlign(tview.AlignLeft)
	table.SetTitleColor(ColorPink)
	table.SetBackgroundColor(ColorBlue)
}

func (a *App) createPipelineHeader(proj config.GitLabProject) *tview.TextView {
	header := tview.NewTextView().
		SetRegions(true).
		SetTextAlign(tview.AlignCenter)

	header.SetBackgroundColor(ColorBlue)

	headerText := fmt.Sprintf(
		"🔧 [::bu]%s[::-] - Pipeline Overview\n[::d]Projekt-ID: %d | Last updated: %s[::-]",
		proj.Name,
		proj.ID,
		time.Now().Format("15:04:05"),
	)

	header.SetText(headerText)

	header.SetBorder(true)
	header.SetBorderColor(ColorOrange)
	header.SetTitle(" 🚀 Pipeline Status ")
	header.SetTitleAlign(tview.AlignCenter)
	header.SetTitleColor(ColorPink)

	return header
}

func (a *App) createPipelineFooter() *tview.TextView {
	footer := tview.NewTextView().
		SetText("Navigation:[::-] ↑/↓ Auswählen | Enter show jobs | B back | R[::-] update | Esc[::-] back to home").
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

func (a *App) handlePipelineSelected(table *tview.Table, projectID int) {
	row, _ := table.GetSelection()
	cell := table.GetCell(row, 0)
	if cell == nil {
		a.showNotification("Keine Pipeline-Daten verfügbar", ColorWarning)
		return
	}

	ref := cell.GetReference()
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

func (a *App) handlePipelineClick(projectID string) *tview.Table {
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)

	table.SetBackgroundColor(ColorBlue)

	loadingCell := tview.NewTableCell("⏳ Lade Pipelines...").
		SetTextColor(ColorPrimary).
		SetSelectable(false)
	table.SetCell(0, 0, loadingCell)

	pipelines, err := gitlab.GetAllPipelines(projectID, a.token, 5)
	if err != nil {
		errorCell := tview.NewTableCell("❌ Fehler beim Laden der Pipelines: " + err.Error()).
			SetTextColor(ColorDanger).
			SetSelectable(false)
		table.SetCell(0, 0, errorCell)
		return table
	}

	table.Clear()

	headerCell := tview.NewTableCell("🔧 Pipelines (" + fmt.Sprint(len(pipelines)) + ")").
		SetTextColor(ColorPink).
		SetSelectable(false).
		SetAttributes(tcell.AttrBold)
	table.SetCell(0, 0, headerCell)

	for i, p := range pipelines {
		pipelineCell := a.createPipelineCell(projectID, p)
		table.SetCell(i+1, 0, pipelineCell)
	}

	if len(pipelines) > 0 {
		table.Select(1, 0) // Erste Pipeline auswählen
	}

	return table
}

func (a *App) createPipelineCell(projectID string, pipeline gitlab.Pipeline) *tview.TableCell {
	commitMessage, err := gitlab.GetCommit(projectID, pipeline.Sha, a.token)
	if err != nil {
		commitMessage = &gitlab.Commit{Message: "Unknown commit message"}
	}

	statusEmoji := gitlab.StatusEmoji(pipeline.Status)
	nodeColor := a.getStatusColor(pipeline.Status)

	message := commitMessage.Message
	if len(message) > 60 {
		message = message[:57] + "..."
	}

	cellText := fmt.Sprintf("%s Pipeline #%d: %s",
		statusEmoji,
		pipeline.ID,
		strings.TrimSpace(message))

	if len(pipeline.Sha) >= 8 {
		shortSha := pipeline.Sha[:8]
		cellText += fmt.Sprintf(" (%s)", shortSha)
	}

	cell := tview.NewTableCell(cellText).
		SetReference(pipeline).
		SetTextColor(nodeColor).
		SetSelectable(true)

	return cell
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

func (a *App) refreshPipelines(table *tview.Table, proj config.GitLabProject) {
	loadingCell := tview.NewTableCell("⏳ Aktualisiere Pipelines...").
		SetTextColor(ColorPrimary).
		SetSelectable(false)

	table.Clear()
	table.SetCell(0, 0, loadingCell)

	go func() {
		time.Sleep(500 * time.Millisecond)

		a.app.QueueUpdateDraw(func() {
			newTable := a.handlePipelineClick(fmt.Sprint(proj.ID))

			// Kopiere Inhalt zur bestehenden Table
			table.Clear()
			for row := 0; row < newTable.GetRowCount(); row++ {
				for col := 0; col < newTable.GetColumnCount(); col++ {
					if cell := newTable.GetCell(row, col); cell != nil {
						table.SetCell(row, col, cell)
					}
				}
			}

			a.showNotification("Pipelines aktualisiert!", ColorSuccess)
		})
	}()
}
