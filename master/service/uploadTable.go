package service

import (
	"bytes"
	"database/sql"
	_ "database/sql"
	"github.com/lib/pq"
	"kaleido/common/message"
	"kaleido/common/oss"
	"kaleido/common/tools"
	"kaleido/master/DB"
	"kaleido/master/model/IPRange"
	"kaleido/master/model/MirrorStation"
	"log"
)

func makeTable() (*KaleidoMessage.KaleidoMessage, error) {
	result := new(KaleidoMessage.KaleidoMessage)
	tx, err := DB.DB.Begin()
	defer tx.Commit()
	if err != nil {
		return result, err
	}
	_, err = tx.Exec(`
	SET TRANSACTION ISOLATION LEVEL REPEATABLE READ;
	`)
	if err != nil {
		return result, err
	}
	result.Mirrors, err = makeMirrors(tx)
	if err != nil {
		return result, err
	}
	result.MirrorStationId_Url, err = makeMirrorStationIdToUrl(tx)
	if err != nil {
		return result, err
	}
	result.Address_AreaISP, err = makeAddressToAreaISP(tx)
	return result, nil
}

func makeAddressToAreaISP(tx *sql.Tx) (map[uint64]uint64, error) {
	result := map[uint64]uint64{}
	ipRanges, err := IPRange.AllWithTransaction(tx)
	if err != nil {
		return nil, err
	}
	for _, ipRange := range ipRanges {
		u64Form, err := ipRange.GetUint64FormatWithTransaction(tx)
		if err != nil {
			return nil, err
		}
		area, isp, err := ipRange.GetAreaISPWithTransaction(tx)
		if err != nil {
			return nil, err
		}
		result[u64Form] = tools.PackUInt32(uint32(area.Id), uint32(isp.Id))
	}
	return result, nil
}

func makeMirrorStationIdToUrl(tx *sql.Tx) (map[uint64]string, error) {
	result := map[uint64]string{}
	stations, err := MirrorStation.AllWithTransaction(tx)
	if err != nil {
		return nil, err
	}
	for _, station := range stations {
		result[station.GetId()], err = station.GetURLWithTransaction(tx)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func makeMirrors(tx *sql.Tx) (map[string]*KaleidoMessage.Mirror, error) {
	result := map[string]*KaleidoMessage.Mirror{}
	rows, err := tx.Query(`
	select mirror_name, station_group, to_area_id, isp_id
	from (select mirror.name                                                                              as mirror_name,
	             mirrorstation_mirror.mirror_id                                                           as mirror_id,
	             array_agg(mirrorstation_mirror.mirrorstation_id)                                         as station_group,
	             area_area.distance +
	             case when isp.id = ispForMirrorStation.id then 0 else 10 end                             as network_distance,
	             area_area.to_id                                                                          as to_area_id,
	             isp.id                                                                                   as isp_id,
	             row_number() over (partition by mirrorstation_mirror.mirror_id, area_area.to_id, isp.id) as rn
	      from mirrorstation_mirror,
	           mirrorstation_iprange,
	           isp as ispForMirrorStation,
	           isp,
	           iprange_area_isp as iprange_area_ispForMirrorstation,
	           area_area,
	           mirror
	      where mirrorstation_mirror.mirrorstation_id = mirrorstation_iprange.mirrorstation_id
	        AND mirrorstation_iprange.iprange_id = iprange_area_ispForMirrorstation.iprange_id
	        AND ispForMirrorStation.id = iprange_area_ispForMirrorstation.isp_id
	        AND area_area.from_id = iprange_area_ispForMirrorstation.area_id
	        AND mirror.id = mirrorstation_mirror.mirror_id
	      group by network_distance, area_area.to_id, isp.id, mirrorstation_mirror.mirror_id, mirror.name
	      order by mirrorstation_mirror.mirror_id, to_area_id, isp_id, network_distance) as all_data
	where rn = 1;
	`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var mirrorName string
		var area uint64
		var isp uint64
		var stationGroup pq.Int64Array
		err = rows.Scan(&mirrorName, &stationGroup, &area, &isp)
		if err != nil {
			return nil, err
		}
		if err != nil {
			return nil, err
		}
		mirrorObject, has := result[mirrorName]
		if !has {
			result[mirrorName] = &KaleidoMessage.Mirror{}
			mirrorObject = result[mirrorName]
		}
		mirrorObject.FallbackMirrorStationId = uint64(stationGroup[0])
		areaISP := tools.PackUInt32(uint32(area), uint32(isp))
		if mirrorObject.AreaISP_MirrorStationGroup == nil {
			mirrorObject.AreaISP_MirrorStationGroup = map[uint64]*KaleidoMessage.MirrorStationGroup{}
		}
		mirrorObject.AreaISP_MirrorStationGroup[areaISP] = &KaleidoMessage.MirrorStationGroup{}
		for _, stationId := range stationGroup {
			mirrorObject.AreaISP_MirrorStationGroup[areaISP].Stations = append(mirrorObject.AreaISP_MirrorStationGroup[areaISP].Stations, uint64(stationId))
		}
	}
	return result, nil
}

func uploadTable() {
	message, err := makeTable()
	if err != nil {
		log.Println(err)
	}
	data, err := message.Marshal()
	if err != nil {
		log.Println(err)
	}
	err = oss.Bucket.PutObject("kaleido-message", bytes.NewBuffer(data))
	if err != nil {
		log.Println(err)
	}
}
