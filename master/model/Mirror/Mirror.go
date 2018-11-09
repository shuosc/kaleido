package Mirror

import (
	"database/sql"
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

func AllWithTransaction(tx *sql.Tx) (result []Mirror, err error) {
	rows, err := tx.Query(`
	SELECT id FROM mirror;
	`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var mirror Mirror
		err := rows.Scan(&mirror.Id)
		if err != nil {
			return nil, err
		}
		result = append(result, mirror)
	}
	return result, err
}

func All() (result []Mirror, err error) {
	rows, err := DB.DB.Query(`
	SELECT id FROM mirror;
	`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var mirror Mirror
		err := rows.Scan(&mirror.Id)
		if err != nil {
			return nil, err
		}
		result = append(result, mirror)
	}
	return result, err
}

func (mirror Mirror) GetName() (result string, err error) {
	row := DB.DB.QueryRow(`
	SELECT name FROM mirror WHERE id=$1;
	`, mirror.Id)
	err = row.Scan(&result)
	return result, err
}

func (mirror Mirror) GetNameWithTransaction(tx *sql.Tx) (result string, err error) {
	row := tx.QueryRow(`
	SELECT name FROM mirror WHERE id=$1;
	`, mirror.Id)
	err = row.Scan(&result)
	return result, err
}

func (mirror Mirror) GetFallbackMirrorStationIdWithTransaction(tx *sql.Tx) (uint64, error) {
	var result uint64
	row := tx.QueryRow(`
	SELECT mirrorstation_id from mirrorstation_mirror WHERE mirror_id=$1 LIMIT 1;
	`, mirror.Id)
	err := row.Scan(&result)
	return result, err
}
