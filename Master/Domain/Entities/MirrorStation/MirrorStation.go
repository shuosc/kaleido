package MirrorStation

import (
	_ "database/sql"
	"encoding/json"
	"kaleido/Master/Infrastructure/DB"
	"log"
	"net/http"
)

type entity interface {
	GetUrl() string
	IsAlive() bool
	GetMirrorList() []string
	CheckAlive()
}

type mirrorStationBase struct {
	url   string
	alive bool
}

func (station *mirrorStationBase) CheckAlive() {
	response, err := http.Head(station.url)
	if err != nil || response.StatusCode != 200 {
		// In case some mirrorstation doesn't support HEAD request
		// e.g http://mirrors.hust.edu.cn
		response, err = http.Get(station.url)
	}
	if err != nil || response.StatusCode != 200 {
		station.alive = false
	} else {
		station.alive = true
	}
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
						url:   url,
						alive: false,
					},
					"https://mirrors.shu.edu.cn/data/mirrors.json",
					ShuJsonStructure{},
				}
			case 3:
				Repo.Entities[id] = &JsonIndexMirrorStation{
					mirrorStationBase{
						url:   url,
						alive: false,
					},
					"https://mirrors.tuna.tsinghua.edu.cn/static/tunasync.json",
					TunaJsonStructure{},
				}
			}
		} else {
			Repo.Entities[id] = &WebPageIndexMirrorStation{
				mirrorStationBase{
					url,
					false,
				},
				selector,
			}
		}
	}
	CronJob()
	log.Println("MirrorStation inited")
}

func CronJob() {
	ch := make(chan int, len(Repo.Entities))
	for _, entity := range Repo.Entities {
		go func(ch chan int) {
			entity.CheckAlive()
			ch <- 1
		}(ch)
	}
	for range Repo.Entities {
		<-ch
	}
}
