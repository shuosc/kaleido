package IPRange

import (
	_ "database/sql"
	"kaleido/master/DB"
	"kaleido/master/model/Area"
	"kaleido/master/model/ISP"
)

type IPRange struct {
	Id uint64
}

func NewIPV4WithIPRange(from string, to string) IPRange {
	var result IPRange
	row := DB.DB.QueryRow(`
	INSERT INTO iprange (ip) VALUES (inet_merge($1,$2)) RETURNING id;
	`, from, to)
	if row.Scan(&result.Id) != nil {
		panic("Cannot create IPRange!")
	}
	return result
}

func (iprange IPRange) SetAreaISP(area Area.Area, isp ISP.ISP) {
	_, err := DB.DB.Exec(`
	INSERT INTO iprange_area_isp(iprange_id, area_id, isp_id) 
	VALUES ($1,$2,$3);
	`, iprange.Id, area.Id, isp.Id)
	if err != nil {
		panic("Failed to set iprange's area or isp")
	}
}
