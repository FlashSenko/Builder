package build

import (
	"builder/pkg/utils"
)

func checkNeedByDate(sourcePath string, destPath string) (bool, error) {
	sourceModtime, err := utils.GetFileModTime(sourcePath)
	if err != nil {
		return false, err
	}

	destModTime, err := utils.GetFileModTime(destPath)
	if err != nil {
		return false, err
	}

	if sourceModtime.After(destModTime) {
		return true, nil
	}

	return false, nil
}

func checkNeedByHash(sourcePath string) bool {
	return false
}
