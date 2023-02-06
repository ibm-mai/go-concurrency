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

func GenerateChunksOfFilenames(slice []string, numGorutine int) [][]string {
	var result [][]string
	for i := 0; i < numGorutine; i++ {
		min := (i * len(slice) / numGorutine)
		max := ((i + 1) * len(slice)) / numGorutine
		result = append(result, slice[min:max])
	}
	return result
}
