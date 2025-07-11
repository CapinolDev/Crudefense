package main

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var (
	mutationChance     = 0.15
	hasRolledMutation  = false
	successfulMutation = false
	mutationType       = 0
)

type UpgradeScreen struct {
	selected int // index of selected upgrade
	options  []string
}

func rollForMutation() bool {
	return rand.Float64() < mutationChance
}
func rollWhichMutation() int {
	return rand.Intn(4)
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
			"Mutation ",
		},
	}
}

func (u *UpgradeScreen) Update() {
	if !hasRolledMutation {
		if rollForMutation() {
			successfulMutation = true
			hasRolledMutation = true
			mutationType = rollWhichMutation()
		} else {
			successfulMutation = false
			hasRolledMutation = true
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		u.applyUpgrade(u.selected)

		currentScene = "Game"
		NewGameplay()
		hasRolledMutation = false // Reset mutation roll for next upgrade screen

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
	//Mutations
	case 3:
		{
			switch mutationType {
			case 0:
				bulletSpeed += 1.5
			case 1:
				bulletAmount += 1
			case 2:
				vampireLvl++
			case 3:
				bouncesAmount++
			}
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
	dMutTXT := &font.Drawer{
		Dst:  screen,
		Src:  imageWhite,
		Face: fontFace,
		Dot:  fixed.P(110, 240),
	}
	dMut := &font.Drawer{
		Dst:  screen,
		Src:  imageWhite,
		Face: fontFace,
		Dot:  fixed.P(110, 270),
	}
	title.DrawString("Choose an Upgrade:")
	dMutTXT.DrawString("Mutations " + fmt.Sprintf("(%.0f%%", mutationChance*100) + " chance) :")
	if successfulMutation {
		switch mutationType {
		case 0:
			dMut.DrawString("Sniper - Increases bullet speed")
		case 1:
			dMut.DrawString("Shotgun - Increases bullet amount")
		case 2:
			dMut.DrawString("Vampire - Heals for each enemy killed")
		case 3:
			dMut.DrawString("Ricochet - Bullets bounce off walls (+1 bounce)")
		}
	} else {
		dMut.DrawString("No Mutation")
	}

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
