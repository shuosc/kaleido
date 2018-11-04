package api

import (
	"encoding/json"
	"kaleido/master/model"
	"net/http"
	"strconv"
)

func MirrorList(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var mirrors []struct {
		Id       uint64   `json:"id"`
		Name     string   `json:"name"`
		Stations []uint64 `json:"stations"`
	}
	for _, mirror := range model.MirrorRepo.GetAll() {
		name, _ := mirror.GetName()
		obj := struct {
			Id       uint64   `json:"id"`
			Name     string   `json:"name"`
			Stations []uint64 `json:"stations"`
		}{Id: mirror.Id, Name: name, Stations: []uint64{}}
		stations, _ := mirror.GetStationsContainThis()
		for _, station := range stations {
			obj.Stations = append(obj.Stations, station.Id)
		}
		mirrors = append(mirrors, obj)
	}
	data, _ := json.Marshal(mirrors)
	w.Write(data)
}

func Mirror(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	queryId, _ := strconv.ParseInt(r.URL.Query()["id"][0], 10, 64)
	m := model.Mirror{}
	m.Id = uint64(queryId)
	name, _ := m.GetName()
	inStations, _ := m.GetStationsContainThis()
	obj := struct {
		Id       uint64   `json:"id"`
		Name     string   `json:"name"`
		Stations []uint64 `json:"stations"`
	}{
		Id:       m.Id,
		Name:     name,
		Stations: []uint64{},
	}
	for _, station := range inStations {
		obj.Stations = append(obj.Stations, station.Id)
	}
	data, _ := json.Marshal(obj)
	w.Write(data)
}
