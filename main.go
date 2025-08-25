package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed settings.json
var settingsData []byte

const sampleRate = 44100

type Settings struct {
	Username   string         `json:"username"`
	Fullscreen bool           `json:"fullscreen"`
	UserStats  map[string]int `json:"user_stats"`
	ShowFps    bool           `json:"show_fps"`
}

var gameplay *Gameplay
var (
	screenWidth     = 640
	screenHeight    = 480
	cursorX         = 0
	cursorY         = 0
	playBtnX        = 50.0
	playBtnY        = 50.0
	playBtnScale    = 0.7
	archBtnX        = 0.0
	archBtnY        = 0.0
	settingsX       = 10.0
	settingsY       = 420.0
	settingsScale   = 0.1
	goBackX         = 0.0
	goBackY         = 0.0
	goBackScale     = 0.3
	fscreenX        = 120.0
	fscreenY        = 165.0
	fscreenScale    = 0.26
	playerX         = 320.0
	playerY         = 240.0
	playerScale     = 0.4
	menuButtonX     = 220.0
	menuButtonY     = 180.0
	menuButtonScale = 0.5

	showFps        bool
	currentFps     = 0.0
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
	mainChar       *ebiten.Image
	menuButton     *ebiten.Image
	background     *ebiten.Image
	audioCtx       *audio.Context
	player         *audio.Player
	fontFace       font.Face
	settings       Settings
)

func boolToOnOff(b bool) string {
	if b {
		return "On"
	}
	return "Off"
}
func getLogicalCursorPosition() (int, int) {
	return ebiten.CursorPosition()
}

func updateMenu() {
	widthP := float64(playButton.Bounds().Dx()) * playBtnScale
	heightP := float64(playButton.Bounds().Dy()) * playBtnScale
	widthS := float64(settingsButton.Bounds().Dx()) * settingsScale
	heightS := float64(settingsButton.Bounds().Dy()) * settingsScale

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		if float64(cursorX) >= playBtnX && float64(cursorX) <= playBtnX+widthP &&
			float64(cursorY) >= playBtnY && float64(cursorY) <= playBtnY+heightP {
			currentScene = "CharSelect"
			player.Rewind()
			player.Play()
		}

		if float64(cursorX) >= settingsX && float64(cursorX) <= settingsX+widthS &&
			float64(cursorY) >= settingsY && float64(cursorY) <= settingsY+heightS {
			currentScene = "Settings"
			player.Rewind()
			player.Play()
		}
	}
}

func updateCharSelect() {
	width := float64(archerButton.Bounds().Dx()) * playBtnScale
	height := float64(archerButton.Bounds().Dy()) * playBtnScale

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		if float64(cursorX) >= archBtnX && float64(cursorX) <= archBtnX+width &&
			float64(cursorY) >= archBtnY && float64(cursorY) <= archBtnY+height {

			resetValues()
			currentScene = "Game"
			NewGameplay()
			player.Rewind()
			player.Play()
		}
	}
}

func updateSettings() {
	widthG := float64(goBack.Bounds().Dx()) * goBackScale
	heightG := float64(goBack.Bounds().Dy()) * goBackScale
	widthF := float64(fscreen.Bounds().Dx()) * fscreenScale
	heightF := float64(fscreen.Bounds().Dy()) * fscreenScale

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		if float64(cursorX) >= goBackX && float64(cursorX) <= goBackX+widthG &&
			float64(cursorY) >= goBackY && float64(cursorY) <= goBackY+heightG {
			player.Rewind()
			player.Play()
			currentScene = "Menu"
		}
		if float64(cursorX) >= fscreenX && float64(cursorX) <= fscreenX+widthF &&
			float64(cursorY) >= fscreenY && float64(cursorY) <= fscreenY+heightF {
			settings.Fullscreen = !settings.Fullscreen
			ebiten.SetFullscreen(settings.Fullscreen)
			player.Rewind()
			player.Play()
		}

	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		userName = userInput
		settings.Username = userInput
		userInput = ""
		if err := SaveSettings("settings.json", settings); err != nil {
			log.Println("Settings failure:", err)
		}
	}
}

func updateGameOver() {
	widthM := float64(menuButton.Bounds().Dx()) * menuButtonScale
	heightM := float64(menuButton.Bounds().Dy()) * menuButtonScale

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		if float64(cursorX) >= menuButtonX && float64(cursorX) <= menuButtonX+widthM &&
			float64(cursorY) >= menuButtonY && float64(cursorY) <= menuButtonY+heightM {
			currentScene = "Menu"
			resetValues()

		}
		player.Rewind()
		player.Play()
	}
}
func drawMenu(screen *ebiten.Image) {
	screen.DrawImage(background, nil)

	opPlay := &ebiten.DrawImageOptions{}
	opPlay.GeoM.Scale(playBtnScale, playBtnScale)
	opPlay.GeoM.Translate(playBtnX, playBtnY)
	screen.DrawImage(playButton, opPlay)

	opSettings := &ebiten.DrawImageOptions{}
	opSettings.GeoM.Scale(settingsScale, settingsScale)
	opSettings.GeoM.Translate(settingsX, settingsY)
	screen.DrawImage(settingsButton, opSettings)
}
func drawCharSelect(screen *ebiten.Image) {
	screen.DrawImage(background, nil)

	opArcher := &ebiten.DrawImageOptions{}
	opArcher.GeoM.Scale(playBtnScale, playBtnScale)
	opArcher.GeoM.Translate(archBtnX, archBtnY)
	screen.DrawImage(archerButton, opArcher)
}
func drawSettings(screen *ebiten.Image) {
	screen.DrawImage(background, nil)

	// Draw Go Back button
	opGoBack := &ebiten.DrawImageOptions{}
	opGoBack.GeoM.Scale(goBackScale, goBackScale)
	opGoBack.GeoM.Translate(goBackX, goBackY)
	screen.DrawImage(goBack, opGoBack)

	// Draw Fullscreen toggle button
	opFScreen := &ebiten.DrawImageOptions{}
	opFScreen.GeoM.Scale(fscreenScale, fscreenScale)
	opFScreen.GeoM.Translate(fscreenX, fscreenY)
	screen.DrawImage(fscreen, opFScreen)

	// Draw Username input box or current name
	ebitenutil.DebugPrintAt(screen, "Enter Name: "+userInput, 20, 60)
	ebitenutil.DebugPrintAt(screen, "Current name: "+settings.Username, 20, 80)

	// Draw Fullscreen label and status
	fullscreenStatus := boolToOnOff(settings.Fullscreen)
	ebitenutil.DebugPrintAt(screen, "Fullscreen: "+fullscreenStatus, int(fscreenX)-100, int(fscreenY))

}
func drawGameOver(screen *ebiten.Image) {
	screen.DrawImage(background, nil)

	opMenu := &ebiten.DrawImageOptions{}
	opMenu.GeoM.Scale(menuButtonScale, menuButtonScale)
	opMenu.GeoM.Translate(menuButtonX, menuButtonY)
	screen.DrawImage(menuButton, opMenu)
}

func loadFont() font.Face {
	ttfFont, err := opentype.Parse(fontBytes)
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

func LoadSettings() Settings {
	var s Settings
	err := json.Unmarshal(settingsData, &s)
	if err != nil {
		log.Println("Failed to load embedded settings:", err)
		// fallback defaults
		return Settings{
			Username:   "",
			Fullscreen: true,
			UserStats:  make(map[string]int),
			ShowFps:    false,
		}
	}
	return s
}

func init() {
	initEnemies()
	var err error
	settings = LoadSettings()

	userName = settings.Username
	showFps = settings.ShowFps

	if settings.Fullscreen {
		ebiten.SetFullscreen(true)
	} else {
		ebiten.SetFullscreen(false)
	}
	fontFace = loadFont()
	audioCtx = audio.NewContext(sampleRate)
	crosshair, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(crosshairPNG))
	playButton, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(playButtonPNG))
	archerButton, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(archerButtonPNG))
	settingsButton, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(settingsButtonPNG))
	goBack, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(goBackPNG))
	fscreen, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(fscreenPNG))
	mainChar, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(mainCharPNG))
	menuButton, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(menuButtonPNG))
	background, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(backgroundPNG))

	if err != nil {
		log.Fatal(err)
	}

	d, err := wav.Decode(audioCtx, bytes.NewReader(clickWav))
	if err != nil {
		log.Fatal(err)
	}

	player, err = audioCtx.NewPlayer(d)
	if err != nil {
		log.Fatal(err)
	}
	gameplay = NewGameplay()
	upgradeScreen = NewUpgradeScreen()
}

type Game struct{}

func (g *Game) Update() error {
	mx, my := getLogicalCursorPosition()
	cursorX = mx
	cursorY = my
	inputRunes = inputRunes[:0]
	inputRunes = ebiten.AppendInputChars(inputRunes)

	for _, r := range inputRunes {
		if r >= 0x20 && r != 0x7F {
			userInput += string(r)
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(userInput) > 0 {
		userInput = userInput[:len(userInput)-1]
	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		if err := SaveSettings("settings.json", settings); err != nil {
			log.Println("Settings failure:", err)
		}
		player.Close()

		log.Fatal("Game closed by user")
	}

	switch currentScene {
	case "Menu":
		updateMenu()
	case "CharSelect":
		updateCharSelect()
	case "Settings":
		updateSettings()
	case "Game":
		gameplay.Update()
	case "Upgrade":
		upgradeScreen.Update()
	case "GameOver":
		updateGameOver()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch currentScene {
	case "Menu":
		drawMenu(screen)
	case "CharSelect":
		drawCharSelect(screen)
	case "Settings":
		drawSettings(screen)
	case "Game":
		gameplay.Draw(screen)
	case "Upgrade":
		upgradeScreen.Draw(screen)
	case "GameOver":
		drawGameOver(screen)
	}

	if crosshair != nil {
		op := &ebiten.DrawImageOptions{}
		scale := 0.02

		op.GeoM.Scale(scale, scale)
		w, h := crosshair.Bounds().Dx(), crosshair.Bounds().Dy()
		op.GeoM.Translate(float64(cursorX)-float64(w)*scale/2, float64(cursorY)-float64(h)*scale/2)
		screen.DrawImage(crosshair, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
func main() {
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	ebiten.SetWindowTitle("Crudefense")
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowResizable(false)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
