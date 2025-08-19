// Package ui is a nice ui interface
// Package ui is a nice ui interface
package ui

import (
	"github.com/Youdontknowme720/Cimonv2/config"
	"github.com/rivo/tview"
)

const (
	PageHome     = "home"
	PageSettings = "settings"
	PagePipeline = "pipelines"
	PageAddProj  = "addProject"
	PageAddToken = "addToken"
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
	home := a.createHomeScreen(a.gitlabProjects)
	a.pages.AddPage(PageHome, home, true, true)
	a.app.SetRoot(a.pages, true)
}
