package HeartBeat

import (
	"kaleido/Master/Domain/Entities/MirrorStation"
	"log"
	"net/http"
)

func testOne(entity *MirrorStation.Entity) {
	entity.Mutex.Lock()
	defer entity.Mutex.Unlock()
	response, err := http.Head(entity.Url)
	if err != nil || response.StatusCode != 200 {
		// In case some mirrorstation doesn't support HEAD request
		// e.g http://mirrors.hust.edu.cn
		response, err = http.Get(entity.Url)
		if err != nil || response.StatusCode != 200 {
			if entity.Alive == true {
				log.Println("Found ", entity.Url, " dead!")
				entity.Alive = false
			}
		} else {
			if entity.Alive == false {
				log.Println("Found ", entity.Url, " come back to life!")
				entity.Alive = true
			}
		}
	} else {
		if entity.Alive == false {
			log.Println("Found ", entity.Url, " come back to life!")
			entity.Alive = true
		}
	}
}

func Test() {
	for _, entity := range MirrorStation.Repo {
		go testOne(entity)
	}
}
