package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"pixel_boot_img/lib"
	"slices"
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
	body := os.Args[1]
	fmt.Println("==============+body+==============")
	fmt.Println(body)
	fmt.Println("==============-body-==============")

	lines := strings.Split(body, "\n")
	lines = slices.DeleteFunc(lines, func(line string) bool { return line == "" })
	lasline := ""
	if len(lines) > 0 {
		lasline = lines[len(lines)-1]
	}

	device, tag := lib.PARSE_DOWNLOAD_URL(lasline)
	if device == "" || tag == "" {
		return
	}

	fmt.Println(device, tag)
	action_output("DEVICE", device)
	action_output("TAG", tag)
}

func action_output(K, V string) {
	echo := fmt.Sprintf(`echo "%s==%s" >> $GITHUB_OUTPUT`, K, V)
	exec.Command("bash", "-c", echo).Run()
}
