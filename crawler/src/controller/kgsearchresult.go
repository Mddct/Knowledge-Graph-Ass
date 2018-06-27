package controller

import (
	"html/template"
	"knowledgegraph"
	"net/http"
)

type KgSearchResult struct {
	Kgs string
}

//func GetResult(name string) []*model.Profile

func (s *KgSearchResult) ServeHTTP(
	w http.ResponseWriter, r *http.Request) {

	// qs := r.FormValue("q")
	res, _ := knowledgegraph.GetResult(s.Kgs, 0)
	if res != nil {
		// TODO 移到vie中

		t, err := template.ParseFiles("view/template/3.html")
		// TODO 统一的错误处理
		if err != nil {
			http.Error(w, "internal error",
				http.StatusInternalServerError)
			return
		}
		if t.Execute(w, res) != nil {
			http.Error(w, "internal error",
				http.StatusInternalServerError)
			return
		}
	} else {
		http.NotFound(w, r)
	}

}
