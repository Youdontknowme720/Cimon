package ui

import (
	"fmt"
	"time"

	"github.com/Youdontknowme720/Cimonv2/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	ColorBlue      = tcell.NewRGBColor(13, 17, 100)
	ColorPurple    = tcell.NewRGBColor(100, 13, 95)
	ColorPink      = tcell.NewRGBColor(234, 34, 100)
	ColorOrange    = tcell.NewRGBColor(247, 141, 96)
	ColorPrimary   = tcell.ColorDodgerBlue
	ColorSecondary = tcell.ColorDarkCyan
	ColorAccent    = tcell.ColorGold
	ColorSuccess   = tcell.ColorDarkGreen
	ColorWarning   = tcell.ColorOrange
	ColorDanger    = tcell.ColorDarkRed
	ColorText      = tcell.ColorWhite
	ColorBorder    = tcell.NewRGBColor(234, 34, 100)
	ColorSelected  = tcell.ColorLightSkyBlue
)

func (a *App) createHomeScreen(projects []config.GitLabProject) *tview.Flex {
	mainContainer := tview.NewFlex().SetDirection(tview.FlexRow)

	header := a.createHeader()

	table := a.createStyledProjectTable(projects)

	mainContainer.
		AddItem(header, 3, 0, false).
		AddItem(table, 0, 2, true)

	return mainContainer
}

func (a *App) createHeader() *tview.TextView {
	header := tview.NewTextView().
		SetText("CIMON v2 - GitLab Pipeline Monitor\n Choose a project or add a new one within the settings").
		SetTextAlign(tview.AlignCenter).
		SetRegions(true)

	header.SetBackgroundColor(ColorBlue)
	header.SetBorder(true).
		SetBorderColor(ColorOrange).
		SetTitle(" Welcome ").
		SetTitleAlign(tview.AlignCenter).
		SetTitleColor(ColorPink)

	return header
}

func (a *App) createStyledProjectTable(projects []config.GitLabProject) *tview.Table {
	table := newSelectableTable()

	table.SetBorder(true).
		SetBorderColor(ColorOrange).
		SetTitle(" HOME ").
		SetTitleAlign(tview.AlignCenter).
		SetTitleColor(ColorPink).
		SetBackgroundColor(ColorBlue)

	table.SetSelectable(true, false).
		SetSelectedStyle(tcell.StyleDefault.
			Background(ColorBlue).
			Foreground(ColorPink).
			Bold(true))

	settingsButtons := a.createSettingsButtons()
	for pos, cell := range settingsButtons {
		table.SetCell(pos[0], pos[1], cell)
	}

	table.SetCell(2, 0, tview.NewTableCell("").SetSelectable(false))

	for i, project := range projects {
		cellText := fmt.Sprintf("- %s ", project.Name)

		cell := tview.NewTableCell(cellText).
			SetAlign(tview.AlignLeft).
			SetSelectable(true).
			SetReference(project).
			SetTextColor(ColorText)

		table.SetCell(3+i, 0, cell)
	}

	table.SetSelectedFunc(func(row, column int) {
		a.handleHomeSelected(row, column, table)
	})

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			a.app.Stop()
			return nil
		}
		return event
	})

	return table
}

func (a *App) createSettingsButtons() map[[2]int]*tview.TableCell {
	buttons := make(map[[2]int]*tview.TableCell)

	addProjectCell := tview.NewTableCell("+ Add Project").
		SetAlign(tview.AlignCenter).
		SetSelectable(true).
		SetReference("AddProj").
		SetTextColor(tcell.ColorWhite).
		SetBackgroundColor(ColorBlue)

	addTokenCell := tview.NewTableCell("+ Add Token").
		SetAlign(tview.AlignCenter).
		SetSelectable(true).
		SetReference("AddToken").
		SetTextColor(tcell.ColorWhite).
		SetBackgroundColor(ColorBlue)

	buttons[[2]int{0, 0}] = addProjectCell
	buttons[[2]int{1, 0}] = addTokenCell

	return buttons
}

func (a *App) handleHomeSelected(row, column int, table *tview.Table) {
	cell := table.GetCell(row, column)
	if cell == nil {
		a.showNotification("Keine Zelle ausgewählt", ColorWarning)
		return
	}

	ref := cell.GetReference()
	if ref == nil {
		a.showNotification("Keine Referenz – Auswahl ignoriert", ColorWarning)
		return
	}

	switch v := ref.(type) {
	case config.GitLabProject:
		a.showNotification(fmt.Sprintf("Lade Pipeline für %s...", v.Name), ColorSuccess)
		page := a.createPipelinePage(v)
		a.pages.AddPage(PagePipeline, page, true, true)
		a.pages.SwitchToPage(PagePipeline)

	case string:
		switch v {
		case "AddProj":
			a.handleAddingProject()
		case "AddToken":
			a.handleAddingToken()
		default:
			a.showNotification("Unbekannter Button – Auswahl ignoriert", ColorDanger)
		}

	default:
		a.showNotification("Unbekannter Typ – Auswahl ignoriert", ColorDanger)
	}
}

func (a *App) showNotification(message string, color tcell.Color) {
	modal := tview.NewModal().
		SetText(message).
		SetTextColor(color).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.pages.HidePage("notification")
		})

	a.pages.AddPage("notification", modal, false, true)

	go func() {
		time.Sleep(2 * time.Second)
		a.app.QueueUpdateDraw(func() {
			a.pages.HidePage("notification")
		})
	}()
}
