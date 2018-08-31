package MirrorStation

import (
	"io/ioutil"
	"net/http"
)

type JsonStructure interface {
	GetNamesForJson(json []byte) []string
}

type JsonIndexMirrorStation struct {
	mirrorStationBase
	MirrorListUrl string
	structure     JsonStructure
}

func (jsonIndexMirrorStation JsonIndexMirrorStation) GetMirrorList() []string {
	response, _ := http.Get(jsonIndexMirrorStation.MirrorListUrl)
	buffer, _ := ioutil.ReadAll(response.Body)
	return jsonIndexMirrorStation.structure.GetNamesForJson(buffer)
}
