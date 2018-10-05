package initDB

import (
	_ "database/sql"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"kaleido/master/infrastructure"
	"math"
	"net/http"
	"strings"
)

func parsePage(url string) {
	response, _ := http.Get(url)
	doc, _ := goquery.NewDocumentFromReader(response.Body)
	lastIndex := strings.LastIndex(url, "/i_")
	ispName := url[lastIndex+3:]
	cityName := ""
	var cityId int
	doc.Find(".list").Each(func(index int, selection *goquery.Selection) {
		selection.Find("dd,dt").Each(func(_ int, selectedLine *goquery.Selection) {
			selected := strings.Trim(selectedLine.Text(), " \n")
			if selected == "" {

			} else if strings.Contains(selected, ".") {
				from := selectedLine.Find(".v_l").Text()
				to := selectedLine.Find(".v_r").Text()
				var ipRangeId int
				row := infrastructure.DB.QueryRow(`
				insert into iprange (ip) values (inet_merge($1,$2)) returning id;
				`, from, to)
				row.Scan(&ipRangeId)
				infrastructure.DB.Exec(`
				insert into iprange_area(iprange_id, area_id) values ($1,$2);
				`, ipRangeId, cityId)
			} else {
				cityName = selected[:strings.Index(selected, "(")]
				row := infrastructure.DB.QueryRow(`
				insert into area (name,isp) values ($1,$2) returning id;
				`, cityName, ispName)
				println(ispName)
				row.Scan(&cityId)
			}
		})
	})
}

func InitIPArea() {
	response, _ := http.Get("http://ipcn.chacuo.net")
	doc, _ := goquery.NewDocumentFromReader(response.Body)
	selection := doc.Find(".section_content > .list a")
	ch := make(chan int, selection.Length())
	selection.Map(func(index int, anchor *goquery.Selection) string {
		go func(channel chan int, idx int) {
			href, _ := anchor.Attr("href")
			parsePage(href)
			ch <- idx
		}(ch, index)
		return ""
	})
	for i := 0; i < selection.Length(); i++ {
		<-ch
	}
	println("init db ok")
}

func rad(degree float64) float64 {
	return degree * math.Pi / 180.0
}

func distance(lat1 float64, lng1 float64, lat2 float64, lng2 float64) float64 {
	radLat1 := rad(lat1)
	radLat2 := rad(lat2)
	a := math.Pow(math.Sin((radLat1-radLat2)/2), 2)
	b := math.Cos(radLat1) * math.Cos(radLat2) * math.Pow(math.Sin((rad(lng1)-rad(lng2))/2), 2)
	return 6378 * 2 * math.Asin(math.Sqrt(a+b))
}

func InitAreaDistance() {
	placeCoordinate := map[string]struct {
		Lng float64 `json:"lng"`
		Lat float64 `json:"lat"`
	}{}
	rows, _ := infrastructure.DB.Query(`
	select distinct name from area;
	`)
	for rows.Next() {
		var areaName string
		rows.Scan(&areaName)
		response, _ := http.Get("http://api.map.baidu.com/geocoder?address=" + areaName +
			"&output=json&key=37492c0ee6f924cb5e934fa08c6b1676&city=北京市")
		body, _ := ioutil.ReadAll(response.Body)
		var result struct {
			Result struct {
				Location struct {
					Lng float64 `json:"lng"`
					Lat float64 `json:"lat"`
				} `json:"location"`
			} `json:"result"`
		}
		json.Unmarshal(body, &result)
		placeCoordinate[areaName] = result.Result.Location
	}
	// 填百度接口的坑
	placeCoordinate["海南"] = struct {
		Lng float64 `json:"lng"`
		Lat float64 `json:"lat"`
	}{Lng: 110.35, Lat: 20.02}
	placeCoordinate["香港"] = struct {
		Lng float64 `json:"lng"`
		Lat float64 `json:"lat"`
	}{Lng: 114.1, Lat: 22.2}
	placeCoordinate["河南"] = struct {
		Lng float64 `json:"lng"`
		Lat float64 `json:"lat"`
	}{Lng: 113.65, Lat: 34.76}
	rows, _ = infrastructure.DB.Query(`
	select *
	from area, area as area_;
	`)
	for rows.Next() {
		var fromId, toId uint64
		var fromName, toName string
		var fromISP, toISP string
		rows.Scan(&fromId, &fromName, &fromISP, &toId, &toName, &toISP)
		fromLocation := placeCoordinate[fromName]
		toLocation := placeCoordinate[toName]
		distance := distance(fromLocation.Lng, fromLocation.Lat, toLocation.Lng, toLocation.Lat)
		if fromISP != toISP {
			distance += 10
		}
		infrastructure.DB.Exec(`
		insert into area_area (from_id, to_id, distance) VALUES ($1,$2,$3);
		`, fromId, toId, uint64(math.Ceil(distance)))
		fmt.Println(fromName, fromISP, "->", toName, toISP, ":", distance)
	}
	println("init distance ok")
}

func InitAll() {
	InitIPArea()
	InitAreaDistance()
}
