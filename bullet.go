package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Bullet struct {
	X, Y        float64
	VelX, VelY  float64
	Speed       float64
	Radius      float64
	BouncesLeft int
}

var (
	bulletSpeed   = 3.0
	bulletImage   *ebiten.Image
	bouncesAmount = 0 // Number of bounces before bullet is removed
)

func init() {
	bulletImage = ebiten.NewImage(5, 5)
	bulletImage.Fill(color.White)
}
func (b *Bullet) IsOffscreen(screenWidth, screenHeight int) bool {
	return b.X < -10 || b.Y < -10 || b.X > float64(screenWidth)+10 || b.Y > float64(screenHeight)+10
}

func NewBullet(x, y, velX, velY float64) *Bullet {
	return &Bullet{
		X:           x,
		Y:           y,
		VelX:        velX,
		VelY:        velY,
		Speed:       bulletSpeed,
		Radius:      2.5,
		BouncesLeft: bouncesAmount,
	}
}

func (b *Bullet) Update(screenWidth, screenHeight int) {
	b.X += b.VelX * b.Speed
	b.Y += b.VelY * b.Speed
	b.Bounce(screenWidth, screenHeight)
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.X-b.Radius, b.Y-b.Radius) // center properly
	screen.DrawImage(bulletImage, op)
}
func (b *Bullet) Bounce(screenWidth, screenHeight int) {
	bounced := false

	if b.X-b.Radius <= 0 || b.X+b.Radius >= float64(screenWidth) {
		b.VelX = -b.VelX
		bounced = true
		if b.X-b.Radius < 0 {
			b.X = b.Radius
		} else if b.X+b.Radius > float64(screenWidth) {
			b.X = float64(screenWidth) - b.Radius
		}
	}

	if b.Y-b.Radius <= 0 || b.Y+b.Radius >= float64(screenHeight) {
		b.VelY = -b.VelY
		bounced = true
		if b.Y-b.Radius < 0 {
			b.Y = b.Radius
		} else if b.Y+b.Radius > float64(screenHeight) {
			b.Y = float64(screenHeight) - b.Radius
		}
	}

	if bounced {
		b.BouncesLeft--
	}
}
