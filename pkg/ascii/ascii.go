package ascii

import (
	"fmt"
	"image"
	"regexp"
	"strconv"
	"strings"

	"github.com/qeesung/image2ascii/convert"
)

const (
	asciiWidth      = 100
	charAspectRatio = 0.5
)

func ConvertToASCIIArt(img image.Image) string {
	convertOptions := convert.DefaultOptions
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	convertOptions.FixedWidth = asciiWidth
	convertOptions.FixedHeight = int(float64(height) * (float64(asciiWidth) / float64(width)) * charAspectRatio)
	convertOptions.Colored = false
	convertOptions.Reversed = true

	converter := convert.NewImageConverter()
	return converter.Image2ASCIIString(img, &convertOptions)
}

func ConvertToColorASCIIArt(img image.Image, format string) string {
	options := convert.DefaultOptions
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	options.FixedWidth = asciiWidth
	options.FixedHeight = int(float64(height) * (float64(asciiWidth) / float64(width)) * charAspectRatio)
	options.Colored = true

	converter := convert.NewImageConverter()

	if format == "html" {
		asciiArt := converter.Image2ASCIIString(img, &options)
		return ansiToHTML(asciiArt)
	}

	return converter.Image2ASCIIString(img, &options)
}

var ansi8bitToRGB = func() map[int]string {
	m := make(map[int]string)
	baseColors := []string{
		"000000", "800000", "008000", "808000", "000080", "800080", "008080", "c0c0c0",
		"808080", "ff0000", "00ff00", "ffff00", "0000ff", "ff00ff", "00ffff", "ffffff",
	}
	for i, c := range baseColors {
		m[i] = "#" + c
	}

	i := 16
	for r := 0; r < 6; r++ {
		for g := 0; g < 6; g++ {
			for b := 0; b < 6; b++ {
				m[i] = fmt.Sprintf("#%02x%02x%02x", r*51, g*51, b*51)
				i++
			}
		}
	}

	for gray := 8; gray <= 238; gray += 10 {
		m[i] = fmt.Sprintf("#%02x%02x%02x", gray, gray, gray)
		i++
	}

	return m
}()

func ansiToHTML(asciiArt string) string {
	reColor := regexp.MustCompile(`\x1b\[38;5;(\d+)m`)
	asciiArt = reColor.ReplaceAllStringFunc(asciiArt, func(match string) string {
		matches := reColor.FindStringSubmatch(match)
		if len(matches) != 2 {
			return ""
		}
		code, _ := strconv.Atoi(matches[1])
		color, ok := ansi8bitToRGB[code]
		if !ok {
			color = "#ffffff"
		}
		return fmt.Sprintf(`<span style="color:%s">`, color)
	})

	asciiArt = strings.ReplaceAll(asciiArt, "\x1b[0m", "</span>")

	asciiArt = regexp.MustCompile(`\x1b\[[0-9;]*m`).ReplaceAllString(asciiArt, "")

	html := "<!DOCTYPE html><html><head><meta charset=\"UTF-8\">" +
		"<style>body { background: black; color: white; font-family: monospace; white-space: pre; }</style></head><body>" +
		asciiArt + "</body></html>"

	return html
}
