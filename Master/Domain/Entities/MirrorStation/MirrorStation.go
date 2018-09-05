package MirrorStation

import (
	_ "database/sql"
	"encoding/json"
	"kaleido/Master/Infrastructure/DB"
	"kaleido/Master/Service/DirtyCheck"
	"log"
	"net/http"
	"sort"
)

type entity interface {
	GetUrl() string
	IsAlive() bool
	GetMirrorList() []string
	UpdateMirrorList()
	CheckAlive()
}

type mirrorStationBase struct {
	url        string
	alive      bool
	mirrorList []string
}

func (station *mirrorStationBase) CheckAlive() {
	response, err := http.Head(station.url)
	if err != nil || response.StatusCode != 200 {
		// In case some mirrorstation doesn't support HEAD request
		// e.g http://mirrors.hust.edu.cn
		response, err = http.Get(station.url)
	}
	lastAlive := station.alive
	if err != nil || response.StatusCode != 200 {
		station.alive = false
	} else {
		station.alive = true
	}
	if station.alive != lastAlive {
		if station.alive == false {
			log.Println("Found station", station.url, "dead")
		} else {
			log.Println("Found station", station.url, "come back to live")
		}
		DirtyCheck.Dirty = true
	}
}

func (station mirrorStationBase) GetMirrorList() []string {
	return station.mirrorList
}

func (station mirrorStationBase) GetUrl() string {
	return station.url
}

func (station mirrorStationBase) IsAlive() bool {
	return station.alive
}

var Repo struct {
	Entities map[uint32]entity
}

type ShuJsonStructure struct {
}

func (shuJsonStructure ShuJsonStructure) GetNamesForJson(buffer []byte) []string {
	var data struct {
		Mirrors []struct {
			Name string `json:"name"`
		} `json:"mirrors"`
	}
	json.Unmarshal(buffer, &data)
	result := make([]string, len(data.Mirrors))
	for index, mirror := range data.Mirrors {
		result[index] = mirror.Name
	}
	sort.Strings(result)
	return result
}

type TunaJsonStructure struct {
}

func (tunaJsonStructure TunaJsonStructure) GetNamesForJson(buffer []byte) []string {
	var data []struct {
		Name string `json:"name"`
	}
	json.Unmarshal(buffer, &data)
	result := make([]string, len(data))
	for index, mirror := range data {
		result[index] = mirror.Name
	}
	sort.Strings(result)
	return result
}

func init() {
	Repo.Entities = map[uint32]entity{}
	rows, _ := DB.DB.Query(`
		select id,url,selector from mirrorstation;
	`)
	for rows.Next() {
		var id uint32
		var url string
		var selector string
		rows.Scan(&id, &url, &selector)
		if selector == "" {
			switch id {
			case 1:
				Repo.Entities[id] = &JsonIndexMirrorStation{
					mirrorStationBase{
						url:        url,
						alive:      false,
						mirrorList: []string{},
					},
					"https://mirrors.shu.edu.cn/data/mirrors.json",
					ShuJsonStructure{},
				}
			case 3:
				Repo.Entities[id] = &JsonIndexMirrorStation{
					mirrorStationBase{
						url:        url,
						alive:      false,
						mirrorList: []string{},
					},
					"https://mirrors.tuna.tsinghua.edu.cn/static/tunasync.json",
					TunaJsonStructure{},
				}
			}
		} else {
			Repo.Entities[id] = &WebPageIndexMirrorStation{
				mirrorStationBase{
					url:        url,
					alive:      false,
					mirrorList: []string{},
				},
				selector,
			}
		}
	}
	log.Println("MirrorStation inited")
}

func CronJob() {
	ch := make(chan int, len(Repo.Entities))
	for _, e := range Repo.Entities {
		go func(entity entity) {
			entity.CheckAlive()
			entity.UpdateMirrorList()
			ch <- 1
		}(e)
	}
	for range Repo.Entities {
		<-ch
	}
}
