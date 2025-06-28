package main

import (
	"image/color"
	_ "image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const sampleRate = 44100

var (
	cursorX       = 0
	cursorY       = 0
	playBtnX      = 50.0
	playBtnY      = 50.0
	playBtnScale  = 0.7
	archBtnX      = 0.0
	archBtnY      = 0.0
	archBtnScaleX = 0.6
	archBtnScaleY = 0.8
	currentScene  = "Menu"
	crosshair     *ebiten.Image
	playButton    *ebiten.Image
	archerButton  *ebiten.Image
	audioCtx      *audio.Context
	player        *audio.Player
)

func init() {
	var err error
	audioCtx = audio.NewContext(sampleRate)
	crosshair, _, err = ebitenutil.NewImageFromFile("./src/gui/crosshair.png")
	if err != nil {
		log.Fatal(err)
	}
	playButton, _, err = ebitenutil.NewImageFromFile("./src/gui/playButton.png")
	if err != nil {
		log.Fatal(err)
	}
	archerButton, _, err = ebitenutil.NewImageFromFile("./src/gui/archerButton.png")
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Open("./src/audio/clickSound.wav")
	if err != nil {
		log.Fatal(err)
	}

	d, err := wav.Decode(audioCtx, f)
	if err != nil {
		log.Fatal(err)
	}

	player, err = audioCtx.NewPlayer(d)
	if err != nil {
		log.Fatal(err)
	}
}

type Game struct{}

func (g *Game) Update() error {
	cursorX, cursorY = ebiten.CursorPosition()

	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		if currentScene == "Menu" {
			width := float64(playButton.Bounds().Dx()) * playBtnScale
			height := float64(playButton.Bounds().Dy()) * playBtnScale

			if float64(cursorX) >= playBtnX && float64(cursorX) <= playBtnX+width &&
				float64(cursorY) >= playBtnY && float64(cursorY) <= playBtnY+height {
				player.Play()
				currentScene = "CharSelect"
			}

		}
		if currentScene == "CharSelect" {
			width := float64(archerButton.Bounds().Dx()) * playBtnScale
			height := float64(archerButton.Bounds().Dy()) * playBtnScale

			if float64(cursorX) >= archBtnX && float64(cursorX) <= archBtnX+width &&
				float64(cursorY) >= archBtnY && float64(cursorY) <= archBtnY+height {
				player.Play()
				currentScene = "Game"
			}
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		log.Fatal("Game closed by user")
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if currentScene == "Menu" {
		screen.Fill(color.RGBA{119, 123, 165, 1})
		playOp := &ebiten.DrawImageOptions{}
		playOp.GeoM.Scale(playBtnScale, playBtnScale)
		playOp.GeoM.Translate(playBtnX, playBtnY)
		screen.DrawImage(playButton, playOp)
	}
	if currentScene == "CharSelect" {
		screen.Fill(color.RGBA{119, 123, 165, 1})
		archOp := &ebiten.DrawImageOptions{}
		archOp.GeoM.Scale(archBtnScaleX, archBtnScaleY)
		archOp.GeoM.Translate(archBtnX, archBtnY)
		screen.DrawImage(archerButton, archOp)

	}
	crossOp := &ebiten.DrawImageOptions{}

	crossW := crosshair.Bounds().Dx()
	crossH := crosshair.Bounds().Dx()

	crossOp.GeoM.Translate(-float64(crossW)/2, -float64(crossH)/2)
	crossOp.GeoM.Scale(0.03, 0.03)
	crossOp.GeoM.Translate(float64(cursorX), float64(cursorY))
	screen.DrawImage(crosshair, crossOp)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	ebiten.SetFullscreen(true)
	ebiten.SetWindowTitle("Crudefense")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
