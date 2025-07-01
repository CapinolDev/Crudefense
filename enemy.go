package main

import (
	"image/color"
	"math"
	"math/rand"

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

var enemyImage *ebiten.Image

func init() {
	enemyImage = ebiten.NewImage(20, 20)
	enemyImage.Fill(color.RGBA{255, 0, 0, 255})
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

	dx := *e.TargetX - e.X
	dy := *e.TargetY - e.Y
	dist := math.Hypot(dx, dy)
	if dist > 1 {
		offset := rand.Float64()*0.5 - 0.25 // So that the enemies wont be inside of eachother
		e.X += (dx/dist)*e.Speed + offset
		e.Y += (dy/dist)*e.Speed + offset
	}
}

func (e *Enemy) Draw(screen *ebiten.Image) {

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(e.X, e.Y)
	screen.DrawImage(enemyImage, op)
}
