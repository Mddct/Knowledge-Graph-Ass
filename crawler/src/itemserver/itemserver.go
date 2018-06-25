package itemserver

import (
	"context"

	"gopkg.in/olivere/elastic.v5"
)

var client *elastic.Client

type Item interface {
	GetID() string
}

func init() {
	cli, err := elastic.NewClient(
		elastic.SetSniff(false),
	)
	if err != nil {
		panic(err)
	}
	client = cli

}
func ItemServer() chan interface{} {

	type inder interface {
		GetID() string
	}
	out := make(chan interface{})
	go func() {
		for {
			item := <-out
			if id, ok := item.(inder); ok {
				save(item, id.GetID())
			} else {
				save(item, "")
			}

		}
	}()
	return out
}
func save(m interface{}, id string) error {
	_, err := client.Index().Index("movie").
		Type("movie1905").
		Id(id).
		BodyJson(m).
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
