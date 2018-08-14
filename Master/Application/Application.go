package Application

import (
	"github.com/robfig/cron"
	"kaleido/Master/Domain/Services/HeartBeat"
	"kaleido/Master/Domain/Services/IPMirrorstationTable"
	"net/http"
	"sync"
)

var data []byte
var mutex sync.RWMutex

func StartCron() {
	c := cron.New()
	c.AddFunc("@every 60s", func() {
		HeartBeat.Test()
		mutex.Lock()
		defer mutex.Unlock()
		data = IPMirrorstationTable.Marshal()
	})
	c.Start()
}

func ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	mutex.RLock()
	defer mutex.RUnlock()
	w.Write(data)
}

func StartServer() {
	mux := http.NewServeMux()
	tableHandler := http.HandlerFunc(ServeHTTP)
	mux.Handle("/", tableHandler)
	http.ListenAndServe(":3000", mux)
}

func init() {
	data = IPMirrorstationTable.Marshal()
}
