package sprite

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"strings"
)

// Render converts PNG image data into colored half-block terminal art at most
// maxWidth characters wide. Transparent pixels become spaces or default background.
// A fully transparent image renders to an empty string.
func Render(data []byte, maxWidth int) (string, error) {
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	minX, minY, maxX, maxY, ok := opaqueBounds(img)
	if !ok {
		return "", nil
	}
	cw := maxX - minX + 1
	ch := maxY - minY + 1

	outW := cw
	if maxWidth > 0 && outW > maxWidth {
		outW = maxWidth
	}
	if outW < 1 {
		outW = 1
	}
	scale := float64(cw) / float64(outW)
	outH := int(float64(ch)/scale + 0.5)
	if outH < 1 {
		outH = 1
	}
	if outH%2 == 1 {
		outH++
	}

	sample := func(ox, oy int) (r, g, b uint8, opaque bool) {
		sx := minX + int(float64(ox)*scale)
		sy := minY + int(float64(oy)*scale)
		if sx > maxX {
			sx = maxX
		}
		if sy > maxY {
			sy = maxY
		}
		cr, cg, cb, ca := img.At(sx, sy).RGBA()
		if ca>>8 < 128 {
			return 0, 0, 0, false
		}
		return uint8(cr >> 8), uint8(cg >> 8), uint8(cb >> 8), true
	}

	var sb strings.Builder
	for row := 0; row < outH/2; row++ {
		for ox := 0; ox < outW; ox++ {
			tr, tg, tb, top := sample(ox, row*2)
			br, bg, bb, bot := sample(ox, row*2+1)
			switch {
			case !top && !bot:
				sb.WriteByte(' ')
			case top && !bot:
				fmt.Fprintf(&sb, "\x1b[38;2;%d;%d;%dm▀", tr, tg, tb)
			case !top && bot:
				fmt.Fprintf(&sb, "\x1b[38;2;%d;%d;%dm▄", br, bg, bb)
			default:
				fmt.Fprintf(&sb, "\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm▀", tr, tg, tb, br, bg, bb)
			}
		}
		sb.WriteString("\x1b[0m\n")
	}
	return sb.String(), nil
}

// opaqueBounds returns the bounding box of pixels with alpha >= 128.
func opaqueBounds(img image.Image) (minX, minY, maxX, maxY int, ok bool) {
	b := img.Bounds()
	minX, minY = b.Max.X, b.Max.Y
	maxX, maxY = b.Min.X-1, b.Min.Y-1
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a>>8 >= 128 {
				if x < minX {
					minX = x
				}
				if y < minY {
					minY = y
				}
				if x > maxX {
					maxX = x
				}
				if y > maxY {
					maxY = y
				}
				ok = true
			}
		}
	}
	return
}
