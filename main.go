package main

import (
	"fmt"
	"net/http"
	"pixel_boot_img/lib"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

func main() {
	collector := colly.NewCollector()
	cookie, _ := http.ParseCookie("cookies_accepted=true; django_language=en; devsite_wall_acks=nexus-image-tos")
	collector.SetCookies("https://developers.google.com", cookie)

	collector.OnHTML("h2", func(h *colly.HTMLElement) {
		d := lib.TryNewDevice(h.Text)
		if d == nil {
			return
		}
		fmt.Printf("[device] %#v\n", *d)
		table := h.DOM.Next().Find("table tbody tr")
		table.Each(func(i int, s *goquery.Selection) {
			// each images
			img := lib.TryNewImageFromTr(s)
			if img == nil {
				return
			}
			d.Images = append(d.Images, img)
			fmt.Printf("%#v\n", img)
		})
		panic("end")
	})
	collector.Visit("https://developers.google.com/android/images")
}
