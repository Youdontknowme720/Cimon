package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// newSelectableTable erstellt eine Tabelle mit einheitlichem Style,
// die in Home/Settings/etc. verwendet werden kann.
func newSelectableTable() *tview.Table {
	table := tview.NewTable()
	table.SetSelectedStyle(
		tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite),
	)
	table.SetSelectable(true, false)
	table.SetBackgroundColor(tcell.ColorDefault)
	table.SetBorder(true)
	return table
}
