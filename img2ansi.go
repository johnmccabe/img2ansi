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
	reader, err := os.Open("images/ryu.png")
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
		fmt.Print("\033[0m\n")
	}
}

func paletteMatch(c color.Color) color.RGBA {
	return color.RGBA{uint8(255), uint8(255), uint8(255), uint8(0)}
}

func rgbToTrueColor(top, bottom color.Color) string {
	var topcode, bottomcode string
	rt, gt, bt, at := top.RGBA()
	rb, gb, bb, ab := bottom.RGBA()
	var symbol string
	if uint8(ab) >= 30 && uint8(at) >= 30 {
		symbol = "▄"
	} else if uint8(ab) >= 30 {
		symbol = "▄"
	} else if uint8(at) >= 30 {
		symbol = "▀"
	} else {
		symbol = " "
	}

	if uint8(ab) >= 30 {
		topcode = fmt.Sprintf("38;2;%d;%d;%d", rb, gb, bb)
	} else {
		topcode = "39"
	}
	if uint8(at) >= 30 {
		bottomcode = fmt.Sprintf("48;2;%d;%d;%d", rt, gt, bt)
	} else {
		bottomcode = "49"
	}
	return fmt.Sprintf("\033[%s;%sm%s", topcode, bottomcode, symbol)
}

func rgbToXterm256(c color.Color) string {
	return "TODO"
}
