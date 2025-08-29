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

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'b', 'B':
				a.showNotification("Zurück zur Startseite...", ColorSuccess)
				a.pages.SwitchToPage(PageHome)
				return nil
			case 'r', 'R':
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
		AddItem(header, 3, 0, false).
		AddItem(table, 0, 1, true)

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

	table.SetBorder(true)
	table.SetBorderColor(ColorOrange)
	table.SetBackgroundColor(ColorBlue)

	loadingCell := tview.NewTableCell("⏳ Lade Pipelines...").
		SetTextColor(ColorPrimary).
		SetSelectable(false)
	table.SetCell(0, 0, loadingCell)

	go func() {
		pipelines, err := gitlab.GetAllPipelines(projectID, a.token, 5)

		a.app.QueueUpdateDraw(func() {
			table.Clear()

			if err != nil {
				errorCell := tview.NewTableCell("❌ Fehler beim Laden der Pipelines: " + err.Error()).
					SetTextColor(ColorDanger).
					SetSelectable(false)
				table.SetCell(0, 0, errorCell)
				return
			}

			headerCell := tview.NewTableCell("🔧 Pipelines (" + fmt.Sprint(len(pipelines)) + ")").
				SetTextColor(tcell.ColorWhite).
				SetSelectable(false).
				SetAttributes(tcell.AttrBold)
			table.SetCell(0, 0, headerCell)

			for i, p := range pipelines {
				pipelineCell := a.createPipelineCell(projectID, p)
				table.SetCell(i+1, 0, pipelineCell)
			}

			if len(pipelines) > 0 {
				table.Select(1, 0)
			}
		})
	}()

	return table
}

func (a *App) createPipelineCell(projectID string, pipeline gitlab.Pipeline) *tview.TableCell {
	statusEmoji := gitlab.StatusEmoji(pipeline.Status)

	placeholderText := fmt.Sprintf("%s Pipeline : ⏳ Lade Commit...", statusEmoji)

	if len(pipeline.Sha) >= 8 {
		shortSha := pipeline.Sha[:8]
		placeholderText += fmt.Sprintf(" (%s)", shortSha)
	}

	cell := tview.NewTableCell(placeholderText).
		SetReference(pipeline).
		SetTextColor(tcell.ColorWhite).
		SetSelectedStyle(tcell.StyleDefault.
			Background(ColorBlue).
			Foreground(ColorPink).
			Bold(true))

	go func(cell *tview.TableCell, sha string) {
		commit, err := gitlab.GetCommit(projectID, sha, a.token)
		message := "Unknown commit message"
		if err == nil {
			message = commit.Message
			if len(message) > 60 {
				message = message[:57] + "..."
			}
		}

		a.app.QueueUpdateDraw(func() {
			newText := fmt.Sprintf("%s Pipeline : %s", statusEmoji, strings.TrimSpace(message))
			if len(sha) >= 8 {
				shortSha := sha[:8]
				newText += fmt.Sprintf(" (%s)", shortSha)
			}
			cell.SetText(newText)
		})
	}(cell, pipeline.Sha)

	return cell
}

func (a *App) refreshPipelines(table *tview.Table, proj config.GitLabProject) {
	loadingCell := tview.NewTableCell("⏳ Aktualisiere Pipelines...").
		SetTextColor(ColorPrimary).
		SetSelectable(false)

	table.Clear()
	table.SetCell(0, 0, loadingCell)

	go func() {
		pipelines, err := gitlab.GetAllPipelines(fmt.Sprint(proj.ID), a.token, 5)

		a.app.QueueUpdateDraw(func() {
			table.Clear()

			if err != nil {
				errorCell := tview.NewTableCell("❌ Fehler beim Laden der Pipelines: " + err.Error()).
					SetTextColor(ColorDanger).
					SetSelectable(false)
				table.SetCell(0, 0, errorCell)
				return
			}

			headerCell := tview.NewTableCell("🔧 Pipelines (" + fmt.Sprint(len(pipelines)) + ")").
				SetTextColor(ColorPink).
				SetSelectable(false).
				SetAttributes(tcell.AttrBold)
			table.SetCell(0, 0, headerCell)

			for i, p := range pipelines {
				pipelineCell := a.createPipelineCell(fmt.Sprint(proj.ID), p)
				table.SetCell(i+1, 0, pipelineCell)
			}

			if len(pipelines) > 0 {
				table.Select(1, 0)
			}
		})
	}()
}
