package MirrorStation

import (
	"database/sql"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"kaleido/master/DB"
	"kaleido/master/model/Mirror"
	"net/http"
	"strings"
)

type WebIndexMirrorStation struct {
	Base
}

func (station WebIndexMirrorStation) getSelector() (string, error) {
	var result string
	row := DB.DB.QueryRow(`
	SELECT selector from webindexedmirrorstation where id=$1;
	`, station.Id)
	if err := row.Scan(&result); err != nil {
		return "", err
	}
	return result, nil
}

func (station WebIndexMirrorStation) fetchNewMirrorList() ([]Mirror.Mirror, error) {
	url, err := station.getURL()
	if err != nil {
		station.setAlive(false)
		return nil, err
	}
	response, err := http.Get(url)
	if err != nil || response.StatusCode != 200 {
		station.setAlive(false)
		return nil, err
	}
	pageDoc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		station.setAlive(false)
		return nil, err
	}
	selector, err := station.getSelector()
	if err != nil {
		station.setAlive(false)
		return nil, err
	}
	station.setAlive(true)
	var result []Mirror.Mirror
	pageDoc.Find(selector).Each(
		func(_ int, selection *goquery.Selection) {
			r, _ := selection.Attr("href")
			name := strings.Split(r, "/")
			var mirrorName string
			if r[len(r)-1] == '/' {
				mirrorName = name[len(name)-2]
			} else {
				mirrorName = name[len(name)-1]
			}
			ignored, err := station.isMirrorIgnored(mirrorName)
			if err != nil {
				return
			}
			if !ignored {
				mirror, err := Mirror.GetOrCreate(mirrorName)
				if err != nil {
					return
				}
				result = append(result, mirror)
			}
		})
	return result, nil
}

func (station WebIndexMirrorStation) SyncMirrorList() (bool, error) {
	mirrors, err := station.fetchNewMirrorList()
	if err != nil {
		return false, err
	}
	return station.analyzeMirrorList(mirrors)
}

func allWebIndexed() ([]WebIndexMirrorStation, error) {
	var result []WebIndexMirrorStation
	rows, err := DB.DB.Query(`
	SELECT id FROM webindexedmirrorstation;
	`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var station WebIndexMirrorStation
		if rows.Scan(&station.Id) != nil {
			return nil, err
		}
		result = append(result, station)
	}
	return result, nil
}

func allWebIndexedWithTransaction(tx *sql.Tx) ([]WebIndexMirrorStation, error) {
	var result []WebIndexMirrorStation
	rows, err := tx.Query(`
	SELECT id FROM webindexedmirrorstation;
	`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var station WebIndexMirrorStation
		if rows.Scan(&station.Id) != nil {
			return nil, err
		}
		result = append(result, station)
	}
	return result, nil
}

func getWebIndexed(id uint64) (JsonIndexMirrorStation, error) {
	var result bool
	row := DB.DB.QueryRow(`
	SELECT exists(select id FROM webindexedmirrorstation where id=$1);
	`, id)
	row.Scan(&result)
	if result {
		return JsonIndexMirrorStation{Base{id}}, nil
	}
	return JsonIndexMirrorStation{}, errors.New("not found")
}

func getWebIndexedWithTransaction(id uint64, tx *sql.Tx) (JsonIndexMirrorStation, error) {
	var result bool
	row := tx.QueryRow(`
	SELECT exists(select id FROM webindexedmirrorstation where id=$1);
	`, id)
	row.Scan(&result)
	if result {
		return JsonIndexMirrorStation{Base{id}}, nil
	}
	return JsonIndexMirrorStation{}, errors.New("not found")
}
