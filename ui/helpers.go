package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func newSelectableTable() *tview.Table {
	table := tview.NewTable()
	table.SetSelectedStyle(
		tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorWhite),
	)
	table.SetSelectable(true, false)
	table.SetBackgroundColor(tcell.ColorBlack)
	table.SetBorder(true)
	return table
}
