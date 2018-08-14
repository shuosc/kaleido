package Init

import (
	"bytes"
	_ "database/sql"
	"github.com/PuerkitoBio/goquery"
	"kaleido/Master/Infrastructure/DB"
	_ "kaleido/Master/Infrastructure/DB"
	"net/http"
	"strconv"
	"strings"
)

func ipv4ToNumberForm(ipv4 string) uint32 {
	splitResult := strings.Split(ipv4, ".")
	var numberForm uint64 = 0
	for _, s := range splitResult {
		num, _ := strconv.ParseUint(s, 10, 32)
		numberForm <<= 8
		numberForm |= num
	}
	return uint32(numberForm)
}

func getMaskForIPRange(from string, to string) int {
	fromNumber := ipv4ToNumberForm(from)
	toNumbers := ipv4ToNumberForm(to)
	for i := 31; i != 0; i-- {
		fromBit := fromNumber & (1 << uint32(i))
		toBit := toNumbers & (1 << uint32(i))
		if fromBit != toBit {
			return 31 - i
		}
	}
	return -1
}

func ipRangeToCidrForm(from string, to string) string {
	var i int
	var v string
	var result bytes.Buffer
	splits := strings.Split(from, ".")
	for i, v = range splits {
		if v == "0" {
			break
		}
	}
	for j := 0; j < i; j++ {
		result.WriteString(splits[j])
		if j != i-1 {
			result.WriteString(".")
		}
	}
	result.WriteString("/")
	result.WriteString(strconv.Itoa(getMaskForIPRange(from, to)))
	return result.String()
}

func getIPRangesOnPage(url string) []string {
	var result []string
	response, _ := http.Get(url)
	pageDoc, _ := goquery.NewDocumentFromReader(response.Body)
	froms := pageDoc.Find(".list dd").Map(
		func(_ int, selection *goquery.Selection) string {
			return selection.Find(".v_l").Text()
		})
	tos := pageDoc.Find(".list dd").Map(
		func(_ int, selection *goquery.Selection) string {
			return selection.Find(".v_r").Text()
		})
	if len(froms) != len(tos) {
		panic("len(froms) should always equals to len(tos)")
	}
	for i := range froms {
		result = append(result, ipRangeToCidrForm(froms[i], tos[i]))
	}
	return result
}

func crawlOneCity(cityName string, cityPageUrl string, finish chan int) {
	var newCityId int
	DB.DB.QueryRow(`
	INSERT INTO area(province) VALUES ($1) RETURNING id;
	`, cityName).Scan(&newCityId)
	for _, ip := range getIPRangesOnPage(cityPageUrl) {
		var newIPId int
		DB.DB.QueryRow(`
		INSERT INTO ip(data) VALUES ($1) RETURNING id;
		`, ip).Scan(&newIPId)
		DB.DB.Exec(`
		INSERT INTO ip_area(ip_id, area_id) VALUES ($1,$2);
		`, newIPId, newCityId)
	}
	finish <- 0
}

func CrawlIP() {
	DB.DB.Exec(`TRUNCATE TABLE area CASCADE;`)
	DB.DB.Exec(`TRUNCATE TABLE ip CASCADE;`)
	response, _ := http.Get("http://ips.chacuo.net")
	indexPageDoc, _ := goquery.NewDocumentFromReader(response.Body)
	allCityNames := indexPageDoc.Find(".list > li > a").Map(
		func(_ int, selection *goquery.Selection) string {
			name := selection.Text()
			return name
		})
	allCityUrls := indexPageDoc.Find(".list > li > a").Map(
		func(_ int, selection *goquery.Selection) string {
			href, _ := selection.Attr("href")
			return href
		})
	if len(allCityNames) != len(allCityUrls) {
		panic("len(allCityNames) should always equals to len(allCityUrls)")
	}
	finish := make(chan int, len(allCityNames))
	for i := range allCityNames {
		go crawlOneCity(allCityNames[i], allCityUrls[i], finish)
	}
	for range allCityNames {
		<-finish
	}
}
