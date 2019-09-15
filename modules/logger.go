package modules

import (
	"fmt"
	"log"
	"os"
)

var (
	ERROR_CHANNEL chan string
	INFO_CHANNEL  chan string
)

func init() {
	ERROR_CHANNEL = make(chan string, 100)
	INFO_CHANNEL = make(chan string, 100)
}

func StartLog() {
	logFile, err := os.OpenFile("./run.log", os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(3)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Llongfile)
	for {
		select {
		case val := <-ERROR_CHANNEL:
			log.SetPrefix("[Error]")
			log.Println(val)
		case val := <-INFO_CHANNEL:
			log.SetPrefix("[Info]")
			log.Println(val)
		}
	}
}

func LogError(msg string) {
	fmt.Println(msg)
	ERROR_CHANNEL <- msg
}

func LogInfo(msg string) {
	fmt.Println(msg)
	INFO_CHANNEL <- msg
}
