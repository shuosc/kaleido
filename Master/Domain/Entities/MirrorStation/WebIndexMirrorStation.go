package MirrorStation

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

type WebPageIndexMirrorStation struct {
	mirrorStationBase
	selector string
}

func (station WebPageIndexMirrorStation) GetMirrorList() []string {
	r, _ := http.Get(station.url)
	pageDoc, _ := goquery.NewDocumentFromReader(r.Body)
	result := pageDoc.Find(station.selector).Map(
		func(_ int, selection *goquery.Selection) string {
			r, _ := selection.Attr("href")
			result := strings.Split(r, "/")
			if r[len(r)-1] == '/' {
				return result[len(result)-2]
			}
			return result[len(result)-1]
		})
	return result
}
