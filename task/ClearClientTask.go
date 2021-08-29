package task

import (
	"fmt"
	"log"
	"time"
)

var (
	ClearClientInterval = time.Second * 30
)

func StartClearClientTask() {
	ticker := time.NewTicker(ClearClientInterval)
	log.Println("启动定时清理客户端任务")
	go tick(ticker)
}

func tick(ticker *time.Ticker) {
	for range ticker.C {
		fmt.Printf("ticked at %v", time.Now())

	}
}
