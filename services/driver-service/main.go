package main

import (
	"log"
	"time"
)

func main() {
	loopNumber := 1
	for {
		log.Printf("hello from driver-service %d", loopNumber)
		time.Sleep(time.Second * 5)
		loopNumber += 1
	}
	
}