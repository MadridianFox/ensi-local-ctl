package src

import (
	"errors"
	"gopkg.in/yaml.v2"
)

type HomeConfigItem struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

type HomeConfig struct {
	Path             string           `yaml:"-"`
	CurrentWorkspace string           `yaml:"current_workspace"`
	UpdateCommand    string           `yaml:"update_command"`
	Workspaces       []HomeConfigItem `yaml:"workspaces"`
}

const defaultUpdateCommand = "curl -sSL https://raw.githubusercontent.com/MadridianFox/ensi-local-ctl/master/get.sh | sudo bash"

func LoadHomeConfig(configPath string) (*HomeConfig, error) {
	yamlFile, err := Pc.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	cfg := &HomeConfig{}
	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		return nil, err
	}
	cfg.Path = configPath
	return cfg, nil
}

func SaveHomeConfig(homeConfig *HomeConfig) error {
	data, err := yaml.Marshal(homeConfig)
	if err != nil {
		return err
	}

	err = Pc.WriteFile(homeConfig.Path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func CheckHomeConfigIsEmpty(configPath string) error {
	if Pc.FileExists(configPath) {
		return nil
	}
	return SaveHomeConfig(&HomeConfig{Path: configPath, UpdateCommand: defaultUpdateCommand})
}

func (hc *HomeConfig) AddWorkspace(name string, path string) error {
	hc.Workspaces = append(hc.Workspaces, HomeConfigItem{Name: name, Path: path})
	return SaveHomeConfig(hc)
}

func (hc *HomeConfig) GetCurrentWsPath() (string, error) {
	if hc.CurrentWorkspace == "" {
		return "", errors.New("current workspace is not set")
	}

	for _, hci := range hc.Workspaces {
		if hci.Name == hc.CurrentWorkspace {
			return hci.Path, nil
		}
	}

	return "", errors.New("current workspace is bad")
}

func (hc *HomeConfig) findWorkspace(name string) *HomeConfigItem {
	for _, workspace := range hc.Workspaces {
		if workspace.Name == name {
			return &workspace
		}
	}

	return nil
}
