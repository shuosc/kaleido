package Application

import (
	"github.com/robfig/cron"
	"io/ioutil"
	"kaleido/Common/Services/IPTools"
	"kaleido/Node/Services/IPMirrorstationTable"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
)

var RedirectTable map[uint32]map[uint32][]string
var mutex sync.RWMutex

func GetMirrorListForIP(testIP string) []string {
	//table := NodeTable.Unmarshal(MasterTable.Marshal())
	for maskBitLenth, maskedIp := range RedirectTable {
		masked := IPTools.MaskIP(testIP, uint8(maskBitLenth))
		result, found := maskedIp[masked]
		if found {
			return result
		}
	}
	return nil
}

var lastTimeFetchSuccess = true

func FetchRedirectTable() {
	response, err := http.Get("http://localhost:3000")
	if err != nil {
		if lastTimeFetchSuccess {
			log.Println("Failed to fetch RedirectTable!")
		}
		lastTimeFetchSuccess = false
		return
	} else {
		if !lastTimeFetchSuccess {
			log.Println("Fetch RedirectTable success again!")
		}
		lastTimeFetchSuccess = true
	}
	result, _ := ioutil.ReadAll(response.Body)
	mutex.Lock()
	defer mutex.Unlock()
	RedirectTable = IPMirrorstationtable.Unmarshal(result)
}

func StartCron() {
	c := cron.New()
	c.AddFunc("@every 5s", func() {
		FetchRedirectTable()
	})
	c.Start()
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	from := r.RemoteAddr
	if r.Header.Get("force-redirect-ip") != "" {
		from = r.Header.Get("force-redirect-ip")
	}
	toList := GetMirrorListForIP(strings.Split(from, ":")[0])
	var to string
	if len(toList) == 0 {
		to = "https://mirrors.shu.edu.cn"
	} else {
		to = toList[rand.Int()%len(toList)]
	}
	log.Println("Redirect request from", from, "to", to)
	http.Redirect(w, r, to+r.URL.Path, http.StatusMovedPermanently)
}

func StartServer() {
	mux := http.NewServeMux()
	tableHandler := http.HandlerFunc(ServeHTTP)
	mux.Handle("/", tableHandler)
	http.ListenAndServe(":8086", mux)
}

func init() {
	FetchRedirectTable()
}
