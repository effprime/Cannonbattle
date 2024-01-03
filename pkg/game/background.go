package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ComponentNameBackground = "background"
)

type BackgroundOptions struct {
	Height int
	Width  int
}

type Background struct {
	opts *BackgroundOptions
	img  *ebiten.Image
}

func NewBackground(opts *BackgroundOptions) *Background {
	img := ebiten.NewImage(opts.Width, opts.Height)
	ground := ebiten.NewImage(opts.Width, int(0.20*float64(opts.Height)))

	img.Fill(color.RGBA{R: 0, G: 100, B: 100, A: 255})
	ground.Fill(color.RGBA{R: 1, G: 50, B: 32, A: 255})

	bgOpts := &ebiten.DrawImageOptions{}
	bgOpts.GeoM.Translate(0, 0.80*float64(opts.Height))
	img.DrawImage(ground, bgOpts)
	return &Background{
		opts: opts,
		img:  img,
	}
}

func (b *Background) Name() string {
	return ComponentNameBackground
}

func (b *Background) Update(g *Game) error {
	return nil
}

func (b *Background) Draw(screen *ebiten.Image) {
	screen.DrawImage(b.img, &ebiten.DrawImageOptions{})
}
