package lib

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Device struct {
	CodeName   string
	DeviceName string
	Images     []*Image
}

type Image struct {
	Build  string
	Date   string
	Url    string
	Sha256 string
}

var (
	TITLE_REGEX = regexp.MustCompile(`^"([^"]+)" for (.+)$`)
	BUILD_REGEX = regexp.MustCompile(`^[^(]+\((.+)\)$`)

	PARSE_BUILD_DATE = func(build_string string) (build, date string) {
		subStr := BUILD_REGEX.FindStringSubmatch(build_string)
		if len(subStr) != 2 {
			return
		}
		subStr = strings.Split(subStr[1], ",")
		if len(subStr) < 2 {
			return
		}
		build, date = strings.Trim(subStr[0], " "), strings.Trim(subStr[1], " ")
		return
	}
)

func TryNewDevice(title string) (d *Device) {
	match := TITLE_REGEX.FindStringSubmatch(title)
	if len(match) != 3 {
		return
	}

	d = &Device{
		CodeName:   match[1],
		DeviceName: match[2],
		Images:     make([]*Image, 0),
	}
	return
}

func TryNewImageFromTr(tr *goquery.Selection) (i *Image) {
	tds := make([]*goquery.Selection, 0)
	tr.Find("td").Each(func(i int, td *goquery.Selection) {
		if !td.Is("td") {
			return
		}
		tds = append(tds, td)
	})
	if len(tds) != 4 {
		return
	}

	Build, Date := PARSE_BUILD_DATE(tds[0].Text())
	if Build == "" || Date == "" {
		return
	}

	Url := tds[2].Find("a").First().AttrOr("href", "")
	Sha256 := tds[3].Text()
	if Url == "" || Sha256 == "" {
		return
	}

	i = &Image{Build, Date, Url, Sha256}
	return
}
