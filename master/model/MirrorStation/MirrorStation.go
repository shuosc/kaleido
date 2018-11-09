package MirrorStation

import (
	"database/sql"
	_ "database/sql"
	"kaleido/master/DB"
	"kaleido/master/model/Mirror"
	"log"
	"sync"
)

type MirrorStation interface {
	SyncMirrorList() (bool, error)
	GetURL() (string, error)
	AddIPRange(ip string) error
	GetId() uint64
	GetURLWithTransaction(tx *sql.Tx) (string, error)
}

type Base struct {
	Id uint64
}

func (station Base) GetId() uint64 {
	return station.Id
}

func (station Base) getName() (string, error) {
	var result string
	row := DB.DB.QueryRow(`
	SELECT name from mirrorstation where id=$1;
	`, station.Id)
	if err := row.Scan(&result); err != nil {
		return "", err
	}
	return result, nil
}

func (station Base) getURL() (string, error) {
	var result string
	row := DB.DB.QueryRow(`
	SELECT url from mirrorstation where id=$1;
	`, station.Id)
	if err := row.Scan(&result); err != nil {
		return "", err
	}
	return result, nil
}

func (station Base) GetURLWithTransaction(tx *sql.Tx) (string, error) {
	var result string
	row := tx.QueryRow(`
	SELECT url from mirrorstation where id=$1;
	`, station.Id)
	if err := row.Scan(&result); err != nil {
		return "", err
	}
	return result, nil
}

func (station Base) addMirror(mirror Mirror.Mirror) error {
	_, err := DB.DB.Exec(`
	INSERT INTO mirrorstation_mirror(mirrorstation_id, mirror_id) VALUES ($1,$2);
	`, station.Id, mirror.Id)
	return err
}

func (station Base) removeMirror(mirror Mirror.Mirror) error {
	_, err := DB.DB.Exec(`
	DELETE FROM mirrorstation_mirror WHERE mirrorstation_id=$1 AND mirror_id=$2;
	`, station.Id, mirror.Id)
	return err
}

func (station Base) getMirrors() ([]Mirror.Mirror, error) {
	var result []Mirror.Mirror
	rows, err := DB.DB.Query(`
		SELECT mirror_id FROM mirrorstation_mirror where mirrorstation_id=$1;
		`, station.Id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var mirror Mirror.Mirror
		rows.Scan(&mirror.Id)
		result = append(result, mirror)
	}
	return result, nil
}

func (station Base) isMirrorIgnored(mirrorName string) (bool, error) {
	var result bool
	row := DB.DB.QueryRow(`
	select exists(
         select * from mirrorignore where mirrorstation_id = $1
                                      and name = $2)
	`, station.Id, mirrorName)
	if err := row.Scan(&result); err != nil {
		return true, err
	}
	return result, nil
}

func (station Base) setAlive(alive bool) error {
	_, err := DB.DB.Exec(`
	UPDATE mirrorstation SET alive=$1 WHERE id=$2;
	`, alive, station.Id)
	return err
}

func (station Base) GetURL() (string, error) {
	var result string
	row := DB.DB.QueryRow(`
	SELECT url from mirrorstation where id=$1;
	`, station.Id)
	if err := row.Scan(&result); err != nil {
		return "", err
	}
	return result, nil
}

func (station Base) AddIPRange(ip string) error {
	_, err := DB.DB.Exec(`
	INSERT INTO mirrorstation_iprange (mirrorstation_id, iprange_id)
	VALUES ($1, (SELECT id FROM iprange WHERE ip >> $2));
	`, station.Id, ip)
	return err
}

func toSet(list []Mirror.Mirror) map[Mirror.Mirror]bool {
	result := map[Mirror.Mirror]bool{}
	for _, mirror := range list {
		result[mirror] = true
	}
	return result
}

func (station Base) analyzeMirrorList(mirrors []Mirror.Mirror) (bool, error) {
	oldMirrors, err := station.getMirrors()
	if err != nil {
		return false, err
	}
	var modified bool
	if len(mirrors) != len(oldMirrors) {
		modified = true
	}
	oldMirrorSet := toSet(oldMirrors)
	newMirrorSet := toSet(mirrors)
	for mirror := range oldMirrorSet {
		// 镜像站删除了镜像
		if !newMirrorSet[mirror] {
			modified = true
			station.removeMirror(mirror)
		}
	}
	for mirror := range newMirrorSet {
		// 镜像站新增了镜像
		if !oldMirrorSet[mirror] {
			modified = true
			station.addMirror(mirror)
		}
	}
	return modified, nil
}

func All() ([]MirrorStation, error) {
	var result []MirrorStation
	webIndexed, err := allWebIndexed()
	if err != nil {
		return nil, err
	}
	jsonIndexed, err := allJsonIndexed()
	if err != nil {
		return nil, err
	}
	for _, station := range webIndexed {
		result = append(result, station)
	}
	for _, station := range jsonIndexed {
		result = append(result, station)
	}
	return result, nil
}

func AllWithTranscation(tx *sql.Tx) ([]MirrorStation, error) {
	var result []MirrorStation
	webIndexed, err := allWebIndexedWithTransaction(tx)
	if err != nil {
		return nil, err
	}
	jsonIndexed, err := allJsonIndexedWithTransaction(tx)
	if err != nil {
		return nil, err
	}
	for _, station := range webIndexed {
		result = append(result, station)
	}
	for _, station := range jsonIndexed {
		result = append(result, station)
	}
	return result, nil
}

func InitialSync() {
	stations, err := All()
	if err != nil {
		panic("Init stations failed!")
	}
	var wg sync.WaitGroup
	for _, station := range stations {
		wg.Add(1)
		go func(s MirrorStation) {
			defer wg.Done()
			s.SyncMirrorList()
		}(station)
	}
	wg.Wait()
	log.Println("Initial crawl success!")
}

func Get(id uint64) (MirrorStation, error) {
	jsonIndexed, err := getJsonIndexed(id)
	if err == nil {
		return jsonIndexed, nil
	}
	webIndexed, err := getWebIndexed(id)
	if err == nil {
		return webIndexed, nil
	}
	return nil, err
}

func GetWithTransaction(id uint64, tx *sql.Tx) (MirrorStation, error) {
	jsonIndexed, err := getJsonIndexedWithTransaction(id, tx)
	if err == nil {
		return jsonIndexed, nil
	}
	webIndexed, err := getWebIndexedWithTransaction(id, tx)
	if err == nil {
		return webIndexed, nil
	}
	return nil, err
}
