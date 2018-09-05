package main

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/robfig/cron"
	"io/ioutil"
	"kaleido/Common/Infrastructure/OSS"
	"kaleido/Common/Services/IPTools"
	"kaleido/Common/Services/KaleidoMessage"
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
	props, err := OSS.Bucket.GetObjectMeta("kaleido-message")
	if err != nil {
		fmt.Println("Fetch Failed", err)
		return
	}
	newUpdateTime := props["Last-Modified"][0]
	if newUpdateTime != lastUpdateTime {
		lastUpdateTime = newUpdateTime
		fmt.Println(newUpdateTime)
		object, _ := OSS.Bucket.GetObject("kaleido-message")
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
	table := message.Mirrors[mirror].AreaToMirrorStationTable
	for mask := 32; mask >= 0; mask-- {
		ipNumberForm, err := IPTools.IPv4ToNumberForm(ip)
		if err != nil {
			break
		}
		masked := IPTools.MaskIP(ipNumberForm, uint8(mask))
		if ipGroup, ok := message.IPGroups[uint32(mask)]; ok {
			if areaId, ok := ipGroup.IPs[masked]; ok {
				if areaId == 0 {
					break
				}
				mirrorStationIdGroup := table[areaId].Stations
				mirrorStationId := mirrorStationIdGroup[spinNum%uint32(len(mirrorStationIdGroup))]
				spinNum++
				return message.MirrorStationUrl[mirrorStationId]
			}
		}
	}
	return message.MirrorStationUrl[message.Mirrors[mirror].DefaultMirrorStationId]
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	from := strings.Split(r.RemoteAddr, ":")[0]
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
	c.AddFunc("@every 5s", func() {
		go fetchTable()
	})
	c.Start()
	mux := http.NewServeMux()
	handler := http.HandlerFunc(ServeHTTP)
	mux.Handle("/", handler)
	http.ListenAndServe(":8080", mux)
	select {}
}
