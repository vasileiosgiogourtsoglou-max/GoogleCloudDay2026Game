package entities

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"olympus2/internal/audio"
	"olympus2/internal/render"
)

type EnemyType int

const (
	TypeHoplite EnemyType = iota
	TypeMinotaur
	TypeGorgon
	TypeHarpy
	TypeTitan
)

type Enemy struct {
	Type     EnemyType
	X, Y     float64
	VX, VY   float64
	Speed    float64
	MaxHP    float64
	HP       float64
	Damage   float64
	Radius   float64
	XPAmount int
	Alive    bool
	HitFlash float64
	ChargeCD float64
	IsBoss   bool
}

func NewEnemy(t EnemyType, x, y float64) *Enemy {
	e := &Enemy{
		Type:  t,
		X:     x,
		Y:     y,
		Alive: true,
	}
	switch t {
	case TypeHoplite:
		e.Speed = 95.0
		e.MaxHP = 18.0
		e.HP = 18.0
		e.Damage = 12.0
		e.Radius = 10.0
		e.XPAmount = 1
	case TypeMinotaur:
		e.Speed = 65.0
		e.MaxHP = 85.0
		e.HP = 85.0
		e.Damage = 25.0
		e.Radius = 16.0
		e.XPAmount = 4
	case TypeGorgon:
		e.Speed = 45.0
		e.MaxHP = 45.0
		e.HP = 45.0
		e.Damage = 15.0
		e.Radius = 12.0
		e.XPAmount = 2
	case TypeHarpy:
		e.Speed = 130.0
		e.MaxHP = 22.0
		e.HP = 22.0
		e.Damage = 14.0
		e.Radius = 11.0
		e.XPAmount = 2
	case TypeTitan:
		e.Speed = 50.0
		e.MaxHP = 1200.0
		e.HP = 1200.0
		e.Damage = 35.0
		e.Radius = 26.0
		e.XPAmount = 30
		e.IsBoss = true
	}
	return e
}

func (e *Enemy) Update(dt float64, targetX, targetY float64) {
	if !e.Alive {
		return
	}
	if e.HitFlash > 0 {
		e.HitFlash -= dt
	}

	dx := targetX - e.X
	dy := targetY - e.Y
	dist := math.Sqrt(dx*dx + dy*dy)

	if dist > 1.0 {
		spd := e.Speed
		// Minotaur charge
		if e.Type == TypeMinotaur {
			e.ChargeCD -= dt
			if e.ChargeCD <= 0 && dist < 170.0 {
				spd *= 2.2
				if e.ChargeCD <= -0.8 {
					e.ChargeCD = 3.5
				}
			}
		}
		e.VX = (dx / dist) * spd
		e.VY = (dy / dist) * spd
	} else {
		e.VX = 0
		e.VY = 0
	}

	e.X += e.VX * dt
	e.Y += e.VY * dt
}

func (e *Enemy) TakeDamage(dmg float64) bool {
	e.HP -= dmg
	e.HitFlash = 0.12
	if e.HP <= 0 {
		e.Alive = false
		audio.PlayExplosion()
		return true
	}
	return false
}

func (e *Enemy) Draw(screen *ebiten.Image, camX, camY float64) {
	if !e.Alive {
		return
	}

	// Shadow
	shOp := &ebiten.DrawImageOptions{}
	shW := e.Radius * 1.6
	shOp.GeoM.Scale(shW/22.0, 1.0)
	shOp.GeoM.Translate(e.X-camX-(shW/2.0), e.Y-camY+e.Radius*0.6)
	screen.DrawImage(render.Atlas.Shadow, shOp)

	var sprite *ebiten.Image
	offsetX, offsetY := 0.0, 0.0
	switch e.Type {
	case TypeHoplite:
		sprite = render.Atlas.EnemyHoplite
		offsetX, offsetY = 12, 12
	case TypeMinotaur:
		sprite = render.Atlas.EnemyMinotaur
		offsetX, offsetY = 17, 17
	case TypeGorgon:
		sprite = render.Atlas.EnemyGorgon
		offsetX, offsetY = 13, 13
	case TypeHarpy:
		sprite = render.Atlas.EnemyHarpy
		offsetX, offsetY = 13, 11
	case TypeTitan:
		sprite = render.Atlas.EnemyTitan
		offsetX, offsetY = 27, 27
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(e.X-camX-offsetX, e.Y-camY-offsetY)
	if e.HitFlash > 0 {
		op.ColorScale.Scale(2.2, 1.8, 1.4, 1.0) // Gold-white flash
	}
	screen.DrawImage(sprite, op)

	// Titan Boss Health Bar (Zero-allocation using pre-rendered textures)
	if e.IsBoss && render.Atlas.BossBarBg != nil {
		barW := 70.0
		bx := e.X - camX - barW/2.0
		by := e.Y - camY - offsetY - 12.0
		pct := math.Max(0.0, e.HP/e.MaxHP)

		bgOp := &ebiten.DrawImageOptions{}
		bgOp.GeoM.Translate(bx, by)
		screen.DrawImage(render.Atlas.BossBarBg, bgOp)

		if pct > 0 {
			fgOp := &ebiten.DrawImageOptions{}
			fgOp.GeoM.Scale(pct, 1.0)
			fgOp.GeoM.Translate(bx, by)
			screen.DrawImage(render.Atlas.BossBarFg, fgOp)
		}
	}
}
