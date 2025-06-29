package main

import (
	"encoding/json"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const sampleRate = 44100

type Settings struct {
	Username   string         `json:"username"`
	Fullscreen bool           `json:"fullscreen"`
	UserStats  map[string]int `json:"user_stats"`
}

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
	fscreenX       = 240.0
	fscreenY       = 165.0
	fscreenScale   = 0.26
	inputRunes     []rune
	currentScene   = "Menu"
	userName       string
	userInput      string
	crosshair      *ebiten.Image
	playButton     *ebiten.Image
	archerButton   *ebiten.Image
	settingsButton *ebiten.Image
	goBack         *ebiten.Image
	fscreen        *ebiten.Image
	audioCtx       *audio.Context
	player         *audio.Player
	fontFace       font.Face
	settings       Settings
)

func loadFont() font.Face {
	ttfBytes, err := os.ReadFile("./src/fonts/Queensides-3z7Ey.ttf")
	if err != nil {
		log.Fatal(err)
	}

	ttfFont, err := opentype.Parse(ttfBytes)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	const fontSize = 24

	fontFace, err := opentype.NewFace(ttfFont, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	return fontFace
}

func LoadSettings(filename string) Settings {
	var settings Settings

	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, return default settings
			return Settings{
				Username:   "",
				Fullscreen: true,
				UserStats:  make(map[string]int),
			}
		} else {
			log.Fatal(err)
		}
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&settings)
	if err != nil {
		log.Fatal("Failed to decode settings:", err)
	}

	return settings
}
func SaveSettings(filename string, settings Settings) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print JSON
	return encoder.Encode(settings)
}

func init() {

	var err error
	settings = LoadSettings("settings.json")

	userName = settings.Username

	if settings.Fullscreen {
		ebiten.SetFullscreen(true)
	} else {
		ebiten.SetFullscreen(false)
	}
	fontFace = loadFont()
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
	fscreen, _, err = ebitenutil.NewImageFromFile("./src/gui/fscreen.png")
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
	inputRunes = inputRunes[:0] // ctrl chars
	inputRunes = ebiten.AppendInputChars(inputRunes)

	for _, r := range inputRunes {
		if r >= 0x20 && r != 0x7F {
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
			widthG := float64(goBack.Bounds().Dx()) * goBackScale
			heightG := float64(goBack.Bounds().Dy()) * goBackScale
			widthF := float64(goBack.Bounds().Dx()) * goBackScale
			heightF := float64(goBack.Bounds().Dy()) * goBackScale

			if float64(cursorX) >= goBackX && float64(cursorX) <= goBackX+widthG &&
				float64(cursorY) >= goBackY && float64(cursorY) <= goBackY+heightG {
				player.Rewind()
				player.Play()
				currentScene = "Menu"
			}
			if float64(cursorX) >= fscreenX && float64(cursorX) <= fscreenX+widthF &&
				float64(cursorY) >= fscreenY && float64(cursorY) <= fscreenY+heightF {
				player.Rewind()
				player.Play()
				settings.Fullscreen = !settings.Fullscreen
				ebiten.SetFullscreen(settings.Fullscreen)
			}
		}

	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		err := SaveSettings("settings.json", settings)
		if err != nil {
			log.Println("Settings failure:", err)

		}
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
		fScreenOp := &ebiten.DrawImageOptions{}
		fScreenOp.GeoM.Scale(fscreenScale, fscreenScale)
		fScreenOp.GeoM.Translate(fscreenX, fscreenY)
		fScreenOp.ColorScale.Reset()
		if settings.Fullscreen {

			fScreenOp.ColorScale.Scale(0, 1, 0, 1)
		} else {

			fScreenOp.ColorScale.Scale(1, 0, 0, 1)
		}

		screen.DrawImage(goBack, goBackOp)
		screen.DrawImage(fscreen, fScreenOp)
		//user input
		dUI := &font.Drawer{
			Dst:  screen,
			Src:  image.NewUniform(color.White),
			Face: fontFace,
			Dot:  fixed.P(80, 80),
		}
		//"Input username" text
		dIUN := &font.Drawer{
			Dst:  screen,
			Src:  image.NewUniform(color.White),
			Face: fontFace,
			Dot:  fixed.P(80, 60),
		}
		//"Current username" text
		dUN := &font.Drawer{
			Dst:  screen,
			Src:  image.NewUniform(color.White),
			Face: fontFace,
			Dot:  fixed.P(80, 130),
		}
		// username
		dUNV := &font.Drawer{
			Dst:  screen,
			Src:  image.NewUniform(color.White),
			Face: fontFace,
			Dot:  fixed.P(80, 150),
		}
		//fullscreen
		dFS := &font.Drawer{
			Dst:  screen,
			Src:  image.NewUniform(color.White),
			Face: fontFace,
			Dot:  fixed.P(80, 180),
		}
		dUI.DrawString(userInput)
		dIUN.DrawString("Input username:")
		dUN.DrawString("Current username:")
		dUNV.DrawString(userName)
		dFS.DrawString("Fullscreen:")
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			userName = userInput
			settings.Username = userInput
			userInput = ""
			err := SaveSettings("settings.json", settings)
			if err != nil {
				log.Println("Settings failure:", err)
			}
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
	ebiten.SetWindowTitle("Crudefense")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
