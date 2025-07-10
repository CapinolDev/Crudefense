package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var enemyImage *ebiten.Image

func initEnemies() {
	const size = 20
	enemyImage = ebiten.NewImage(size, size)
	enemyImage.Fill(color.RGBA{255, 0, 0, 255})
}

type Enemy struct {
	X, Y    float64
	Speed   float64
	Color   color.Color
	TargetX *float64
	TargetY *float64
	Radius  float64
	HP      int
	MaxHP   int
}

func NewEnemy(x, y float64, hp int, targetX, targetY *float64) *Enemy {
	return &Enemy{
		X:       x,
		Y:       y,
		Speed:   1.5,
		Radius:  10,
		Color:   color.RGBA{255, 0, 0, 255},
		TargetX: targetX,
		TargetY: targetY,
		HP:      hp,
		MaxHP:   hp,
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
	var clr color.RGBA

	switch e.MaxHP {
	case 1:
		clr = color.RGBA{0, 255, 0, 255} // green
	case 2:
		clr = color.RGBA{0, 0, 255, 255} // blue
	case 3:
		clr = color.RGBA{255, 165, 0, 255} // orange
	case 4:
		clr = color.RGBA{255, 0, 0, 255} // red
	default:
		clr = color.RGBA{255, 255, 255, 255} // white for unknown
	}

	img := ebiten.NewImage(int(e.Radius*2), int(e.Radius*2))
	img.Fill(clr)

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(e.X-e.Radius, e.Y-e.Radius)

	screen.DrawImage(img, opts)
	hpText := fmt.Sprintf("%d", e.HP)
	d := &font.Drawer{
		Dst:  screen,
		Src:  image.NewUniform(color.White),
		Face: fontFace,
		// Position the text roughly centered horizontally above the enemy
		Dot: fixed.P(int(e.X)-int(e.Radius/2), int(e.Y)-int(e.Radius)-5),
	}
	d.DrawString(hpText)
}
