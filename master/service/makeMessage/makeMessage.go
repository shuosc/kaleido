package makeMessage

import (
	_ "database/sql"
	"github.com/lib/pq"
	"kaleido/common/service"
	"kaleido/common/service/message"
	"kaleido/master/infrastructure"
)

func makeStationUrls() (map[uint32]string, error) {
	result := map[uint32]string{}
	rows, err := infrastructure.DB.Query(`
	select id,url from mirrorstation;
	`)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var id uint32
		var url string
		rows.Scan(&id, &url)
		result[id] = url
	}
	return result, nil
}

func makeMaskAddressAreaID() (map[uint32]*KaleidoMessage.Address_AreaId, error) {
	result := map[uint32]*KaleidoMessage.Address_AreaId{}
	rows, err := infrastructure.DB.Query(`
	select masklen(iprange.ip), host(iprange.ip), iprange_area.area_id
	from iprange,iprange_area
	where iprange_area.iprange_id = iprange.id;
	`)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var mask uint32
		var host string
		var area uint32
		rows.Scan(&mask, &host, &area)
		addressAreaId, has := result[mask]
		if !has {
			addressAreaId = &KaleidoMessage.Address_AreaId{
				Address_AreaId: map[uint32]uint32{},
			}
			result[mask] = addressAreaId
		}
		hostNumberForm := service.IPv4ToNumberForm(host)
		addressAreaId.Address_AreaId[hostNumberForm] = area
	}
	return result, nil
}

func castIntoUint32(array pq.Int64Array) []uint32 {
	var data []uint32
	for _, num := range array {
		data = append(data, uint32(num))
	}
	return data
}

func makeMirrors() (map[string]*KaleidoMessage.Mirror, error) {
	result := map[string]*KaleidoMessage.Mirror{}
	rows, err := infrastructure.DB.Query(`
	select distinct on (mirror.id,
                   area.id) mirror.name,
                   area.id,
                   array_agg(
                     mirrorstation.id) over (partition by area.id, mirror.id, area_area.distance order by area_area.distance)
	from area,
	     mirror,
	     mirrorstation,
	     mirrorstation_mirror,
	     mirrorstation_area,
	     area_area
	where mirrorstation_mirror.mirrorstation_id = mirrorstation.id
	  and mirrorstation_mirror.mirror_id = mirror.id
	  and mirrorstation_area.mirrorstation_id = mirrorstation.id
	  and area_area.from_id = mirrorstation_area.area_id
	  and area_area.to_id = area.id
	group by mirror.id, area.id, mirrorstation.id, area_area.distance
	order by mirror.id, area.id, area_area.distance;
	`)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var mirrorName string
		var areaId uint32
		var mirrorStationId pq.Int64Array
		rows.Scan(&mirrorName, &areaId, &mirrorStationId)
		mirror, has := result[mirrorName]
		if !has {
			mirror = &KaleidoMessage.Mirror{
				DefaultMirrorStationId:    uint32(mirrorStationId[0]),
				AreaId_MirrorStationGroup: map[uint32]*KaleidoMessage.MirrorStationGroup{},
			}
			result[mirrorName] = mirror
		}
		mirror.AreaId_MirrorStationGroup[areaId] = &KaleidoMessage.MirrorStationGroup{
			Stations: castIntoUint32(mirrorStationId),
		}
	}
	return result, nil
}

func MakeMessage() (KaleidoMessage.KaleidoMessage, error) {
	result := KaleidoMessage.KaleidoMessage{}
	mirrorStationUrl, err := makeStationUrls()
	if err != nil {
		return result, err
	}
	result.MirrorStationId_Url = mirrorStationUrl
	maskAddressAreaID, err := makeMaskAddressAreaID()
	if err != nil {
		return result, err
	}
	result.Mask_Address_AreaID = maskAddressAreaID
	mirrors, err := makeMirrors()
	if err != nil {
		return result, err
	}
	result.Mirrors = mirrors
	return result, nil
}
