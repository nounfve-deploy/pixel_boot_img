package lib

import (
	"fmt"
	"os"
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

func (d *Device) FindImage(build string) (i *Image) {
	for _, img := range d.Images {
		if img.Build == build {
			i = img
			break
		}
	}
	return
}

func (d *Device) GetBuildStorePath(build string) (path string) {
	root := os.Getenv("gitRoot")
	if root == "" {
		root = "."
	}
	path = fmt.Sprintf("%s/storage/%s/%s/", root, d.CodeName, build)
	err := os.MkdirAll(path, 0755)
	if err != nil {
		panic(err)
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

type ImageZip struct {
	Image
	path string
}

func (i Image) NewImageZip(path_hash []string) (iz *ImageZip) {
	if len(path_hash) != 2 {
		return
	}
	path, hash := path_hash[0], path_hash[1]
	if hash != i.Sha256 {
		return
	}
	iz = &ImageZip{i, path}
	return
}

func (iz *ImageZip) Unzip_boot_img() (bootImgs []string) {
	inner_zip_pattern := fmt.Sprintf("*%s*.zip", strings.ToLower(iz.Build))
	inner_zip := Sh_unzip_ls.RunReturnLastLineSplit(iz.path, inner_zip_pattern)

	if len(inner_zip) != 1 {
		fmt.Printf("%#v\n", inner_zip)
		panic("multiple image found")
	}

	boot_img_pattern := fmt.Sprintf("%s*boot.img", "")
	bootImgs = Sh_unzip_ls.RunReturnLastLineSplit(inner_zip[0], boot_img_pattern)
	return
}

var (
	TITLE_REGEX   = regexp.MustCompile(`^"([^"]+)" for (.+)$`)
	BUILD_REGEX   = regexp.MustCompile(`^[^(]+\((.+)\)$`)
	ZIP_URL_REGEX = regexp.MustCompile(`^.+\/([^-]+)-([^-]+)-.+\.zip$`)

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

	PARSE_DOWNLOAD_URL = func(url string) (deviceCode, build string) {
		subStr := ZIP_URL_REGEX.FindStringSubmatch(url)
		if len(subStr) != 3 {
			return
		}
		deviceCode, build = subStr[1], subStr[2]
		return
	}
)
