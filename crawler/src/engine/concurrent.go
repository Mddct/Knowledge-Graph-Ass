package engine

import (
	"fetcher"
	"log"
	"scheduler"
	"types"
)

type ConcurrentEngine struct {
	Scheduler   scheduler.Scheduler
	WorkerCount int
	ItemChan    chan interface{}
}

func (e *ConcurrentEngine) Run(seeds ...types.Request) {
	out := make(chan types.ParseResult)
	e.Scheduler.Run()

	for i := 0; i < e.WorkerCount; i++ {
		createWorker(e.Scheduler.WorkerChan(), out, e.Scheduler)
	}

	for _, r := range seeds {
		e.Scheduler.Submit(r)
	}

	for {
		result := <-out
		for _, item := range result.Items {
			go func() {
				e.ItemChan <- item
			}()
		}

		for _, request := range result.Requests {
			e.Scheduler.Submit(request)
		}
	}
}

func createWorker(in chan types.Request,
	out chan types.ParseResult,
	ready scheduler.ReadyNotifier) {
	go func() {
		for {
			// tell scheduler i am ready
			ready.WorkerReady(in)
			request := <-in
			result, err := worker(request)
			if err != nil {
				continue
			}
			out <- result
		}
	}()
}
func worker(r types.Request) (types.ParseResult, error) {
	// log.Printf("Fetching %s", r.Url)
	body, err := fetcher.Fetch(r.Url)
	if err != nil {
		log.Printf("Fetcher: error fetching url %s: %v", r.Url, err)
		return types.ParseResult{}, err
	}
	return r.ParseFunc(body), nil
}
