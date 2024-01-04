package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type HealthBar struct {
	opts    *HealthBarOptions
	current int
	img     *ebiten.Image
	pos     Position
}

type HealthBarOptions struct {
	Total int
	Width int
}

func NewHealthBar(opts *HealthBarOptions) *HealthBar {
	bgImg := ebiten.NewImage(opts.Width, 10)
	bgImg.Fill(color.RGBA{R: 0, G: 255, B: 0, A: 255})
	return &HealthBar{
		opts:    opts,
		current: opts.Total,
		img:     bgImg,
	}
}

func (h *HealthBar) Add(amount int) {
	defer h.updateImage()
	newCurrent := h.current + amount
	if newCurrent > h.opts.Total {
		h.current = h.opts.Total
	} else {
		h.current = newCurrent
	}
}

func (h *HealthBar) Subtract(amount int) {
	defer h.updateImage()
	newCurrent := h.current - amount
	if newCurrent < 0 {
		h.current = 0
	} else {
		h.current = newCurrent
	}
}

func (h *HealthBar) SetPosition(p Position) {
	h.pos = p
}

func (h *HealthBar) updateImage() {
	bgImg := ebiten.NewImage(h.opts.Width, 5)
	bgImg.Fill(color.RGBA{R: 255, G: 0, B: 0, A: 255})
	healthWidth := int(float64(h.opts.Width) * float64(h.current) / float64(h.opts.Total))
	if healthWidth != 0 {
		healthBar := ebiten.NewImage(healthWidth, 5)
		healthBar.Fill(color.RGBA{R: 0, G: 255, B: 0, A: 255})
		bgImg.DrawImage(healthBar, &ebiten.DrawImageOptions{})
	}
	h.img = bgImg
}

func (h *HealthBar) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(h.pos.X, h.pos.Y)
	screen.DrawImage(h.img, opts)
}
