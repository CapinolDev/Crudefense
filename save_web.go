//go:build js

package main

import "log"

func SaveSettings(filename string, settings Settings) error {
	log.Println("Skipping SaveSettings on browser build")
	return nil
}
