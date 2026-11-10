package render

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// TextureAtlas holds all pre-rendered Greek mythological pixel art
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
	GreekTile     *ebiten.Image
	GreekPillar   *ebiten.Image
	Shadow        *ebiten.Image

	// Pre-allocated UI Elements (Zero allocation in Draw loops)
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
		DrachmaSmall:   renderDrachma(color.NRGBA{212, 175, 55, 255}, 6), // Bronze Drachma
		DrachmaMed:     renderDrachma(color.NRGBA{255, 215, 0, 255}, 8),  // Golden Stater
		DrachmaLarge:   renderDrachma(color.NRGBA{255, 140, 0, 255}, 10), // Ambrosia Gem
		SpearProj:      renderSpear(),
		OwlProj:        renderOwl(),
		BronzeAegis:    renderBronzeAegis(),
		GreekTile:      renderGreekMarbleTile(),
		GreekPillar:    renderGreekPillar(),
		Shadow:         renderShadow(22, 9),
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

// Athena: Classical Greek Goddess with crested Corinthian helmet, red plume, white peplos & bronze breastplate
func renderAthena(dash bool) *ebiten.Image {
	img := ebiten.NewImage(30, 30)

	gold := color.NRGBA{220, 180, 50, 255}
	redPlume := color.NRGBA{200, 30, 30, 255}
	darkPlume := color.NRGBA{140, 20, 20, 255}
	whitePeplos := color.NRGBA{245, 245, 240, 255}
	shadowPeplos := color.NRGBA{185, 180, 170, 255}
	bronzeArmor := color.NRGBA{190, 130, 45, 255}
	skin := color.NRGBA{245, 210, 175, 255}
	doryWood := color.NRGBA{110, 70, 35, 255}

	if dash {
		whitePeplos = color.NRGBA{255, 255, 255, 210}
	}

	// 1. Red Horsehair Plume (Crest on helmet)
	for y := 2; y <= 7; y++ {
		for x := 12; x <= 18; x++ {
			if y <= 4 && (x == 14 || x == 15 || x == 16) {
				img.Set(x, y, redPlume)
			} else if y >= 5 && (x >= 12 && x <= 18) {
				if x == 12 || x == 18 {
					img.Set(x, y, darkPlume)
				} else {
					img.Set(x, y, redPlume)
				}
			}
		}
	}

	// 2. Corinthian Bronze Helmet
	for y := 7; y <= 13; y++ {
		for x := 12; x <= 18; x++ {
			if y == 7 || x == 12 || x == 18 {
				img.Set(x, y, gold)
			} else if y >= 9 && y <= 10 && (x == 14 || x == 16) {
				img.Set(x, y, color.NRGBA{40, 20, 10, 255}) // Eye slits
			} else if y >= 11 && x == 15 {
				img.Set(x, y, gold) // Nose guard
			} else if y == 12 && (x == 14 || x == 16) {
				img.Set(x, y, skin)
			} else {
				img.Set(x, y, bronzeArmor)
			}
		}
	}

	// 3. Torso & Bronze Cuirass
	for y := 14; y <= 18; y++ {
		for x := 11; x <= 19; x++ {
			if x == 11 || x == 19 {
				img.Set(x, y, skin) // Arms
			} else if y <= 16 {
				img.Set(x, y, bronzeArmor) // Bronze breastplate
			} else {
				img.Set(x, y, gold) // Golden belt / girdle
			}
		}
	}

	// 4. White Peplos Toga (Flowing skirt)
	for y := 19; y <= 27; y++ {
		width := float64(y-17) * 0.75
		for x := 15 - int(width); x <= 15 + int(width); x++ {
			if y == 27 || x == 15-int(width) || x == 15+int(width) {
				img.Set(x, y, gold) // Greek key golden hem
			} else if (x-15)%3 == 0 {
				img.Set(x, y, shadowPeplos) // Pleats
			} else {
				img.Set(x, y, whitePeplos)
			}
		}
	}

	// 5. Holding the Spear (Dory) in Right Hand
	for y := 4; y <= 26; y++ {
		if y <= 7 {
			img.Set(21, y, gold) // Bronze Spearhead
		} else {
			img.Set(21, y, doryWood) // Wooden shaft
		}
	}

	// 6. Holding the Golden Aegis Shield on Left Arm
	for y := 14; y <= 22; y++ {
		for x := 6; x <= 10; x++ {
			dx := float64(x - 8)
			dy := float64(y - 18)
			if dx*dx+dy*dy <= 12.0 {
				if dx*dx+dy*dy <= 3.0 {
					img.Set(x, y, color.NRGBA{140, 20, 20, 255}) // Gorgoneion center
				} else {
					img.Set(x, y, gold) // Bronze shield boss
				}
			}
		}
	}

	return img
}

// Hoplite: Corrupted Skeletal Spartan with bronze Corinthian helmet, red plume & round aspis shield
func renderHoplite() *ebiten.Image {
	img := ebiten.NewImage(24, 24)
	bone := color.NRGBA{225, 220, 200, 255}
	bronze := color.NRGBA{180, 120, 40, 255}
	redPlume := color.NRGBA{180, 40, 40, 255}
	wood := color.NRGBA{90, 60, 30, 255}

	// Helmet & Red Plume
	for x := 10; x <= 14; x++ {
		img.Set(x, 2, redPlume)
		img.Set(x, 3, redPlume)
		img.Set(x, 4, bronze)
		img.Set(x, 5, bronze)
	}
	// Skull eyes
	img.Set(11, 6, color.NRGBA{20, 10, 10, 255})
	img.Set(13, 6, color.NRGBA{20, 10, 10, 255})

	// Skeletal Ribcage & Bronze Shield (Aspis)
	for y := 7; y <= 16; y++ {
		for x := 6; x <= 18; x++ {
			// Round Bronze Shield on left
			dx := float64(x - 9)
			dy := float64(y - 12)
			if dx*dx+dy*dy <= 16.0 {
				if dx*dx+dy*dy <= 4.0 {
					img.Set(x, y, color.NRGBA{220, 180, 50, 255}) // Lambda symbol
				} else {
					img.Set(x, y, bronze)
				}
			} else if x >= 11 && x <= 15 && y <= 14 {
				// Ribs
				if y%2 == 0 {
					img.Set(x, y, bone)
				}
			}
		}
	}

	// Spear
	for y := 2; y <= 21; y++ {
		if y <= 5 {
			img.Set(17, y, bronze)
		} else {
			img.Set(17, y, wood)
		}
	}

	// Skeletal legs
	img.Set(11, 17, bone); img.Set(11, 18, bone); img.Set(11, 19, bronze)
	img.Set(14, 17, bone); img.Set(14, 18, bone); img.Set(14, 19, bronze)

	return img
}

// Minotaur: Mythological Bull of Crete with terracotta fur, massive curved ivory horns & bronze bracers
func renderMinotaur() *ebiten.Image {
	img := ebiten.NewImage(34, 34)
	fur := color.NRGBA{120, 65, 30, 255}
	darkFur := color.NRGBA{80, 40, 18, 255}
	ivoryHorn := color.NRGBA{245, 240, 220, 255}
	bronze := color.NRGBA{195, 140, 50, 255}
	snout := color.NRGBA{170, 110, 75, 255}
	redEye := color.NRGBA{220, 30, 30, 255}

	for y := 4; y < 30; y++ {
		for x := 4; x < 30; x++ {
			dx := float64(x - 17)
			dy := float64(y - 17)
			d := math.Sqrt(dx*dx + dy*dy)

			// Massive muscular chest
			if d <= 10 {
				if y <= 15 && math.Abs(dx) <= 5 {
					img.Set(x, y, snout)
				} else if d <= 6 {
					img.Set(x, y, darkFur)
				} else {
					img.Set(x, y, fur)
				}
			}

			// Eyes & Bronze Nose Ring
			if y == 13 && (x == 14 || x == 20) {
				img.Set(x, y, redEye)
			}
			if y == 18 && (x >= 16 && x <= 18) {
				img.Set(x, y, bronze)
			}

			// Giant Sweeping Horns
			if y >= 4 && y <= 11 {
				// Left Horn
				if (x == 10 && y >= 8) || (x == 8 && y >= 6) || (x == 6 && y >= 5) || (x == 5 && y == 4) || (x == 4 && y == 4) {
					img.Set(x, y, ivoryHorn)
				}
				// Right Horn
				if (x == 24 && y >= 8) || (x == 26 && y >= 6) || (x == 28 && y >= 5) || (x == 29 && y == 4) || (x == 30 && y == 4) {
					img.Set(x, y, ivoryHorn)
				}
			}

			// Bronze shoulder pauldrons
			if (y >= 14 && y <= 17) && (x == 6 || x == 7 || x == 27 || x == 28) {
				img.Set(x, y, bronze)
			}
		}
	}
	return img
}

// Gorgon (Medusa): Coiling green snake hair, pale marble skin & petrifying amber eyes
func renderGorgon() *ebiten.Image {
	img := ebiten.NewImage(26, 26)
	snakeGreen := color.NRGBA{45, 125, 50, 255}
	lightSnake := color.NRGBA{80, 175, 80, 255}
	paleSkin := color.NRGBA{215, 220, 215, 255}
	amberEye := color.NRGBA{255, 185, 0, 255}
	peplos := color.NRGBA{90, 70, 95, 255}

	for y := 2; y < 24; y++ {
		for x := 2; x < 24; x++ {
			dx := float64(x - 13)
			dy := float64(y - 12)
			d := math.Sqrt(dx*dx + dy*dy)

			// Face
			if d <= 5 {
				img.Set(x, y, paleSkin)
				if y == 11 && (x == 11 || x == 15) {
					img.Set(x, y, amberEye) // Glowing petrifying gaze
				}
			} else if d > 5 && d <= 9 && y <= 15 {
				// Coiling Snake Hair
				if (x+y)%2 == 0 {
					img.Set(x, y, snakeGreen)
				} else {
					img.Set(x, y, lightSnake)
				}
				// Snake tongue hints
				if (y == 3 || y == 4) && (x == 9 || x == 17) {
					img.Set(x, y, color.NRGBA{220, 40, 40, 255})
				}
			}

			// Dark Greek robe
			if y >= 17 && y <= 24 && math.Abs(dx) <= float64(y-15)*0.8 {
				img.Set(x, y, peplos)
			}
		}
	}
	return img
}

// Stymphalian Bird / Harpy: Mythological bronze feathered raptor with golden talons
func renderHarpy() *ebiten.Image {
	img := ebiten.NewImage(26, 22)
	bronzeFeather := color.NRGBA{190, 130, 45, 255}
	darkFeather := color.NRGBA{130, 80, 25, 255}
	goldTalon := color.NRGBA{255, 200, 40, 255}
	eye := color.NRGBA{230, 40, 40, 255}

	for y := 2; y < 20; y++ {
		for x := 2; x < 24; x++ {
			// Wing spread
			wingY := int(math.Abs(float64(x-13))*0.55) + 6
			if y >= wingY && y <= wingY+2 {
				if y == wingY {
					img.Set(x, y, bronzeFeather)
				} else {
					img.Set(x, y, darkFeather)
				}
			}
			// Raptor body & beak
			if math.Abs(float64(x-13)) <= 2 && y >= 6 && y <= 15 {
				if y == 6 || y == 7 {
					img.Set(x, y, goldTalon) // Bronze beak
				} else {
					img.Set(x, y, bronzeFeather)
				}
				if y == 8 && (x == 12 || x == 14) {
					img.Set(x, y, eye)
				}
			}
			// Talons
			if y >= 16 && y <= 18 && (x == 11 || x == 12 || x == 14 || x == 15) {
				img.Set(x, y, goldTalon)
			}
		}
	}
	return img
}

// Titan Kronos (Boss): Colossal mythological Titan wielding volcanic earth, divine bronze armor & burning eyes
func renderTitan() *ebiten.Image {
	img := ebiten.NewImage(54, 54)
	stone := color.NRGBA{85, 75, 70, 255}
	bronze := color.NRGBA{200, 140, 40, 255}
	gold := color.NRGBA{255, 215, 0, 255}
	magma := color.NRGBA{255, 60, 20, 255}
	crimson := color.NRGBA{180, 25, 25, 255}

	for y := 2; y < 52; y++ {
		for x := 2; x < 52; x++ {
			dx := float64(x - 27)
			dy := float64(y - 27)
			d := math.Sqrt(dx*dx + dy*dy)

			if d <= 20 {
				if d <= 7 {
					img.Set(x, y, magma) // Molten Titan core
				} else if d <= 13 {
					img.Set(x, y, stone)
				} else if d <= 17 {
					img.Set(x, y, bronze) // Ornate bronze plate
				} else {
					img.Set(x, y, crimson)
				}
			}

			// Titan Horns & Golden Crown
			if y <= 12 {
				if (math.Abs(dx) == 8 && y >= 6) || (math.Abs(dx) == 12 && y >= 4) || (math.Abs(dx) == 16 && y == 3) {
					img.Set(x, y, gold)
				}
			}

			// Glowing Magma Cracks
			if (x == 27 || y == 27 || int(dx+dy)%7 == 0) && d <= 18 {
				img.Set(x, y, magma)
			}
		}
	}
	return img
}

// Golden Drachma / Ambrosia Gem
func renderDrachma(c color.NRGBA, size int) *ebiten.Image {
	img := ebiten.NewImage(size, size)
	center := float64(size) / 2.0
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx := float64(x) - center + 0.5
			dy := float64(y) - center + 0.5
			if dx*dx+dy*dy <= (center-0.5)*(center-0.5) {
				if dx*dx+dy*dy <= (center*0.4)*(center*0.4) {
					img.Set(x, y, color.NRGBA{255, 255, 240, 255}) // Owl crest highlight
				} else {
					img.Set(x, y, c)
				}
			}
		}
	}
	return img
}

// Spear of Athena (Dory Javelin)
func renderSpear() *ebiten.Image {
	img := ebiten.NewImage(18, 5)
	gold := color.NRGBA{255, 215, 0, 255}
	bronze := color.NRGBA{210, 150, 45, 255}
	wood := color.NRGBA{115, 75, 40, 255}
	gleam := color.NRGBA{255, 255, 255, 255}

	for x := 0; x < 18; x++ {
		if x >= 13 {
			// Leaf-shaped bronze spearhead
			img.Set(x, 2, gleam)
			img.Set(x, 1, gold)
			img.Set(x, 3, gold)
		} else {
			// Wooden shaft
			img.Set(x, 2, bronze)
			img.Set(x, 1, wood)
			img.Set(x, 3, wood)
		}
	}
	return img
}

// Sacred Owl of Athena Familiar
func renderOwl() *ebiten.Image {
	img := ebiten.NewImage(10, 10)
	white := color.NRGBA{250, 250, 245, 255}
	gold := color.NRGBA{255, 215, 0, 255}
	dark := color.NRGBA{90, 75, 60, 255}

	for y := 1; y < 9; y++ {
		for x := 1; x < 9; x++ {
			if (y >= 3 && y <= 7) && (x >= 2 && x <= 8) {
				img.Set(x, y, white)
			}
			// Big golden wisdom eyes
			if y == 4 && (x == 3 || x == 7) {
				img.Set(x, y, gold)
			}
			// Beak
			if y == 5 && x == 5 {
				img.Set(x, y, gold)
			}
			// Wing tips
			if (x == 1 || x == 9) && (y >= 4 && y <= 6) {
				img.Set(x, y, dark)
			}
		}
	}
	return img
}

// Whirling Bronze Aegis Shield
func renderBronzeAegis() *ebiten.Image {
	img := ebiten.NewImage(16, 16)
	gold := color.NRGBA{255, 215, 0, 255}
	bronze := color.NRGBA{190, 130, 40, 255}
	gorgonFace := color.NRGBA{130, 30, 20, 255}

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			dx := float64(x - 8)
			dy := float64(y - 8)
			distSq := dx*dx + dy*dy
			if distSq <= 49.0 {
				if distSq >= 36.0 {
					img.Set(x, y, gold) // Outer golden rim
				} else if distSq <= 9.0 {
					img.Set(x, y, gorgonFace) // Sacred Gorgoneion
				} else {
					img.Set(x, y, bronze) // Hammered bronze plate
				}
			}
		}
	}
	return img
}

// Classical Greek Marble Temple Floor Tile (Warm Ivory Marble with Sandstone Seams & Greek Key meander border)
func renderGreekMarbleTile() *ebiten.Image {
	img := ebiten.NewImage(32, 32)
	marbleIvory := color.NRGBA{238, 232, 222, 255}
	marbleLight := color.NRGBA{245, 240, 232, 255}
	marbleShadow := color.NRGBA{220, 212, 198, 255}
	seam := color.NRGBA{175, 165, 150, 255}
	terracotta := color.NRGBA{165, 85, 45, 255}
	goldRune := color.NRGBA{200, 160, 60, 255}

	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			// Stone mortar boundary
			if x == 0 || y == 0 || x == 31 || y == 31 {
				img.Set(x, y, seam)
			} else if (x == 1 || y == 1) && (x%8 == 0 || y%8 == 0) {
				// Terracotta corner accents
				img.Set(x, y, terracotta)
			} else if (x == 16 && y == 16) || (x == 4 && y == 4) || (x == 28 && y == 28) {
				img.Set(x, y, goldRune) // Greek rosette dot
			} else if (x+y)%7 == 0 {
				img.Set(x, y, marbleShadow) // Natural marble grain vein
			} else if (x*3+y)%11 == 0 {
				img.Set(x, y, marbleLight)
			} else {
				img.Set(x, y, marbleIvory)
			}
		}
	}
	return img
}

// Classical Fluted Corinthian Greek Pillar
func renderGreekPillar() *ebiten.Image {
	img := ebiten.NewImage(32, 64)
	marbleWhite := color.NRGBA{245, 240, 235, 255}
	marbleShadow := color.NRGBA{185, 175, 165, 255}
	gold := color.NRGBA{225, 185, 60, 255}
	terracotta := color.NRGBA{175, 80, 40, 255}

	for y := 0; y < 64; y++ {
		for x := 0; x < 32; x++ {
			// Stepped Plinth Base (Bottom)
			if y >= 56 {
				if y == 56 || y == 63 || x == 0 || x == 31 {
					img.Set(x, y, terracotta)
				} else if y >= 58 && y <= 60 {
					img.Set(x, y, gold)
				} else {
					img.Set(x, y, marbleWhite)
				}
			} else if y <= 10 {
				// Corinthian Capital (Top) with carved acanthus leaf scrollwork
				if y == 0 || y == 1 || y == 10 {
					img.Set(x, y, terracotta)
				} else if y >= 3 && y <= 7 && (x%5 == 0 || x == 16) {
					img.Set(x, y, gold) // Golden scroll
				} else {
					img.Set(x, y, marbleWhite)
				}
			} else {
				// Fluted Column Shaft
				if x == 2 || x == 29 {
					img.Set(x, y, terracotta) // Column outer edge
				} else if x%5 == 0 {
					img.Set(x, y, marbleShadow) // Fluting groove
				} else {
					img.Set(x, y, marbleWhite)
				}
			}
		}
	}
	return img
}

func renderShadow(w, h int) *ebiten.Image {
	img := ebiten.NewImage(w, h)
	c := color.NRGBA{40, 30, 20, 95}
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
