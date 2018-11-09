package service

import (
	"database/sql"
	_ "database/sql"
	"kaleido/common/message"
	"kaleido/common/tools"
	"kaleido/master/DB"
	"kaleido/master/model/Area"
	"kaleido/master/model/IPRange"
	"kaleido/master/model/ISP"
	"kaleido/master/model/Mirror"
	"kaleido/master/model/MirrorStation"
)

func makeTable() (KaleidoMessage.KaleidoMessage, error) {
	var result KaleidoMessage.KaleidoMessage
	tx, err := DB.DB.Begin()
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
	stations, err := MirrorStation.AllWithTranscation(tx)
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
	mirrors, err := Mirror.AllWithTransaction(tx)
	if err != nil {
		return nil, err
	}
	for _, mirror := range mirrors {
		name, err := mirror.GetNameWithTransaction(tx)
		if err != nil {
			return nil, err
		}
		result[name], err = makeMirror(tx, mirror)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func makeMirror(tx *sql.Tx, mirror Mirror.Mirror) (result *KaleidoMessage.Mirror, err error) {
	result = new(KaleidoMessage.Mirror)
	result.FallbackMirrorStationId, err = mirror.GetFallbackMirrorStationIdWithTransaction(tx)
	if err != nil {
		return nil, err
	}
	result.AreaISP_MirrorStationGroup, err = GetAreaISPToMirrorStationGroupMap(mirror, tx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetAreaISPToMirrorStationGroupMap(mirror Mirror.Mirror, tx *sql.Tx) (map[uint64]*KaleidoMessage.MirrorStationGroup, error) {
	result := map[uint64]*KaleidoMessage.MirrorStationGroup{}
	areas, err := Area.AllWithTransaction(tx)
	if err != nil {
		return nil, err
	}
	isps, err := ISP.AllWithTranscation(tx)
	if err != nil {
		return nil, err
	}
	for _, area := range areas {
		for _, isp := range isps {
			key := tools.PackUInt32(uint32(area.Id), uint32(isp.Id))
			result[key] = &KaleidoMessage.MirrorStationGroup{}
			stations, err := mirror.GetMirrorStationsForAreaAndISPWithTransaction(area, isp, tx)
			if err != nil {
				return nil, err
			}
			for _, station := range stations {
				result[key].Stations = append(result[key].Stations, station.GetId())
			}
		}
	}
	return result, nil
}

func uploadTable() {

}
