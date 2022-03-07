package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"os"

	"github.com/zachomedia/go-bdf"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var (
	outputPngFile       = flag.String("output_png_file", "out.png", "output png file")
	outputAnimationFile = flag.String("output_animation_file", "out.animation", "output png file")
)

func main() {
	flag.Parse()

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer f.Close()

	buf, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("failed to read file: %s", err)
	}

	fnt, err := bdf.Parse(buf)
	if err != nil {
		log.Fatalf("failed to parse font: %s", err)
	}

	img := image.NewRGBA(image.Rect(0, 0, 200, 1000))

	x := 0
	y := 0

	outAnimationF, err := os.Create(*outputAnimationFile)
	if err != nil {
		log.Fatalf("failed to open png file for writing: %s", err)
	}
	defer outAnimationF.Close()

	for _, c := range fnt.Characters {
		face := fnt.NewFace()

		if _, err := fmt.Fprintf(outAnimationF, "animation state=\"U%06X\"\n", c.Encoding); err != nil {
			log.Fatalf("failed to open animation file for writing: %s", err)
		}

		s := string(c.Encoding)
		_, advance := font.BoundString(face, s)

		if x+advance.Round() >= img.Rect.Dx() {
			x = 0
			y += fnt.Size + 1
		}

		d := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(color.RGBA{255, 255, 255, 255}),
			Face: face,
			Dot:  fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y + fnt.Size)},
		}

		d.DrawString(s)

		if _, err := fmt.Fprintf(outAnimationF, "frame duration=\"0\" x=\"%d\" y=\"%d\" width=\"%d\" height=\"%d\" originx=\"0\" originy=\"0\" flipx=\"0\" flipy=\"0\"\n\n", x, y, advance.Round(), fnt.Size); err != nil {
			log.Fatalf("failed to open animation file for writing: %s", err)
		}

		x += advance.Round() + 1
	}

	outPngF, err := os.Create(*outputPngFile)
	if err != nil {
		log.Fatalf("failed to open png file for writing: %s", err)
	}
	defer outPngF.Close()

	if err := png.Encode(outPngF, img); err != nil {
		log.Fatalf("failed to encode png: %s", err)
	}
}
