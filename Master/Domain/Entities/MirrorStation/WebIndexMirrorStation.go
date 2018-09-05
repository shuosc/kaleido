package MirrorStation

import (
	"github.com/PuerkitoBio/goquery"
	"kaleido/Master/Service/DirtyCheck"
	"log"
	"net/http"
	"sort"
	"strings"
)

type WebPageIndexMirrorStation struct {
	mirrorStationBase
	selector string
}

func (station *WebPageIndexMirrorStation) UpdateMirrorList() {
	response, err := http.Get(station.url)
	if err != nil || response.StatusCode != 200 {
		log.Println("Failed to fetch MirrorListUrl from ", station.url)
		station.mirrorList = []string{}
		return
	}
	pageDoc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Println("Failed to fetch MirrorListUrl from ", station.url)
		station.mirrorList = []string{}
		return
	}
	result := pageDoc.Find(station.selector).Map(
		func(_ int, selection *goquery.Selection) string {
			r, _ := selection.Attr("href")
			result := strings.Split(r, "/")
			if r[len(r)-1] == '/' {
				return result[len(result)-2]
			}
			return result[len(result)-1]
		})
	sort.Strings(result)
	oldListContent := strings.Join(station.mirrorList, "")

	station.mirrorList = result
	if oldListContent != strings.Join(station.mirrorList, "") {
		log.Println("Found content in ", station.url, " changed")
		DirtyCheck.Dirty = true
	}
}
