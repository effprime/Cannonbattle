package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	ComponentNameCannon = "cannon"
)

type CannonOptions struct {
	Length int
	Width  int
	Angle  float64
	Speed  float64
}

type Cannon struct {
	opts          *CannonOptions
	shotImg       *ebiten.Image
	pathMarkerImg *ebiten.Image
	img           *ebiten.Image
	shots         []Shot
}

func NewCannon(opts *CannonOptions) *Cannon {
	img := ebiten.NewImage(opts.Width, opts.Length)
	img.Fill(color.Black)

	shot := ebiten.NewImage(20, 20)
	shot.Fill(color.Black)

	pathMarker := ebiten.NewImage(10, 10)
	pathMarker.Fill(color.RGBA{R: 255, A: 255})

	return &Cannon{
		opts:          opts,
		img:           img,
		shotImg:       shot,
		pathMarkerImg: pathMarker,
		shots:         []Shot{},
	}
}

func (c *Cannon) Name() string {
	return ComponentNameCannon
}

func (c *Cannon) Update(g *Game) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		c.opts.Angle -= 3 * math.Pi / 180
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		c.opts.Angle += 3 * math.Pi / 180
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		c.opts.Speed -= 10
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		c.opts.Speed += 10
	} else if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		c.Fire()
	}

	dt := 1.0 / ebiten.ActualTPS()
	remainingShots := []Shot{}
	for idx, s := range c.shots {
		s.Update(dt)
		if s.Elapsed > PPS {
			s.AddPathMarker()
			s.Elapsed = 0
			c.shots[idx] = s
		} else {
			s.Elapsed += dt
		}
		if s.Pos.Y <= 0.8*float64(GameHeight) {
			remainingShots = append(remainingShots, s)
		}
	}
	c.shots = remainingShots
	return nil
}

func (c *Cannon) Draw(screen *ebiten.Image) {
	cannonOpts := &ebiten.DrawImageOptions{}
	cannonOpts.GeoM.Rotate(c.opts.Angle + math.Pi)
	cannonOpts.GeoM.Translate(10, 0.80*float64(GameHeight))
	screen.DrawImage(c.img, cannonOpts)

	for _, s := range c.shots {
		screen.DrawImage(c.shotImg, s.DrawOptions())
		for _, p := range s.Path {
			pathMarkerOpts := &ebiten.DrawImageOptions{}
			pathMarkerOpts.GeoM.Translate(p.X, p.Y)
			screen.DrawImage(c.pathMarkerImg, pathMarkerOpts)
		}
	}
}

func (c *Cannon) Fire() {
	c.shots = append(c.shots, Shot{
		VelY:  -math.Sin(math.Pi/2-c.opts.Angle) * c.opts.Speed,
		Speed: c.opts.Speed,
		Angle: math.Pi/2 - c.opts.Angle,
		Pos: Position{
			X: float64(c.opts.Length) * math.Cos(math.Pi/2-c.opts.Angle),
			Y: -float64(c.opts.Length)*math.Sin(math.Pi/2-c.opts.Angle) + 0.80*float64(GameHeight),
		},
	})
}

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
