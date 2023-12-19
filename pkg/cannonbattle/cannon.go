package cannonbattle

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Cannon struct {
	ShotImg *ebiten.Image
	Image   *ebiten.Image
	Length  int
	Width   int
	Angle   float64
	Speed   float64
	Shots   []Shot
}

func (c *Cannon) CreateImage() {
	img := ebiten.NewImage(c.Width, c.Length)
	img.Fill(color.Black)
	c.Image = img
}

func (c *Cannon) DrawOptions() *ebiten.DrawImageOptions {
	cannonOpts := &ebiten.DrawImageOptions{}
	cannonOpts.GeoM.Rotate(c.Angle + math.Pi)
	cannonOpts.GeoM.Translate(10, 0.80*float64(Height))
	return cannonOpts
}

func (c *Cannon) Fire() {
	c.Shots = append(c.Shots, Shot{
		VelY:  -math.Sin(math.Pi/2-c.Angle) * c.Speed,
		Speed: c.Speed,
		Angle: math.Pi/2 - c.Angle,
		Pos: Position{
			X: float64(c.Length) * math.Cos(math.Pi/2-c.Angle),
			Y: -float64(c.Length)*math.Sin(math.Pi/2-c.Angle) + 0.80*float64(Height),
		},
	})
}
