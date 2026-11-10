package render

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type TextureAtlas struct {
	AthenaIdle    *ebiten.Image
	AthenaDash    *ebiten.Image
	EnemyHoplite  *ebiten.Image
	EnemyMinotaur *ebiten.Image
	EnemyGorgon   *ebiten.Image
	EnemyHarpy    *ebiten.Image
	EnemyTitan    *ebiten.Image
	DrachmaSmall  *ebiten.Image
	DrachmaMed    *ebiten.Image
	DrachmaLarge  *ebiten.Image
	SpearProj     *ebiten.Image
	OwlProj       *ebiten.Image
	BronzeAegis   *ebiten.Image
	FloorTileA    *ebiten.Image
	FloorTileB    *ebiten.Image
	FloorTileC    *ebiten.Image
	GreekPillar   *ebiten.Image
	Brazier       *ebiten.Image
	Shadow        *ebiten.Image
	RingVFX       *ebiten.Image

	// Pre-allocated UI Elements
	BossBarBg      *ebiten.Image
	BossBarFg      *ebiten.Image
	CardBorderSel  *ebiten.Image
	CardBorderNorm *ebiten.Image
	CardInner      *ebiten.Image
	ScrimDark      *ebiten.Image
	ScrimRed       *ebiten.Image
	ScrimGreen     *ebiten.Image
}

var Atlas *TextureAtlas

func InitAtlas() {
	if Atlas != nil {
		return
	}

	bossBg := ebiten.NewImage(70, 6)
	bossBg.Fill(color.NRGBA{30, 20, 15, 240})
	bossFg := ebiten.NewImage(70, 6)
	bossFg.Fill(color.NRGBA{255, 60, 20, 255})

	cardSel := ebiten.NewImage(184, 210)
	cardSel.Fill(color.NRGBA{255, 215, 0, 255})
	cardNorm := ebiten.NewImage(184, 210)
	cardNorm.Fill(color.NRGBA{145, 110, 50, 255})
	cardIn := ebiten.NewImage(178, 204)
	cardIn.Fill(color.NRGBA{32, 24, 18, 255})

	sDark := ebiten.NewImage(640, 360)
	sDark.Fill(color.NRGBA{18, 14, 10, 220})
	sRed := ebiten.NewImage(640, 360)
	sRed.Fill(color.NRGBA{25, 15, 10, 220})
	sGreen := ebiten.NewImage(640, 360)
	sGreen.Fill(color.NRGBA{15, 25, 15, 220})

	Atlas = &TextureAtlas{
		AthenaIdle:     renderAthena(false),
		AthenaDash:     renderAthena(true),
		EnemyHoplite:   renderHoplite(),
		EnemyMinotaur:  renderMinotaur(),
		EnemyGorgon:    renderGorgon(),
		EnemyHarpy:     renderHarpy(),
		EnemyTitan:     renderTitan(),
		DrachmaSmall:   renderDrachma(color.NRGBA{212, 175, 55, 255}, 7),
		DrachmaMed:     renderDrachma(color.NRGBA{255, 215, 0, 255}, 9),
		DrachmaLarge:   renderDrachma(color.NRGBA{255, 120, 0, 255}, 11),
		SpearProj:      renderSpear(),
		OwlProj:        renderOwl(),
		BronzeAegis:    renderBronzeAegis(),
		FloorTileA:     renderGrandMarbleSlab(0),
		FloorTileB:     renderGrandMarbleSlab(1),
		FloorTileC:     renderGrandMarbleSlab(2),
		GreekPillar:    renderGrandCorinthianColumn(),
		Brazier:        renderOlympicBrazier(),
		Shadow:         renderShadow(24, 10),
		RingVFX:        renderRingVFX(24),
		BossBarBg:      bossBg,
		BossBarFg:      bossFg,
		CardBorderSel:  cardSel,
		CardBorderNorm: cardNorm,
		CardInner:      cardIn,
		ScrimDark:      sDark,
		ScrimRed:       sRed,
		ScrimGreen:     sGreen,
	}
}

// Athena: Classical Greek Goddess with flowing crimson plume, Corinthian helmet, peplos & Aegis
func renderAthena(dash bool) *ebiten.Image {
	img := ebiten.NewImage(32, 32)

	gold := color.NRGBA{235, 195, 50, 255}
	brightGold := color.NRGBA{255, 225, 100, 255}
	redPlume := color.NRGBA{220, 35, 35, 255}
	darkPlume := color.NRGBA{150, 20, 20, 255}
	whitePeplos := color.NRGBA{248, 248, 244, 255}
	shadowPeplos := color.NRGBA{190, 185, 175, 255}
	bronzeArmor := color.NRGBA{200, 135, 45, 255}
	skin := color.NRGBA{245, 210, 175, 255}
	doryWood := color.NRGBA{115, 75, 38, 255}

	if dash {
		whitePeplos = color.NRGBA{255, 255, 255, 210}
		gold = brightGold
	}

	// 1. Majestic Red Horsehair Crest (Plume)
	for y := 2; y <= 8; y++ {
		for x := 13; x <= 20; x++ {
			if y <= 4 && (x >= 15 && x <= 18) {
				img.Set(x, y, redPlume)
			} else if y >= 5 {
				if x == 13 || x == 20 {
					img.Set(x, y, darkPlume)
				} else {
					img.Set(x, y, redPlume)
				}
			}
		}
	}

	// 2. Corinthian Bronze Helmet
	for y := 8; y <= 15; y++ {
		for x := 13; x <= 19; x++ {
			if y == 8 || x == 13 || x == 19 {
				img.Set(x, y, gold)
			} else if (y == 10 || y == 11) && (x == 15 || x == 17) {
				img.Set(x, y, color.NRGBA{35, 18, 10, 255}) // Visor slit
			} else if y == 12 && x == 16 {
				img.Set(x, y, gold) // Nose guard
			} else if y >= 13 && (x == 15 || x == 17) {
				img.Set(x, y, skin)
			} else {
				img.Set(x, y, bronzeArmor)
			}
		}
	}

	// 3. Bronze Cuirass & Golden Girdle
	for y := 16; y <= 20; y++ {
		for x := 12; x <= 20; x++ {
			if x == 12 || x == 20 {
				img.Set(x, y, skin) // Arms
			} else if y <= 18 {
				img.Set(x, y, bronzeArmor)
			} else {
				img.Set(x, y, gold) // Ornate gold belt
			}
		}
	}

	// 4. White Peplos Toga (Flowing skirt with gold hem)
	for y := 21; y <= 29; y++ {
		width := float64(y-19) * 0.75
		for x := 16 - int(width); x <= 16 + int(width); x++ {
			if y == 29 || x == 16-int(width) || x == 16+int(width) {
				img.Set(x, y, gold) // Greek key golden hem
			} else if (x-16)%3 == 0 {
				img.Set(x, y, shadowPeplos) // Pleats
			} else {
				img.Set(x, y, whitePeplos)
			}
		}
	}

	// 5. Golden Dory Spear in Right Hand
	for y := 3; y <= 28; y++ {
		if y <= 7 {
			img.Set(23, y, brightGold) // Leaf-shaped spearhead
			img.Set(22, y, gold)
			img.Set(24, y, gold)
		} else {
			img.Set(23, y, doryWood) // Ash-wood shaft
		}
	}

	// 6. Golden Aegis Shield on Left Arm
	for y := 16; y <= 25; y++ {
		for x := 6; x <= 12; x++ {
			dx := float64(x - 9)
			dy := float64(y - 20)
			if dx*dx+dy*dy <= 16.0 {
				if dx*dx+dy*dy <= 4.0 {
					img.Set(x, y, color.NRGBA{160, 25, 20, 255}) // Gorgoneion sun motif
				} else if dx*dx+dy*dy >= 10.0 {
					img.Set(x, y, brightGold)
				} else {
					img.Set(x, y, gold)
				}
			}
		}
	}

	return img
}

// Skeletal Spartan Hoplite with round Aspis shield & spear
func renderHoplite() *ebiten.Image {
	img := ebiten.NewImage(28, 28)
	bone := color.NRGBA{230, 225, 210, 255}
	bronze := color.NRGBA{190, 130, 40, 255}
	redPlume := color.NRGBA{200, 40, 40, 255}
	wood := color.NRGBA{100, 65, 30, 255}
	eyeGlow := color.NRGBA{255, 30, 30, 255}

	// Red Corinthian Plume
	for x := 12; x <= 16; x++ {
		img.Set(x, 2, redPlume)
		img.Set(x, 3, redPlume)
		img.Set(x, 4, bronze)
		img.Set(x, 5, bronze)
	}
	// Glowing red skull eyes
	img.Set(13, 6, eyeGlow)
	img.Set(15, 6, eyeGlow)

	// Round Bronze Shield (Aspis with Spartan Lambda 'Λ')
	for y := 8; y <= 20; y++ {
		for x := 6; x <= 20; x++ {
			dx := float64(x - 11)
			dy := float64(y - 14)
			if dx*dx+dy*dy <= 24.0 {
				// Spartan Lambda 'Λ' in red & gold
				isLambda := (y-8 == int(math.Abs(float64(x-11))*1.8)) && y <= 16
				if isLambda {
					img.Set(x, y, redPlume)
				} else if dx*dx+dy*dy >= 18.0 {
					img.Set(x, y, color.NRGBA{220, 170, 50, 255}) // Rim
				} else {
					img.Set(x, y, bronze)
				}
			} else if x >= 13 && x <= 17 && y <= 17 {
				if y%2 == 0 {
					img.Set(x, y, bone) // Ribs
				}
			}
		}
	}

	// Spear
	for y := 3; y <= 25; y++ {
		if y <= 6 {
			img.Set(20, y, bronze)
		} else {
			img.Set(20, y, wood)
		}
	}

	// Greaves
	img.Set(13, 21, bone); img.Set(13, 22, bronze); img.Set(13, 23, bronze)
	img.Set(16, 21, bone); img.Set(16, 22, bronze); img.Set(16, 23, bronze)

	return img
}

// Cretan Minotaur: Horned beast wielding a massive bronze double-headed Labrys axe
func renderMinotaur() *ebiten.Image {
	img := ebiten.NewImage(38, 38)
	fur := color.NRGBA{115, 60, 28, 255}
	darkFur := color.NRGBA{75, 38, 16, 255}
	ivoryHorn := color.NRGBA{250, 245, 225, 255}
	bronze := color.NRGBA{200, 140, 45, 255}
	snout := color.NRGBA{165, 105, 70, 255}
	redEye := color.NRGBA{230, 30, 30, 255}

	for y := 4; y < 34; y++ {
		for x := 4; x < 34; x++ {
			dx := float64(x - 18)
			dy := float64(y - 19)
			d := math.Sqrt(dx*dx + dy*dy)

			if d <= 12 {
				if y <= 17 && math.Abs(dx) <= 6 {
					img.Set(x, y, snout)
				} else if d <= 7 {
					img.Set(x, y, darkFur)
				} else {
					img.Set(x, y, fur)
				}
			}

			// Red Eyes & Bronze Nose Ring
			if y == 15 && (x == 15 || x == 21) {
				img.Set(x, y, redEye)
			}
			if y == 20 && (x >= 17 && x <= 19) {
				img.Set(x, y, bronze)
			}

			// Curved Ivory Horns
			if y >= 4 && y <= 13 {
				if (x == 11 && y >= 9) || (x == 9 && y >= 7) || (x == 7 && y >= 6) || (x == 5 && y == 5) || (x == 4 && y == 4) {
					img.Set(x, y, ivoryHorn)
				}
				if (x == 25 && y >= 9) || (x == 27 && y >= 7) || (x == 29 && y >= 6) || (x == 31 && y == 5) || (x == 32 && y == 4) {
					img.Set(x, y, ivoryHorn)
				}
			}

			// Bronze pauldrons
			if (y >= 16 && y <= 20) && (x <= 8 || x >= 28) {
				img.Set(x, y, bronze)
			}
		}
	}
	return img
}

// Gorgon Medusa: Slithering green snakes & petrifying gaze
func renderGorgon() *ebiten.Image {
	img := ebiten.NewImage(28, 28)
	snakeGreen := color.NRGBA{40, 130, 50, 255}
	lightSnake := color.NRGBA{90, 185, 90, 255}
	paleSkin := color.NRGBA{220, 225, 218, 255}
	amberEye := color.NRGBA{255, 190, 0, 255}
	crimsonPeplos := color.NRGBA{110, 30, 45, 255}

	for y := 2; y < 26; y++ {
		for x := 2; x < 26; x++ {
			dx := float64(x - 14)
			dy := float64(y - 13)
			d := math.Sqrt(dx*dx + dy*dy)

			if d <= 6 {
				img.Set(x, y, paleSkin)
				if y == 12 && (x == 12 || x == 16) {
					img.Set(x, y, amberEye) // Glowing petrifying gaze
				}
			} else if d > 6 && d <= 10 && y <= 16 {
				if (x+y)%2 == 0 {
					img.Set(x, y, snakeGreen)
				} else {
					img.Set(x, y, lightSnake)
				}
				if (y == 3 || y == 4) && (x == 10 || x == 18) {
					img.Set(x, y, color.NRGBA{230, 40, 40, 255})
				}
			}

			// Peplos
			if y >= 18 && y <= 26 && math.Abs(dx) <= float64(y-16)*0.85 {
				img.Set(x, y, crimsonPeplos)
			}
		}
	}
	return img
}

// Stymphalian Bird / Harpy
func renderHarpy() *ebiten.Image {
	img := ebiten.NewImage(28, 24)
	bronzeFeather := color.NRGBA{195, 135, 45, 255}
	darkFeather := color.NRGBA{135, 85, 25, 255}
	goldTalon := color.NRGBA{255, 205, 45, 255}
	eye := color.NRGBA{230, 40, 40, 255}

	for y := 2; y < 22; y++ {
		for x := 2; x < 26; x++ {
			wingY := int(math.Abs(float64(x-14))*0.55) + 6
			if y >= wingY && y <= wingY+2 {
				if y == wingY {
					img.Set(x, y, bronzeFeather)
				} else {
					img.Set(x, y, darkFeather)
				}
			}
			if math.Abs(float64(x-14)) <= 2 && y >= 6 && y <= 16 {
				if y <= 7 {
					img.Set(x, y, goldTalon)
				} else {
					img.Set(x, y, bronzeFeather)
				}
				if y == 9 && (x == 13 || x == 15) {
					img.Set(x, y, eye)
				}
			}
			if y >= 17 && y <= 20 && (x == 12 || x == 13 || x == 15 || x == 16) {
				img.Set(x, y, goldTalon)
			}
		}
	}
	return img
}

// Titan Kronos (Boss)
func renderTitan() *ebiten.Image {
	img := ebiten.NewImage(58, 58)
	stone := color.NRGBA{80, 70, 65, 255}
	bronze := color.NRGBA{210, 145, 40, 255}
	gold := color.NRGBA{255, 215, 0, 255}
	magma := color.NRGBA{255, 60, 20, 255}
	crimson := color.NRGBA{190, 25, 25, 255}

	for y := 2; y < 56; y++ {
		for x := 2; x < 56; x++ {
			dx := float64(x - 29)
			dy := float64(y - 29)
			d := math.Sqrt(dx*dx + dy*dy)

			if d <= 22 {
				if d <= 8 {
					img.Set(x, y, magma)
				} else if d <= 14 {
					img.Set(x, y, stone)
				} else if d <= 18 {
					img.Set(x, y, bronze)
				} else {
					img.Set(x, y, crimson)
				}
			}

			// Titan Horns & Golden Crown
			if y <= 13 {
				if (math.Abs(dx) == 9 && y >= 6) || (math.Abs(dx) == 13 && y >= 4) || (math.Abs(dx) == 17 && y == 3) {
					img.Set(x, y, gold)
				}
			}

			if (x == 29 || y == 29 || int(dx+dy)%7 == 0) && d <= 20 {
				img.Set(x, y, magma)
			}
		}
	}
	return img
}

// Grand Greek Palace Marble Floor (64x64 tiles with rich variations)
func renderGrandMarbleSlab(variation int) *ebiten.Image {
	img := ebiten.NewImage(64, 64)

	ivoryBase := color.NRGBA{238, 230, 218, 255}
	ivoryLight := color.NRGBA{248, 242, 234, 255}
	ivoryVein := color.NRGBA{218, 208, 192, 255}
	terracotta := color.NRGBA{165, 65, 35, 255}
	goldMeander := color.NRGBA{212, 175, 55, 255}
	mortar := color.NRGBA{150, 138, 122, 255}

	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			// Mortar border
			if x == 0 || y == 0 || x == 63 || y == 63 {
				img.Set(x, y, mortar)
				continue
			}

			// Classical Greek-Key / Meander inner frame (2px in from border)
			if (x == 2 || x == 61 || y == 2 || y == 61) && (x > 2 && x < 61 && y > 2 && y < 61) {
				if (x+y)%4 == 0 {
					img.Set(x, y, terracotta)
				} else {
					img.Set(x, y, goldMeander)
				}
				continue
			}

			// Tile Centerpiece Variations:
			if variation == 1 {
				// Central Olympian Mosaic Rosette
				dx := float64(x - 32)
				dy := float64(y - 32)
				dist := math.Sqrt(dx*dx + dy*dy)
				if dist <= 14.0 {
					if dist <= 4.0 {
						img.Set(x, y, goldMeander)
					} else if dist >= 10.0 && dist <= 12.0 {
						img.Set(x, y, terracotta)
					} else if int(dist)%3 == 0 {
						img.Set(x, y, goldMeander)
					} else {
						img.Set(x, y, ivoryBase)
					}
					continue
				}
			} else if variation == 2 {
				// Corner Terracotta Rosettes
				if (x <= 8 && y <= 8) || (x >= 55 && y <= 8) || (x <= 8 && y >= 55) || (x >= 55 && y >= 55) {
					if (x == 5 && y == 5) || (x == 58 && y == 5) || (x == 5 && y == 58) || (x == 58 && y == 58) {
						img.Set(x, y, goldMeander)
						continue
					} else if (x+y)%2 == 0 {
						img.Set(x, y, terracotta)
						continue
					}
				}
			}

			// Natural Marble Grain
			if (x*2+y)%13 == 0 || (x-y*2)%17 == 0 {
				img.Set(x, y, ivoryVein)
			} else if (x+y)%9 == 0 {
				img.Set(x, y, ivoryLight)
			} else {
				img.Set(x, y, ivoryBase)
			}
		}
	}
	return img
}

// Corinthian Greek Pillar with carved capital and shadow
func renderGrandCorinthianColumn() *ebiten.Image {
	img := ebiten.NewImage(36, 72)
	marbleWhite := color.NRGBA{248, 244, 238, 255}
	marbleShadow := color.NRGBA{180, 168, 155, 255}
	gold := color.NRGBA{230, 185, 55, 255}
	terracotta := color.NRGBA{170, 70, 35, 255}

	for y := 0; y < 72; y++ {
		for x := 0; x < 36; x++ {
			// Stepped Plinth Base
			if y >= 64 {
				if y == 64 || y == 71 || x == 0 || x == 35 {
					img.Set(x, y, terracotta)
				} else if y >= 66 && y <= 68 {
					img.Set(x, y, gold)
				} else {
					img.Set(x, y, marbleWhite)
				}
			} else if y <= 12 {
				// Carved Corinthian Capital
				if y == 0 || y == 1 || y == 12 {
					img.Set(x, y, terracotta)
				} else if y >= 4 && y <= 8 && (x%6 == 0 || x == 18) {
					img.Set(x, y, gold)
				} else {
					img.Set(x, y, marbleWhite)
				}
			} else {
				// Fluted Column Shaft
				if x == 2 || x == 33 {
					img.Set(x, y, terracotta)
				} else if x%5 == 0 {
					img.Set(x, y, marbleShadow)
				} else {
					img.Set(x, y, marbleWhite)
				}
			}
		}
	}
	return img
}

// Olympic Brazier: Sacred burning altar on bronze tripod
func renderOlympicBrazier() *ebiten.Image {
	img := ebiten.NewImage(24, 30)
	bronze := color.NRGBA{195, 130, 40, 255}
	fireCore := color.NRGBA{255, 240, 120, 255}
	fireOrange := color.NRGBA{255, 120, 20, 255}
	fireRed := color.NRGBA{220, 40, 20, 255}

	// Sacred Flame
	for y := 2; y <= 14; y++ {
		for x := 6; x <= 17; x++ {
			dx := math.Abs(float64(x) - 11.5)
			dy := float64(14 - y)
			if dx <= dy*0.6 {
				if dx <= dy*0.25 && y >= 6 {
					img.Set(x, y, fireCore)
				} else if dx <= dy*0.45 {
					img.Set(x, y, fireOrange)
				} else {
					img.Set(x, y, fireRed)
				}
			}
		}
	}

	// Bronze Bowl
	for y := 14; y <= 19; y++ {
		for x := 4; x <= 19; x++ {
			img.Set(x, y, bronze)
		}
	}

	// Tripod Legs
	for y := 20; y <= 28; y++ {
		img.Set(5, y, bronze)
		img.Set(11, y, bronze)
		img.Set(18, y, bronze)
	}

	return img
}

func renderDrachma(c color.NRGBA, size int) *ebiten.Image {
	img := ebiten.NewImage(size, size)
	center := float64(size) / 2.0
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx := float64(x) - center + 0.5
			dy := float64(y) - center + 0.5
			if dx*dx+dy*dy <= (center-0.5)*(center-0.5) {
				if dx*dx+dy*dy <= (center*0.4)*(center*0.4) {
					img.Set(x, y, color.NRGBA{255, 255, 240, 255})
				} else {
					img.Set(x, y, c)
				}
			}
		}
	}
	return img
}

func renderSpear() *ebiten.Image {
	img := ebiten.NewImage(20, 6)
	gold := color.NRGBA{255, 220, 50, 255}
	bronze := color.NRGBA{210, 150, 45, 255}
	wood := color.NRGBA{115, 75, 40, 255}
	gleam := color.NRGBA{255, 255, 255, 255}

	for x := 0; x < 20; x++ {
		if x >= 14 {
			img.Set(x, 2, gleam)
			img.Set(x, 3, gleam)
			img.Set(x, 1, gold)
			img.Set(x, 4, gold)
		} else {
			img.Set(x, 2, bronze)
			img.Set(x, 3, bronze)
			img.Set(x, 1, wood)
			img.Set(x, 4, wood)
		}
	}
	return img
}

func renderOwl() *ebiten.Image {
	img := ebiten.NewImage(12, 12)
	white := color.NRGBA{250, 250, 245, 255}
	gold := color.NRGBA{255, 215, 0, 255}
	dark := color.NRGBA{90, 75, 60, 255}

	for y := 1; y < 11; y++ {
		for x := 1; x < 11; x++ {
			if (y >= 3 && y <= 8) && (x >= 2 && x <= 9) {
				img.Set(x, y, white)
			}
			if y == 5 && (x == 4 || x == 8) {
				img.Set(x, y, gold)
			}
			if y == 6 && x == 6 {
				img.Set(x, y, gold)
			}
			if (x == 1 || x == 10) && (y >= 5 && y <= 7) {
				img.Set(x, y, dark)
			}
		}
	}
	return img
}

func renderBronzeAegis() *ebiten.Image {
	img := ebiten.NewImage(18, 18)
	gold := color.NRGBA{255, 220, 50, 255}
	bronze := color.NRGBA{195, 135, 45, 255}
	gorgonFace := color.NRGBA{150, 30, 20, 255}

	for y := 0; y < 18; y++ {
		for x := 0; x < 18; x++ {
			dx := float64(x - 9)
			dy := float64(y - 9)
			distSq := dx*dx + dy*dy
			if distSq <= 64.0 {
				if distSq >= 49.0 {
					img.Set(x, y, gold)
				} else if distSq <= 12.0 {
					img.Set(x, y, gorgonFace)
				} else {
					img.Set(x, y, bronze)
				}
			}
		}
	}
	return img
}

func renderRingVFX(radius int) *ebiten.Image {
	size := radius * 2
	img := ebiten.NewImage(size, size)
	c := color.NRGBA{255, 220, 60, 200}
	center := float64(radius)
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx := float64(x) - center
			dy := float64(y) - center
			d := math.Sqrt(dx*dx + dy*dy)
			if d >= float64(radius-2) && d <= float64(radius) {
				img.Set(x, y, c)
			}
		}
	}
	return img
}

func renderShadow(w, h int) *ebiten.Image {
	img := ebiten.NewImage(w, h)
	c := color.NRGBA{35, 25, 15, 95}
	rx := float64(w) / 2.0
	ry := float64(h) / 2.0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dx := (float64(x) - rx) / rx
			dy := (float64(y) - ry) / ry
			if dx*dx+dy*dy <= 1.0 {
				img.Set(x, y, c)
			}
		}
	}
	return img
}
