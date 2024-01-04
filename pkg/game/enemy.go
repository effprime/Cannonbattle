package game

import (
	"math"
	"time"

	"github.com/effprime/cannonbattle/pkg/util"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ComponentNameEnemy    = "enemy"
	EnemyHitPointsDefault = 1000
)

type Enemy struct {
	health    Hitpoints
	healthBar *HealthBar
	img       *ebiten.Image
	pos       Position
	scale     float64
	lastMoved time.Time
}

func NewEnemy(path string) (*Enemy, error) {
	img, err := util.TransparentImageFromPath(path)
	if err != nil {
		return nil, err
	}

	health := Hitpoints{
		Current: EnemyHitPointsDefault,
		Max:     EnemyHitPointsDefault,
	}
	healthBar := NewHealthBar(&HealthBarOptions{
		Total: health.Max,
		Width: (70 * img.Bounds().Dx()) / 100,
	})

	enemy := &Enemy{
		img:       img,
		scale:     0.7,
		health:    health,
		healthBar: healthBar,
	}

	enemy.NewPosition()
	enemy.healthBar.SetPosition(enemy.pos)

	return enemy, nil
}

func (e *Enemy) Name() string {
	return ComponentNameEnemy
}

func (e *Enemy) Update(g *Game) error {
	if time.Since(e.lastMoved) > 5*time.Second {
		e.NewPosition()
	}
	return nil
}

func (e *Enemy) Draw(screen *ebiten.Image) {
	enemyOpts := &ebiten.DrawImageOptions{}
	enemyOpts.GeoM.Scale(e.scale, e.scale)
	enemyOpts.GeoM.Translate(e.pos.X, e.pos.Y)

	screen.DrawImage(e.img, enemyOpts)
	e.healthBar.Draw(screen)
}

func (e *Enemy) NewPosition() {
	midpoint := GameWidth / 2.0
	min := int(EnemySpawnShiftPercent / 100.0 * midpoint)
	max := int(((EnemySpawnShiftPercent / 100.0) + 1.0) * midpoint)
	e.pos = Position{
		X: float64(util.RandomRangeInt(min, max)),
		Y: 0.8*GameHeight - e.scale*float64(e.img.Bounds().Dy()),
	}
	e.lastMoved = time.Now()
	e.healthBar.SetPosition(e.pos)
}

func (e *Enemy) InHitbox(p Position) bool {
	centerX := e.pos.X
	centerY := e.pos.Y + e.scale*float64(e.img.Bounds().Dy())/2
	radius := e.scale * float64(e.img.Bounds().Dy()) / 2
	distance := math.Sqrt(math.Pow(centerX-p.X, 2) + math.Pow(centerY-p.Y, 2))
	return distance <= radius
}

func (e *Enemy) TakeDamage(amount int) {
	if e.health.Current == 0 {
		return
	}
	if amount >= e.health.Current {
		amount = e.health.Current
	}
	e.health.Current = e.health.Current - amount
	e.healthBar.Subtract(amount)
}

type Hitpoints struct {
	Current int
	Max     int
}
