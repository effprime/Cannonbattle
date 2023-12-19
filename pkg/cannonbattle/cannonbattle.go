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
	Width                  = 1280
	Height                 = 720
	PPS                    = 0.5
	Gravity                = 10
	MaxPath                = 10
	EnemySpawnShiftPercent = 10.0
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
	Pool ImagePool
}

type ImagePool struct {
	Cannon     *Cannon
	Enemy      *Enemy
	Shot       *ebiten.Image
	PathMarker *ebiten.Image
	Background *ebiten.Image
}

func NewGame() (*Game, error) {
	bg := ebiten.NewImage(Width, Height)
	ground := ebiten.NewImage(Width, int(0.20*float64(Height)))

	bg.Fill(color.RGBA{R: 0, G: 100, B: 100, A: 255})
	ground.Fill(color.RGBA{R: 1, G: 50, B: 32, A: 255})

	bgOpts := &ebiten.DrawImageOptions{}
	bgOpts.GeoM.Translate(0, 0.80*float64(Height))
	bg.DrawImage(ground, bgOpts)

	shot := ebiten.NewImage(20, 20)
	shot.Fill(color.Black)

	pathMarker := ebiten.NewImage(10, 10)
	pathMarker.Fill(color.RGBA{R: 255, A: 255})

	cannon := &Cannon{
		Speed:  150,
		Length: 200,
		Width:  50,
		Angle:  0.78,
	}
	cannon.CreateImage()

	enemy, err := NewEnemy("../../assets/enemy.png")
	if err != nil {
		return nil, err
	}

	return &Game{
		Pool: ImagePool{
			Background: bg,
			Cannon:     cannon,
			Shot:       shot,
			PathMarker: pathMarker,
			Enemy:      enemy,
		},
	}, nil
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.Pool.Cannon.Angle -= 3 * math.Pi / 180
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		g.Pool.Cannon.Angle += 3 * math.Pi / 180
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		g.Pool.Cannon.Speed -= 10
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		g.Pool.Cannon.Speed += 10
	} else if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.Pool.Cannon.Fire()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.Pool.Background, &ebiten.DrawImageOptions{})
	screen.DrawImage(g.Pool.Cannon.Image, g.Pool.Cannon.DrawOptions())

	enemyOpts := &ebiten.DrawImageOptions{}
	enemyOpts.GeoM.Scale(g.Pool.Enemy.Scale, g.Pool.Enemy.Scale)
	enemyOpts.GeoM.Translate(g.Pool.Enemy.Position.X, g.Pool.Enemy.Position.Y)
	screen.DrawImage(g.Pool.Enemy.Image, enemyOpts)

	dt := 1.0 / ebiten.ActualTPS()
	remainingShots := []Shot{}
	for idx, s := range g.Pool.Cannon.Shots {
		screen.DrawImage(g.Pool.Shot, s.DrawOptions())
		for _, p := range s.Path {
			pathMarkerOpts := &ebiten.DrawImageOptions{}
			pathMarkerOpts.GeoM.Translate(p.X, p.Y)
			screen.DrawImage(g.Pool.PathMarker, pathMarkerOpts)
		}
		s.Update(dt)
		if s.Elapsed > PPS {
			s.AddPathMarker()
			g.Pool.Cannon.Shots[idx] = s
			s.Elapsed = 0
		} else {
			s.Elapsed += dt
		}
		if g.Pool.Enemy.InHitbox(s.Pos) {
			g.Pool.Enemy.TakeDamage(10)
			if g.Pool.Enemy.Health.Current == 0 {
				fmt.Println("Enemy is dead!")
			} else {
				fmt.Printf("10 dmg to enemy! Health remaining: %v\n", g.Pool.Enemy.Health.Current)
			}
		} else if s.Pos.Y <= 0.8*float64(Height) {
			remainingShots = append(remainingShots, s)
		}
	}
	g.Pool.Cannon.Shots = remainingShots

	text.Draw(screen, fmt.Sprintf("Speed: %.2f", g.Pool.Cannon.Speed), mplusNormalFont, 10, Height-10, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return Width, Height
}

type Position struct {
	X float64
	Y float64
}
