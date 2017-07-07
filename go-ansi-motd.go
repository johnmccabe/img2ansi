package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"os"
	"reflect"
)

// ANSI256Palette TODO
var ANSI256Palette = []color.Color{
	color.RGBA{0x00, 0x00, 0x44, 0xff},
}

func main() {
	fmt.Println("Running go-ansi-motd")
	reader, err := os.Open("images/mark_small.png")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	m, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	bounds := m.Bounds()
	fmt.Printf("Bounds: %v\n", bounds)
	fmt.Printf("Space:  %s\n", reflect.TypeOf(m.ColorModel()).String())

	for y := bounds.Min.Y; y < bounds.Max.Y; y += 2 {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			top := m.At(x, y)
			bottom := m.At(x, y+1)
			fmt.Print(rgbToTrueColor(top, bottom))
		}
		fmt.Print("\n")
	}
}

func paletteMatch(c color.Color) color.RGBA {
	return color.RGBA{uint8(255), uint8(255), uint8(255), uint8(0)}
}

func rgbToTrueColor(top, bottom color.Color) string {
	rt, gt, bt, _ := top.RGBA()
	rb, gb, bb, _ := bottom.RGBA()
	return fmt.Sprintf("\033[38;2;%d;%d;%d;48;2;%d;%d;%dm▄", uint8(rb), uint8(gb), uint8(bb), uint8(rt), uint8(gt), uint8(bt))
}

func rgbToXterm256(c color.Color) string {
	return "TODO"
}
