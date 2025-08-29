package ui

import (
	"fmt"
	"time"

	"github.com/Youdontknowme720/Cimonv2/gitlab"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *App) createJobPage(projectID int, pipelineID int) tview.Primitive {
	container := tview.NewFlex().SetDirection(tview.FlexRow)

	header := a.createJobHeader(projectID, pipelineID)

	table := a.handleJobClick(fmt.Sprint(projectID), pipelineID)
	a.styleJobTable(table, pipelineID)

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'b', 'B':
				a.showNotification("Zurück zur Pipeline-Ansicht...", ColorSuccess)
				a.pages.SwitchToPage(PagePipeline)
				return nil
			case 'r', 'R':
				a.refreshJobs(table, projectID, pipelineID)
				return nil
			}
		case tcell.KeyEsc:
			a.pages.SwitchToPage(PagePipeline)
			return nil
		}
		return event
	})

	table.SetSelectedFunc(func(row, column int) {
		cell := table.GetCell(row, column)
		if cell != nil {
			ref := cell.GetReference()
			if job, ok := ref.(gitlab.Job); ok {
				a.showJobDetailsModal(job, projectID, pipelineID)
			}
		}
	})

	container.
		AddItem(header, 4, 0, false).
		AddItem(table, 0, 1, true)

	return container
}

func (a *App) showJobDetailsModal(job gitlab.Job, projectID int, pipelineID int) {
	panic("unimplemented")
}

func (a *App) createJobHeader(projectID int, pipelineID int) *tview.TextView {
	header := tview.NewTextView().
		SetRegions(true).
		SetTextAlign(tview.AlignCenter)
	header.SetBackgroundColor(ColorBlue)

	headerText := fmt.Sprintf(
		"⚙️ [::bu]Jobs for pipeline #%d[::-]\n[::d]Projekt-ID: %d | last updated: %s[::-]",
		pipelineID,
		projectID,
		time.Now().Format("15:04:05"),
	)

	header.SetText(headerText)

	header.SetBorder(true)
	header.SetBorderColor(ColorOrange)
	header.SetTitle(" 🔨 Job Übersicht ")
	header.SetTitleAlign(tview.AlignCenter)
	header.SetTitleColor(ColorPink)

	return header
}

func (a *App) styleJobTable(table *tview.Table, pipelineID int) {
	table.SetBorder(true)
	table.SetBorderColor(ColorOrange)
	table.SetTitle(fmt.Sprintf(" 📋 Jobs for Pipeline #%d ", pipelineID))
	table.SetTitleAlign(tview.AlignLeft)
	table.SetTitleColor(ColorPink)
	table.SetBackgroundColor(ColorBlue)
}

func (a *App) handleJobClick(projectID string, pipelineID int) *tview.Table {
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)

	table.SetBackgroundColor(ColorBlue)

	loadingCell := tview.NewTableCell("⏳ Lade Jobs...").
		SetTextColor(tcell.ColorWhite).
		SetSelectable(false)
	table.SetCell(0, 0, loadingCell)

	go func() {
		jobs, err := gitlab.GetJobDetails(projectID, pipelineID, a.token)

		a.app.QueueUpdateDraw(func() {
			table.Clear()

			if err != nil {
				errorCell := tview.NewTableCell("❌ Fehler beim Laden der Jobs: " + err.Error()).
					SetTextColor(ColorDanger).
					SetSelectable(false)
				table.SetCell(0, 0, errorCell)
				return
			}

			headerCell := tview.NewTableCell(fmt.Sprintf("Jobs (%d)", len(jobs))).
				SetTextColor(tcell.ColorWhite).
				SetSelectable(false).
				SetAttributes(tcell.AttrBold)
			table.SetCell(0, 0, headerCell)

			for i, job := range jobs {
				cell := a.createJobTableCell(job)
				table.SetCell(i+1, 0, cell)
			}
		})
	}()

	return table
}

func (a *App) createJobTableCell(job gitlab.Job) *tview.TableCell {
	statusEmoji := gitlab.StatusEmoji(job.Status)

	cellText := fmt.Sprintf("%s %s ⏳ Lade Details...", statusEmoji, job.Name)
	if job.Duration > 0 {
		duration := time.Duration(job.Duration) * time.Second
		cellText += fmt.Sprintf(" [gray](%v)[white]", duration.Round(time.Second))
	}
	if job.Stage != "" {
		cellText += fmt.Sprintf(" [darkgray][%s][white]", job.Stage)
	}

	cell := tview.NewTableCell(cellText).
		SetReference(job).
		SetTextColor(tcell.ColorWhite).
		SetSelectable(true).
		SetSelectedStyle(tcell.StyleDefault.
			Background(ColorBlue).
			Foreground(ColorPink).
			Bold(true))

	go func(cell *tview.TableCell, job gitlab.Job) {
		a.app.QueueUpdateDraw(func() {
			newText := fmt.Sprintf("%s %s", statusEmoji, job.Name)
			if job.Duration > 0 {
				duration := time.Duration(job.Duration) * time.Second
				newText += fmt.Sprintf(" [gray](%v)[white]", duration.Round(time.Second))
			}
			if job.Stage != "" {
				newText += fmt.Sprintf(" [darkgray][%s][white]", job.Stage)
			}
			cell.SetText(newText)
		})
	}(cell, job)

	return cell
}

func (a *App) refreshJobs(table *tview.Table, projectID int, pipelineID int) {
	table.Clear()
	loadingCell := tview.NewTableCell("⏳ Aktualisiere Jobs...").
		SetTextColor(tcell.ColorWhite).
		SetSelectable(false)
	table.SetCell(0, 0, loadingCell)

	go func() {
		time.Sleep(300 * time.Millisecond)

		jobs, err := gitlab.GetJobDetails(fmt.Sprint(projectID), pipelineID, a.token)

		a.app.QueueUpdateDraw(func() {
			table.Clear()

			if err != nil {
				errorCell := tview.NewTableCell("❌ Fehler beim Laden der Jobs: " + err.Error()).
					SetTextColor(ColorDanger).
					SetSelectable(false)
				table.SetCell(0, 0, errorCell)
				return
			}

			headerCell := tview.NewTableCell(fmt.Sprintf("⚙️ Jobs (%d)", len(jobs))).
				SetTextColor(ColorPrimary).
				SetSelectable(false).
				SetAttributes(tcell.AttrBold)
			table.SetCell(0, 0, headerCell)

			for i, job := range jobs {
				cell := a.createJobTableCell(job)
				table.SetCell(i+1, 0, cell)
			}
		})
	}()
}
