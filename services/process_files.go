package services

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

func ProcessFile(wg *sync.WaitGroup, numGoroutine int, fileList []string) {
	defer wg.Done()
	fmt.Println("Thread:", numGoroutine, "is processing files ... ", fileList)
	time.Sleep(time.Second * 2)
}

func ProcessFile2(numGoroutine int, fileList []string) {
	fmt.Println("Thread:", numGoroutine, "is processing files ... ", fileList)
	time.Sleep(time.Second * 2)
}

func ProcessFileAndInsertDB(db *sql.DB, goRoutineId int, fileList []string) {
	start := time.Now()
	for _, file := range fileList {
		fmt.Printf("[Goroutine %d] Start Processing file ... %s\n", goRoutineId, file)

		file, err := os.Open("./input/" + file)
		if err != nil {
			fmt.Println("An error occurred:", err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
			msisdn := strings.Split(scanner.Text(), "|")
			switch msisdn[0] {
			case "01":
				continue
			case "09":
				continue
			default:
				tx, err := db.Begin()
				if err != nil {
					fmt.Println("Error starting transaction:", err)
					return
				}

				stmt, err := tx.Prepare("INSERT INTO PRIVUSER.CUSTOMER_PROFILE_TEST_BATCH_OPTIMIZE (MSISDN,BILLING_ACCOUNT) VALUES (?,?)")
				if err != nil {
					fmt.Println("Error preparing statement:", err)
					return
				}
				defer stmt.Close()

				_, err = stmt.Exec(msisdn[4], "9999")
				if err != nil {
					fmt.Println("Error executing statement:", err)
					tx.Rollback()
					return
				}

				err = tx.Commit()
				if err != nil {
					fmt.Println("Error committing transaction:", err)
					tx.Rollback()
					return
				}
			}
		}
	}
	elapsed := time.Since(start)
	fmt.Printf("[Goroutine %d] Done! Time used: %s \n", goRoutineId, elapsed)
}
