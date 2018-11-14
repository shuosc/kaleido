package MirrorStation

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"kaleido/master/DB"
	"kaleido/master/model/Mirror"
	"net/http"
	"regexp"
)

type JsonIndexMirrorStation struct {
	Base
}

func (station JsonIndexMirrorStation) getIndexURL() (string, error) {
	var result string
	row := DB.DB.QueryRow(`
			select case
         			when substr(indexurl, 1, 4) = 'http' then indexurl
         			else concat(url, indexurl)
           		   end
			from jsonindexedmirrorstation
			where id = $1;`, station.Id)
	if err := row.Scan(&result); err != nil {
		return "", err
	}
	return result, nil
}

func (station JsonIndexMirrorStation) fetchNewMirrorList() ([]Mirror.Mirror, error) {
	url, err := station.getIndexURL()
	if err != nil {
		station.setAlive(false)
		return nil, err
	}
	response, err := http.Get(url)
	if err != nil {
		station.setAlive(false)
		return nil, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		station.setAlive(false)
		return nil, err
	}
	station.setAlive(true)
	expression := regexp.MustCompile(`"name":\s*"(\S+?)"`)
	matches := expression.FindAllStringSubmatch(string(body), -1)
	var result []Mirror.Mirror
	for _, match := range matches {
		ignored, err := station.isMirrorIgnored(match[1])
		if err != nil {
			return nil, err
		}
		if !ignored {
			mirror, err := Mirror.GetOrCreate(match[1])
			if err != nil {
				return nil, err
			}
			result = append(result, mirror)
		}
	}
	return result, nil
}

func (station JsonIndexMirrorStation) SyncMirrorList() (bool, error) {
	mirrors, err := station.fetchNewMirrorList()
	if err != nil {
		return false, err
	}
	return station.analyzeMirrorList(mirrors)
}

func allJsonIndexed() ([]JsonIndexMirrorStation, error) {
	var result []JsonIndexMirrorStation
	rows, err := DB.DB.Query(`
	SELECT id FROM jsonindexedmirrorstation;
	`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var station JsonIndexMirrorStation
		if rows.Scan(&station.Id) != nil {
			return nil, err
		}
		result = append(result, station)
	}
	return result, nil
}

func allJsonIndexedWithTransaction(tx *sql.Tx) ([]JsonIndexMirrorStation, error) {
	var result []JsonIndexMirrorStation
	rows, err := tx.Query(`
	SELECT id FROM jsonindexedmirrorstation;
	`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var station JsonIndexMirrorStation
		if rows.Scan(&station.Id) != nil {
			return nil, err
		}
		result = append(result, station)
	}
	return result, nil
}

func getJsonIndexed(id uint64) (JsonIndexMirrorStation, error) {
	var result bool
	row := DB.DB.QueryRow(`
	SELECT exists(select id FROM jsonindexedmirrorstation where id=$1);
	`, id)
	row.Scan(&result)
	if result {
		return JsonIndexMirrorStation{Base{id}}, nil
	}
	return JsonIndexMirrorStation{}, errors.New("not found")
}

func getJsonIndexedWithTransaction(id uint64, tx *sql.Tx) (JsonIndexMirrorStation, error) {
	var result bool
	row := tx.QueryRow(`
	SELECT exists(select id FROM jsonindexedmirrorstation where id=$1);
	`, id)
	row.Scan(&result)
	if result {
		return JsonIndexMirrorStation{Base{id}}, nil
	}
	return JsonIndexMirrorStation{}, errors.New("not found")
}
