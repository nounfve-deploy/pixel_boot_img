package main

import (
	"fmt"
	"log/slog"
	"os"
	"pixel_boot_img/lib"
	"strings"
)

func init() {
	if strings.ToUpper(os.Getenv("ALT_MAIN")) != "PARSE_ISSUE" {
		return
	}
	parse_issue()
	os.Exit(0)
}

func parse_issue() {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		slog.Error("Recovered from panic", "error", r)
		fmt.Println("usage: ALT_MAIN=PARSE_ISSUE pixel_boot_img [zip_url]")
	}()
	input := os.Args[1]
	fmt.Println(input)
	device, build := lib.PARSE_DOWNLOAD_URL(input)
	fmt.Println(device, build)
}
