package services

import (
	"bufio"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"main.go/utils"
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

type CustomerProfile struct {
	Msisdn         string `db:"msisdn"`
	BillingAccount int    `db:"billing_account"`
}

func ProcessFileAndInsertDB(db *sqlx.DB, goRoutineId int, fileList []string, inputPath string, batchSize int, msisdnDuplicatePath string) {
	config, err := utils.LoadConfig("./config")
	if err != nil {
		logrus.Infof("Cannot load config: %s", err)
	}
	if len(fileList) == 0 {
		logrus.Debugf("[Goroutine %d] fileList is empty", goRoutineId)
		return
	}

	logrus.Infof("Start import data to: %s", config.Db.CustomerProfileTableName)
	start := time.Now()
	for _, fileName := range fileList {
		logrus.Infof("[Goroutine %d] Start Processing filename: %s", goRoutineId, fileName)
		file, err := os.Open(inputPath + "/" + fileName)
		if err != nil {
			logrus.Error("An error occurred:", err)
			return
		}
		defer file.Close()

		profiles := make([]CustomerProfile, batchSize)
		reader := bufio.NewReader(file)
		count := 0
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			data := strings.Split(line, "|")
			switch data[0] {
			case "01":
				continue
			case "09":
				continue
			default:
				// Add the Person object to the profile
				profiles[count%batchSize] = CustomerProfile{Msisdn: data[4], BillingAccount: 9999}
				count++
				// If the profile is full, insert it into the database
				if count%batchSize == 0 {
					tx, err := db.Beginx()
					if err != nil {
						logrus.Error("Error starting transaction:", err)
						return
					}
					stmt, err := tx.PrepareNamed(`INSERT INTO PRIVUSER.` + config.Db.CustomerProfileTableName + ` (msisdn, billing_account) VALUES (:msisdn, :billing_account)`)
					if err != nil {
						logrus.Error("Error preparing statement:", err)
						return
					}
					defer stmt.Close()

					// Execute the SQL statement with named parameters for each Person object in the array
					for _, profile := range profiles {
						_, err := stmt.Exec(profile)
						if err != nil {
							file, err := os.OpenFile(msisdnDuplicatePath+fmt.Sprintf("/dup-msisdn-goroutine-%d.txt", goRoutineId), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
							if err != nil {
								logrus.Error("Error open file:", err)
								return
							}
							defer file.Close()
							_, err = file.WriteString(profile.Msisdn + "\n")
							if err != nil {
								logrus.Error("Error write file:", err)
								return
							}
							continue
						}
					}

					err = tx.Commit()
					if err != nil {
						logrus.Error("Error committing transaction:", err)
						tx.Rollback()
						return
					}
					logrus.Infof("[Goroutine %d] Commit successfully: %d records", goRoutineId, count)
				}
			}
		}
		// Insert any remaining lines into the database
		if count%batchSize != 0 {
			tx, err := db.Beginx()
			if err != nil {
				logrus.Error("Error starting transaction:", err)
				return
			}

			stmt, err := tx.PrepareNamed("INSERT INTO PRIVUSER.CUSTOMER_PROFILE_TEST (msisdn, billing_account) VALUES (:msisdn, :billing_account)")
			if err != nil {
				logrus.Error("Error preparing statement:", err)
				return
			}
			defer stmt.Close()

			for _, profile := range profiles[:count%batchSize] {
				_, err := stmt.Exec(profile)
				if err != nil {
					//logrus.Infof("[Remaining] Skip MSISDN: %s", profile.Msisdn)
					//file, err := os.OpenFile(msisdnDuplicatePath+fmt.Sprintf("/dup-msisdn-goroutine-%d.txt", goRoutineId), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
					//if err != nil {
					//	logrus.Error("Error open file:", err)
					//}
					//defer file.Close()
					//_, err = file.WriteString(profile.Msisdn + "\n")
					//if err != nil {
					//	logrus.Error("Error write file:", err)
					//	return
					//}
					continue
				}
			}

			err = tx.Commit()
			if err != nil {
				panic(err)
			}
			logrus.Infof("[Goroutine %d] Commit successfully: %d records (remaining)", goRoutineId, count)
		}
		elapsed := time.Since(start)
		logrus.Infof("[Goroutine %d] Done! Time used: %s ", goRoutineId, elapsed)

		// move files to done folder
		logrus.Infof("[Goroutine %d] Start moving files %s ", goRoutineId, fileName)
		oldLocation := inputPath + "/" + fileName
		newLocation := "./done/" + fileName
		err = os.Rename(oldLocation, newLocation)
		if err != nil {
			logrus.Fatal(err)
		}
	}
}
