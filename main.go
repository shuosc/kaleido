package main

import (
	"net/http"
	"log"
	"database/sql"
	"strings"
	"strconv"
	_ "github.com/lib/pq"
)

var connStr = "user=mirrorRedirector password=test123456 dbname=postgres sslmode=disable"
var db *sql.DB

func IPv4ToNumberForm(ipv4 string) uint64 {
	splitResult := strings.Split(ipv4, ".")
	var numberForm uint64 = 0
	for _, s := range splitResult {
		num, _ := strconv.ParseUint(s, 10, 64)
		numberForm <<= 8
		numberForm |= num
	}
	return numberForm
}

func GetMirrorURLByIPv4(ipv4 string) string {
	numberForm := IPv4ToNumberForm(ipv4)
	row := db.QueryRow(`
	select mirror.url from 
		area,
		ipv4_area,
		ipv4_ranges,
		area_mirror,
		mirror
	where ipv4_ranges.data @> $1 :: numeric
	and ipv4_area.id = ipv4_ranges.id
	and area.id = ipv4_area.area_id
	and area_id = area_mirror.id
	and mirror.id = area_mirror.id;`, numberForm)
	var result string
	err := row.Scan(&result)
	if err != nil {
		return "https://mirrors.shu.edu.cn"
	}
	return result
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	from := r.RemoteAddr
	to := GetMirrorURLByIPv4(strings.Split(from, ":")[0])
	log.Println("Redirect request from", from, "to", to)
	http.Redirect(w, r, to+r.URL.Path, http.StatusMovedPermanently)
}

func main() {
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/", Redirect)

	log.Println("Listening...")
	http.ListenAndServe(":8386", mux)
}
