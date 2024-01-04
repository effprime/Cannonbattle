package game

import (
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	GameWidth              = 1920
	GameHeight             = 1080
	PPS                    = 0.5
	Gravity                = 50.0
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
	ebiten.SetWindowSize(GameWidth, GameHeight)
	ebiten.SetWindowTitle("Cannonbattle")
	ebiten.SetWindowResizable(true)

	game, err := New()
	if err != nil {
		return err
	}
	return ebiten.RunGame(game)
}

type Game struct {
	Components *ComponentState
}

func New() (*Game, error) {
	state := NewComponentState()

	background := NewBackground(&BackgroundOptions{
		Height: GameHeight,
		Width:  GameWidth,
	})
	state.RegisterComponent(background)

	cannon := NewCannon(&CannonOptions{
		Length: 200,
		Width:  50,
		Angle:  0.78,
		Speed:  150,
	})
	state.RegisterComponent(cannon)

	enemy, err := NewEnemy("../../assets/enemy.png")
	if err != nil {
		return nil, err
	}
	state.RegisterComponent(enemy)

	return &Game{
		Components: state,
	}, nil
}

func (g *Game) Update() error {
	return g.Components.UpdateAll(g)
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Components.DrawAll(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return GameWidth, GameHeight
}

type Position struct {
	X float64
	Y float64
}
