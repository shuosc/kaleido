package MirrorStation

import (
	_ "database/sql"
	"kaleido/Master/Infrastructure/DB"
	"sync"
)

type Entity struct {
	Url   string
	Alive bool
	Mutex sync.RWMutex
}

var Repo map[uint32]*Entity

func init() {
	Repo = make(map[uint32]*Entity)
	rows, _ := DB.DB.Query(`
	SELECT id,url FROM mirrorstation;
	`)
	for rows.Next() {
		var id uint32
		var entity Entity
		rows.Scan(&id, &entity.Url)
		entity.Alive = true
		Repo[id] = &entity
	}
}
