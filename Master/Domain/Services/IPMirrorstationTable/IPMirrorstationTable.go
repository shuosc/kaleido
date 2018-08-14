package IPMirrorstationTable

import (
	_ "database/sql"
	"github.com/gogo/protobuf/proto"
	"kaleido/Common/Services/IPMirrorstationTableMessage"
	. "kaleido/Common/Services/IPTools"
	"kaleido/Master/Domain/Entities/IPRange"
	_ "kaleido/Master/Domain/Entities/IPRange"
	"kaleido/Master/Domain/Entities/MirrorStation"
	"kaleido/Master/Infrastructure/DB"
)

func Marshal() []byte {
	var table IPMirrorStationtableMessages.Table
	for ipId, ipEntity := range IPRange.Repo {
		rows, _ := DB.DB.Query(`
			select mirrorstation.id
			from ip_area,
     			area_mirrorstation,
     			mirrorstation
			where $1 = ip_area.ip_id
  			and ip_area.area_id = area_mirrorstation.area_id
  			and area_mirrorstation.mirrorstation_id = mirrorstation.id;`,
			ipId)
		relation := IPMirrorStationtableMessages.Relation{
			Ip:            IPv4ToNumberForm(ipEntity.Address),
			MaskBitLength: uint32(ipEntity.MaskLength),
		}
		for rows.Next() {
			var mirrorStationID uint32
			rows.Scan(&mirrorStationID)
			mirrorStationEntity := MirrorStation.Repo[mirrorStationID]
			if mirrorStationEntity.Alive {
				mirrorStationEntity.Mutex.Lock()
				relation.Url = append(relation.Url, mirrorStationEntity.Url)
				mirrorStationEntity.Mutex.Unlock()
			}
		}
		table.Relations = append(table.Relations, &relation)
	}
	result, _ := proto.Marshal(&table)
	return result
}
