package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type UpgradeScreen struct {
	selected int // example: index of selected upgrade
	options  []string
}

var (
	imageWhite  = image.NewUniform(color.White)
	imageYellow = image.NewUniform(color.RGBA{255, 255, 0, 255})
)

func NewUpgradeScreen() *UpgradeScreen {
	return &UpgradeScreen{
		selected: 0,
		options: []string{
			"Increase HP",
			"Increase Speed",
			"Faster Shooting",
		},
	}
}

func (u *UpgradeScreen) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		u.applyUpgrade(u.selected)

		currentScene = "Game"
		NewGameplay()

	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) || inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		u.selected++
		if u.selected >= len(u.options) {
			u.selected = 0
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyW) || inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		u.selected--
		if u.selected < 0 {
			u.selected = len(u.options) - 1
		}
	}

}

func (u *UpgradeScreen) applyUpgrade(index int) {
	switch index {
	case 0:
		playerHp += 2
	case 1:
		playerSpeed += 0.5
	case 2:
		if shootCooldown > 0.2 {
			shootCooldown -= 0.1
		}
	}
}

func (u *UpgradeScreen) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{50, 50, 50, 255})

	// Draw title
	title := &font.Drawer{
		Dst:  screen,
		Src:  imageWhite,
		Face: fontFace,
		Dot:  fixed.P(100, 50),
	}
	title.DrawString("Choose an Upgrade:")

	// Draw options
	for i, opt := range u.options {
		d := &font.Drawer{
			Dst:  screen,
			Src:  imageWhite,
			Face: fontFace,
			Dot:  fixed.P(120, 100+i*30),
		}
		if i == u.selected {
			d.Src = imageYellow // Highlight selected option
		}
		d.DrawString(opt)
	}
}
