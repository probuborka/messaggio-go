package domain

import (
	"fmt"
	"os"
	"path/filepath"
)

func Path(dir string) (string, error) {
	dirPath, err := os.Getwd()
	if err != nil {
		return "", err
	}
	// search dir
	for {
		retPath := filepath.Join(dirPath, dir)
		_, err = os.Stat(retPath)
		if err != nil {
			if os.IsNotExist(err) {
				sPath, _ := filepath.Split(dirPath)
				dirPath = sPath[:len(sPath)-1]
			} else {
				return "", err
			}
			if dirPath == "" {
				return "", fmt.Errorf("directory not found")
			}
			continue
		}
		return retPath, nil
	}
}
