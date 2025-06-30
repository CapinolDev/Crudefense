package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Enemy struct {
	X, Y    float64
	Speed   float64
	Color   color.Color
	TargetX *float64
	TargetY *float64
	Radius  float64
}

func NewEnemy(x, y float64, targetX, targetY *float64) *Enemy {
	return &Enemy{
		X:       x,
		Y:       y,
		Speed:   1.5,
		Radius:  10,
		Color:   color.RGBA{255, 0, 0, 255},
		TargetX: targetX,
		TargetY: targetY,
	}
}

func (e *Enemy) Update() {
	// Move towards player (target)
	dx := *e.TargetX - e.X
	dy := *e.TargetY - e.Y
	dist := math.Hypot(dx, dy)
	if dist > 1 {
		e.X += (dx / dist) * e.Speed
		e.Y += (dy / dist) * e.Speed
	}
}

func (e *Enemy) Draw(screen *ebiten.Image) {
	const size = 20
	enemyImage := ebiten.NewImage(size, size)
	enemyImage.Fill(e.Color)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(e.X, e.Y)
	screen.DrawImage(enemyImage, op)
}
