package api

import (
	"encoding/json"
	"kaleido/master/model"
	"net/http"
	"strconv"
)

func MirrorStationList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var stations []struct {
		Id      uint64   `json:"id"`
		Name    string   `json:"name"`
		Url     string   `json:"url"`
		Alive   bool     `json:"alive"`
		Mirrors []uint64 `json:"mirrors"`
	}
	for _, station := range model.MirrorStationRepo.GetAll() {
		m := model.MirrorStation{}
		m.Id = uint64(station.GetId())
		name, _ := m.GetName()
		url, _ := m.GetUrl()
		alive, _ := m.IsAlive()
		mirrors, _ := m.GetMirrors()
		obj := struct {
			Id      uint64   `json:"id"`
			Name    string   `json:"name"`
			Url     string   `json:"url"`
			Alive   bool     `json:"alive"`
			Mirrors []uint64 `json:"mirrors"`
		}{
			Id:      station.GetId(),
			Name:    name,
			Url:     url,
			Alive:   alive,
			Mirrors: []uint64{},
		}
		for _, mirror := range mirrors {
			obj.Mirrors = append(obj.Mirrors, mirror.Id)
		}
		stations = append(stations, obj)
	}
	data, _ := json.Marshal(stations)
	w.Write(data)
}

func MirrorStation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	queryId, _ := strconv.ParseInt(r.URL.Query()["id"][0], 10, 64)
	m := model.MirrorStation{}
	m.Id = uint64(queryId)
	name, _ := m.GetName()
	alive, _ := m.IsAlive()
	url, _ := m.GetUrl()
	mirrors, _ := m.GetMirrors()
	for alive && len(mirrors) == 0 {
		alive, _ = m.IsAlive()
		mirrors, _ = m.GetMirrors()
	}
	obj := struct {
		Id      uint64   `json:"id"`
		Name    string   `json:"name"`
		Url     string   `json:"url"`
		Alive   bool     `json:"alive"`
		Mirrors []uint64 `json:"mirrors"`
	}{
		Id:      m.Id,
		Name:    name,
		Url:     url,
		Alive:   alive,
		Mirrors: []uint64{},
	}
	for _, mirror := range mirrors {
		obj.Mirrors = append(obj.Mirrors, mirror.Id)
	}
	data, _ := json.Marshal(obj)
	w.Write(data)
}
