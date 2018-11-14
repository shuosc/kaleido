package main

import (
	"kaleido/master/service"
	"kaleido/master/service/initDB"
	_ "net/http/pprof"
)

func main() {
	initDB.InitAll()
	go service.StartCronJobs()
	go service.StartGraphQLServer()
	select {}
}
