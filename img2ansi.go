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

var debugBuffer string

// type ANSIBuffer struct {
// 	Character [][]string
// }

const RESET = "\033[0m"
const ALPHATHRESHOLD = 30

// const FG = 38
// const BG = 48

func main() {
	fmt.Println("Running img2ansi")
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

	// var buffer ANSIBuffer

	for y := bounds.Min.Y; y < bounds.Max.Y; y += 2 {
		debugBuffer += fmt.Sprintf("Image Row: %d\n", y)
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			top := m.At(x, y)
			bottom := m.At(x, y+1)
			debugBuffer += fmt.Sprintf("  X= %2d: ", x)
			// buffer[x][y/2] = rgbToTrueColor(top, bottom)
			fmt.Print(rgbToTrueColor(top, bottom))
		}
		fmt.Printf("%s\n", RESET)
	}
	fmt.Println("")
	fmt.Println("DEBUG============")
	fmt.Println(debugBuffer)
}

func paletteMatch(c color.Color) color.RGBA {
	return color.RGBA{uint8(255), uint8(255), uint8(255), uint8(0)}
}

func rgbToTrueColor(top, bottom color.Color) string {
	var topcode, bottomcode string
	rt, gt, bt, at := top.RGBA()
	rb, gb, bb, ab := bottom.RGBA()
	var symbol string
	if uint8(ab) >= ALPHATHRESHOLD && uint8(at) >= ALPHATHRESHOLD {
		topcode = fmt.Sprintf("38;2;%d;%d;%d", uint8(rt), uint8(gt), uint8(bt))    // FG
		bottomcode = fmt.Sprintf("48;2;%d;%d;%d", uint8(rb), uint8(gb), uint8(bb)) // BG
		symbol = "▀"
	} else if uint8(ab) >= ALPHATHRESHOLD {
		topcode = fmt.Sprintf("49")
		bottomcode = fmt.Sprintf("38;2;%d;%d;%d", uint8(rb), uint8(gb), uint8(bb)) //FG
		symbol = "▄"
	} else if uint8(at) >= ALPHATHRESHOLD {
		topcode = fmt.Sprintf("38;2;%d;%d;%d", uint8(rt), uint8(gt), uint8(bt)) // FG
		bottomcode = fmt.Sprintf("49")                                          // BG
		symbol = "▀"
	} else {
		symbol = " "
	}

	debugBuffer += fmt.Sprintf("glyph: %s, top: %s, bottom: %s %s", symbol, topcode, bottomcode, RESET)

	debugBuffer += fmt.Sprintf("top {%3d, %3d, %3d, %3d}, bottom  {%3d, %3d, %3d, %3d}\n", uint8(rt), uint8(gt), uint8(bt), uint8(at), uint8(rb), uint8(gb), uint8(bb), uint8(ab))

	return fmt.Sprintf("\033[%s;%sm%s", topcode, bottomcode, symbol)
}

func rgbToXterm256(c color.Color) string {
	return "TODO"
}
