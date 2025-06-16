package build

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"builder/pkg/utils"
)

// Loop runs when the destination is directory.
func buildLoopForDDir(destPath string, config *Config, execution *Execution) error {

	for source := range execution.SourcesIter(config) {
		var (
			sourceFileName = filepath.Base(source)
			targetPath     = destPath + (strings.TrimSuffix(sourceFileName, filepath.Ext(sourceFileName))) + "." + execution.To

			need bool = true
			err  error
		)

		if config.CheckDate {
			need, err = checkNeedByDate(source, targetPath)

			if os.IsNotExist(err) {
				need = true
			} else if err != nil {
				return err
			}
		}

		if need && config.CheckHash {
			need = checkNeedByHash(source)
		}

		if need {
			err := execute(execution.Command, source, targetPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Loop runs when the destination is a file.
func buildLoopForDFile(destPath string, config *Config, execution *Execution) (string, time.Time, error) {
	var err error

	destModTime, err := utils.GetFileModTime(destPath)
	if err != nil && !os.IsNotExist(err) {
		return "", time.Time{}, err
	}

	var (
		sourceModTime time.Time = time.Time{}
		sourcesArg    string    = ""
		need          bool      = false
		source        string
	)

	if !config.CheckDate && !config.CheckHash {
		need = true
	}

	for source, err = range execution.SourcesIter(config) {
		if err != nil {
			return "", time.Time{}, err
		}

		sourcesArg += " " + source

		if need {
			continue
		}

		if config.CheckDate {
			sourceModTime, err = utils.GetFileModTime(source)

			if os.IsNotExist(err) {
				need = true
			} else if err != nil {
				return "", time.Time{}, err
			}

			if sourceModTime.After(destModTime) {
				need = true
			}
		}

		if need && config.CheckHash {
			need = checkNeedByHash(source)
		}
	}

	if need {
		err := execute(execution.Command, sourcesArg[1:], destPath)
		if err != nil {
			return "", time.Time{}, err
		}

		return destPath, sourceModTime, nil
	}

	return "", time.Time{}, nil
}

func Build(config *Config) error {
	for _, execution := range config.Executions {
		var compiledFilePath string = ""
		var minimumDate time.Time = time.Time{}

		for dest := range execution.DestsIter(config) {
			isDir, err := utils.IsDir(dest)
			if err != nil {
				return err
			}

			if isDir {
				err := buildLoopForDDir(dest, config, &execution)
				if err != nil {
					return err
				}
				continue
			}

			if minimumDate.IsZero() {
				compiledFilePath, minimumDate, err = buildLoopForDFile(dest, config, &execution)
				if err != nil {
					return err
				}
				continue
			} else {
				destModTime, err := utils.GetFileModTime(dest)
				if err != nil {
					return err
				}

				if destModTime.After(minimumDate) {
					continue
				}
			}

			if compiledFilePath != "" {
				utils.Copy(compiledFilePath, dest)
			}
		}
	}

	return nil
}
