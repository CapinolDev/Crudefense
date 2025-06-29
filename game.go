package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var (
	playerHp     = 10
	playerSpeed  = 3.0
	currentRound = 1.0
)
var enemies []*Enemy
var bullets []*Bullet

type Gameplay struct {
}

func NewGameplay() *Gameplay {
	enemies = []*Enemy{}
	for len(enemies) < int(math.Round(currentRound*1.5)) {
		enemy := NewEnemy(rand.Float64()*700, rand.Float64()*700, &playerX, &playerY)
		enemies = append(enemies, enemy)
	}
	return &Gameplay{}
}

func (gp *Gameplay) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		playerY = playerY - playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		playerY = playerY + playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		playerX = playerX - playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		playerX = playerX + playerSpeed
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		dx := float64(mx) - playerX
		dy := float64(my) - playerY
		dist := math.Hypot(dx, dy)
		if dist == 0 {
			dist = 1 // avoid division by zero
		}
		velX := dx / dist
		velY := dy / dist

		bullet := NewBullet(playerX, playerY, velX, velY)
		bullets = append(bullets, bullet)
	}

	// Update bullets
	for _, b := range bullets {
		b.Update()
	}

	// Update enemies (your existing)
	for _, enemy := range enemies {
		enemy.Update()

		for _, enemy := range enemies {
			enemy.Update()
		}
	}
}

func (gp *Gameplay) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{136, 174, 224, 1})
	dCHR := &font.Drawer{
		Dst:  screen,
		Src:  image.NewUniform(color.White),
		Face: fontFace,
		Dot:  fixed.P(int(playerX-20), int(playerY)-10),
	}
	dHP := &font.Drawer{
		Dst:  screen,
		Src:  image.NewUniform(color.RGBA{255, 0, 0, 255}),
		Face: fontFace,
		Dot:  fixed.P(10, 20),
	}
	for _, enemy := range enemies {
		enemy.Draw(screen)
	}
	for _, b := range bullets {
		b.Draw(screen)
	}
	dHP.DrawString(fmt.Sprintf("HP: %d", playerHp))
	dCHR.DrawString(userName)
	charOps := &ebiten.DrawImageOptions{}
	charOps.GeoM.Scale(playerScale, playerScale)
	charOps.GeoM.Translate(playerX, playerY)
	screen.DrawImage(mainChar, charOps)

}
