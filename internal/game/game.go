package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Scene interface {
	Update(dt float64) (StateType, bool)
	Draw(screen *ebiten.Image)
}

type Game struct {
	CurrentState StateType
	titleScene   Scene
	playScene    Scene
}

func NewGame(title Scene, play Scene) *Game {
	return &Game{
		CurrentState: StateTitle,
		titleScene:   title,
		playScene:    play,
	}
}

func (g *Game) Update() error {
	dt := 1.0 / float64(TargetTPS)

	// F11 / Alt+Enter Fullscreen toggle per ebitengineer guide
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) || (ebiten.IsKeyPressed(ebiten.KeyAlt) && inpututil.IsKeyJustPressed(ebiten.KeyEnter)) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	switch g.CurrentState {
	case StateTitle:
		nextState, change := g.titleScene.Update(dt)
		if change {
			g.CurrentState = nextState
		}
	case StateGameplay:
		nextState, change := g.playScene.Update(dt)
		if change {
			g.CurrentState = nextState
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.CurrentState {
	case StateTitle:
		g.titleScene.Draw(screen)
	case StateGameplay:
		g.playScene.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return VirtualWidth, VirtualHeight
}
