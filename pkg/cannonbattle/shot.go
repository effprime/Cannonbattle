package cannonbattle

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Shot struct {
	Elapsed float64
	Speed   float64
	VelY    float64
	Angle   float64
	Pos     Position
	Path    []Position
}

func (s *Shot) Update(dt float64) {
	s.VelY = s.VelY + Gravity*dt
	s.Pos.X = s.Pos.X + s.Speed*math.Cos(s.Angle)*dt
	s.Pos.Y = s.Pos.Y + s.VelY*dt
}

func (s *Shot) AddPathMarker() {
	s.Path = append(s.Path, s.Pos)
	if len(s.Path) > MaxPath {
		s.Path = s.Path[1:]
	}
}

func (s *Shot) DrawOptions() *ebiten.DrawImageOptions {
	shotOpts := &ebiten.DrawImageOptions{}
	shotOpts.GeoM.Rotate(math.Pi/4 + math.Pi)
	shotOpts.GeoM.Translate(float64(s.Pos.X), float64(s.Pos.Y))
	return shotOpts
}
