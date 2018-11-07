package initDB

import (
	_ "database/sql"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"kaleido/master/model/Area"
	"kaleido/master/model/IPRange"
	"kaleido/master/model/ISP"
	"math"
	"net/http"
	"strings"
	"sync"
)

type ispPageInfo struct {
	Doc  *goquery.Document
	Name string
}

func parseISPPage(pageInfo ispPageInfo) {
	isp := ISP.New(pageInfo.Name)
	var area Area.Area
	pageInfo.Doc.Find(".list").Each(func(index int, selection *goquery.Selection) {
		selection.Find("dd,dt").Each(func(_ int, selectedLine *goquery.Selection) {
			selected := strings.Trim(selectedLine.Text(), " \n")
			if selected != "" {
				if strings.Contains(selected, ".") {
					// is an IPRange
					from := selectedLine.Find(".v_l").Text()
					to := selectedLine.Find(".v_r").Text()
					ipRange := IPRange.NewIPV4WithIPRange(from, to)
					ipRange.SetAreaISP(area, isp)
				} else {
					// is name of a city
					areaName := selected[:strings.Index(selected, "(")]
					area = Area.GetOrCreate(areaName)
				}
			}
		})
	})
}

func downloadISPPage(url string) ispPageInfo {
	response, _ := http.Get(url)
	doc, _ := goquery.NewDocumentFromReader(response.Body)
	ispName := url[strings.LastIndex(url, "/i_")+3:]
	return ispPageInfo{doc, ispName}
}

func initAreas() {
	response, _ := http.Get("http://ipcn.chacuo.net")
	doc, _ := goquery.NewDocumentFromReader(response.Body)
	selection := doc.Find(".section_content > .list a")
	channels := make(chan ispPageInfo, selection.Size())
	selection.Map(func(index int, anchor *goquery.Selection) string {
		href, _ := anchor.Attr("href")
		go func(ch chan ispPageInfo) {
			channels <- downloadISPPage(href)
		}(channels)
		return ""
	})
	for i := 0; i < selection.Size(); i++ {
		info := <-channels
		parseISPPage(info)
	}
}

func getAreaCoordinateURL(areaName string) string {
	return "http://api.map.baidu.com/geocoder?address=" + areaName +
		"&output=json&key=37492c0ee6f924cb5e934fa08c6b1676&city=北京市"
}

type coordinate struct {
	Lng float64 `json:"lng"`
	Lat float64 `json:"lat"`
}

func initAreaDistance() {
	areas := Area.All()
	areaCoordinate := fetchAreaCoordinates(areas)
	for _, areaFrom := range areas {
		for _, areaTo := range areas {
			if areaFrom == areaTo {
				Area.SetDistance(areaFrom, areaTo, 0)
			} else {
				fromCoordinate := areaCoordinate[areaFrom]
				toCoordinate := areaCoordinate[areaTo]
				distance := Distance(fromCoordinate.Lng, fromCoordinate.Lat, toCoordinate.Lng, toCoordinate.Lat)
				Area.SetDistance(areaFrom, areaTo, uint64(math.Ceil(distance)))
			}
		}
	}
}

func fetchAreaCoordinates(areas []Area.Area) map[Area.Area]coordinate {
	areaCoordinate := map[Area.Area]coordinate{}
	var areaCoordinateMutex sync.Mutex
	channel := make(chan int, len(areas))
	for _, area := range areas {
		go func(area_ Area.Area) {
			response, _ := http.Get(getAreaCoordinateURL(area_.GetName()))
			body, _ := ioutil.ReadAll(response.Body)
			var result struct {
				Result struct {
					Location coordinate `json:"location"`
				} `json:"result"`
			}
			json.Unmarshal(body, &result)
			areaCoordinateMutex.Lock()
			defer areaCoordinateMutex.Unlock()
			areaCoordinate[area_] = result.Result.Location
			channel <- 1
		}(area)
	}
	// 填百度接口的坑
	areaCoordinate[Area.GetOrCreate("海南")] = coordinate{Lng: 110.35, Lat: 20.02}
	areaCoordinate[Area.GetOrCreate("香港")] = coordinate{Lng: 114.1, Lat: 22.2}
	areaCoordinate[Area.GetOrCreate("河南")] = coordinate{Lng: 113.65, Lat: 34.76}
	for range areas {
		<-channel
	}
	return areaCoordinate
}

func InitAll() {
	initAreas()
	initAreaDistance()
	println("DB init success!")
}
