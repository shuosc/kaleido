package main

import (
	"kaleido/Master/Application"
	_ "kaleido/Master/Domain/Entities/IPRange"
	"log"
)

func main() {
	go Application.StartCron()
	log.Println("Cron started")
	go Application.StartServer()
	log.Println("Master booted")
	select {}
}
