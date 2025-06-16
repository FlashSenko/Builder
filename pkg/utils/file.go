package utils

import (
	"fmt"
	"io"
	"os"
	"time"
)

func GetFileModTime(path string) (time.Time, error) {
	sourceInfo, err := os.Stat(path)

	if os.IsNotExist(err) {
		return time.Time{}, nil
	} else if err != nil {
		return time.Time{}, err
	}

	return sourceInfo.ModTime(), nil
}

func Copy(sourcePath, destPath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("cannnot open source file: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("cannot create destination file: %w", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("occur an error while copying: %w", err)
	}

	err = destFile.Sync()
	if err != nil {
		return fmt.Errorf("occur an error while writting files into disk: %w", err)
	}

	return nil
}
