package game

// // const (
// // 	EnemyHitPointsDefault = 50
// // )

// // func NewEnemy(path string) (*Enemy, error) {
// // 	img, err := util.TransparentImageFromPath(path)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	enemy := &Enemy{
// // 		Image: img,
// // 		Scale: 0.7,
// // 		Health: Hitpoints{
// // 			Current: EnemyHitPointsDefault,
// // 			Max:     EnemyHitPointsDefault,
// // 		},
// // 	}
// // 	enemy.NewPosition()
// // 	return enemy, nil
// // }

// // type Enemy struct {
// // 	Health   Hitpoints
// // 	Image    *ebiten.Image
// // 	Position Position
// // 	Scale    float64
// // }

// // type Hitpoints struct {
// // 	Current int
// // 	Max     int
// // }

// // func (e *Enemy) NewPosition() {
// // 	midpoint := Width / 2.0
// // 	min := int(EnemySpawnShiftPercent / 100.0 * midpoint)
// // 	max := int(((EnemySpawnShiftPercent / 100.0) + 1.0) * midpoint)
// // 	e.Position = Position{
// // 		X: float64(util.RandomRangeInt(min, max)),
// // 		Y: 0.8*Height - e.Scale*float64(e.Image.Bounds().Dy()),
// // 	}
// // }

// // func (e *Enemy) InHitbox(p Position) bool {
// // 	centerX := e.Position.X
// // 	centerY := e.Position.Y + e.Scale*float64(e.Image.Bounds().Dy())/2
// // 	radius := e.Scale * float64(e.Image.Bounds().Dy()) / 2
// // 	distance := math.Sqrt(math.Pow(centerX-p.X, 2) + math.Pow(centerY-p.Y, 2))
// // 	return distance <= radius
// // }

// // func (e *Enemy) TakeDamage(amount int) {
// // 	if e.Health.Current == 0 {
// // 		return
// // 	}
// // 	if amount >= e.Health.Current {
// // 		amount = e.Health.Current
// // 	}
// // 	e.Health.Current = e.Health.Current - amount
// // }
