package model

import (
	_ "database/sql"
	"kaleido/common/model"
	"kaleido/master/infrastructure"
)

type Mirror struct {
	model.Entity
}

func (mirror Mirror) GetName() (string, error) {
	row := infrastructure.DB.QueryRow(`
	select name from mirror where id=$1;
	`, mirror.Id)
	var result string
	err := row.Scan(&result)
	return result, err
}

func (mirror Mirror) GetStationsContainThis() ([]MirrorStation, error) {
	rows, err := infrastructure.DB.Query(`
	select mirrorstation_id from mirrorstation_mirror where mirror_id=$1 order by mirrorstation_id;
	`, mirror.Id)
	if err != nil {
		return nil, err
	}
	var result []MirrorStation
	for rows.Next() {
		station := MirrorStation{}
		err := rows.Scan(&station.Id)
		if err != nil {
			return nil, err
		}
		result = append(result, station)
	}
	return result, nil
}

type mirrorRepo struct {
}

func (repo mirrorRepo) GetAll() []Mirror {
	rows, _ := infrastructure.DB.Query(`
	select id from mirror order by name;
	`)
	var result []Mirror
	for rows.Next() {
		var id uint64
		rows.Scan(&id)
		result = append(result, Mirror{
			model.Entity{
				Id: id,
			},
		})
	}
	return result
}

var MirrorRepo = mirrorRepo{}
