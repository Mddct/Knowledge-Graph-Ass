package engine

import (
	"parser"
	"scheduler"
	"testing"
	"types"
)

func TestConcurrentEngine_simpleScheduler(t *testing.T) {
	e := ConcurrentEngine{
		Scheduler:   &scheduler.SimpleScheduler{},
		WorkerCount: 100,
	}
	e.Run(types.Request{
		Url:       "http://www.zhenai.com/zhenghun",
		ParseFunc: parser.ParseCityList,
	})
}

func TestConcurrentEngine_queuedScheduler(t *testing.T) {
	e := ConcurrentEngine{
		Scheduler:   &scheduler.QueuedScheduler{},
		WorkerCount: 100,
	}
	e.Run(types.Request{
		Url:       "http://www.zhenai.com/zhenghun",
		ParseFunc: parser.ParseCityList,
	})
}
