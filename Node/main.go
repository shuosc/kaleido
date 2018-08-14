package main

import (
	"kaleido/Node/Application"
	_ "kaleido/Node/Application"
	"log"
)

func main() {
	go Application.StartCron()
	go Application.StartServer()
	log.Println("Node booted")
	select {}
}
