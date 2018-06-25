package controller

import (
	"context"
	"model"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"view"

	"gopkg.in/olivere/elastic.v5"
)

var el *elastic.Client

func init() {
	e, err := elastic.NewClient(
		elastic.SetSniff(false),
	)
	if err != nil {
		panic(err)
	}
	el = e
}

type SearchResultHandler struct {
	view   *view.SearchResultView
	client *elastic.Client
}

func CreateSearchResultHandle(
	template string,
) *SearchResultHandler {
	return &SearchResultHandler{
		view:   view.CreateSearchResultView(template),
		client: el,
	}

}

//xxxx/search?q= & from=20
// 每页最多十个
func (s *SearchResultHandler) ServeHTTP(
	w http.ResponseWriter, r *http.Request) {

	q := strings.TrimSpace(r.FormValue("q"))
	from, err := strconv.Atoi(
		r.FormValue("from"),
	)
	if err != nil {
		from = 0
	}

	page, err := s.getSearchResult(q, from)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = s.view.Render(w, page)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (s *SearchResultHandler) getSearchResult(
	q string, from int) (*model.SearchResult, error) {
	var res model.SearchResult
	res.SearchString = q
	resp, err := s.client.Search("movie").
		Query(elastic.NewSimpleQueryStringQuery(q)).
		From(from).Do(context.Background())
	if err != nil {
		return nil, err
	}
	res.Hits = int(resp.TotalHits())
	res.Start = from + 1
	for _, v := range resp.Each(
		reflect.TypeOf(&model.Profile{})) {
		item := v.(*model.Profile)
		res.Items = append(res.Items, item)
	}
	return &res, nil
}
