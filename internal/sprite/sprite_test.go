package sprite

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"strings"
	"testing"
)

func encodePNG(t *testing.T, img image.Image) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestRenderColors(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	red := color.RGBA{255, 0, 0, 255}
	blue := color.RGBA{0, 0, 255, 255}
	img.Set(0, 0, red)
	img.Set(1, 0, red)
	img.Set(0, 1, blue)
	img.Set(1, 1, blue)

	out, err := Render(encodePNG(t, img), 2)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "▀") {
		t.Error("expected a half-block char")
	}
	if !strings.Contains(out, "38;2;255;0;0") {
		t.Error("expected red foreground code")
	}
	if !strings.Contains(out, "48;2;0;0;255") {
		t.Error("expected blue background code")
	}
}

func TestRenderTransparent(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 2, 2)) // all-zero alpha
	out, err := Render(encodePNG(t, img), 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "38;2;") {
		t.Errorf("transparent image must have no color codes: %q", out)
	}
}

func TestRenderResetsBackgroundOnTransparency(t *testing.T) {
	// A column with an opaque pixel directly above a transparent pixel must reset
	// the background on that half-block so it cannot inherit a neighbor's color.
	// Opaque pixels above and below keep the transparent pixel inside the crop box.
	img := image.NewRGBA(image.Rect(0, 0, 1, 3))
	img.Set(0, 0, color.RGBA{255, 0, 0, 255})
	// (0,1) left transparent
	img.Set(0, 2, color.RGBA{255, 0, 0, 255})

	out, err := Render(encodePNG(t, img), 1)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "\x1b[38;2;255;0;0m\x1b[49m▀") {
		t.Errorf("expected red top half with a background reset, got %q", out)
	}
}

func TestRenderMalformed(t *testing.T) {
	if _, err := Render([]byte("not a png"), 10); err == nil {
		t.Fatal("expected an error for malformed input")
	}
}
