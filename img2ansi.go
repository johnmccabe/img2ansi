package img2ansi

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/johnmccabe/img2ansi/palette"
)

const alphaThreshold = 30

// ansiBlock
type ansiBlock struct {
	FG    string
	BG    string
	Glyph string
}

// RenderANSI16 TODO
func RenderANSI16(i image.Image) (string, error) {
	return renderANSI(i, palette.Xterm16)
}

// RenderANSI256 TODO
func RenderANSI256(i image.Image) (string, error) {
	return renderANSI(i, palette.Xterm256)
}

// RenderTrueColor TODO
func RenderTrueColor(i image.Image) (string, error) {
	return renderANSI(i, nil)
}

func renderANSI(i image.Image, p color.Palette) (string, error) {
	bounds := i.Bounds()
	rowbuffer := make([][]ansiBlock, (bounds.Max.Y+1)/2)

	for row := 0; row < bounds.Max.Y; row += 2 {
		rowbuffer[row/2] = make([]ansiBlock, bounds.Max.X)
		for x := 0; x < bounds.Max.X; x++ {
			top := i.At(x, row)
			bottom := i.At(x, row+1)
			if p == nil {
				rowbuffer[row/2] = append(rowbuffer[row/2], rgbToTrueColor(top, bottom))

			} else {
				rowbuffer[row/2] = append(rowbuffer[row/2], rgbToXterm(top, bottom, p))
			}
		}
	}
	return bufferToString(rowbuffer), nil
}

func bufferToString(rowbuffer [][]ansiBlock) string {
	var s string
	for y := 0; y < len(rowbuffer); y++ {
		for x := 0; x < len(rowbuffer[0]); x++ {
			if x == 0 {
				s += rowbuffer[y][x].printCompact(ansiBlock{})
			}
			if x > 0 {
				s += rowbuffer[y][x].printCompact(rowbuffer[y][x-1])
			}
		}
		// Reset escape codes and add a newline
		s += "\033[0m"
		s += "\n"
	}
	return s
}

func (b ansiBlock) printCompact(previous ansiBlock) (out string) {
	if (b.FG == previous.FG) && (b.BG == previous.BG) {
		out = fmt.Sprintf("%s", b.Glyph)
	} else if b.FG == previous.FG {
		out = fmt.Sprintf("\033[%sm%s", b.BG, b.Glyph)
	} else if b.BG == previous.BG {
		out = fmt.Sprintf("\033[%sm%s", b.FG, b.Glyph)
	} else {
		out = fmt.Sprintf("\033[%s;%sm%s", b.BG, b.FG, b.Glyph)
	}
	return out
}

func rgbToTrueColor(top, bottom color.Color) ansiBlock {
	var fg, bg string
	rt, gt, bt, at := top.RGBA()
	rb, gb, bb, ab := bottom.RGBA()
	var symbol string
	if uint8(ab) >= alphaThreshold && uint8(at) >= alphaThreshold {
		symbol = "▀"
		fg = fmt.Sprintf("38;2;%d;%d;%d", uint8(rt), uint8(gt), uint8(bt))
		bg = fmt.Sprintf("48;2;%d;%d;%d", uint8(rb), uint8(gb), uint8(bb))
	} else if uint8(ab) >= alphaThreshold {
		symbol = "▄"
		fg = fmt.Sprintf("38;2;%d;%d;%d", uint8(rb), uint8(gb), uint8(bb))
		bg = fmt.Sprintf("49")
	} else if uint8(at) >= alphaThreshold {
		symbol = "▀"
		fg = fmt.Sprintf("38;2;%d;%d;%d", uint8(rt), uint8(gt), uint8(bt))
		bg = fmt.Sprintf("49")
	} else {
		symbol = " "
		fg = "39"
		bg = "49"
	}
	return ansiBlock{FG: fg, BG: bg, Glyph: symbol}
}

func rgbToXterm(top, bottom color.Color, p color.Palette) ansiBlock {
	var fg, bg string
	_, _, _, at := top.RGBA()
	_, _, _, ab := bottom.RGBA()
	var symbol string
	if uint8(ab) >= alphaThreshold && uint8(at) >= alphaThreshold {
		symbol = "▀"
		fg = fmt.Sprintf("38;5;%d", p.Index(top))
		bg = fmt.Sprintf("48;5;%d", p.Index(bottom))
	} else if uint8(ab) >= alphaThreshold {
		symbol = "▄"
		fg = fmt.Sprintf("38;5;%d", p.Index(bottom))
		bg = fmt.Sprintf("49")
	} else if uint8(at) >= alphaThreshold {
		symbol = "▀"
		fg = fmt.Sprintf("38;5;%d", p.Index(top))
		bg = fmt.Sprintf("49")
	} else {
		symbol = " "
		fg = "39"
		bg = "49"
	}
	return ansiBlock{FG: fg, BG: bg, Glyph: symbol}
}
