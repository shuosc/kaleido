package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/robfig/cron"
	"io/ioutil"
	"kaleido/common/iptools"
	"kaleido/common/message"
	"kaleido/common/oss"
	"log"
	"net/http"
	"strings"
	"sync"
)

var message KaleidoMessage.KaleidoMessage

var lastUpdateTime string

var mutex sync.RWMutex

var spinNum uint32

func fetchTable() {
	props, err := oss.Bucket.GetObjectMeta("kaleido-message")
	if err != nil {
		fmt.Println("Fetch Failed", err)
		return
	}
	newUpdateTime := props["Last-Modified"][0]
	if newUpdateTime != lastUpdateTime {
		lastUpdateTime = newUpdateTime
		fmt.Println(newUpdateTime)
		object, _ := oss.Bucket.GetObject("kaleido-message")
		buffer, _ := ioutil.ReadAll(object)
		newMessage := KaleidoMessage.KaleidoMessage{}
		proto.Unmarshal(buffer, &newMessage)
		mutex.Lock()
		defer mutex.Unlock()
		message = newMessage
		log.Println("Table fetched")
	}
}

func GetRedirectToStation(mirror string, ip string) string {
	mutex.RLock()
	defer mutex.RUnlock()
	table := message.Mirrors[mirror].AreaISP_MirrorStationGroup
	for mask := 32; mask >= 0; mask-- {
		masked, err := iptools.MergedMaskedIP(ip, uint8(mask))
		fmt.Println(ip, mask, masked)
		if err != nil {
			continue
		}
		if areaISP, ok := message.Address_AreaISP[masked]; ok {
			mirrorStationIdGroup := table[areaISP].Stations
			mirrorStationId := mirrorStationIdGroup[spinNum%uint32(len(mirrorStationIdGroup))]
			spinNum++
			return message.MirrorStationId_Url[mirrorStationId]
		}
	}
	return message.MirrorStationId_Url[message.Mirrors[mirror].FallbackMirrorStationId]
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	from := r.Header["Remote_addr"][0]
	if r.Header.Get("force-redirect-ip") != "" {
		from = r.Header.Get("force-redirect-ip")
	}
	mirror := strings.Split(r.URL.Path, "/")[1]
	to := GetRedirectToStation(mirror, from)
	log.Println("Redirect request from ", from, " for mirror", mirror, "to", to)
	http.Redirect(w, r, to+r.URL.Path, http.StatusMovedPermanently)
}

func main() {
	c := cron.New()
	c.AddFunc("@every 10s", func() {
		go fetchTable()
	})
	c.Start()
	mux := http.NewServeMux()
	handler := http.HandlerFunc(ServeHTTP)
	mux.Handle("/", handler)
	http.ListenAndServe(":8080", mux)
	select {}
}
