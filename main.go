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
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const sampleRate = 44100

var (
	screenWidth    = 640
	screenHeight   = 480
	cursorX        = 0
	cursorY        = 0
	playBtnX       = 50.0
	playBtnY       = 50.0
	playBtnScale   = 0.7
	archBtnX       = 0.0
	archBtnY       = 0.0
	archBtnScaleX  = 0.6
	archBtnScaleY  = 0.8
	settingsX      = 10.0
	settingsY      = 420.0
	settingsScale  = 0.1
	goBackX        = 0.0
	goBackY        = 0.0
	goBackScale    = 0.3
	inputRunes     []rune
	currentScene   = "Menu"
	userName       string
	userInput      string
	crosshair      *ebiten.Image
	playButton     *ebiten.Image
	archerButton   *ebiten.Image
	settingsButton *ebiten.Image
	goBack         *ebiten.Image
	audioCtx       *audio.Context
	player         *audio.Player
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
	settingsButton, _, err = ebitenutil.NewImageFromFile("./src/gui/cogwheel.png")
	if err != nil {
		log.Fatal(err)
	}
	goBack, _, err = ebitenutil.NewImageFromFile("./src/gui/goBack.png")
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
	inputRunes = inputRunes[:0]
	inputRunes = ebiten.AppendInputChars(inputRunes)

	for _, r := range inputRunes {
		if r >= 0x20 && r != 0x7F { // ignore control characters
			userInput += string(r)
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		if currentScene == "Menu" {
			widthP := float64(playButton.Bounds().Dx()) * playBtnScale
			heightP := float64(playButton.Bounds().Dy()) * playBtnScale
			widthS := float64(settingsButton.Bounds().Dx()) * settingsScale
			heightS := float64(settingsButton.Bounds().Dy()) * settingsScale

			if float64(cursorX) >= playBtnX && float64(cursorX) <= playBtnX+widthP &&
				float64(cursorY) >= playBtnY && float64(cursorY) <= playBtnY+heightP {
				player.Rewind()
				player.Play()
				currentScene = "CharSelect"
			}

			if float64(cursorX) >= settingsX && float64(cursorX) <= settingsX+widthS &&
				float64(cursorY) >= settingsY && float64(cursorY) <= settingsY+heightS {
				player.Rewind()
				player.Play()
				currentScene = "Settings"
			}

		}
		if currentScene == "CharSelect" {
			width := float64(archerButton.Bounds().Dx()) * playBtnScale
			height := float64(archerButton.Bounds().Dy()) * playBtnScale

			if float64(cursorX) >= archBtnX && float64(cursorX) <= archBtnX+width &&
				float64(cursorY) >= archBtnY && float64(cursorY) <= archBtnY+height {
				player.Rewind()
				player.Play()
				currentScene = "Game"
			}
		}
		if currentScene == "Settings" {
			width := float64(goBack.Bounds().Dx()) * goBackScale
			height := float64(goBack.Bounds().Dy()) * goBackScale

			if float64(cursorX) >= goBackX && float64(cursorX) <= goBackX+width &&
				float64(cursorY) >= goBackY && float64(cursorY) <= goBackY+height {
				player.Rewind()
				player.Play()
				currentScene = "Menu"
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
		settOp := &ebiten.DrawImageOptions{}
		settOp.GeoM.Scale(settingsScale, settingsScale)
		settOp.GeoM.Translate(settingsX, settingsY)
		screen.DrawImage(playButton, playOp)
		screen.DrawImage(settingsButton, settOp)
	}
	if currentScene == "CharSelect" {
		screen.Fill(color.RGBA{119, 123, 165, 1})
		archOp := &ebiten.DrawImageOptions{}
		archOp.GeoM.Scale(archBtnScaleX, archBtnScaleY)
		archOp.GeoM.Translate(archBtnX, archBtnY)
		screen.DrawImage(archerButton, archOp)

	}
	if currentScene == "Settings" {
		screen.Fill(color.RGBA{119, 123, 165, 1})
		goBackOp := &ebiten.DrawImageOptions{}
		goBackOp.GeoM.Scale(goBackScale, goBackScale)
		goBackOp.GeoM.Translate(goBackX, goBackY)
		screen.DrawImage(goBack, goBackOp)
		text.Draw(screen, userInput)
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			userName = userInput
			userInput = ""
		}

	}
	crossOp := &ebiten.DrawImageOptions{}

	crossW := crosshair.Bounds().Dx()
	crossH := crosshair.Bounds().Dy()

	crossOp.GeoM.Translate(-float64(crossW)/2, -float64(crossH)/2)
	crossOp.GeoM.Scale(0.03, 0.03)
	crossOp.GeoM.Translate(float64(cursorX), float64(cursorY))
	screen.DrawImage(crosshair, crossOp)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	ebiten.SetFullscreen(true)
	ebiten.SetWindowTitle("Crudefense")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
