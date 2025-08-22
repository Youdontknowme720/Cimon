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

	footer := a.createFooter()

	mainContainer.
		AddItem(header, 5, 0, false).
		AddItem(table, 0, 1, true).
		AddItem(footer, 5, 0, false)

	return mainContainer
}

func (a *App) createHeader() *tview.TextView {
	header := tview.NewTextView().
		SetText("🚀 CIMON v2 - GitLab Pipeline Monitor\n Choose a project or add a new one within the settings").
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

	settingsCell := tview.NewTableCell(" ⚙️  Settings").
		SetAlign(tview.AlignLeft).
		SetSelectable(true).
		SetReference("Settings").
		SetTextColor(tcell.ColorWhite).
		SetBackgroundColor(ColorBlue)

	table.SetCell(0, 0, settingsCell)

	separatorCell := tview.NewTableCell("").
		SetAlign(tview.AlignLeft).
		SetSelectable(false)
	table.SetCell(1, 0, separatorCell)
	table.SetCell(2, 0, separatorCell)

	for i, project := range projects {
		icon := a.getProjectIcon(project)
		cellText := fmt.Sprintf(" %s  %s", icon, project.Name)

		cell := tview.NewTableCell(cellText).
			SetAlign(tview.AlignLeft).
			SetSelectable(true).
			SetReference(project).
			SetTextColor(ColorText)

		table.SetCell(2*i+3, 0, cell)
	}

	table.SetSelectedFunc(func(row, column int) {
		a.handleHomeSelected(row, column, table)
	})

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			a.app.Stop()
			return nil
		case tcell.KeyCtrlR:
			return nil
		}
		return event
	})

	return table
}

func (a *App) createFooter() *tview.TextView {
	footer := tview.NewTextView().
		SetText("Tastenkürzel:↑/↓ Navigation | Enter[::-] Auswählen Esc[::-] Beenden Ctrl+R[::-] Aktualisieren").
		SetTextAlign(tview.AlignCenter).
		SetRegions(true)

	footer.SetBackgroundColor(ColorBlue)

	footer.SetBorder(true).
		SetBorderColor(ColorOrange).
		SetTitle(" Hilfe ").
		SetTitleAlign(tview.AlignCenter).
		SetTitleColor(ColorPink)

	return footer
}

func (a *App) getProjectIcon(project config.GitLabProject) string {
	icons := []string{"🔧", "🏗️", "📦", "🌐", "🔬", "📱", "💻", "⚡"}

	hash := 0
	for _, char := range project.Name {
		hash += int(char)
	}

	return icons[hash%len(icons)]
}

func (a *App) handleHomeSelected(row, column int, table *tview.Table) {
	cell := table.GetCell(row, column)
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
		if v == "Settings" {
			a.showNotification("Öffne Einstellungen...", ColorSuccess)
			page := a.createSettingsPage()
			a.pages.AddPage(PageSettings, page, true, true)
			a.pages.SwitchToPage(PageSettings)
		}

	default:
		a.showNotification("Unbekannter Typ – Auswahl ignoriert", ColorDanger)
	}
}

func (a *App) showNotification(message string, color tcell.Color) {
	modal := tview.NewModal().
		SetText(message).
		SetTextColor(color).
		AddButtons([]string{"OK"}).
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
