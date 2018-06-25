package main

import (
	"engine"
	"itemserver"
	"movie1905/parser"
	"scheduler"
	"types"
)

func main() {
	e := engine.ConcurrentEngine{
		Scheduler:   &scheduler.QueuedScheduler{},
		WorkerCount: 100,
		ItemChan:    itemserver.ItemServer(),
	}
	e.Run(types.Request{
		Url:       "http://www.1905.com/mdb/film/list/year-2018/o0d0p1.html",
		ParseFunc: parser.ParseMovieOnePage,
	})
}
