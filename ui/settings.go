package ui

import (
	"fmt"
	"strconv"

	"github.com/Youdontknowme720/Cimonv2/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *App) handleAddingToken() {
	form := tview.NewForm().
		AddInputField("Token", "", 20, nil, nil)
	form.SetFieldTextColor(tcell.ColorBlack)
	form.SetBorder(true).SetTitle("Editing ...").SetTitleAlign(tview.AlignCenter)
	form.SetBorderColor(ColorOrange)

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
		AddInputField("ProjectName", "", 0, nil, nil).
		AddInputField("ProjectID", "", 0, nil, nil)

	form.SetBorder(true).SetTitle("... Editing ...").SetTitleAlign(tview.AlignCenter)
	form.SetFieldBackgroundColor(ColorOrange)
	form.SetFieldTextColor(tcell.ColorBlack)
	form.SetLabelColor(tcell.ColorWhite)
	form.SetTitleColor(ColorPink)
	form.SetBorderColor(ColorOrange)
	form.SetBackgroundColor(ColorBlue)
	form.SetButtonBackgroundColor(ColorOrange)
	form.SetButtonTextColor(tcell.ColorWhite)

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
		SetTitleColor(ColorPink).
		SetBorder(true).
		SetBorderColor(ColorOrange).
		SetBackgroundColor(ColorBlue)

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
