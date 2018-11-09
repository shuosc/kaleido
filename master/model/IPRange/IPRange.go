package IPRange

import (
	"database/sql"
	_ "database/sql"
	"kaleido/common/iptools"
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

func (iprange IPRange) GetAreaISPWithTransaction(tx *sql.Tx) (area Area.Area, isp ISP.ISP, err error) {
	row := tx.QueryRow(`
	SELECT area_id,isp_id FROM iprange_area_isp WHERE iprange_id=$1;
	`, iprange.Id)
	err = row.Scan(&area.Id, &isp.Id)
	return area, isp, err
}

func (iprange IPRange) GetUint64FormatWithTransaction(tx *sql.Tx) (uint64, error) {
	var ip string
	var mask uint8
	row := tx.QueryRow(`
	SELECT host(ip),masklen(ip) FROM iprange where id=$1;
	`, iprange.Id)
	if err := row.Scan(&ip, &mask); err != nil {
		return 0, err
	}
	return iptools.MergedMaskedIP(ip, mask)
}

func AllWithTransaction(tx *sql.Tx) (result []IPRange, err error) {
	rows, err := tx.Query(`
	SELECT id FROM iprange;
	`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var iprange IPRange
		err := rows.Scan(&iprange.Id)
		if err != nil {
			return nil, err
		}
		result = append(result, iprange)
	}
	return result, err
}
