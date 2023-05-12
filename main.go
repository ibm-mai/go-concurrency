package main

import (
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"main.go/services"
	"main.go/utils"
	"os"
	"sync"
	"time"
)

func main() {
	utils.InitLogRus()

	var wg sync.WaitGroup

	config, err := utils.LoadConfig("./config")
	if err != nil {
		logrus.Infof("Cannot load config: %s", err)
	}
	if err := os.Mkdir(config.GenData.MsisdnDuplicatePath, os.ModePerm); err != nil {
		logrus.Error("Error creating folder", err)
	}
	start := time.Now()

	// InitDB
	connString := utils.GetConnectionString(config)
	db, err := sqlx.Open(config.Db.Driver, connString)
	if err != nil {
		logrus.Error("Error opening database:", err)
		return
	}
	defer db.Close()

	totalFiles := utils.GetTotalFiles(config.GenData.InputPath)
	logrus.Infof("Number of total files: %d, concurrency: %d", totalFiles, config.GenData.Concurrence)
	fileLists, _ := utils.GetFilesList(config.GenData.InputPath)
	// Create array of chunks
	chunks := utils.GenerateChunksOfFilenames(fileLists, config.GenData.Concurrence)
	logrus.Info("Chunks to process: ", chunks)

	// Way 1: Create annonymous function that create goroutines
	for i, chunck := range chunks {
		wg.Add(1)
		go func(goRoutineId int, chuncks []string) {
			defer wg.Done()
			services.ProcessFileAndInsertDB(db, goRoutineId, chuncks, config.GenData.InputPath, config.GenData.CommitSize, config.GenData.MsisdnDuplicatePath)
		}(i, chunck)
	}
	// Way 2: Pass the named function that create goroutines
	/*
		for i := 0; i < numConcurrent; i++ {
			wg.Add(1)
			go processFile(&wg, i, chunks[i])
		}
	*/
	wg.Wait()
	logrus.Info("Successfully import data")
	elapsed := time.Since(start)
	logrus.Info("Total Time used to process ", totalFiles, " files = ", elapsed)
}
