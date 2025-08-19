// Package config is used for excuting config configurations
package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Token    string          `yaml:"token"`
	Projects []GitLabProject `yaml:"projects"`
}

type GitLabProject struct {
	ID   int    `yaml:"id"`
	Name string `yaml:"name"`
}

func ReadConfig() Config {
	data, err := os.ReadFile("config/config.yml")
	if err != nil {
		log.Fatal("Error during reading config.yml")
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal("Couldn't read data")
	}
	return cfg
}

func GetProjectData() (string, []GitLabProject) {
	var activeProject []GitLabProject
	cfgData := ReadConfig()
	activeProject = append(activeProject, cfgData.Projects...)
	return cfgData.Token, activeProject
}

func AddNewProject(projectID int, projectName string) {
	cfgData := ReadConfig()
	newGitlabProject := GitLabProject{projectID, projectName}
	cfgData.Projects = append(cfgData.Projects, newGitlabProject)

	newData, err := yaml.Marshal(cfgData)
	if err != nil {
		fmt.Println("Couldn't write to config file")
		panic(err)
	}

	err = os.WriteFile("config/config.yml", newData, 0644)
	if err != nil {
		panic(err)
	}
}

func AddNewToken(token string) {
	cfgData := ReadConfig()
	cfgData.Token = token

	newData, err := yaml.Marshal(cfgData)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("config/config.yml", newData, 0644)
	if err != nil {
		panic(err)
	}
}
