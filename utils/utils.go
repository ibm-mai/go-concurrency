package utils

import (
	"fmt"
	"io/ioutil"
)

func GetTotalFiles(direcotry string) int {
	files, _ := ioutil.ReadDir(direcotry)
	return len(files)
}

func GetFilesList(directory string) ([]string, error) {
	filesList := []string{}
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	for _, f := range files {
		filesList = append(filesList, f.Name())
	}
	return filesList, err
}
