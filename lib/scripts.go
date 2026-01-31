package lib

import (
	"bytes"
	_ "embed"
	"fmt"
	"os/exec"
	"slices"
	"strings"
)

//go:embed scripts/curl_sha256.sh
var curl_sha256 string

//go:embed scripts/unzip_ls.sh
var unzip_ls string

//go:embed scripts/cp.sh
var cp string

var (
	Sh_curl_sha256 = &ShellScript{&curl_sha256}
	Sh_unzip_ls    = &ShellScript{&unzip_ls}
	Sh_cp          = &ShellScript{&cp}
)

type ShellScript struct {
	*string
}

func (ss *ShellScript) Run(args ...string) (output []byte) {
	cmd := exec.Command("/bin/sh")
	cmd.Stdin = strings.NewReader(SetString(args...) + *ss.string)
	fmt.Println("[script-start]")
	output, err := cmd.CombinedOutput()
	fmt.Printf("%s\n", string(output))
	fmt.Println("[script-end]")
	if err != nil {
		fmt.Printf("Error running script: %v\n", err)
		panic("script error")
	}
	return
}

func (ss *ShellScript) RunReturnLastLine(args ...string) (output []byte) {
	output = ss.Run(args...)
	lines := bytes.Split(output, []byte{'\n'})
	slices.Reverse(lines)
	for _, line := range lines {
		if len(line) != 0 {
			output = line
			break
		}
	}
	return
}

func (ss *ShellScript) RunReturnLastLineSplit(args ...string) (output []string) {
	outputByte := ss.RunReturnLastLine(args...)
	output = strings.Split(string(outputByte), " ")
	output = slices.DeleteFunc(output, func(e string) bool { return e == "" })
	return
}

func SetString(args ...string) (setter string) {
	if len(args) == 0 {
		return
	}
	setter = strings.Join(args, `" "`)
	setter = `set -- "` + setter + `" ` + "\n"
	return
}
