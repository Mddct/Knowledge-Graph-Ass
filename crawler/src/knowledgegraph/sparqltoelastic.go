package knowledgegraph

import (
	"fetcher"
	"fmt"
	"model"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/knakk/sparql"
)

type SearchType int

const (
	SearchMovie SearchType = iota
	SearchMovieRecommend
	SearchBook
	SearchBookRecommend
	SearchGame
	SearchGameRecommend
)
const moviePatern = `
PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
PREFIX dbo: <http://dbpedia.org/ontology/>
PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
PREFIX foaf: <http://xmlns.com/foaf/0.1/>
PREFIX dc: <http://purl.org/dc/elements/1.1/>

select distinct ?name ?image ?film ?director ?wikiLink ?language ?country group_concat(' ',?actor) AS ?actor ?comment where{
?film 
	  rdf:type 	 dbo:Film;
      rdfs:comment ?comment;
      rdfs:label ?name;
      dbp:language ?language;
      dbp:country ?country;
      dbo:director ?director;
      dbo:starring ?actor;
dbp:caption ?image;
      dbo:wikiPageExternalLink ?wikiLink
      
FILTER( (lang(?name)="zh") && regex(str(?name),"%s") && (lang(?comment)="zh") ).

}GROUP BY ?name ?film ?image ?comment ?language ?director ?wikiLink ?country
LIMIT 5`

const movieRecommendPatern = `

prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#>
prefix dbo: <http://dbpedia.org/ontology/>
prefix dbp: <http://dbpedia.org/property/>
prefix foaf: <http://xmlns.com/foaf/0.1/>

SELECT ?name2 ?wikiLink ?comment ?image COUNT(?movie2) SAMPLE(?movie2)
        FROM <http://en.dbpedia.org>
        WHERE
        {
		  ?movie1 dbp:title ?name.
		  ?movie1 dct:subject ?o.
		  ?movie2 dct:subject ?o.
		  ?movie2 dbo:wikiPageExternalLink ?wikiLink.
				  ?movie2 rdfs:comment ?comment.
				  ?movie2 dbp:caption ?image.
				  ?movie2 dbp:title ?name2.
			
		  FILTER (?movie2 != ?movie1) 
		  FILTER( (lang(?name2)="zh") && regex(str(?name),"%s")).
        } 
GROUP BY ?movie2 ?wikiLink ?name ?comment ?image ?name2
ORDER BY DESC(COUNT(?movie2))
limit 5
`

const bookPatern = ` PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
PREFIX dbo: <http://dbpedia.org/ontology/>
PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
PREFIX dc: <http://purl.org/dc/elements/1.1/>

SELECT ?name ?book ?author ?genre ?abs
WHERE 
{{
	?book rdf:type dbo:Book .
	?book rdfs:label ?name .
	?book dbo:author ?author .
	?book dbo:literaryGenre ?genre .
	?book dbo:abstract ?abs
	FILTER (  (lang(?name)="zh")&&(lang(?abs)="zh")&&(regex(?name, "%s"))  ).
}}LIMIT 3`

const bookRecommendPattern = `
PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
PREFIX dbo: <http://dbpedia.org/ontology/>
PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
PREFIX dc: <http://purl.org/dc/elements/1.1/>

SELECT ?name ?book
WHERE 
{{
	?book rdf:type dbo:Book .
	?book rdfs:label ?name .
	?book dbo:author dbr:{} .
	FILTER ((lang(?name)="zh")&&(?name != '{}'@zh)) .
	MINUS 
	{{
	?game rdf:type dbo:Book .
	?game rdfs:label ?name .
	FILTER (  (lang(?name)="zh")&&(lang(?abs)="zh")&&(regex(?name, "%s"))  ).
	}}
}}LIMIT 5
`
const DbpediaPath = "http://dbpedia.org/sparql"

func GetResult(name string, category SearchType) ([]*model.Profile, error) {
	if category == SearchMovie {
		return searchMovie(name, moviePatern)
	} else if category == SearchMovieRecommend {
		return searchMovie(name, movieRecommendPatern)
	}
	return nil, nil
}
func searchMovie(name string, pattern string) ([]*model.Profile, error) {

	qstring := fmt.Sprintf(moviePatern, name)
	repo, err := sparql.NewRepo(DbpediaPath)
	if err != nil {
		return nil, err
	}

	res, err := repo.Query(qstring)
	if err != nil {
		return nil, err
	}

	ret := make([]*model.Profile, 0, 10)

	for _, k := range res.Solutions() {
		mi := &model.Profile{
			Link:     k["wikiLink"].String(),
			Name:     k["name"].String(),
			Abstract: k["comment"].String(),
		}
		ret = append(ret, mi)
	}

	var wg sync.WaitGroup
	for i := range ret {
		wg.Add(1)
		go func(i int, imgSrc *string) {
			defer wg.Done()
			iui := fmt.Sprintf(`https://www.douban.com/search?q=%s`,
				ret[i].Name)
			contents, err := fetcher.Fetch(iui)
			if err != nil {
				return
			}
			var dom, _ = goquery.NewDocumentFromReader(
				strings.NewReader(string(contents)))

			ok := false
			*imgSrc, ok = dom.Find(".pic").Find("img").Attr("src")
			if !ok {
				return
			}
		}(i, &(ret[i].ImageSrc))

	}

	wg.Wait()

	return ret, nil
}
