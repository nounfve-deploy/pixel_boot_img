package lib

import (
	"bytes"
	_ "embed"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"slices"
	"strings"
)

var (
	//go:embed scripts/curl_sha256.sh
	curl_sha256 string
	//go:embed scripts/unzip_ls.sh
	unzip_ls string
	//go:embed scripts/cp.sh
	cp string
)

var (
	Sh_curl_sha256 = &ShellScript{&curl_sha256}
	Sh_unzip_ls    = &ShellScript{&unzip_ls}
	Sh_cp          = &ShellScript{&cp}
)

type ShellScript struct {
	*string
}

func (ss *ShellScript) Run(args ...string) (output []byte) {
	var (
		captureBuffer = bytes.Buffer{}
		SetIn         = strings.NewReader(SetString(args...) + *ss.string)
		TeeOut        = io.MultiWriter(os.Stdout, &captureBuffer)
		TeeErr        = io.MultiWriter(os.Stdout, &captureBuffer)
	)

	cmd := exec.Command("/bin/sh")
	cmd.Stdin = SetIn
	cmd.Stdout = TeeOut
	cmd.Stderr = TeeErr

	slog.Debug("[script-start]")
	defer slog.Debug("[script-end]")
	err := cmd.Run()
	output = captureBuffer.Bytes()
	if err != nil {
		slog.Error("Error running script", "err", err)
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
