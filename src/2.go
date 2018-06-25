package main

import (
	"fmt"

	"github.com/knakk/sparql"
)

const DbpediaPath = "https://dbpedia.org/sparql"
const WikiDataPath = `https://query.wikidata.org/sparql`
const CNDbepdiaPath = "http://shuyantech.com/api/cndbpedia/ment2ent"
const qstring = `
PREFIX xsd: <http://www.w3.org/2001/XMLSchema#>

#Movies released in 2017
SELECT DISTINCT ?item ?itemLabel WHERE {
  ?item wdt:P31 wd:Q11424.
  ?item wdt:P577 ?pubdate.
  ?item wdt:P495 wd:Q148.
  SERVICE wikibase:label { bd:serviceParam wikibase:language "[AUTO_LANGUAGE],zh-cn". }
  FILTER((?pubdate >= "2018-01-01T00:00:00Z"^^xsd:dateTime) && (?pubdate <= "2018-6-31T00:00:00Z"^^xsd:dateTime))
}
`

func main() {
	// url := url.Values{}
	// url.Set("format", "json")
	// url.Set("query", qstring)

	// res, err := http.Get(WikiDataPath + "?" + url.Encode())
	// if err != nil {
	//	panic(err)
	// }
	// defer res.Body.Close()
	// io.Copy(os.Stdout, res.Body)
	repo, _ := sparql.NewRepo(WikiDataPath)
	res, err := repo.Query(qstring)
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Solutions())
}
