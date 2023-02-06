package main

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
	"main.go/services"
	"main.go/utils"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	config, err := utils.LoadConfig("./config")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}
	start := time.Now()

	connString := utils.GetConnectionString(config)
	fmt.Println(connString)
	db, err := sql.Open(config.Db.Driver, connString)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	totalFiles := utils.GetTotalFiles("./input")
	fmt.Println("Number of total files:", totalFiles)
	fileLists, _ := utils.GetFilesList("./input")
	// Create array of chunks
	chuncks := utils.GenerateChunksOfFilenames(fileLists, config.Concurrence)
	fmt.Println("Chunks to process:", chuncks)

	// Way 1: Create annonymous function that create goroutines
	for i, chunck := range chuncks {
		wg.Add(1)
		go func(goRoutineId int, chuncks []string) {
			defer wg.Done()
			services.ProcessFileAndInsertDB(db, goRoutineId, chuncks)
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
	elapsed := time.Since(start)
	fmt.Println("Total Time used to process", totalFiles, "=", elapsed)
}
