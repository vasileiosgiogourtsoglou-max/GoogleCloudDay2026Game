package scenes

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"olympus2/internal/audio"
	"olympus2/internal/entities"
	"olympus2/internal/game"
	"olympus2/internal/render"
)

type UpgradeCard struct {
	Type     entities.UpgradeType
	Title    string
	Lines    []string
	LevelStr string
}

type LevelUpOverlay struct {
	Cards     []UpgradeCard
	Selected  int
	darkScrim *ebiten.Image
}

func NewLevelUpOverlay(p *entities.Player) *LevelUpOverlay {
	pool := []UpgradeCard{
		{
			Type:     entities.UpgradeOrbitalAegis,
			Title:    "DIVINE AEGIS",
			Lines:    []string{"+1 Whirling Bronze Shield", "Strikes Nearby Foes"},
			LevelStr: "Rank +1",
		},
		{
			Type:     entities.UpgradeLaser,
			Title:    "SPEAR OF ATHENA",
			Lines:    []string{"Faster Javelin Throws", "+35% Pierce Damage"},
			LevelStr: "Rank +1",
		},
		{
			Type:     entities.UpgradeCyberOwl,
			Title:    "OWL OF WISDOM",
			Lines:    []string{"Sacred Owl Companion", "Fires Divine Starlight"},
			LevelStr: "Rank +1",
		},
		{
			Type:     entities.UpgradeSpeed,
			Title:    "HERMES' SANDALS",
			Lines:    []string{"+20% Sacred Move Speed", "-15% Dash Cooldown"},
			LevelStr: "Divine Speed",
		},
		{
			Type:     entities.UpgradeMaxHP,
			Title:    "AMBROSIA OF LIFE",
			Lines:    []string{"+30 Maximum Health", "Restores Full Vitality"},
			LevelStr: "Immortal Grace",
		},
	}

	rand.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
	chosen := pool[:3]

	scrim := ebiten.NewImage(game.VirtualWidth, game.VirtualHeight)
	scrim.Fill(color.NRGBA{18, 14, 10, 220}) // Warm dark bronze scrim

	return &LevelUpOverlay{
		Cards:     chosen,
		darkScrim: scrim,
	}
}

func (l *LevelUpOverlay) Update(p *entities.Player) bool {
	choice := -1
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		choice = 0
	} else if inpututil.IsKeyJustPressed(ebiten.Key2) {
		choice = 1
	} else if inpututil.IsKeyJustPressed(ebiten.Key3) {
		choice = 2
	} else if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		choice = l.Selected
	}

	mx, my := ebiten.CursorPosition()
	cardW := 184
	cardH := 210
	cardGap := 16
	startX := (game.VirtualWidth - (3*cardW + 2*cardGap)) / 2
	cardY := 85

	for i := 0; i < len(l.Cards); i++ {
		cx := startX + i*(cardW+cardGap)
		if mx >= cx && mx <= cx+cardW && my >= cardY && my <= cardY+cardH {
			l.Selected = i
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				choice = i
			}
		}
	}

	if choice >= 0 && choice < len(l.Cards) {
		l.applyUpgrade(p, l.Cards[choice].Type)
		audio.PlayLevelUp()
		return true
	}
	return false
}

func (l *LevelUpOverlay) applyUpgrade(p *entities.Player, up entities.UpgradeType) {
	switch up {
	case entities.UpgradeOrbitalAegis:
		p.OrbitalLvl++
	case entities.UpgradeLaser:
		p.LaserLvl++
	case entities.UpgradeCyberOwl:
		p.CyberOwlLvl++
	case entities.UpgradeSpeed:
		p.Speed *= 1.20
		p.DashCooldown *= 0.85
	case entities.UpgradeMaxHP:
		p.MaxHP += 30.0
		p.HP = p.MaxHP
		p.Shield = p.MaxShield
	}
}

func (l *LevelUpOverlay) Draw(screen *ebiten.Image) {
	screen.DrawImage(l.darkScrim, &ebiten.DrawImageOptions{})

	gold := color.NRGBA{255, 215, 0, 255}
	warmWhite := color.NRGBA{245, 240, 230, 255}
	rankGold := color.NRGBA{220, 175, 60, 255}

	render.DrawCentered(screen, "FAVOR OF THE OLYMPIANS", render.Fonts.TitleLarge, float64(game.VirtualWidth/2), 34, gold)
	render.DrawCentered(screen, "Select a Divine Blessing from Mount Olympus", render.Fonts.TitleSmall, float64(game.VirtualWidth/2), 58, warmWhite)

	cardW := 184
	cardH := 210
	cardGap := 16
	startX := (game.VirtualWidth - (3*cardW + 2*cardGap)) / 2
	cardY := 85

	borderSel := ebiten.NewImage(cardW, cardH)
	borderSel.Fill(color.NRGBA{255, 215, 0, 255}) // Radiant Olympian Gold

	borderNorm := ebiten.NewImage(cardW, cardH)
	borderNorm.Fill(color.NRGBA{145, 110, 50, 255}) // Ancient Bronze

	cardInner := ebiten.NewImage(cardW-6, cardH-6)
	cardInner.Fill(color.NRGBA{32, 24, 18, 255}) // Dark Grecian Stone Tile

	for i, c := range l.Cards {
		cx := startX + i*(cardW+cardGap)

		bOp := &ebiten.DrawImageOptions{}
		bOp.GeoM.Translate(float64(cx), float64(cardY))
		if i == l.Selected {
			screen.DrawImage(borderSel, bOp)
		} else {
			screen.DrawImage(borderNorm, bOp)
		}

		inOp := &ebiten.DrawImageOptions{}
		inOp.GeoM.Translate(float64(cx+3), float64(cardY+3))
		screen.DrawImage(cardInner, inOp)

		numStr := string(rune('1' + i))
		numColor := gold
		if i == l.Selected {
			numColor = color.NRGBA{255, 255, 220, 255}
		}
		render.DrawCentered(screen, "["+numStr+"]", render.Fonts.TitleMedium, float64(cx+cardW/2), float64(cardY+28), numColor)
		render.DrawCentered(screen, c.Title, render.Fonts.CardTitle, float64(cx+cardW/2), float64(cardY+65), warmWhite)
		render.DrawCentered(screen, c.LevelStr, render.Fonts.CardTag, float64(cx+cardW/2), float64(cardY+95), rankGold)

		for lineIdx, line := range c.Lines {
			lineY := float64(cardY + 130 + lineIdx*22)
			render.DrawCentered(screen, line, render.Fonts.CardBody, float64(cx+cardW/2), lineY, color.NRGBA{230, 225, 215, 255})
		}
	}
}
