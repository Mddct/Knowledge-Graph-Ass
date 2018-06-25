package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/knakk/sparql"
)

var homePage []byte
var limitedPage []byte

func init() {
	h, err := os.Open("template/1.html")
	if err != nil {
		panic(err)
	}
	hp, err := ioutil.ReadAll(h)
	if err != nil {
		panic(err)
	}
	homePage = hp
	l, err := os.Open("template/2.html")
	if err != nil {
		panic(err)
	}
	lp, err := ioutil.ReadAll(l)
	if err != nil {
		panic(err)
	}
	limitedPage = lp
}

type Results struct {
	TotalCount int
	Items      []*MovieInfo
}

func lessWorlds(s string) string {
	var builder strings.Builder
	c := 0
	for i := 0; i < len(s); {
		// ommit the encoding error
		r, d := utf8.DecodeRuneInString(s[i:])
		i += d
		c++
		builder.WriteRune(r)
		if c > 80 {
			break
		}
	}
	builder.WriteString("...")
	return builder.String()
}

// q represent search
// newwindow=1
// maxworldlen = 32
// 搜索不到 提供关键字
// 超出限制
func Search(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// io.WriteString(w, r.FormValue("q"))
	text := r.FormValue("q")
	if len(text) < 32 {
		results := template.Must(template.New("resultslist").Funcs(template.FuncMap{
			"lessWorlds": lessWorlds,
		}).ParseFiles("./template/2.html"))

		ret := getResult(text)
		res := Results{
			Items:      ret,
			TotalCount: len(ret),
		}
		results.ExecuteTemplate(w, `temp`, res)
		return
	}
	fmt.Println(r.FormValue("q"))
}
func HomePage(w http.ResponseWriter, r *http.Request) {
	//	w.
	// fmt.Fprint(w, "hello")
	w.Write(homePage)
}
func main() {

	// fileserverHandler := http.FileServer(http.Dir("static/"))
	// http.Handle("/asset/", http.StripPrefix("/asset/", fileserverHandler))
	// http.HandleFunc("/", HomePage)
	// http.HandleFunc("/search", Search)
	// http.ListenAndServe(":8080", nil)
	getBooks("鲁迅")
}

type MovieInfo struct {
	Link        string
	Name        string
	Countrys    string
	Description string
	Year        string
	Thumbnail   string
}

type BookInfo struct {
	FirstLine string
	// topic's main category
	TopicMC string
	Title   string
	Author  string

	// country of origin

	CoO             string
	PublicationDate string
}

func getBooks(authorName string) []*BookInfo {
	const WikiDataPath = "https://query.wikidata.org/sparql"
	qstring := fmt.Sprintf(`
#Books by a given Author including genres, series, and publication year
#added before 2016-10
 SELECT ?book ?bookLabel ?authorLabel ?genre_label ?series_label ?publicationDate ?firstline
WHERE
{
	?author ?label "%s"@zh .
	?book wdt:P31 wd:Q571 .
	?book wdt:P50 ?author .
        OPTIONAL{?book wdt:P1922 ?firstline}
	OPTIONAL {
		?book wdt:P136 ?genre .
		?genre rdfs:label ?genre_label filter (lang(?genre_label) = "zh-cn").
	}
	OPTIONAL {
		?book wdt:P179 ?series .
		?series rdfs:label ?series_label filter (lang(?series_label) = "zh-cn").
	}
	OPTIONAL {
		?book wdt:P577 ?publicationDate .
	}
	SERVICE wikibase:label {
		bd:serviceParam wikibase:language "zh" .
	}
}

`, authorName)
	repo, _ := sparql.NewRepo(WikiDataPath)
	res, err := repo.Query(qstring)
	if err != nil {
		panic(err)
	}

	ret := make([]*BookInfo, 0, 10)
	for _, k := range res.Solutions() {
		mi := &BookInfo{}
		v := reflect.ValueOf(mi).Elem()
		if _, ok := k["title"]; ok {
			v.FieldByName("Title").Set(reflect.ValueOf(k["title"].String()))
		}
		if _, ok := k["authorLabel"]; ok {
			v.FieldByName("Author").Set(reflect.ValueOf(k["authorLabel"].String()))
		}
		if _, ok := k["firstline"]; ok {
			v.FieldByName("FirstLine").Set(reflect.ValueOf(k["firstline"].String()))
		}

		if _, ok := k["publicationDate"]; ok {
			v.FieldByName("PublicationDate").Set(reflect.ValueOf(k["publicationDate"].String()))
		}
		ret = append(ret, mi)
	}

	fmt.Println(ret[0])
	return nil

}

func getResult(name string) []*MovieInfo {
	const DbpediaPath = "https://dbpedia.org/sparql"
	qstring := fmt.Sprintf(`
    PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
    PREFIX dbo: <http://dbpedia.org/ontology/>
    PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
    PREFIX dct: <http://purl.org/dc/terms/>

select ?link ?name ?countrys ?abstract ?year ?thumbnail
where{
?link rdf:type dbo:Film;
       rdfs:label ?name;
       dbp:country ?countrys;
      rdfs:comment ?abstract;
      dbo:runtime ?year
filter regex(?name,"%s")
filter(LANG(?abstract)="zh" ).
optional {
?link dbo:thumbnail ?thumbnail
}
}
   `, name)
	repo, _ := sparql.NewRepo(DbpediaPath)
	res, err := repo.Query(qstring)
	if err != nil {
		panic(err)
	}

	ret := make([]*MovieInfo, 0, 10)
	for _, k := range res.Solutions() {
		mi := &MovieInfo{
			Link:        k["link"].String(),
			Name:        k["name"].String(),
			Description: k["abstract"].String(),
			Year:        k["year"].String(),
			Countrys:    k["countrys"].String(),
		}
		if _, ok := k["thumbnail"]; ok {
			mi.Thumbnail = k["thumbnail"].String()
		}
		ret = append(ret, mi)
	}

	// // fmt.Println(res.Solutions())
	// url := url.Values{}
	// url.Set("format", "json")
	// url.Set("query", qstring)

	// res, err := http.Get(DbpediaPath + "?" + url.Encode())
	// if err != nil {
	//	panic(err)
	// }
	// defer res.Body.Close()
	// io.Copy(os.Stdout, res.Body)
	return ret
}
