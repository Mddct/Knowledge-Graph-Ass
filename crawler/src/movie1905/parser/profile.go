package parser

import (
	"movie1905/model"
	"regexp"
	"strings"
	"types"

	"github.com/PuerkitoBio/goquery"
)

var abstractRe = regexp.MustCompile(``)

func ParseProfile(
	contents []byte, name string,
	link string) types.ParseResult {
	var profile model.Profile
	profile.Name = name
	profile.Link = link

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
