package Area

import (
	_ "database/sql"
	"kaleido/master/DB"
)

type Area struct {
	Id uint64
}

func (area Area) GetName() string {
	var result string
	row := DB.DB.QueryRow(`
	SELECT name from area where id=$1;
	`, area.Id)
	if row.Scan(&result) != nil {
		panic("Unable to get name!")
	}
	return result
}

func SetDistance(areaFrom Area, areaTo Area, distance uint64) {
	_, err := DB.DB.Exec(`
			INSERT INTO area_area(from_id, to_id, distance) VALUES ($1,$2,$3);
			`, areaFrom.Id, areaTo.Id, distance)
	if err != nil {
		panic("Cannot set distance!")
	}
}

func GetOrCreate(name string) Area {
	var result Area
	row := DB.DB.QueryRow(`
	SELECT id FROM area WHERE name=$1;
	`, name)
	if err := row.Scan(&result.Id); err != nil {
		return New(name)
	}
	return result
}

func New(name string) Area {
	var result Area
	row := DB.DB.QueryRow(`
	INSERT INTO area(name) VALUES ($1) RETURNING id;
	`, name)
	if row.Scan(&result.Id) != nil {
		panic("Cannot create Area!")
	}
	return result
}

func All() []Area {
	var result []Area
	rows, err := DB.DB.Query(`
	SELECT id from area;
	`)
	if err != nil {
		panic("Unable to get all areas!")
	}
	for rows.Next() {
		var item Area
		rows.Scan(&item.Id)
		result = append(result, item)
	}
	return result
}
