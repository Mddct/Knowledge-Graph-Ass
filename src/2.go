package main

import (
	"fmt"

	"github.com/knakk/sparql"
)

const DbpediaPath = "https://dbpedia.org/sparql"
const WikiDataPath = `https://query.wikidata.org/sparql`
const CNDbepdiaPath = "http://shuyantech.com/api/cndbpedia/ment2ent"
const qstring1 = `
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

const qstring = `
#Number of handed out academy awards per award type
#added before 2016-10

SELECT ?awardCount ?award ?awardLabel WHERE {
	{
		SELECT (COUNT(?award) AS ?awardCount) ?award
		WHERE
		{
			{
				SELECT (SAMPLE(?human) AS ?human) ?award ?awardWork (SAMPLE(?director) AS ?director) (SAMPLE(?awardEdition) AS ?awardEdition) (SAMPLE(?time) AS ?time) WHERE {
					?award wdt:P31 wd:Q19020 .			# All items that are instance of(P31) of Academy awards (Q19020)
					{
						?human p:P166 ?awardStat .              # Humans with an awarded(P166) statement
						?awardStat ps:P166 ?award .     	 # ... that has any of the values of ?award
						?awardStat pq:P805 ?awardEdition . # Get the award edition (which is "subject of" XXth Academy Awards)
						?awardStat pq:P1686 ?awardWork . # The work they have been awarded for
						?human wdt:P31 wd:Q5 . 				# Humans
					} UNION {
						?awardWork wdt:P31 wd:Q11424 . # Films
						?awardWork p:P166 ?awardStat . # ... with an awarded(P166) statement
						?awardStat ps:P166 ?award .     	 # ... that has any of the values of ?award
						?awardStat pq:P805 ?awardEdition . # Get the award edition (which is "subject of" XXth Academy Awards)
					}
					OPTIONAL {
						?awardEdition wdt:P585 ?time . # the "point of time" of the Academy Award
						?awardWork wdt:P57 ?director .
					}
				}
				GROUP BY ?awardWork ?award # We only want every movie once for a category (a 'random' person is selected)
			}
		} GROUP BY ?award
		ORDER BY ASC(?awardCount)
	}
	SERVICE wikibase:label {            # ... include the labels
		bd:serviceParam wikibase:language "[AUTO_LANGUAGE],zh" .
	}
}
`

func main() {

	repo, _ := sparql.NewRepo(WikiDataPath)
	res, err := repo.Query(qstring)
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Head)
	fmt.Println(res.Solutions())
	fmt.Println(res.Bindings)
}
