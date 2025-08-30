// Package config is used for excuting config configurations
package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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
	configDir := "config"
	configPath := filepath.Join(configDir, "config.yml")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return createDefaultConfig(configDir, configPath)
	}
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

func createDefaultConfig(configDir, configPath string) Config {
	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Printf("Warning: Could not create config directory: %v", err)
	}

	defaultConfig := Config{
		Token:    "",
		Projects: []GitLabProject{},
	}

	defaultYAML := `token: ""
projects:
`

	if err := os.WriteFile(configPath, []byte(defaultYAML), 0644); err != nil {
		log.Printf("Warning: Could not create default config file: %v", err)
	} else {
		log.Printf("Created default config at: %s", configPath)
	}

	return defaultConfig
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
