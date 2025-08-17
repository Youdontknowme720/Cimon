// Package ui is a nice ui interface
package ui

import (
	"fmt"
	"strconv"

	"github.com/Youdontknowme720/Cimonv2/config"
	"github.com/Youdontknowme720/Cimonv2/gitlab"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type App struct {
	app            *tview.Application
	pages          *tview.Pages
	gitlabProjects []config.GitLabProject
	token          string
}

func NewApp() *App {
	token, projects := config.GetProjectData()
	app := &App{tview.NewApplication(), tview.NewPages(), projects, token}
	return app
}

func (a *App) Run() error {
	return a.app.Run()
}

func (a *App) Setup() {
	homeHeader := []config.GitLabProject{}
	homeHeader = append(homeHeader, a.gitlabProjects...)
	home := a.createHomeScreen(homeHeader)
	home.SetTitle("Home Menue")
	a.pages.AddPage("home", home, true, true)
	a.app.SetRoot(a.pages, true)
}

func (a *App) createHomeScreen(projects []config.GitLabProject) *tview.Table {
	table := tview.NewTable()
	table.SetSelectedStyle(tcell.StyleDefault.
		Background(tcell.ColorBlue).
		Foreground(tcell.ColorWhite)).
		SetSelectable(true, false)
	table.SetBackgroundColor(tcell.ColorDefault)
	table.SetBorder(true)
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
		cell := table.GetCell(row, column)
		ref := cell.GetReference()

		if ref == nil {
			fmt.Println("Keine Referenz – Auswahl ignoriert")
			return
		}

		switch v := ref.(type) {

		case config.GitLabProject:
			page := a.createPipelinePage(v)
			a.pages.AddPage("pipelines", page, true, true)
			a.pages.SwitchToPage("pipelines")

		case string:
			if v == "Settings" {
				page := a.createSettingsPage()
				a.pages.AddPage("settings", page, true, true)
				a.pages.SwitchToPage("settings")
			}

		default:
			fmt.Println("Unbekannter Typ – Auswahl ignoriert")
		}
	})

	return table
}

func (a *App) createPipelinePage(proj config.GitLabProject) tview.Primitive {
	tree := a.handlePipelineClick(fmt.Sprint(proj.ID))
	tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == 'b' {
			a.pages.SwitchToPage("home")
			return nil
		}
		return event
	})
	return tree
}

func (a *App) handlePipelineClick(projectID string) *tview.TreeView {
	pipelines, err := gitlab.GetAllPipelines(projectID, a.token, 3)
	if err != nil {
		panic(err)
	}
	root := tview.NewTreeNode("Pipelines").
		SetColor(tcell.ColorGreen)

	for _, p := range pipelines {
		pipelineNode := tview.NewTreeNode(fmt.Sprint(p.ID)).
			SetReference(p).
			SetSelectable(true)
		root.AddChild(pipelineNode)
	}

	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	return tree
}

func (a *App) createSettingsPage() tview.Primitive {
	table := a.handleSettingsClick()
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == 'b' {
			a.pages.SwitchToPage("home")
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
				fmt.Println("Adding")
			case "Conf":
				fmt.Println("Adding")
			}
		}
	})
	return table
}

func (a *App) handleAddingProject() {
	form := tview.NewForm().
		AddInputField("ProjectID", "", 20, nil, nil).
		AddInputField("ProjectName", "", 20, nil, nil)

	form.AddButton("Save", func() {
		name := form.GetFormItemByLabel("ProjectName").(*tview.InputField).GetText()
		idStr := form.GetFormItemByLabel("ProjectID").(*tview.InputField).GetText()
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return
		}
		config.AddNewProject(id, name)

		a.pages.SwitchToPage("home")
	})
	form.AddButton("Abbort", func() {
		a.pages.SwitchToPage("home")
	})
	form.SetBorder(true).SetTitle("Adding new Project").SetTitleAlign(tview.AlignLeft)
	a.pages.AddPage("addProject", form, true, true)
	a.pages.SwitchToPage("addProject")
}

func (a *App) handleSettingsClick() *tview.Table {
	tableHeader := map[string]string{
		"Add new project":   "Add",
		"Delete project":    "Del",
		"Configure project": "Conf",
	}
	table := tview.NewTable()
	table.SetSelectedStyle(tcell.StyleDefault.
		Background(tcell.ColorBlue).
		Foreground(tcell.ColorWhite)).
		SetSelectable(true, false)
	table.SetBackgroundColor(tcell.ColorDefault)
	table.SetBorder(true)
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
