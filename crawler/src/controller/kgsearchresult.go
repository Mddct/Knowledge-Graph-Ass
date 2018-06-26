package controller

import (
	"fmt"
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

	fmt.Println("calling")
	// qs := r.FormValue("q")
	res := knowledgegraph.GetResult(s.Kgs)
	if res != nil {
		// TODO 移到vie中
		t, err := template.ParseFiles("view/template/3.html")
		// TODO 统一的错误处理
		if err != nil {
			http.Error(w, "internal error",
				http.StatusInternalServerError)
		}
		if t.Execute(w, res) != nil {
			http.Error(w, "internal error",
				http.StatusInternalServerError)
		}
	}
	http.NotFound(w, r)
}
