package ui

import (
	"github.com/rivo/tview"
)

func StartView(){
	app := tview.NewApplication()

	text := tview.NewTextView().
		SetText("Hello a new View").
		SetTextAlign(tview.AlignCenter)

	if err := app.SetRoot(text, true); err != nil{
		panic(err)
	}
}