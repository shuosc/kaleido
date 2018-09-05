package Application

import (
	"bytes"
	"kaleido/Common/Infrastructure/OSS"
	"kaleido/Common/Services/KaleidoMessage"
	"kaleido/Master/Domain/Entities/Area"
	"kaleido/Master/Domain/Entities/IP"
	"kaleido/Master/Domain/Entities/Mirror"
	"kaleido/Master/Domain/Entities/MirrorStation"
	"kaleido/Master/Service/DirtyCheck"
	"log"
	"net/http"
	"sync"
)

var message []byte

var mutex sync.RWMutex

func CronJob() {
	DirtyCheck.Dirty = false
	MirrorStation.CronJob()
	Mirror.CronJob()
	if DirtyCheck.Dirty {
		rawMessage := GenerateMessage()
		msg, _ := rawMessage.Marshal()
		mutex.Lock()
		defer mutex.Unlock()
		message = msg
		go uploadMessage()
		log.Println("Sync with stations success")
	}
}

func uploadMessage() {
	mutex.RLock()
	defer mutex.RUnlock()
	OSS.Bucket.PutObject("kaleido-message", bytes.NewBuffer(message))
	log.Println("New message uploaded")
}

func GetTableHandler(w http.ResponseWriter, r *http.Request) {
	mutex.RLock()
	defer mutex.RUnlock()
	w.Write(message)
}

func GenerateMessage() KaleidoMessage.KaleidoMessage {
	result := KaleidoMessage.KaleidoMessage{}
	result.MirrorStationUrl = map[uint32]string{}
	for id, mirrorStation := range MirrorStation.Repo.Entities {
		result.MirrorStationUrl[id] = mirrorStation.GetUrl()
	}
	result.IPGroups = map[uint32]*KaleidoMessage.IP_Area{}
	for _, ip := range IP.Repo.Entities {
		if _, exist := result.IPGroups[uint32(ip.GetMaskBitLength())]; !exist {
			newIPArea := KaleidoMessage.IP_Area{}
			newIPArea.IPs = map[uint32]uint32{}
			result.IPGroups[uint32(ip.GetMaskBitLength())] = &newIPArea
		}
		result.IPGroups[uint32(ip.GetMaskBitLength())].IPs[ip.GetAddressNumberForm()] = ip.GetAreaId()
	}
	result.Mirrors = map[string]*KaleidoMessage.Mirror{}
	for mirrorName, mirror := range Mirror.Repo.Entities {
		newMirror := KaleidoMessage.Mirror{}
		newMirror.DefaultMirrorStationId = mirror.GetDefaultMirrorStationId()
		newMirror.AreaToMirrorStationTable = map[uint32]*KaleidoMessage.MirrorStationGroup{}
		for areaId := range Area.Repo.Entities {
			group := KaleidoMessage.MirrorStationGroup{}
			group.Stations = mirror.GetAreaToMirrorStationTable(areaId)
			newMirror.AreaToMirrorStationTable[areaId] = &group
		}
		result.Mirrors[mirrorName] = &newMirror
	}
	return result
}
