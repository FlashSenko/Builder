package build

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"builder/pkg/utils"
)

var configs []*Config

type Config struct {
	CheckDate bool `json:"check_date"`
	CheckHash bool `json:"check_hash"`

	Executions []Execution    `json:"executions"`
	MapsDict   map[string]Map `json:"maps_dict"`
}

type Execution struct {
	Command string `json:"command"`
	Maps    []any  `json:"maps"`
	To      string `json:"to"`
}

type Map struct {
	Sources      []string `json:"sources"`
	Destinations []string `json:"destinations"`
}

func (e *Execution) MapsIter(config *Config) func(yield func(*Map, error) bool) {
	return func(yield func(*Map, error) bool) {
		for _, mapValue := range e.Maps {
			mapStr, ok := mapValue.(string)
			if ok {
				mapStruct, ok := config.MapsDict[mapStr]
				if !ok {
					yield(nil, fmt.Errorf("map key not found: %s", mapStr))
					return
				}

				if !yield(&mapStruct, nil) {
					return
				}
			} else {
				mapMap := mapValue.(map[string]any)

				if !yield(&Map{
					Sources:      utils.ToStringSlice(mapMap["sources"].([]any)),
					Destinations: utils.ToStringSlice(mapMap["destinations"].([]any)),
				}, nil) {
					return
				}
			}
		}
	}
}

func (e *Execution) SourcesIter(config *Config) func(yield func(string, error) bool) {
	return func(yield func(string, error) bool) {
		for mapStruct, err := range e.MapsIter(config) {
			if err != nil {
				yield("", err)
				return
			}

			for _, source := range mapStruct.Sources {
				isDir, err := utils.IsDir(source)
				if err != nil {
					yield("", err)
					return
				}

				if isDir {
					if source[len(source)-1] != '/' {
						source = source + "/"
					}

					source = source + "*"
				}

				matches, err := filepath.Glob(source)
				if err != nil {
					yield("", err)
					return
				}

				for _, path := range matches {
					isDir, err = utils.IsDir(path)
					if isDir {
						continue
					}

					if !yield(path, err) {
						return
					}
				}

			}
		}
	}
}

func (e *Execution) DestsIter(config *Config) func(yield func(string, error) bool) {
	return func(yield func(string, error) bool) {
		for mapStruct, err := range e.MapsIter(config) {
			if err != nil {
				yield("", err)
				return
			}

			for _, destination := range mapStruct.Destinations {
				matches, err := filepath.Glob(destination)

				if len(matches) == 0 {
					if !yield(destination, err) {
						return
					}
				}

				for _, path := range matches {
					if !yield(path, err) {
						return
					}
				}
			}
		}
	}
}

func LoadConfig(configPath string, save bool) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config *Config = &Config{}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	if save {
		configs = append(configs, config)
	}

	return config, nil
}
