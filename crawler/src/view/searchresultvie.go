package view

import (
	"html/template"
	"io"
	"model"
	"strings"
	"unicode/utf8"
)

type SearchResultView struct {
	template *template.Template
}

func lessWorlds(s string) string {
	if s == "" {
		return ""
	}
	var builder strings.Builder
	c := 0

	for i := 0; i < len(s); {
		// ommit the encoding error
		r, d := utf8.DecodeRuneInString(s[i:])
		i += d
		c++
		builder.WriteRune(r)
		if c > 100 {
			break
		}
	}
	return builder.String() + "..."
}
func realRange(data interface{}) []*model.Profile {
	d, ok := data.(*model.SearchResult)
	if !ok {
		return nil
	}

	return d.Items
}
func realRecommendRange(data interface{}) []*model.Profile {
	d, ok := data.([]*model.Profile)
	if !ok {
		return nil
	}

	return d
}
func CreateSearchResultView(filename string) *SearchResultView {
	return &SearchResultView{
		template: template.Must(template.New("tmp").Funcs(template.FuncMap{
			"lessWorlds":         lessWorlds,
			"realRange":          realRange,
			"realRecommendRange": realRecommendRange,
		}).ParseFiles(filename))}
}

func (s *SearchResultView) Render(
	w io.Writer, data *model.SearchResult) error {
	return s.template.Execute(w, map[string]interface{}{
		"data":      data,
		"recommend": data.Items[1 : len(data.Items)-1],
	})

}
