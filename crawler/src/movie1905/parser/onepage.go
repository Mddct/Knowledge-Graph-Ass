package parser

import (
	"fmt"
	"regexp"
	"strings"
	"types"

	"github.com/PuerkitoBio/goquery"
)

const movieBasePath = "http://www.1905.com"

var movieRe = regexp.MustCompile(`<a target="_blank" title="([^"]+)" href="(/mdb/film/[0-9]+/)">[^<]+</a>`)
var year = 2018

func ParseMovieOnePage(contents []byte) types.ParseResult {
	var dom, _ = goquery.NewDocumentFromReader(
		strings.NewReader(string(contents)))
	matches := movieRe.FindAllSubmatch(contents, -1)

	result := types.ParseResult{}
	for _, m := range matches {
		name := string(m[1])
		link := string(m[2])

		result.Requests = append(
			result.Requests,
			types.Request{
				Url: movieBasePath + string(m[2]),
				ParseFunc: func(bytes []byte) types.ParseResult {
					return ParseProfile(
						bytes, name, movieBasePath+link)
				},
			})

	}

	// find nextpage
	np := dom.Find("#new_page").
		Find("a").Last()

	if np.Text() == "下一页" {
		nextPageLink, ok := np.Attr("href")
		if ok {
			result.Requests = append(result.Requests,
				types.Request{
					Url:       movieBasePath + nextPageLink,
					ParseFunc: ParseMovieOnePage,
				})
		}
	} else if year > 2000 { // 当前年完成
		year--
		result.Requests = append(result.Requests,
			types.Request{

				Url:       fmt.Sprintf("http://www.1905.com/mdb/film/list/year-%d/o0d0p1.html", year),
				ParseFunc: ParseMovieOnePage,
			})

	}
	return result
}
