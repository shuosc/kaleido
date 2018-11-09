package ISP

import (
	"database/sql"
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

func GetOrCreate(name string) ISP {
	var result ISP
	row := DB.DB.QueryRow(`
	SELECT id FROM isp WHERE name=$1;
	`, name)
	if err := row.Scan(&result.Id); err != nil {
		return New(name)
	}
	return result
}

func AllWithTranscation(tx *sql.Tx) ([]ISP, error) {
	var result []ISP
	rows, err := tx.Query(`
	SELECT id from isp;
	`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var item ISP
		rows.Scan(&item.Id)
		result = append(result, item)
	}
	return result, nil
}
