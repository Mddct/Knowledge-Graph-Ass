package knowledgegraph

import (
	"fmt"
	"model"

	"github.com/knakk/sparql"
)

const moviePattern = `
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
`

// TODO support other language
const boolPatern = `
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
}`

func GetResult(name string) []*model.Profile {
	const DbpediaPath = "https://dbpedia.org/sparql"
	qstring := fmt.Sprintf(moviePattern, name)

	repo, _ := sparql.NewRepo(DbpediaPath)
	res, err := repo.Query(qstring)
	if err != nil {
		panic(err)
	}

	ret := make([]*model.Profile, 0, 10)
	for _, k := range res.Solutions() {
		mi := &model.Profile{
			Link:     k["link"].String(),
			Name:     k["name"].String(),
			Abstract: k["abstract"].String(),
			Year:     k["year"].String(),
		}
		if _, ok := k["thumbnail"]; ok {
			mi.ImageSrc = k["thumbnail"].String()
		}
		ret = append(ret, mi)
	}

	return ret
}
