package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"pixel_boot_img/lib"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

func main() {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		slog.Error("Recovered from panic", "error", r)
		fmt.Println("usage: pixel_boot_img [target] [build]")
	}()
	target, build := os.Args[1], os.Args[2]
	slog.Info("cli input", "target", target, "build", build)

	devices := collect_device_images()
	d, ok := devices[target]
	if !ok {
		slog.Error("unmatched device")
		panic(-1)
	}
	img := d.FindImage(build)
	if img == nil {
		slog.Error("unmatched build")
		panic(-1)
	}
	slog.Info("image to process", "image", img)
	storagePath := d.GetBuildStorePath(img.Build)
	if dir, err := os.ReadDir(storagePath); len(dir) != 0 || err != nil {
		slog.Error("build storage not empty", "dir", dir, "err", err)
		return
	}

	// download and get hash
	lastline := lib.Sh_curl_sha256.RunReturnLastLineSplit(img.Url)
	zipFile := img.NewImageZip(lastline)
	if zipFile == nil {
		slog.Error("failed to download zip or verify sha256")
		panic(-1)
	}

	_ = zipFile.Unzip_boot_img()
	lib.Sh_cp.Run("init_boot.img", storagePath)
}

func collect_device_images() (devices map[string]*lib.Device) {
	collector := colly.NewCollector()
	devices = make(map[string]*lib.Device)

	collector.OnHTML("h2", func(h *colly.HTMLElement) {
		d := lib.TryNewDevice(h.Text)
		if d == nil {
			return
		}

		defer func() {
			devices[d.CodeName] = d
			// fmt.Printf("[device] %#v\n", d)
		}()
		table := h.DOM.Next().Find("table tbody tr")
		table.Each(func(i int, s *goquery.Selection) {
			// each images
			img := lib.TryNewImageFromTr(s)
			if img == nil {
				return
			}
			d.Images = append(d.Images, img)
			// fmt.Printf("%#v\n", img)
		})
	})

	cookie, _ := http.ParseCookie("cookies_accepted=true; django_language=en; devsite_wall_acks=nexus-image-tos")
	collector.SetCookies("https://developers.google.com", cookie)
	collector.Visit("https://developers.google.com/android/images")

	return
}
