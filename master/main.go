package main

import (
	"bytes"
	"fmt"
	"kaleido/common/infrastructure"
	"kaleido/master/model"
	"kaleido/master/service/api"
	"kaleido/master/service/makeMessage"
	"log"
	"net/http"
	"time"
)

func checkStation(mirrorStationChannel chan model.MirrorSupplier,
	changedChannel chan int) {
	for {
		station := <-mirrorStationChannel
		if station.Check() {
			select {
			case x := <-changedChannel:
				changedChannel <- x
			default:
				changedChannel <- 1
			}
		}
		mirrorStationChannel <- station
		time.Sleep(time.Second * 10)
	}
}

func main() {
	stations := model.MirrorStationRepo.GetAll()
	channel := make(chan model.MirrorSupplier, len(stations))
	makeMessageRequest := make(chan int)
	for _, station := range stations {
		channel <- station
	}
	for range stations {
		go checkStation(channel, makeMessageRequest)
	}
	go func(ch chan int) {
		for {
			<-ch
			msg, err := makeMessage.MakeMessage()
			data, err := msg.Marshal()
			if err == nil {
				infrastructure.Bucket.PutObject("kaleido-message-v2", bytes.NewBuffer(data))
				fmt.Println("New message Uploaded")
			}
		}
	}(makeMessageRequest)
	http.HandleFunc("/api/mirror_stations", api.MirrorStationList)
	http.HandleFunc("/api/mirrors", api.MirrorList)
	http.HandleFunc("/api/mirror_station", api.MirrorStation)
	http.HandleFunc("/api/mirror", api.Mirror)

	log.Fatal(http.ListenAndServe(":8086", nil))

	select {}
}
