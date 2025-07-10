package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Bullet struct {
	X, Y       float64
	VelX, VelY float64
	Speed      float64
	Radius     float64
}

var bulletImage *ebiten.Image

func init() {
	bulletImage = ebiten.NewImage(5, 5)
	bulletImage.Fill(color.White)
}
func (b *Bullet) IsOffscreen(screenWidth, screenHeight int) bool {
	return b.X < -10 || b.Y < -10 || b.X > float64(screenWidth)+10 || b.Y > float64(screenHeight)+10
}

func NewBullet(x, y, velX, velY float64) *Bullet {
	return &Bullet{
		X:      x,
		Y:      y,
		VelX:   velX,
		VelY:   velY,
		Speed:  5,
		Radius: 2.5,
	}
}

func (b *Bullet) Update() {
	b.X += b.VelX * b.Speed
	b.Y += b.VelY * b.Speed
	filtered := bullets[:0]
	for _, b := range bullets {
		if !b.IsOffscreen(screenWidth, screenHeight) {
			filtered = append(filtered, b)
		}
	}
	bullets = filtered

}

func (b *Bullet) Draw(screen *ebiten.Image) {

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.X, b.Y)
	screen.DrawImage(bulletImage, op)

}
