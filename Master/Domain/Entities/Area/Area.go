package Area

import (
	_ "database/sql"
	"kaleido/Master/Infrastructure/DB"
)

type Entity struct {
}

var Repo struct {
	Entities map[uint32]Entity
}

func init() {
	Repo.Entities = map[uint32]Entity{}
	rows, _ := DB.DB.Query(`
		SELECT id from area;
	`)
	for rows.Next() {
		var id uint32
		rows.Scan(&id)
		Repo.Entities[id] = Entity{}
	}
}
