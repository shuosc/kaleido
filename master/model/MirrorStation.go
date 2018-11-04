package model

import (
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"kaleido/common/model"
	"kaleido/master/infrastructure"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

type MirrorSupplier interface {
	Check() bool
	GetName() (string, error)
	GetId() uint64
}

type MirrorStation struct {
	model.Entity
	WasAlive      bool
	oldMirrorList []string
}

func (mirrorStation MirrorStation) GetId() uint64 {
	return mirrorStation.Id
}

func (mirrorStation MirrorStation) GetName() (string, error) {
	row := infrastructure.DB.QueryRow(`
	select name from mirrorstation where id=$1
	`, mirrorStation.Id)
	var result string
	err := row.Scan(&result)
	return result, err
}

func (mirrorStation MirrorStation) GetUrl() (string, error) {
	row := infrastructure.DB.QueryRow(`
	select url from mirrorstation where id=$1
	`, mirrorStation.Id)
	var result string
	err := row.Scan(&result)
	return result, err
}

func (mirrorStation MirrorStation) SetAlive(value bool) error {
	_, err := infrastructure.DB.Exec(`
	update mirrorstation
	set alive=$1
	where id=$2;
	`, value, mirrorStation.Id)
	return err
}

func (mirrorStation MirrorStation) IsAlive() (bool, error) {
	row := infrastructure.DB.QueryRow(`
	select alive from mirrorstation where id=$1
	`, mirrorStation.Id)
	var result bool
	err := row.Scan(&result)
	return result, err
}

func (mirrorStation MirrorStation) ContainsMirror(mirrorId uint64) bool {
	row := infrastructure.DB.QueryRow(`
		select exists(
			select from mirrorstation_mirror
			where mirrorstation_id=$1 and mirror_id=$2)
	`, mirrorStation.Id, mirrorId)
	var exists bool
	row.Scan(&exists)
	return exists
}

func (mirrorStation MirrorStation) GetMirrors() ([]Mirror, error) {
	rows, err := infrastructure.DB.Query(`
	select mirror_id from mirrorstation_mirror where mirrorstation_id=$1 order by mirror_id;
	`, mirrorStation.Id)
	if err != nil {
		return nil, err
	}
	var result []Mirror
	for rows.Next() {
		mirror := Mirror{}
		err := rows.Scan(&mirror.Id)
		if err != nil {
			return nil, err
		}
		result = append(result, mirror)
	}
	return result, nil
}

func (mirrorStation MirrorStation) addMirror(mirrorName string) {
	var mirrorId uint64
	row := infrastructure.DB.QueryRow(`
		select id 
		from mirror
		where name=$1;
		`, mirrorName)
	err := row.Scan(&mirrorId)
	if err != nil {
		log.Println("Find new mirror:", mirrorName)
		row = infrastructure.DB.QueryRow(`
			insert into mirror(name) values ($1) returning id;
			`, mirrorName)
		row.Scan(&mirrorId)
	}
	_, err = infrastructure.DB.Exec(`
		insert into mirrorstation_mirror (mirrorstation_id, mirror_id) values ($1,$2);
		`, mirrorStation.Id, mirrorId)
}

func (mirrorStation MirrorStation) isMirrorIgnored(name string) bool {
	var result bool
	row := infrastructure.DB.QueryRow(`
	select exists(
         select * from mirrorignore where mirrorstationid = $1
                                      and mirrorname = $2
           ) 
	`, mirrorStation.Id, name)
	row.Scan(&result)
	return result
}

type webIndexMirrorStation struct {
	MirrorStation
}

func (mirrorStation webIndexMirrorStation) GetSelector() (string, error) {
	row := infrastructure.DB.QueryRow(`
	select selector from mirrorstation where id=$1
	`, mirrorStation.Id)
	var result string
	err := row.Scan(&result)
	return result, err
}

func (mirrorStation *webIndexMirrorStation) Check() bool {
	infrastructure.DB.Exec(`
	delete from mirrorstation_mirror where mirrorstation_id=$1;
	`, mirrorStation.Id)
	url, err := mirrorStation.GetUrl()
	if err != nil {
		log.Fatal("DB connection lost!")
		return false
	}
	response, err := http.Get(url)
	if err != nil || response.StatusCode != 200 {
		log.Println("Failed to fetch MirrorListUrl from ", url)
		mirrorStation.SetAlive(false)
		if mirrorStation.WasAlive {
			mirrorStation.WasAlive = false
			return true
		}
		return false
	}
	pageDoc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Println("Failed to fetch MirrorListUrl from ", url)
		mirrorStation.SetAlive(false)
		if mirrorStation.WasAlive {
			mirrorStation.WasAlive = false
			return true
		}
		return false
	}
	selector, err := mirrorStation.GetSelector()
	if err != nil {
		log.Fatal("DB connection lost!")
		return false
	}
	changed := !mirrorStation.WasAlive
	mirrorStation.WasAlive = true
	mirrorStation.SetAlive(true)
	var newMirrorList []string
	pageDoc.Find(selector).Each(
		func(_ int, selection *goquery.Selection) {
			r, _ := selection.Attr("href")
			result := strings.Split(r, "/")
			var mirrorName string
			if r[len(r)-1] == '/' {
				mirrorName = result[len(result)-2]
			} else {
				mirrorName = result[len(result)-1]
			}
			if !mirrorStation.isMirrorIgnored(mirrorName) {
				mirrorStation.addMirror(mirrorName)
				newMirrorList = append(newMirrorList, mirrorName)
			}
		})
	sort.Strings(newMirrorList)
	if len(newMirrorList) != len(mirrorStation.oldMirrorList) {
		changed = true
	} else {
		for idx := range newMirrorList {
			if newMirrorList[idx] != mirrorStation.oldMirrorList[idx] {
				changed = false
				break
			}
		}
	}
	mirrorStation.oldMirrorList = newMirrorList
	if changed {
		log.Println("Find change in ", url)
		return true
	}
	return false
}

type mirrorStationRepo struct {
}

type jsonIndexMirrorStation struct {
	MirrorStation
}

func (mirrorStation *jsonIndexMirrorStation) Check() bool {
	tx, _ := infrastructure.DB.Begin()
	tx.Exec(`
	delete from mirrorstation_mirror where mirrorstation_id=$1;
	`, mirrorStation.Id)
	url, err := mirrorStation.GetJsonUrl()
	if err != nil {
		log.Fatal("DB connection lost!")
		tx.Commit()
		return false
	}
	response, err := http.Get(url)
	if err != nil || response.StatusCode != 200 {
		log.Println("Failed to fetch MirrorListUrl from ", url)
		mirrorStation.SetAlive(false)
		tx.Commit()
		if mirrorStation.WasAlive {
			mirrorStation.WasAlive = false
			return true
		}
		return false
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Failed to fetch MirrorListUrl from ", url)
		mirrorStation.SetAlive(false)
		tx.Commit()
		if mirrorStation.WasAlive {
			mirrorStation.WasAlive = false
			return true
		}
		return false
	}
	changed := !mirrorStation.WasAlive
	mirrorStation.WasAlive = true
	mirrorStation.SetAlive(true)
	expression := regexp.MustCompile(`"name":\s*"(\S+?)"`)
	matches := expression.FindAllStringSubmatch(string(body), -1)
	var newMirrorList []string
	for _, match := range matches {
		if !mirrorStation.isMirrorIgnored(match[1]) {
			mirrorStation.addMirror(match[1])
			newMirrorList = append(newMirrorList, match[1])
		}
	}
	sort.Strings(newMirrorList)
	if len(newMirrorList) != len(mirrorStation.oldMirrorList) {
		changed = true
	} else {
		for idx := range newMirrorList {
			if newMirrorList[idx] != mirrorStation.oldMirrorList[idx] {
				changed = false
				break
			}
		}
	}
	mirrorStation.oldMirrorList = newMirrorList
	tx.Commit()
	if changed {
		log.Println("Find change in ", url)
		return true
	}
	return false
}

func (mirrorStation jsonIndexMirrorStation) GetJsonUrl() (string, error) {
	var url string
	row := infrastructure.DB.QueryRow(`
		select case
         when substr(jsonurl, 1, 4) = 'http' then jsonurl
         else concat(url, jsonurl)
           end
	from mirrorstation
	where id = $1;`, mirrorStation.Id)
	err := row.Scan(&url)
	return url, err
}

func (repo mirrorStationRepo) GetAll() []MirrorSupplier {
	rows, _ := infrastructure.DB.Query(`
	select id from mirrorstation where selector is not null;
	`)
	var result []MirrorSupplier
	for rows.Next() {
		var id uint64
		rows.Scan(&id)
		result = append(result, &webIndexMirrorStation{MirrorStation{
			model.Entity{Id: id},
			false,
			[]string{},
		}})
	}
	rows, _ = infrastructure.DB.Query(`
	select id from mirrorstation where jsonurl is not null;
	`)
	for rows.Next() {
		var id uint64
		rows.Scan(&id)
		result = append(result, &jsonIndexMirrorStation{MirrorStation{
			model.Entity{Id: id},
			false,
			[]string{},
		}})
	}
	return result
}

var MirrorStationRepo = mirrorStationRepo{}
