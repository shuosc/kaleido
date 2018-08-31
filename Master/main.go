package main

import (
	"github.com/robfig/cron"
	"kaleido/Master/Application"
	"log"
	"net/http"
)

func main() {
	Application.CronJob()
	c := cron.New()
	c.AddFunc("@every 5s", func() {
		Application.CronJob()
		log.Println("Sync with stations success")
	})
	c.Start()
	mux := http.NewServeMux()
	tableHandler := http.HandlerFunc(Application.GetTableHandler)
	mux.Handle("/", tableHandler)
	http.ListenAndServe(":8086", mux)
}
