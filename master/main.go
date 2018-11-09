package main

import (
	"kaleido/master/service"
	_ "net/http/pprof"
)

func main() {
	//initDB.InitAll()
	service.StartCronJobs()
}
