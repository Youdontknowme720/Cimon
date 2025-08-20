package ui

import (
	"fmt"
	"time"

	"github.com/Youdontknowme720/Cimonv2/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Farbschema definieren
var (
	ColorPrimary   = tcell.ColorDodgerBlue
	ColorSecondary = tcell.ColorDarkCyan
	ColorAccent    = tcell.ColorGold
	ColorSuccess   = tcell.ColorGreen
	ColorWarning   = tcell.ColorOrange
	ColorDanger    = tcell.ColorRed
	ColorText      = tcell.ColorWhite
	ColorBorder    = tcell.ColorSteelBlue
	ColorSelected  = tcell.ColorLightSkyBlue
)

func (a *App) createHomeScreen(projects []config.GitLabProject) *tview.Flex {
	// Hauptcontainer mit Border und Titel
	mainContainer := tview.NewFlex().SetDirection(tview.FlexRow)

	// Header mit Titel und Info
	header := a.createHeader()

	// Projekt-Tabelle erstellen
	table := a.createStyledProjectTable(projects)

	// Footer mit Hilfetext
	footer := a.createFooter()

	// Alles zusammenfügen
	mainContainer.
		AddItem(header, 3, 0, false).
		AddItem(table, 0, 1, true).
		AddItem(footer, 2, 0, false)

	return mainContainer
}

func (a *App) createHeader() *tview.TextView {
	header := tview.NewTextView().
		SetText("🚀 [::bu]CIMON v2[::-] - GitLab Pipeline Monitor\n[::d]Wählen Sie ein Projekt oder Settings aus[::-]").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetRegions(true)

	header.SetBorder(true).
		SetBorderColor(ColorPrimary).
		SetTitle(" Welcome ").
		SetTitleAlign(tview.AlignCenter).
		SetTitleColor(ColorAccent)

	return header
}

func (a *App) createStyledProjectTable(projects []config.GitLabProject) *tview.Table {
	table := newSelectableTable()

	// Tabellen-Styling
	table.SetBorder(true).
		SetBorderColor(ColorBorder).
		SetTitle(" 📁 Projekte ").
		SetTitleAlign(tview.AlignLeft).
		SetTitleColor(ColorPrimary)

	// Tabellen-Eigenschaften
	table.SetSelectable(true, false).
		SetSelectedStyle(tcell.StyleDefault.
			Background(ColorSelected).
			Foreground(tcell.ColorBlack).
			Bold(true))

	// Settings-Zeile mit Icon und Styling
	settingsCell := tview.NewTableCell(" ⚙️  Settings").
		SetAlign(tview.AlignLeft).
		SetSelectable(true).
		SetReference("Settings").
		SetTextColor(ColorAccent).
		SetBackgroundColor(tcell.ColorDefault)

	table.SetCell(0, 0, settingsCell)

	// Separator-Zeile
	separatorCell := tview.NewTableCell("─────────────────────").
		SetAlign(tview.AlignLeft).
		SetSelectable(false).
		SetTextColor(ColorBorder)
	table.SetCell(1, 0, separatorCell)

	// Projekt-Zeilen mit Icons und Farben
	for i, project := range projects {
		icon := a.getProjectIcon(project)
		cellText := fmt.Sprintf(" %s  %s", icon, project.Name)

		cell := tview.NewTableCell(cellText).
			SetAlign(tview.AlignLeft).
			SetSelectable(true).
			SetReference(project).
			SetTextColor(ColorText)

		table.SetCell(i+2, 0, cell)
	}

	// Selection Handler
	table.SetSelectedFunc(func(row, column int) {
		a.handleHomeSelected(row, column, table)
	})

	// Input Handler für zusätzliche Tastenkürzel
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			a.app.Stop()
			return nil
		case tcell.KeyCtrlR:
			// Refresh-Funktion hier implementieren
			return nil
		}
		return event
	})

	return table
}

func (a *App) createFooter() *tview.TextView {
	footer := tview.NewTextView().
		SetText("[::b]Tastenkürzel:[::-] [yellow]↑/↓[::-] Navigation | [yellow]Enter[::-] Auswählen | [yellow]Esc[::-] Beenden | [yellow]Ctrl+R[::-] Aktualisieren").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetRegions(true)

	footer.SetBorder(true).
		SetBorderColor(ColorSecondary).
		SetTitle(" Hilfe ").
		SetTitleAlign(tview.AlignCenter).
		SetTitleColor(ColorAccent)

	return footer
}

func (a *App) getProjectIcon(project config.GitLabProject) string {
	// Basierend auf Projektname oder anderen Eigenschaften
	// können Sie verschiedene Icons zurückgeben
	icons := []string{"🔧", "🏗️", "📦", "🌐", "🔬", "📱", "💻", "⚡"}

	// Einfache Hash-basierte Icon-Auswahl
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

// Hilfsfunktion für Benachrichtigungen
func (a *App) showNotification(message string, color tcell.Color) {
	// Modal für kurze Benachrichtigungen
	modal := tview.NewModal().
		SetText(message).
		SetTextColor(color).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.pages.HidePage("notification")
		})

	a.pages.AddPage("notification", modal, false, true)

	// Auto-Hide nach 2 Sekunden (optional)
	go func() {
		time.Sleep(2 * time.Second)
		a.app.QueueUpdateDraw(func() {
			a.pages.HidePage("notification")
		})
	}()
}

// Erweiterte Tabellen-Erstellung mit mehr Optionen

// Zusätzliche Styling-Funktionen
func (a *App) applyGlobalStyling() {
	// Globale App-Styling-Einstellungen
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault
	tview.Styles.ContrastBackgroundColor = tcell.ColorBlue
	tview.Styles.MoreContrastBackgroundColor = tcell.ColorGreen
	tview.Styles.BorderColor = ColorBorder
	tview.Styles.TitleColor = ColorAccent
	tview.Styles.GraphicsColor = ColorSecondary
	tview.Styles.PrimaryTextColor = ColorText
	tview.Styles.SecondaryTextColor = tcell.ColorGray
	tview.Styles.TertiaryTextColor = tcell.ColorDarkGray
	tview.Styles.InverseTextColor = tcell.ColorBlack
	tview.Styles.ContrastSecondaryTextColor = tcell.ColorWhite
}
