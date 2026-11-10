package scenes

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"olympus2/internal/audio"
	"olympus2/internal/entities"
	"olympus2/internal/game"
	"olympus2/internal/render"
)

type GameplayScene struct {
	Player        *entities.Player
	Enemies       []*entities.Enemy
	Projectiles   []*entities.Projectile
	Pickups       []*entities.Pickup
	SpawnTimer    float64
	BossSpawned   bool
	IsGameOver    bool
	IsVictory     bool
	LevelUpActive *LevelUpOverlay

	CamX, CamY float64

	// Pre-rendered HUD elements
	hpBarBg  *ebiten.Image
	hpBarFg  *ebiten.Image
	shieldFg *ebiten.Image
	xpBarBg  *ebiten.Image
	xpBarFg  *ebiten.Image
	dashIcon *ebiten.Image
}

func NewGameplayScene() *GameplayScene {
	p := entities.NewPlayer(0, 0)
	s := &GameplayScene{
		Player:      p,
		Enemies:     make([]*entities.Enemy, 0, 200),
		Projectiles: make([]*entities.Projectile, 0, 100),
		Pickups:     make([]*entities.Pickup, 0, 300),

		hpBarBg:  ebiten.NewImage(120, 10),
		hpBarFg:  ebiten.NewImage(116, 6),
		shieldFg: ebiten.NewImage(116, 2),
		xpBarBg:  ebiten.NewImage(game.VirtualWidth, 6),
		xpBarFg:  ebiten.NewImage(game.VirtualWidth, 6),
		dashIcon: ebiten.NewImage(16, 16),
	}
	s.hpBarBg.Fill(color.NRGBA{40, 28, 18, 240})
	s.hpBarFg.Fill(color.NRGBA{220, 50, 40, 255})    // Spartan Crimson
	s.shieldFg.Fill(color.NRGBA{255, 215, 0, 255})   // Olympian Gold
	s.xpBarBg.Fill(color.NRGBA{30, 22, 14, 200})
	s.xpBarFg.Fill(color.NRGBA{230, 180, 50, 255})   // Golden Glory
	s.dashIcon.Fill(color.NRGBA{255, 215, 0, 255})

	return s
}

func (s *GameplayScene) Update(dt float64) (game.StateType, bool) {
	if s.LevelUpActive != nil {
		if s.LevelUpActive.Update(s.Player) {
			s.LevelUpActive = nil
		}
		return game.StateGameplay, false
	}

	if s.IsGameOver || s.IsVictory {
		if inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			return game.StateTitle, true
		}
		return game.StateGameplay, false
	}

	inX, inY := 0.0, 0.0
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		inY -= 1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		inY += 1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		inX -= 1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		inX += 1.0
	}

	dash := inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
		inpututil.IsKeyJustPressed(ebiten.KeyShift) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)

	s.Player.Update(dt, inX, inY, dash)

	targetCamX := s.Player.X - float64(game.VirtualWidth)/2.0
	targetCamY := s.Player.Y - float64(game.VirtualHeight)/2.0
	s.CamX += (targetCamX - s.CamX) * (1.0 - math.Exp(-10.0*dt))
	s.CamY += (targetCamY - s.CamY) * (1.0 - math.Exp(-10.0*dt))

	s.updateWeapons(dt)
	s.updateSpawner(dt)
	s.updateEnemies(dt)
	s.updateProjectiles(dt)
	s.updatePickups(dt)

	if s.Player.HP <= 0 {
		s.IsGameOver = true
		audio.PlayExplosion()
	}

	return game.StateGameplay, false
}

func (s *GameplayScene) updateWeapons(dt float64) {
	var closest *entities.Enemy
	minDist := 330.0

	for _, e := range s.Enemies {
		if !e.Alive {
			continue
		}
		d := math.Hypot(e.X-s.Player.X, e.Y-s.Player.Y)
		if d < minDist {
			minDist = d
			closest = e
		}
	}

	// 1. Spear of Athena
	s.Player.LaserTimer -= dt
	fireRate := math.Max(0.14, 0.38-float64(s.Player.LaserLvl)*0.04)
	if s.Player.LaserTimer <= 0 && closest != nil {
		s.Player.LaserTimer = fireRate
		spearDmg := 15.0 + float64(s.Player.LaserLvl)*5.0
		s.Projectiles = append(s.Projectiles, entities.NewSpear(s.Player.X, s.Player.Y, closest.X, closest.Y, spearDmg))
		audio.PlayLaser()
	}

	// 2. Owl of Athena
	if s.Player.CyberOwlLvl > 0 {
		s.Player.OwlTimer -= dt
		if s.Player.OwlTimer <= 0 && closest != nil {
			s.Player.OwlTimer = 0.8
			dartDmg := 22.0 + float64(s.Player.CyberOwlLvl)*6.0
			s.Projectiles = append(s.Projectiles, entities.NewOwlDart(s.Player.X, s.Player.Y, dartDmg))
		}
	}

	// 3. Whirling Bronze Aegis
	if s.Player.OrbitalLvl > 0 {
		count := s.Player.OrbitalLvl + 1
		radius := 36.0
		for i := 0; i < count; i++ {
			ang := s.Player.OrbitalAngle + (float64(i) * (2.0 * math.Pi / float64(count)))
			ox := s.Player.X + math.Cos(ang)*radius
			oy := s.Player.Y + math.Sin(ang)*radius

			for _, e := range s.Enemies {
				if !e.Alive {
					continue
				}
				if math.Hypot(e.X-ox, e.Y-oy) <= e.Radius+8.0 {
					if e.TakeDamage(26.0 * dt) {
						s.onEnemyKilled(e)
					}
				}
			}
		}
	}
}

func (s *GameplayScene) updateSpawner(dt float64) {
	s.SpawnTimer -= dt
	spawnInterval := math.Max(0.24, 1.2-s.Player.SurvivalTime*0.003)

	if s.SpawnTimer <= 0 {
		s.SpawnTimer = spawnInterval

		spawnAngle := rand.Float64() * math.Pi * 2.0
		spawnDist := 380.0
		ex := s.Player.X + math.Cos(spawnAngle)*spawnDist
		ey := s.Player.Y + math.Sin(spawnAngle)*spawnDist

		t := entities.TypeHoplite
		roll := rand.Float64()
		if s.Player.SurvivalTime > 90.0 && roll < 0.25 {
			t = entities.TypeMinotaur
		} else if s.Player.SurvivalTime > 50.0 && roll < 0.45 {
			t = entities.TypeGorgon
		} else if s.Player.SurvivalTime > 25.0 && roll < 0.65 {
			t = entities.TypeHarpy
		}

		s.Enemies = append(s.Enemies, entities.NewEnemy(t, ex, ey))
	}

	if !s.BossSpawned && s.Player.SurvivalTime >= 180.0 {
		s.BossSpawned = true
		s.Enemies = append(s.Enemies, entities.NewEnemy(entities.TypeTitan, s.Player.X+280, s.Player.Y))
	}
}

func (s *GameplayScene) updateEnemies(dt float64) {
	for _, e := range s.Enemies {
		if !e.Alive {
			continue
		}
		e.Update(dt, s.Player.X, s.Player.Y)

		dist := math.Hypot(e.X-s.Player.X, e.Y-s.Player.Y)
		if dist <= e.Radius+10.0 {
			if s.Player.TakeDamage(e.Damage) {
				s.IsGameOver = true
			}
		}
	}
}

func (s *GameplayScene) updateProjectiles(dt float64) {
	for _, p := range s.Projectiles {
		if !p.Active {
			continue
		}
		p.Update(dt, s.Player.X, s.Player.Y)

		for _, e := range s.Enemies {
			if !e.Alive {
				continue
			}
			if math.Hypot(e.X-p.X, e.Y-p.Y) <= e.Radius+p.Radius {
				if e.TakeDamage(p.Damage) {
					s.onEnemyKilled(e)
				}
				p.Penetrate--
				if p.Penetrate <= 0 {
					p.Active = false
					break
				}
			}
		}
	}
}

func (s *GameplayScene) updatePickups(dt float64) {
	for _, p := range s.Pickups {
		if p.Collected {
			continue
		}
		if p.Update(dt, s.Player.X, s.Player.Y) {
			if s.Player.AddXP(p.Value) {
				s.LevelUpActive = NewLevelUpOverlay(s.Player)
			}
		}
	}
}

func (s *GameplayScene) onEnemyKilled(e *entities.Enemy) {
	s.Player.TotalKills++
	s.Pickups = append(s.Pickups, entities.NewPickup(e.X, e.Y, e.XPAmount))
	if e.IsBoss {
		s.IsVictory = true
		audio.PlayLevelUp()
	}
}

func (s *GameplayScene) Draw(screen *ebiten.Image) {
	tileW := 32.0
	startX := math.Floor(s.CamX/tileW)*tileW - tileW
	startY := math.Floor(s.CamY/tileW)*tileW - tileW
	endX := s.CamX + float64(game.VirtualWidth) + tileW
	endY := s.CamY + float64(game.VirtualHeight) + tileW

	for ty := startY; ty <= endY; ty += tileW {
		for tx := startX; tx <= endX; tx += tileW {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(tx-s.CamX, ty-s.CamY)
			screen.DrawImage(render.Atlas.GreekTile, op)
		}
	}

	for _, p := range s.Pickups {
		p.Draw(screen, s.CamX, s.CamY)
	}
	for _, e := range s.Enemies {
		e.Draw(screen, s.CamX, s.CamY)
	}
	s.Player.Draw(screen, s.CamX, s.CamY)
	for _, pr := range s.Projectiles {
		pr.Draw(screen, s.CamX, s.CamY)
	}

	s.drawHUD(screen)

	if s.LevelUpActive != nil {
		s.LevelUpActive.Draw(screen)
	} else if s.IsGameOver {
		s.drawGameOver(screen)
	} else if s.IsVictory {
		s.drawVictory(screen)
	}
}

func (s *GameplayScene) drawHUD(screen *ebiten.Image) {
	screen.DrawImage(s.xpBarBg, &ebiten.DrawImageOptions{})
	xpPct := math.Min(1.0, float64(s.Player.XP)/float64(s.Player.XPToNext))
	if xpPct > 0 {
		xpOp := &ebiten.DrawImageOptions{}
		xpOp.GeoM.Scale(xpPct, 1.0)
		screen.DrawImage(s.xpBarFg, xpOp)
	}

	hpX, hpY := 14.0, 14.0
	bgOp := &ebiten.DrawImageOptions{}
	bgOp.GeoM.Translate(hpX, hpY)
	screen.DrawImage(s.hpBarBg, bgOp)

	hpPct := math.Max(0.0, s.Player.HP/s.Player.MaxHP)
	if hpPct > 0 {
		fgOp := &ebiten.DrawImageOptions{}
		fgOp.GeoM.Scale(hpPct, 1.0)
		fgOp.GeoM.Translate(hpX+2, hpY+2)
		screen.DrawImage(s.hpBarFg, fgOp)
	}

	shieldPct := math.Max(0.0, s.Player.Shield/s.Player.MaxShield)
	if shieldPct > 0 {
		shOp := &ebiten.DrawImageOptions{}
		shOp.GeoM.Scale(shieldPct, 1.0)
		shOp.GeoM.Translate(hpX+2, hpY+7)
		screen.DrawImage(s.shieldFg, shOp)
	}

	lvlText := fmt.Sprintf("FAVOR %d", s.Player.Level)
	render.DrawCentered(screen, lvlText, render.Fonts.HUD, 44, 34, color.NRGBA{212, 175, 55, 255})

	mins := int(s.Player.SurvivalTime) / 60
	secs := int(s.Player.SurvivalTime) % 60
	timerText := fmt.Sprintf("%02d:%02d", mins, secs)
	killsText := fmt.Sprintf("⚔ %d", s.Player.TotalKills)
	render.DrawCentered(screen, timerText, render.Fonts.HUD, float64(game.VirtualWidth-45), 18, color.NRGBA{45, 30, 15, 255})
	render.DrawCentered(screen, killsText, render.Fonts.HUD, float64(game.VirtualWidth-45), 34, color.NRGBA{165, 55, 30, 255})

	dashX := float64(game.VirtualWidth/2 - 8)
	dashY := float64(game.VirtualHeight - 24)
	dOp := &ebiten.DrawImageOptions{}
	dOp.GeoM.Translate(dashX, dashY)
	if s.Player.DashTimer > 0 {
		dOp.ColorScale.Scale(0.35, 0.3, 0.25, 0.6)
	}
	screen.DrawImage(s.dashIcon, dOp)
}

func (s *GameplayScene) drawGameOver(screen *ebiten.Image) {
	scrim := ebiten.NewImage(game.VirtualWidth, game.VirtualHeight)
	scrim.Fill(color.NRGBA{25, 15, 10, 220})
	screen.DrawImage(scrim, &ebiten.DrawImageOptions{})

	render.DrawCentered(screen, "ATHENA HAS FALLEN", render.Fonts.TitleLarge, float64(game.VirtualWidth/2), 95, color.NRGBA{220, 50, 40, 255})
	render.DrawCentered(screen, "Mount Olympus Has Fallen to Mythological Beasts", render.Fonts.TitleSmall, float64(game.VirtualWidth/2), 135, color.NRGBA{235, 225, 215, 255})

	stats := fmt.Sprintf("Survival Time: %02d:%02d  |  Foes Vanquished: %d", int(s.Player.SurvivalTime)/60, int(s.Player.SurvivalTime)%60, s.Player.TotalKills)
	render.DrawCentered(screen, stats, render.Fonts.TitleSmall, float64(game.VirtualWidth/2), 180, color.NRGBA{255, 215, 0, 255})

	render.DrawCentered(screen, "[ PRESS SPACE OR R TO ARISE AGAIN ]", render.Fonts.TitleMedium, float64(game.VirtualWidth/2), 245, color.NRGBA{212, 175, 55, 255})
}

func (s *GameplayScene) drawVictory(screen *ebiten.Image) {
	scrim := ebiten.NewImage(game.VirtualWidth, game.VirtualHeight)
	scrim.Fill(color.NRGBA{15, 25, 15, 220})
	screen.DrawImage(scrim, &ebiten.DrawImageOptions{})

	render.DrawCentered(screen, "TITAN KRONOS VANQUISHED", render.Fonts.TitleLarge, float64(game.VirtualWidth/2), 95, color.NRGBA{255, 215, 0, 255})
	render.DrawCentered(screen, "Olympus is Victorious! Glory to the Goddess Athena!", render.Fonts.TitleSmall, float64(game.VirtualWidth/2), 135, color.NRGBA{245, 240, 230, 255})

	stats := fmt.Sprintf("Total Battle Time: %02d:%02d  |  Total Foes Slain: %d", int(s.Player.SurvivalTime)/60, int(s.Player.SurvivalTime)%60, s.Player.TotalKills)
	render.DrawCentered(screen, stats, render.Fonts.TitleSmall, float64(game.VirtualWidth/2), 180, color.NRGBA{255, 215, 0, 255})

	render.DrawCentered(screen, "[ PRESS SPACE OR R TO REJOICE ]", render.Fonts.TitleMedium, float64(game.VirtualWidth/2), 245, color.NRGBA{255, 215, 0, 255})
}
