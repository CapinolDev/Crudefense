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
	playerHp             = 1
	playerSpeed          = 0.0
	playerRadius         = 00.0
	playerInvincible     = false
	playerInvincibleFor  = 0
	playerInvincibleTick = 0
	currentRound         = 1.0
	timeSinceLastShot    = 0.0
	shootCooldown        = 0.0
	attackType           string
	class                string
)
var enemies []*Enemy
var bullets []*Bullet

type Gameplay struct {
}

func getClassStats() {
	if class == "Archer" {
		playerHp = 10
		playerSpeed = 3.0
		playerRadius = 20.0
		playerInvincible = false
		playerInvincibleFor = 60
		playerInvincibleTick = 0
		timeSinceLastShot = 0.0
		shootCooldown = 0.0
		attackType = "Bullet"

	}
}

func resetValues() {
	playerHp = 1
	playerSpeed = 0.0
	playerRadius = 00.0
	playerInvincible = false
	playerInvincibleFor = 0
	playerInvincibleTick = 0
	currentRound = 1.0
	timeSinceLastShot = 0.0
	shootCooldown = 0.0
	attackType = ""
	class = ""
}
func isColliding(x1, y1, r1, x2, y2, r2 float64) bool {
	dx := x1 - x2
	dy := y1 - y2
	return math.Hypot(dx, dy) < r1+r2
}
func NewGameplay() *Gameplay {
	if currentRound < 1 {
		currentRound = 1
	}
	enemies = []*Enemy{}

	margin := 50.0 // how far off-screen they spawn

	for len(enemies) < int(math.Round(currentRound*1.5)) {
		var x, y float64

		edge := rand.Intn(4) // 0 = top, 1 = right, 2 = bottom, 3 = left

		switch edge {
		case 0: // Top
			x = rand.Float64() * float64(screenWidth)
			y = -margin
		case 1: // Right
			x = float64(screenWidth) + margin
			y = rand.Float64() * float64(screenHeight)
		case 2: // Bottom
			x = rand.Float64() * float64(screenWidth)
			y = float64(screenHeight) + margin
		case 3: // Left
			x = -margin
			y = rand.Float64() * float64(screenHeight)
		}

		enemy := NewEnemy(x, y, &playerX, &playerY)
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
		if attackType == "Bullet" {
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

			timeSinceLastShot = 0

		}
	}

	for _, b := range bullets {
		b.Update()
	}

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
	dFPS := &font.Drawer{
		Dst:  screen,
		Src:  image.NewUniform(color.RGBA{255, 0, 0, 255}),
		Face: fontFace,
		Dot:  fixed.P(10, 60),
	}
	for _, enemy := range enemies {
		enemy.Draw(screen)
	}
	for _, b := range bullets {
		b.Draw(screen)
	}
	if showFps {
		dFPS.DrawString(fmt.Sprintf("Fps: %v", currentFps))
	}

	dHP.DrawString(fmt.Sprintf("HP: %d", playerHp))
	dCHR.DrawString(userName)
	charOps := &ebiten.DrawImageOptions{}
	charOps.GeoM.Scale(playerScale, playerScale)
	charOps.GeoM.Translate(playerX, playerY)
	screen.DrawImage(mainChar, charOps)

}
