package controller

import (
	"context"
	"knowledgegraph"
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

	// ommit error
	r.ParseForm()

	q := strings.TrimSpace(r.Form.Get("q"))
	from, err := strconv.Atoi(
		r.FormValue("from"),
	)
	if err != nil {
		from = 0
	}

	type baseChan struct {
		data *model.SearchResult
		err  error
	}
	type recommendChan struct {
		data []*model.Profile
		err  error
	}
	bc := make(chan baseChan)
	re := make(chan recommendChan)
	go func() {
		defer close(bc)
		page, err := s.getSearchResult(q, from)
		bc <- baseChan{page, err}
	}()
	go func() {
		defer close(re)
		rec, err := knowledgegraph.GetResult(q,
			knowledgegraph.SearchMovieRecommend)
		re <- recommendChan{rec, err}
	}()

	b := <-bc
	nr := <-re
	page, err := b.data, b.err

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	pageRe, err := nr.data, nr.err

	if err == nil {
		if len(pageRe) > 0 {
			if len(page.Items) > 0 {
				page.Items[0] = pageRe[0]
			} else {
				page.Items = pageRe
				pageRe = nil
			}

		}
		page.RecommendItems = pageRe
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
