package parser

import (
	"model"
	"regexp"
	"strings"
	"types"

	"github.com/PuerkitoBio/goquery"
)

const movieBasePath = "http://www.1905.com"

func ParseProfile(
	contents []byte, name string,
	link string) types.ParseResult {
	var profile model.Profile
	profile.Name = name
	profile.Link = movieBasePath + link

	// /mdb/film/2245683/
	l := strings.Split(link, "/")
	profile.ID = l[len(l)-2]
	// profile.Director = ""
	doc, _ := goquery.NewDocumentFromReader(
		strings.NewReader(string(contents)))
	imageSrc, ok := doc.Find(
		".container-left").Find("img.poster").Attr("src")
	if ok {
		profile.ImageSrc = imageSrc
	}
	profile.Abstract = doc.Find(".container-left").
		Find(".plot").Find("p").Text()
	profile.Year = doc.Find(".container-right").
		Find("h1").Find("span").Text()
	var ret types.ParseResult
	ret.Items = append(ret.Items, profile)
	return ret
}

func extractString(contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)
	if len(match) >= 2 {
		return string(match[1])
	}
	return ""
}
