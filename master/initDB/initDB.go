package initDB

import (
	_ "database/sql"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"kaleido/master/DB"
	"kaleido/master/model/Area"
	"kaleido/master/model/IPRange"
	"kaleido/master/model/ISP"
	"kaleido/master/model/MirrorStation"
	"math"
	"net"
	"net/http"
	"strings"
	"sync"
)

type ispPageInfo struct {
	Doc  *goquery.Document
	Name string
}

func execSQLFile(path string) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	requests := strings.Split(string(file), ";")
	for _, request := range requests {
		_, err := DB.DB.Exec(request)
		if err != nil {
			panic(err)
		}
	}
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
	// 部分IP段（尤其是镜像站的IP段）网页上查找不到
	// 修复
	iprange := IPRange.NewIPV4WithIPRange("202.141.160.0", "202.141.160.255")
	area := Area.GetOrCreate("安徽")
	isp := ISP.GetOrCreate("CHINANET")
	iprange.SetAreaISP(area, isp)

	iprange = IPRange.NewIPV4WithIPRange("40.73.103.0", "40.73.103.255")
	area = Area.GetOrCreate("上海")
	isp = ISP.GetOrCreate("UNICOM")
	iprange.SetAreaISP(area, isp)

	iprange = IPRange.NewIPV4WithIPRange("59.111.0.0", "59.111.0.255")
	area = Area.GetOrCreate("浙江")
	isp = ISP.GetOrCreate("CMNET")
	iprange.SetAreaISP(area, isp)
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
	areas, err := Area.All()
	if err != nil {
		panic("Unable to fetch all areas!")
	}
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
	var wg sync.WaitGroup
	for _, area := range areas {
		wg.Add(1)
		go func(area_ Area.Area) {
			defer wg.Done()
			name, _ := area_.GetName()
			response, _ := http.Get(getAreaCoordinateURL(name))
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
		}(area)
	}
	// 填百度接口的坑
	areaCoordinate[Area.GetOrCreate("海南")] = coordinate{Lng: 110.35, Lat: 20.02}
	areaCoordinate[Area.GetOrCreate("香港")] = coordinate{Lng: 114.1, Lat: 22.2}
	areaCoordinate[Area.GetOrCreate("河南")] = coordinate{Lng: 113.65, Lat: 34.76}
	wg.Wait()
	return areaCoordinate
}

func initMirrorStationIPRange() {
	stations, err := MirrorStation.All()
	if err != nil {
		panic("Cannot get all stations!")
	}
	for _, station := range stations {
		url, err := station.GetURL()
		if err != nil {
			panic("Cannot get url!")
		}
		ips, err := net.LookupIP(url[strings.Index(url, "//")+2:])
		for _, ip := range ips {
			station.AddIPRange(ip.String())
		}
	}
}

func InitAll() {
	execSQLFile("./tables.sql")
	execSQLFile("./data.sql")
	initAreas()
	initAreaDistance()
	initMirrorStationIPRange()
	println("DB init success!")
}
