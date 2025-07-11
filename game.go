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
	MaxHp                = 10
	playerSpeed          = 3.0
	playerRadius         = 20.0
	playerInvincible     = false
	playerInvincibleFor  = 30 // ticks of invincibility
	playerInvincibleTick = 0
	currentRound         = 1.0
	timeSinceLastShot    = 0.0
	shootCooldown        = 0.7
	upgradeScreen        *UpgradeScreen
	waveInProgress       = false
	bulletAmount         = 1
	vampireLvl           = 0
)
var enemies []*Enemy
var bullets []*Bullet

type Gameplay struct {
}

func resetValues() {
	playerHp = 10
	MaxHp = 10
	playerX = 320.0
	playerY = 240.0
	playerSpeed = 3.0
	playerRadius = 20.0
	playerInvincible = false
	playerInvincibleTick = 0
	bulletAmount = 1
	bulletSpeed = 3.0
	vampireLvl = 0
	bouncesAmount = 0 // Reset bounces amount

	timeSinceLastShot = 0.0
	shootCooldown = 0.7
	// Reset enemies and bullets
	enemies = []*Enemy{}
	bullets = []*Bullet{}
	//reset round and wave state
	currentRound = 1.0
	waveInProgress = false
}
func vampHeal() {
	if vampireLvl > 0 {
		playerHp += vampireLvl
		if playerHp > MaxHp {
			playerHp = MaxHp // cap hp at MaxHp
		}
	}
}
func randomOffscreenPosition() (x, y float64) {
	screenW := float64(screenWidth)
	screenH := float64(screenHeight)

	edge := rand.Intn(4) // 0 = top, 1 = bottom, 2 = left, 3 = right

	switch edge {
	case 0: // top
		x = rand.Float64() * screenW
		y = -30
	case 1: // bottom
		x = rand.Float64() * screenW
		y = screenH + 30
	case 2: // left
		x = -30
		y = rand.Float64() * screenH
	case 3: // right
		x = screenW + 30
		y = rand.Float64() * screenH
	}
	return
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
	waveInProgress = true

	spawnEnemies := func(count, level int) {
		for i := 0; i < count; i++ {
			ex, ey := randomOffscreenPosition()
			enemy := NewEnemy(ex, ey, level, &playerX, &playerY)
			enemies = append(enemies, enemy)
		}
	}

	switch currentRound {
	case 1:
		spawnEnemies(2, 1)
	case 2:
		spawnEnemies(4, 2)
	case 3:
		spawnEnemies(10, 1)
	case 4:
		spawnEnemies(8, 2)
	case 5:
		spawnEnemies(2, 4)
	case 6:
		spawnEnemies(6, 3)
	default:
		// if i didnt specify a round, spawn enemies based on the current round
		spawnEnemies(int(currentRound*2), int(currentRound/2))
	}

	return &Gameplay{}
}

func (gp *Gameplay) Update() {
	if playerHp <= 0 {
		currentScene = "GameOver"
	}
	if waveInProgress && len(enemies) == 0 {
		currentScene = "Upgrade"
		currentRound++
		waveInProgress = false // stop triggering Upgrade every frame
	}
	if playerInvincible {
		playerInvincibleTick++
		if playerInvincibleTick >= playerInvincibleFor {
			playerInvincible = false
			playerInvincibleTick = 0
		}
	}
	var dx, dy float64

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		dy -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		dy += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		dx -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		dx += 1
	}

	// normalize if moving diagonally
	if dx != 0 || dy != 0 {
		length := math.Hypot(dx, dy)
		dx /= length
		dy /= length
	}

	playerX += dx * playerSpeed
	playerY += dy * playerSpeed

	screenW := float64(screenWidth)
	screenH := float64(screenHeight)
	playerW := float64(mainChar.Bounds().Dx()) * playerScale
	playerH := float64(mainChar.Bounds().Dy()) * playerScale

	if playerX < 0 {
		playerX = 0
	}
	if playerY < 0 {
		playerY = 0
	}
	if playerX+playerW > screenW {
		playerX = screenW - playerW
	}
	if playerY+playerH > screenH {
		playerY = screenH - playerH
	}
	for i := 0; i < len(bullets); i++ {
		b := bullets[i]
		for j := 0; j < len(enemies); j++ {
			e := enemies[j]
			if isColliding(b.X, b.Y, b.Radius, e.X, e.Y, e.Radius) {
				bullets = append(bullets[:i], bullets[i+1:]...)
				i-- // Adjust index after deletion

				e.HP--
				if e.HP <= 0 {
					enemies = append(enemies[:j], enemies[j+1:]...) // Remove enemy
					j--                                             // Adjust index after deletion if needed
				}
				vampHeal() // Heal if vampire mutation is active
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
		for i := 0; i < bulletAmount; i++ {
			// Spread bullets slightly
			spread := float64(i)*0.1 - float64(bulletAmount-1)*0.05
			bullet := NewBullet(playerX, playerY, velX+spread, velY+spread)
			bullets = append(bullets, bullet)

		}

		timeSinceLastShot = 0

	}
	filtered := bullets[:0]
	for _, b := range bullets {
		b.Update(screenWidth, screenHeight)
		if b.BouncesLeft >= 0 {
			filtered = append(filtered, b)
		}
	}
	bullets = filtered
	for _, b := range bullets {
		b.Update(640, 480)
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

}
