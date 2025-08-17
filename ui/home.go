package ui

import (
	"fmt"

	"github.com/Youdontknowme720/Cimonv2/config"
	"github.com/rivo/tview"
)

func (a *App) createHomeScreen(projects []config.GitLabProject) *tview.Table {
	table := newSelectableTable()

	settingsCell := tview.NewTableCell("Settings").
		SetAlign(tview.AlignLeft).
		SetSelectable(true).
		SetReference("Settings")
	table.SetCell(0, 0, settingsCell)

	for i, project := range projects {
		cell := tview.NewTableCell(project.Name).
			SetAlign(tview.AlignLeft).
			SetSelectable(true).
			SetReference(project)
		table.SetCell(i+1, 0, cell)
	}

	table.SetSelectedFunc(func(row, column int) {
		a.handleHomeSelected(row, column, table)
	})

	return table
}

func (a *App) handleHomeSelected(row, column int, table *tview.Table) {
	cell := table.GetCell(row, column)
	ref := cell.GetReference()

	if ref == nil {
		fmt.Println("Keine Referenz – Auswahl ignoriert")
		return
	}

	switch v := ref.(type) {
	case config.GitLabProject:
		page := a.createPipelinePage(v)
		a.pages.AddPage(PagePipeline, page, true, true)
		a.pages.SwitchToPage(PagePipeline)

	case string:
		if v == "Settings" {
			page := a.createSettingsPage()
			a.pages.AddPage(PageSettings, page, true, true)
			a.pages.SwitchToPage(PageSettings)
		}

	default:
		fmt.Println("Unbekannter Typ – Auswahl ignoriert")
	}
}
