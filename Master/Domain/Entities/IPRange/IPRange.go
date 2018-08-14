package IPRange

import (
	_ "database/sql"
	"kaleido/Common/Services/IPTools"
	"kaleido/Master/Infrastructure/DB"
	_ "kaleido/Master/Infrastructure/DB"
	"strconv"
)

type Entity struct {
	Address    string
	MaskLength uint8
}

type RepoType map[uint32]*Entity

var Repo RepoType

var JHashMap map[string][]*Entity

func jhash(arr []uint32) string {
	var h uint32 = 8388617
	l := len(arr)
	for i := 0; i < l; i++ {
		h = ((h<<1 | h>>30) & 0x7fffffff) ^ arr[i]
	}
	return strconv.FormatUint(
		uint64(h), 36,
	)
}

func (entity Entity) NumberFormAddress() uint32 {
	return IPTools.IPv4ToNumberForm(entity.Address)
}

func (entity Entity) Jhash() string {
	return jhash([]uint32{IPTools.MaskIP(entity.Address, entity.MaskLength), uint32(entity.MaskLength)})
}

func init() {
	Repo = make(map[uint32]*Entity)
	JHashMap = make(map[string][]*Entity)
	rows, _ := DB.DB.Query(`
	select id,host(ip.data), masklen(ip.data) from ip;
	`)
	for rows.Next() {
		var id uint32
		var entity Entity
		rows.Scan(&id, &entity.Address, &entity.MaskLength)
		Repo[id] = &entity
		JHashMap[entity.Jhash()] = append(JHashMap[entity.Jhash()], &entity)
	}
}
