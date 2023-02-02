package services

import (
	"fmt"
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
