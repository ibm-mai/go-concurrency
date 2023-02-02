package main

import (
	"fmt"
	"main.go/services"
	"main.go/utils"
	"sync"
)

func generateChunksOfFilenames(slice []string, numGorutine int) [][]string {
	var result [][]string
	for i := 0; i < numGorutine; i++ {
		min := (i * len(slice) / numGorutine)
		max := ((i + 1) * len(slice)) / numGorutine
		result = append(result, slice[min:max])
	}
	return result
}

func main() {
	numConcurrent := 2

	var wg sync.WaitGroup
	//wg.Add(numConcurrent)
	totalFiles := utils.GetTotalFiles("./input")
	fmt.Println("Number of total files:", totalFiles)
	fileLists, _ := utils.GetFilesList("./input")

	// Create array of chunks
	chuncks := generateChunksOfFilenames(fileLists, numConcurrent)
	fmt.Println("Chunks to process:", chuncks)

	// Way 1: Create annonymous function that create goroutines
	for i, chunck := range chuncks {
		wg.Add(1)
		go func(goRoutine int, chuncks []string) {
			defer wg.Done()
			services.ProcessFile2(goRoutine, chuncks)
		}(i, chunck)
	}
	// Way 2: Pass the named function that create goroutines
	/*
		for i := 0; i < numConcurrent; i++ {
			wg.Add(1)
			go processFile(&wg, i, chuncks[i])
		}
	*/
	wg.Wait()
	fmt.Println("Successfully import data")
}
