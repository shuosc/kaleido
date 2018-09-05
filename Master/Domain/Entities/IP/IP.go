package IP

import (
	_ "database/sql"
	"kaleido/Common/Services/IPTools"
	"kaleido/Master/Infrastructure/DB"
)

type entity interface {
	GetMaskBitLength() uint8
	GetAddressNumberForm() uint32
	GetAreaId() uint32
}

type ip struct {
	maskBitLength uint8
	address       string
	areaId        uint32
}

func (ipEntity ip) GetMaskBitLength() uint8 {
	return ipEntity.maskBitLength
}

func (ipEntity ip) GetAddressNumberForm() uint32 {
	result, _ := IPTools.IPv4ToNumberForm(ipEntity.address)
	return result
}

func (ipEntity ip) GetAreaId() uint32 {
	return ipEntity.areaId
}

var Repo struct {
	Entities map[uint32]entity
}

func init() {
	Repo.Entities = map[uint32]entity{}
	rows, _ := DB.DB.Query(`
		SELECT ip.id, host(ip.data), masklen(ip.data), ip_area.area_id
		FROM ip, ip_area
		WHERE ip_area.ip_id = ip.id;
	`)
	for rows.Next() {
		var id uint32
		var address string
		var mask uint8
		var areaId uint32
		rows.Scan(&id, &address, &mask, &areaId)
		Repo.Entities[id] = ip{
			mask,
			address,
			areaId,
		}
	}
}
