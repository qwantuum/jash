package evaluator

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"strings"
)

const asciiChars = "@%#*+=-:. "

func imageASCIIFunc(args ...Object) Object {
	if len(args) != 1 {
		return &Error{Message: "image.ascii() requires exactly 1 argument: file path or URL"}
	}

	source, ok := args[0].(*String)
	if !ok {
		return &Error{Message: "image.ascii() argument must be a string (file path or URL)"}
	}

	path := source.Value

	var reader io.ReadCloser
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		resp, err := http.Get(path)
		if err != nil {
			return &Error{Message: fmt.Sprintf("failed to fetch image URL: %s", err)}
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return &Error{Message: fmt.Sprintf("HTTP %d fetching image URL", resp.StatusCode)}
		}
		reader = resp.Body
	} else {
		f, err := os.Open(path)
		if err != nil {
			return &Error{Message: fmt.Sprintf("failed to open image file: %s", err)}
		}
		reader = f
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		return &Error{Message: fmt.Sprintf("failed to decode image: %s", err)}
	}

	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	termWidth := 80
	aspect := float64(h) / float64(w)
	charWidth := termWidth
	charHeight := int(float64(charWidth) * aspect * 0.45)
	if charHeight < 1 {
		charHeight = 1
	}

	var result strings.Builder
	for y := 0; y < charHeight; y++ {
		for x := 0; x < charWidth; x++ {
			srcX := x * w / charWidth
			srcY := y * h / charHeight
			r, g, b, _ := img.At(srcX, srcY).RGBA()
			brightness := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
			brightness = brightness / 65535.0
			idx := int(brightness * float64(len(asciiChars)-1))
			if idx < 0 {
				idx = 0
			}
			if idx >= len(asciiChars) {
				idx = len(asciiChars) - 1
			}
			result.WriteByte(asciiChars[idx])
		}
		result.WriteByte('\n')
	}

	return &String{Value: result.String()}
}
