package MirrorStation

import (
	"io/ioutil"
	"kaleido/Master/Service/DirtyCheck"
	"log"
	"net/http"
	"strings"
)

type JsonStructure interface {
	GetNamesForJson(json []byte) []string
}

type JsonIndexMirrorStation struct {
	mirrorStationBase
	MirrorListUrl string
	structure     JsonStructure
}

func (jsonIndexMirrorStation *JsonIndexMirrorStation) UpdateMirrorList() {
	response, err := http.Get(jsonIndexMirrorStation.MirrorListUrl)
	if err != nil {
		log.Println("Failed to fetch MirrorListUrl from ", jsonIndexMirrorStation.MirrorListUrl)
		jsonIndexMirrorStation.mirrorList = []string{}
		return
	}
	buffer, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Failed to fetch MirrorListUrl from ", jsonIndexMirrorStation.MirrorListUrl)
		jsonIndexMirrorStation.mirrorList = []string{}
		return
	}
	oldListContent := strings.Join(jsonIndexMirrorStation.mirrorList, "")
	jsonIndexMirrorStation.mirrorList = jsonIndexMirrorStation.structure.GetNamesForJson(buffer)
	if oldListContent != strings.Join(jsonIndexMirrorStation.mirrorList, "") {
		log.Println("Found content in ", jsonIndexMirrorStation.url, " changed")
		DirtyCheck.Dirty = true
	}
}
