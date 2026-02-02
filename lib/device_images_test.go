package lib_test

import (
	"pixel_boot_img/lib"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTitleRegex(t *testing.T) {
	rango := `"rango" for Pixel 10 Pro Fold`
	match := lib.TITLE_REGEX.FindStringSubmatch(rango)
	assert.Equal(t, len(match), 3)
}

func TestBuildRegex(t *testing.T) {
	build_string := `16.0.0 (BD3A.251005.003.F1, Oct 2025, Japan)`
	build, date := lib.PARSE_BUILD_DATE(build_string)
	assert.Equal(t, build, "BD3A.251005.003.F1")
	assert.Equal(t, date, "Oct 2025")
}

func TestUrlRegex(t *testing.T) {
	url := `https://dl.google.com/dl/android/aosp/rango-bp4a.260105.004.e1-factory-0765b33a.zip`
	device, build := lib.PARSE_DOWNLOAD_URL(url)
	assert.Equal(t, device, "rango")
	assert.Equal(t, build, "bp4a.260105.004.e1")
}
