package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

var homePage []byte

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
}

// q represent search
func Search(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	io.WriteString(w, r.FormValue("q"))
}
func HomePage(w http.ResponseWriter, r *http.Request) {
	//	w.
	// fmt.Fprint(w, "hello")
	w.Write(homePage)
}
func main() {
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/search", Search)
	http.ListenAndServe(":8080", nil)
}
