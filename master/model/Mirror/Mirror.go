package Mirror

import (
	_ "database/sql"
	"errors"
	"kaleido/master/DB"
)

type Mirror struct {
	Id uint64
}

func New(name string) (Mirror, error) {
	var result Mirror
	row := DB.DB.QueryRow(`
	INSERT INTO mirror(name) VALUES ($1) RETURNING id;
	`, name)
	if row.Scan(&result.Id) != nil {
		return result, errors.New("cannot create Mirror")
	}
	return result, nil
}

func GetOrCreate(name string) (Mirror, error) {
	var result Mirror
	row := DB.DB.QueryRow(`
	SELECT id FROM mirror WHERE name=$1;
	`, name)
	if err := row.Scan(&result.Id); err != nil {
		return New(name)
	}
	return result, nil
}
