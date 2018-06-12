package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"unicode/utf8"
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
	Items      []*Item
}

type Item struct {
	Link        string
	Description string
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
		if c > 50 {
			break
		}
	}
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
	if len(r.FormValue("q")) < 32 {
		results := template.Must(template.New("resultslist").Funcs(template.FuncMap{
			"lessWorlds": lessWorlds,
		}).ParseFiles("./template/2.html"))

		res := Results{
			Items: []*Item{
				&Item{"www.baidu.com", "nihao上京东,百万商品在线,满减+闪电送!购爽,购High!nihao到京东,各款商品给力大促,High 翻天!"},
				&Item{"www.baidu.com", "nihao上京东,百万商品在线,满减+闪电送!购爽,购High!nihao到京东,各款商品给力大促,High 翻天!"},
				&Item{"www.baidu.com", "nihao上京东,百万商品在线,满减+闪电送!购爽,购High!nihao到京东,各款商品给力大促,High 翻天!"},
				&Item{"www.baidu.com", "nihao上京东,百万商品在线,满减+闪电送!购爽,购High!nihao到京东,各款商品给力大促,High 翻天!"},
				&Item{"www.baidu.com", "nihao上京东,百万商品在线,满减+闪电送!购爽,购High!nihao到京东,各款商品给力大促,High 翻天!"},
				&Item{"www.baidu.com", "nihao上京东,百万商品在线,满减+闪电送!购爽,购High!nihao到京东,各款商品给力大促,High 翻天!"},
				&Item{"www.baidu.com", "nihao上京东,百万商品在线,满减+闪电送!购爽,购High!nihao到京东,各款商品给力大促,High 翻天!"},
				&Item{"www.baidu.com", "nihao上京东,百万商品在线,满减+闪电送!购爽,购High!nihao到京东,各款商品给力大促,High 翻天!"},
			},
			TotalCount: 4,
		}
		// if err := results.Execute(os.Stdout, res); err != nil {
		// 	log.Fatal(err)
		// }
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
	fileserverHandler := http.FileServer(http.Dir("static/"))
	http.Handle("/asset/", http.StripPrefix("/asset/", fileserverHandler))
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/search", Search)
	http.ListenAndServe(":8080", nil)
}
