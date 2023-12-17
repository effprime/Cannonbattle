package cannonbattle

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	Width   = 1280
	Height  = 720
	PPS     = 0.5
	Gravity = 10
	MaxPath = 10
)

var (
	mplusNormalFont font.Face
)

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		panic(err)
	}
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		panic(err)
	}
}

func Run() error {
	ebiten.SetWindowSize(Width, Height)
	ebiten.SetWindowTitle("Cannonbattle")
	ebiten.SetWindowResizable(true)

	game, err := NewGame()
	if err != nil {
		return err
	}
	return ebiten.RunGame(game)
}

type Game struct {
	Cannon     *Cannon
	Shot       *ebiten.Image
	PathMarker *ebiten.Image
	Background *ebiten.Image
}

func NewGame() (*Game, error) {
	sky := ebiten.NewImage(Width, Height)
	ground := ebiten.NewImage(Width, int(0.20*float64(Height)))

	sky.Fill(color.RGBA{R: 0, G: 100, B: 100, A: 255})
	ground.Fill(color.RGBA{R: 1, G: 50, B: 32, A: 255})

	skyOpts := &ebiten.DrawImageOptions{}
	skyOpts.GeoM.Translate(0, 0.80*float64(Height))
	sky.DrawImage(ground, skyOpts)

	cannon := Cannon{
		Speed:  150,
		Length: 200,
		Width:  50,
		Angle:  0.78,
	}

	shot := ebiten.NewImage(20, 20)
	shot.Fill(color.Black)

	pathMarker := ebiten.NewImage(10, 10)
	pathMarker.Fill(color.RGBA{R: 255, A: 255})
	return &Game{
		Background: sky,
		Cannon:     &cannon,
		Shot:       shot,
		PathMarker: pathMarker,
	}, nil
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.Cannon.Angle -= 3 * math.Pi / 180
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		g.Cannon.Angle += 3 * math.Pi / 180
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		g.Cannon.Speed -= 10
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		g.Cannon.Speed += 10
	} else if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.Cannon.Fire()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.Background, &ebiten.DrawImageOptions{})
	screen.DrawImage(g.Cannon.Draw(), g.Cannon.DrawOptions())

	dt := 1.0 / ebiten.ActualTPS()
	remainingShots := []Shot{}
	for idx, s := range g.Cannon.Shots {
		screen.DrawImage(g.Shot, s.DrawOptions())
		for _, p := range s.Path {
			pathMarkerOpts := &ebiten.DrawImageOptions{}
			pathMarkerOpts.GeoM.Translate(p.X, p.Y)
			screen.DrawImage(g.PathMarker, pathMarkerOpts)
		}
		s.Update(dt)
		if s.Elapsed > PPS {
			s.AddPathMarker()
			g.Cannon.Shots[idx] = s
			s.Elapsed = 0
		} else {
			s.Elapsed += dt
		}
		if s.PosY <= 0.8*float64(Height) {
			remainingShots = append(remainingShots, s)
		}
	}
	g.Cannon.Shots = remainingShots

	text.Draw(screen, fmt.Sprintf("Speed: %.2f", g.Cannon.Speed), mplusNormalFont, 10, Height-10, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return Width, Height
}

type Cannon struct {
	Image  *ebiten.Image
	Length int
	Width  int
	Angle  float64
	Speed  float64
	Shots  []Shot
}

func (c *Cannon) Draw() *ebiten.Image {
	if c.Image == nil {
		c.Image = ebiten.NewImage(c.Width, c.Length)
		c.Image.Fill(color.Black)
	}
	return c.Image
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
		PosX:  float64(c.Length) * math.Cos(math.Pi/2-c.Angle),
		PosY:  -float64(c.Length)*math.Sin(math.Pi/2-c.Angle) + 0.80*float64(Height),
	})
}

type Shot struct {
	Elapsed float64
	Speed   float64
	VelY    float64
	Angle   float64
	PosX    float64
	PosY    float64
	Path    []Position
}

type Position struct {
	X float64
	Y float64
}

func (s *Shot) Update(dt float64) {
	s.VelY = s.VelY + Gravity*dt
	s.PosX = s.PosX + s.Speed*math.Cos(s.Angle)*dt
	s.PosY = s.PosY + s.VelY*dt
}

func (s *Shot) AddPathMarker() {
	s.Path = append(s.Path, Position{
		X: s.PosX, Y: s.PosY,
	})
	if len(s.Path) > MaxPath {
		s.Path = s.Path[1:]
	}
}

func (s *Shot) DrawOptions() *ebiten.DrawImageOptions {
	shotOpts := &ebiten.DrawImageOptions{}
	shotOpts.GeoM.Rotate(math.Pi/4 + math.Pi)
	shotOpts.GeoM.Translate(float64(s.PosX), float64(s.PosY))
	return shotOpts
}
