package scenes

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"olympus2/internal/audio"
	"olympus2/internal/game"
	"olympus2/internal/render"
)

type TitleScene struct {
	timeElapsed float64
}

func NewTitleScene() *TitleScene {
	return &TitleScene{}
}

func (s *TitleScene) Update(dt float64) (game.StateType, bool) {
	s.timeElapsed += dt

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
		inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		audio.PlayLevelUp()
		return game.StateGameplay, true
	}
	return game.StateTitle, false
}

func (s *TitleScene) Draw(screen *ebiten.Image) {
	// Background: Warm Classical Greek Marble Tiles
	tilesX := game.VirtualWidth/32 + 1
	tilesY := game.VirtualHeight/32 + 1
	for ty := 0; ty < tilesY; ty++ {
		for tx := 0; tx < tilesX; tx++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(tx*32), float64(ty*32))
			screen.DrawImage(render.Atlas.GreekTile, op)
		}
	}

	// Classical Fluted Corinthian Pillars
	pOpL := &ebiten.DrawImageOptions{}
	pOpL.GeoM.Scale(1.8, 2.8)
	pOpL.GeoM.Translate(25, 60)
	screen.DrawImage(render.Atlas.GreekPillar, pOpL)

	pOpR := &ebiten.DrawImageOptions{}
	pOpR.GeoM.Scale(1.8, 2.8)
	pOpR.GeoM.Translate(float64(game.VirtualWidth-80), 60)
	screen.DrawImage(render.Atlas.GreekPillar, pOpR)

	// Athena Goddess Sprite
	athenaOp := &ebiten.DrawImageOptions{}
	athenaOp.GeoM.Scale(3.2, 3.2)
	athenaOp.GeoM.Translate(float64(game.VirtualWidth/2-48), 105+math.Sin(s.timeElapsed*3.0)*7.0)
	screen.DrawImage(render.Atlas.AthenaIdle, athenaOp)

	gold := color.NRGBA{212, 175, 55, 255}
	terracotta := color.NRGBA{175, 60, 35, 255}
	darkText := color.NRGBA{45, 30, 15, 255}

	render.DrawCentered(screen, "OLYMPUS", render.Fonts.TitleLarge, float64(game.VirtualWidth/2), 42, gold)
	render.DrawCentered(screen, "AEGIS OF ATHENA", render.Fonts.TitleMedium, float64(game.VirtualWidth/2), 72, terracotta)

	if int(s.timeElapsed*2.5)%2 == 0 {
		render.DrawCentered(screen, "[ PRESS SPACE OR CLICK TO DEFEND OLYMPUS ]", render.Fonts.TitleSmall, float64(game.VirtualWidth/2), 275, darkText)
	}

	render.DrawCentered(screen, "WASD / Arrows to Move | Spacebar to Dash | Divine Auto-Attack", render.Fonts.CardBody, float64(game.VirtualWidth/2), 325, color.NRGBA{90, 70, 50, 255})
}
