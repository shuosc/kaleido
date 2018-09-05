package Mirror

import (
	"database/sql"
	_ "database/sql"
	"github.com/lib/pq"
	"kaleido/Master/Domain/Entities/MirrorStation"
	"kaleido/Master/Infrastructure/DB"
	"log"
)

type entity interface {
	GetDefaultMirrorStationId() uint32
	GetAreaToMirrorStationTable(areaId uint32) []uint32
}

type Mirror struct {
	defaultMirrorStationId uint32
	stationsHaveThis       map[uint32]bool
}

func (m Mirror) GetDefaultMirrorStationId() uint32 {
	return m.defaultMirrorStationId
}

func (m Mirror) GetAreaToMirrorStationTable(areaId uint32) []uint32 {
	var mirrorStationHaveThis []uint32
	for mirrorStationId := range MirrorStation.Repo.Entities {
		if m.stationsHaveThis[mirrorStationId] {
			mirrorStationHaveThis = append(mirrorStationHaveThis, mirrorStationId)
		}
	}
	row := DB.DB.QueryRow(`
			select array_agg(mirrorstation_id)
			from area_mirrorstation
			where area_id = $1
				and mirrorstation_id = ANY($2)
			group by priority
			order by priority
			LIMIT 1;
		`, areaId, pq.Array(mirrorStationHaveThis))
	var nullInt64Result []sql.NullInt64
	row.Scan(pq.Array(&nullInt64Result))
	var result []uint32
	for _, n := range nullInt64Result {
		result = append(result, uint32(n.Int64))
	}
	return result
}

var Repo struct {
	Entities map[string]entity
}

func init() {
	CronJob()
	log.Println("Mirror inited")
}

func CronJob() {
	Repo.Entities = map[string]entity{}
	for stationId, station := range MirrorStation.Repo.Entities {
		if station.IsAlive() {
			for _, mirrorName := range station.GetMirrorList() {
				if mirror, has := Repo.Entities[mirrorName]; !has || mirror == nil {
					Repo.Entities[mirrorName] = Mirror{
						stationId,
						map[uint32]bool{},
					}
				}
				Repo.Entities[mirrorName].(Mirror).stationsHaveThis[stationId] = true
			}
		}
	}
}
