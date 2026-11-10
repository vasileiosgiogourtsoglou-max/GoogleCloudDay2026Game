package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"olympus2/internal/audio"
	"olympus2/internal/game"
	"olympus2/internal/render"
	"olympus2/internal/scenes"
)

func main() {
	// Initialize Classical Greek art, typography, procedural sound & VFX
	render.InitAtlas()
	render.InitFonts()
	render.InitVFX()
	audio.InitAudio()

	title := scenes.NewTitleScene()
	play := scenes.NewGameplayScene()
	g := game.NewGame(title, play)

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Olympus: Aegis of Athena")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	log.Printf("boot: Olympus initialized with Classical Greek typography & art. Virtual: %dx%d", game.VirtualWidth, game.VirtualHeight)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatalf("game crashed: %v", err)
	}
}
