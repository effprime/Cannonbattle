package util

import (
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

func ImageFromPath(path string) (*ebiten.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	ebitenImage := ebiten.NewImageFromImage(img)
	return ebitenImage, nil
}

func TransparentImageFromPath(path string) (*ebiten.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	rgba, ok := img.(*image.RGBA)
	if !ok {
		b := img.Bounds()
		rgba = image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				rgba.Set(x, y, img.At(x, y))
			}
		}
	}

	makeTransparent(rgba)

	ebitenImage := ebiten.NewImageFromImage(rgba)
	return ebitenImage, nil
}

// makeTransparent processes the image to make white pixels transparent.
func makeTransparent(img *image.RGBA) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// Adjust the threshold according to your definition of "white".
			if r >= 0xffff && g >= 0xffff && b >= 0xffff {
				img.SetRGBA(x, y, color.RGBA{0, 0, 0, 0})
			}
		}
	}
}

func RandomRangeInt(min, max int) int {
	return rand.Intn(max-min) + min
}
