package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Bullet struct {
	X, Y       float64
	VelX, VelY float64
	Speed      float64
}

func NewBullet(x, y, velX, velY float64) *Bullet {
	return &Bullet{
		X:     x,
		Y:     y,
		VelX:  velX,
		VelY:  velY,
		Speed: 5, // adjust speed
	}
}

func (b *Bullet) Update() {
	b.X += b.VelX * b.Speed
	b.Y += b.VelY * b.Speed
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	const size = 5
	bulletImage := ebiten.NewImage(size, size)
	bulletImage.Fill(color.White) // or any color you want

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.X, b.Y)
	screen.DrawImage(bulletImage, op)

}
