package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const MovieBasePath = "http://www.1905.com/mdb/film/list/year-2018/"

func main() {
	// seedPage := "o0d0p1.html"
	numPages := 1
	for numPages < 51 {
		resp, err := http.Get(
			MovieBasePath + fmt.Sprintf("o0d0p%d.html", numPages))
		if err != nil {
			panic(err)
		}
		if resp.StatusCode != http.StatusOK {
			break
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		mathMovies(b)
		numPages++
		resp.Body.Close()

	}
	fmt.Println(numPages)
}
func mathMovies(b []byte) {

}
