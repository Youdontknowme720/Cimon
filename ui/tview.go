package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
)

func StartView() {
	app := tview.NewApplication()

	text := tview.NewBox().
		SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetTitle("My First TviewBox")

	app.SetRoot(text, true)

	if err := app.Run(); err != nil {
		log.Fatalf("Run failed: %v", err)
	}
}
