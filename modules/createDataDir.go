package modules

import (
	"os"
)

func CreateDataDir() {
	exist, err := pathExists("./data")

	if err != nil {
		LogError(err.Error())
		os.Exit(1)
	}

	if !exist {
		err := os.Mkdir("./data", os.ModePerm)
		if err != nil {
			LogError(err.Error())
			os.Exit(1)
		} else {
			LogInfo("mkdir ./data")
		}
	}
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
