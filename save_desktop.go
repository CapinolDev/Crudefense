//go:build !js

package main

import (
	"encoding/json"
	"os"
)

func SaveSettings(filename string, settings Settings) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(settings)
}
