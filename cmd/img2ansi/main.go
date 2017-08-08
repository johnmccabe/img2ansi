package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"path"

	"github.com/johnmccabe/img2ansi"
)

var imageFile string
var border int

func main() {
	if len(os.Args) != 2 {
		usage := `Usage:
	%s <path to image>`
		fmt.Println(fmt.Sprintf(usage, path.Base(os.Args[0])))
		os.Exit(1)
	}

	imageFile := os.Args[1]
	f, err := os.Open(imageFile)
	if err != nil {
		log.Fatalf("Error opening image [%s]: %v", imageFile, err)
	}
	defer f.Close()

	i, _, err := image.Decode(f)
	if err != nil {
		log.Fatalf("Error parsing image: %v", err)
	}

	ansi, err := img2ansi.RenderANSI256(i)
	if err != nil {
		log.Fatalf("Error converting image to ANSI: %v", err)
	}
	fmt.Print(ansi)
}
