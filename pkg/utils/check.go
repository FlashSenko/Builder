package utils

import "os"

func IsDir(path string) (bool, error) {
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false, nil
	} else if _, ok := err.(*os.PathError); ok {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return info.IsDir(), nil
}
