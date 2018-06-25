package main

import (
	"controller"
	"io"
	"log"
	"net/http"
	"os"
)

var homePage, _ = os.Open("view/template/1.html")

func HomePage(w http.ResponseWriter, r *http.Request) {
	io.Copy(w, homePage)
}
func main() {
	http.Handle("/search",
		controller.
			CreateSearchResultHandle("view/template/2.html"),
	)
	http.HandleFunc("/", HomePage)
	log.Println(http.ListenAndServe(":8080", nil))
}
