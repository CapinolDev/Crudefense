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

var (
	playerHp             = 10
	playerSpeed          = 3.0
	playerRadius         = 20.0
	playerInvincible     = false
	playerInvincibleFor  = 60
	playerInvincibleTick = 0
	currentRound         = 1.0
	timeSinceLastShot    = 0.0
	shootCooldown        = 0.7
)
var enemies []*Enemy
var bullets []*Bullet

type Gameplay struct {
}

func resetValues() {
	playerHp = 10
	playerSpeed = 3.0
	playerRadius = 20.0
	playerInvincible = false
	playerInvincibleFor = 60
	playerInvincibleTick = 0
	currentRound = 1.0
	timeSinceLastShot = 0.0
	shootCooldown = 0.7
}
func isColliding(x1, y1, r1, x2, y2, r2 float64) bool {
	dx := x1 - x2
	dy := y1 - y2
	return math.Hypot(dx, dy) < r1+r2
}
func NewGameplay() *Gameplay {
	enemies = []*Enemy{}
	for len(enemies) < int(math.Round(currentRound*1.5)) {
		enemy := NewEnemy((rand.Float64()+0.2)*700, (rand.Float64()+0.2)*700, &playerX, &playerY)
		enemies = append(enemies, enemy)
	}
	return &Gameplay{}
}

func (gp *Gameplay) Update() {
	if playerHp <= 0 {
		currentScene = "GameOver"
	}
	if playerInvincible {
		playerInvincibleTick++
		if playerInvincibleTick >= playerInvincibleFor {
			playerInvincible = false
			playerInvincibleTick = 0
		}
	}
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
	for i := 0; i < len(bullets); i++ {
		b := bullets[i]
		for j := 0; j < len(enemies); j++ {
			e := enemies[j]
			if isColliding(b.X, b.Y, b.Radius, e.X, e.Y, e.Radius) {
				// Remove bullet and enemy
				bullets = append(bullets[:i], bullets[i+1:]...)
				i-- // Adjust index after deletion
				enemies = append(enemies[:j], enemies[j+1:]...)
				break
			}
		}
	}
	for _, enemy := range enemies {
		if isColliding(enemy.X, enemy.Y, enemy.Radius, playerX, playerY, playerRadius) {
			if !playerInvincible {
				playerHp--
				playerInvincible = true
				playerInvincibleTick = 0
			}
		}
	}
	timeSinceLastShot += 1.0 / ebiten.ActualTPS()

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && timeSinceLastShot >= shootCooldown {
		mx, my := ebiten.CursorPosition()
		dx := float64(mx) - playerX
		dy := float64(my) - playerY
		dist := math.Hypot(dx, dy)
		if dist == 0 {
			dist = 1
		}
		velX := dx / dist
		velY := dy / dist

		bullet := NewBullet(playerX, playerY, velX, velY)
		bullets = append(bullets, bullet)

		timeSinceLastShot = 0 // reset timer
	}

	// Update bullets
	for _, b := range bullets {
		b.Update()
	}

	// Update enemies (your existing)
	for _, enemy := range enemies {
		enemy.Update()

	}
}

func (gp *Gameplay) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{136, 174, 224, 1})
	if !playerInvincible || (playerInvincible && playerInvincibleTick%10 < 5) {
		charOps := &ebiten.DrawImageOptions{}
		charOps.GeoM.Scale(playerScale, playerScale)
		charOps.GeoM.Translate(playerX, playerY)
		screen.DrawImage(mainChar, charOps)
	}
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
