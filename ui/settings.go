package ui

import (
	"fmt"
	"strconv"

	"github.com/Youdontknowme720/Cimonv2/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *App) createSettingsPage() tview.Primitive {
	table := a.handleSettingsClick()

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == 'b' {
			a.pages.SwitchToPage(PageHome)
			return nil
		}
		return event
	})

	table.SetSelectedFunc(func(row, column int) {
		cell := table.GetCell(row, column)
		ref := cell.GetReference()

		if ref == nil {
			fmt.Println("Keine Referenz – Auswahl ignoriert")
			return
		}

		if refStr, ok := ref.(string); ok {
			switch refStr {
			case "Add":
				a.handleAddingProject()
			case "Del":
				fmt.Println("Deleting project...")
			case "Conf":
				fmt.Println("Configuring project...")
			case "Tok":
				a.handleAddingToken()
			}
		}
	})
	return table
}

func (a *App) handleAddingToken() {
	form := tview.NewForm().
		AddInputField("Token", "", 20, nil, nil)
	form.SetFieldTextColor(tcell.ColorBlack)
	form.SetBorder(true).SetTitle("Editing ...").SetTitleAlign(tview.AlignCenter)

	saveFunc := func() {
		newToken := form.GetFormItemByLabel("Token").(*tview.InputField).GetText()
		config.AddNewToken(newToken)
		a.token = newToken
		a.pages.SwitchToPage(PageHome)
	}

	abortFunc := func() {
		a.pages.SwitchToPage(PageHome)
	}

	form.AddButton("Save", saveFunc)
	form.AddButton("Abort", abortFunc)

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Key() == tcell.KeyCtrlS:
			saveFunc()
			return nil
		case event.Key() == tcell.KeyCtrlB || event.Key() == tcell.KeyEsc:
			abortFunc()
			return nil
		}
		return event
	})
	a.pages.AddPage(PageAddToken, form, true, true)
	a.pages.SwitchToPage(PageAddToken)
}

func (a *App) handleAddingProject() {
	form := tview.NewForm().
		AddInputField("ProjectID", "", 20, nil, nil).
		AddInputField("ProjectName", "", 20, nil, nil)

	form.SetFieldTextColor(tcell.ColorBlack)
	form.SetBorder(true).SetTitle("Editing ...").SetTitleAlign(tview.AlignCenter)

	configOverviewTable := tview.NewTable()
	_, activeProjects := config.GetProjectData()
	for idx, project := range activeProjects {
		configOverviewTable.SetCell(idx, 0,
			tview.NewTableCell(fmt.Sprint(project.ID)).SetAlign(tview.AlignLeft).SetSelectable(false))
		configOverviewTable.SetCell(idx, 1,
			tview.NewTableCell(project.Name).SetAlign(tview.AlignLeft).SetSelectable(false))
	}
	configOverviewTable.SetTitle("Current projects definied in config").
		SetTitleAlign(tview.AlignCenter).
		SetBorder(true)

	saveFunc := func() {
		name := form.GetFormItemByLabel("ProjectName").(*tview.InputField).GetText()
		idStr := form.GetFormItemByLabel("ProjectID").(*tview.InputField).GetText()
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return
		}
		config.AddNewProject(id, name)
		a.pages.RemovePage(PageHome)
		_, projects := config.GetProjectData()
		a.gitlabProjects = projects
		home := a.createHomeScreen(a.gitlabProjects)
		a.pages.AddPage(PageHome, home, true, true)
		a.pages.SwitchToPage(PageHome)
	}

	abortFunc := func() {
		a.pages.SwitchToPage(PageHome)
	}

	form.AddButton("Save", saveFunc)
	form.AddButton("Abort", abortFunc)

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Key() == tcell.KeyCtrlS:
			saveFunc()
			return nil
		case event.Key() == tcell.KeyCtrlB || event.Key() == tcell.KeyEsc:
			abortFunc()
			return nil
		}
		return event
	})

	flex := tview.NewFlex().
		AddItem(form, 0, 1, true).
		AddItem(configOverviewTable, 0, 1, false)
	a.pages.AddPage(PageAddProj, flex, true, true)
	a.pages.SwitchToPage(PageAddProj)
}

func (a *App) handleSettingsClick() *tview.Table {
	tableHeader := map[string]string{
		"Add new project":   "Add",
		"Delete project":    "Del",
		"Configure project": "Conf",
		"Add new Token":     "Tok",
	}
	table := newSelectableTable()
	cnt := 0
	for key, value := range tableHeader {
		cell := tview.NewTableCell(key).
			SetAlign(tview.AlignLeft).
			SetSelectable(true).
			SetReference(value)
		table.SetCell(cnt, 0, cell)
		cnt++
	}
	return table
}
