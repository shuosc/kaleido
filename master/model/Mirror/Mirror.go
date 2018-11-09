package Mirror

import (
	"database/sql"
	_ "database/sql"
	"errors"
	"github.com/lib/pq"
	"kaleido/master/DB"
	"kaleido/master/model/Area"
	"kaleido/master/model/ISP"
	"kaleido/master/model/MirrorStation"
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

func (mirror Mirror) GetMirrorStationsForAreaAndISPWithTransaction(area Area.Area, isp ISP.ISP, tx *sql.Tx) ([]MirrorStation.MirrorStation, error) {
	row := tx.QueryRow(`
	select array_agg(mirrorstation_id)
	from (select distinct on (mirrorstation_mirror.mirrorstation_id) mirrorstation_mirror.mirrorstation_id,
	                         area_area.distance +
	                         case when isp.id = ispForMirrorStation.id then 0 else 10 end as network_distance
	      from mirrorstation_mirror,
	           mirrorstation_iprange,
	           area as areaForMirrorstation,
	           isp as ispForMirrorStation,
	           isp,
	           iprange_area_isp as iprange_area_ispForMirrorstation,
	           area_area
	      where mirrorstation_mirror.mirror_id = $1
	        AND mirrorstation_mirror.mirrorstation_id = mirrorstation_iprange.mirrorstation_id
	        AND mirrorstation_iprange.iprange_id = iprange_area_ispForMirrorstation.iprange_id
	        AND areaForMirrorstation.id = iprange_area_ispForMirrorstation.area_id
	        AND ispForMirrorStation.id = iprange_area_ispForMirrorstation.isp_id
	        AND area_area.from_id = iprange_area_ispForMirrorstation.area_id
	        and area_area.to_id = $2
	        and isp.id = $3
	      order by mirrorstation_mirror.mirrorstation_id) as distances
	group by distances.network_distance
	order by distances.network_distance
	LIMIT 1;
	`, mirror.Id, area.Id, isp.Id)
	queryResult := pq.Int64Array{}
	err := row.Scan(&queryResult)
	if err != nil {
		return nil, err
	}
	var result []MirrorStation.MirrorStation
	for _, id := range queryResult {
		station, err := MirrorStation.GetWithTransaction(uint64(id), tx)
		if err != nil {
			return nil, err
		}
		result = append(result, station)
	}
	return result, nil
}
