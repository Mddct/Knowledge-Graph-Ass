package model

type SearchResult struct {
	SearchString   string
	Hits           int
	Start          int
	Items          []*Profile
	RecommendItems []*Profile
}
