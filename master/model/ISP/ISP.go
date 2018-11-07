package ISP

import (
	_ "database/sql"
	"kaleido/master/DB"
)

type ISP struct {
	Id uint64
}

func New(name string) ISP {
	var result ISP
	row := DB.DB.QueryRow(`
	INSERT INTO isp(name) VALUES ($1) RETURNING id;
	`, name)
	if row.Scan(&result.Id) != nil {
		panic("Cannot create ISP!")
	}
	return result
}
