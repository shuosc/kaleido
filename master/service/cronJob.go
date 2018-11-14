package service

import (
	"kaleido/master/model/MirrorStation"
	"log"
	"runtime"
	"time"
)

func StartCronJobs() {
	syncStationChannel := make(chan MirrorStation.MirrorStation, runtime.NumCPU())
	uploadTableChannel := make(chan int)
	go func() {
		for {
			<-uploadTableChannel
			uploadTable()
			log.Println("new table uploaded")
			time.Sleep(5 * time.Second)
		}
	}()
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				station := <-syncStationChannel
				changed, err := station.SyncMirrorList()
				if err != nil {
					log.Println(err)
				}
				if changed {
					select {
					case uploadTableChannel <- 1:
					default:
					}
				}
			}
		}()
	}
	uploadTableChannel <- 1
	for {
		stations, err := MirrorStation.All()
		if err != nil {
			log.Println(err)
		}
		for _, mirror := range stations {
			syncStationChannel <- mirror
		}
		time.Sleep(5 * time.Second)
	}
}
