package plog

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
)

func Init() {
	logFile, err := os.OpenFile("./server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open plog file failed, err:", err)
		return
	}

	//log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
	log.SetPrefix("[server] ")
	writer := io.MultiWriter(logFile, os.Stdout)

	log.SetOutput(writer)
	// gin记录到文件。
	gin.DefaultWriter = writer
	gin.DefaultErrorWriter = writer
}
