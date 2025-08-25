package main

import (
	_ "embed"
)

// Fonts
//
//go:embed src/fonts/Queensides-3z7Ey.ttf
var fontBytes []byte

// Images
//
//go:embed src/gui/crosshair.png
var crosshairPNG []byte

//go:embed src/gui/playButton.png
var playButtonPNG []byte

//go:embed src/gui/archerButton.png
var archerButtonPNG []byte

//go:embed src/gui/cogwheel.png
var settingsButtonPNG []byte

//go:embed src/gui/goBack.png
var goBackPNG []byte

//go:embed src/gui/fscreen.png
var fscreenPNG []byte

//go:embed src/entities/mainChar.png
var mainCharPNG []byte

//go:embed src/gui/menuButton.png
var menuButtonPNG []byte

//go:embed src/gui/background.png
var backgroundPNG []byte

// Audio
//
//go:embed src/audio/clickSound.wav
var clickWav []byte
