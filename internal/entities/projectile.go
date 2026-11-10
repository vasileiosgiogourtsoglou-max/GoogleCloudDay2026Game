package entities

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"olympus2/internal/render"
)

type ProjectileType int

const (
	ProjSpear ProjectileType = iota
	ProjOwlDart
)

type Projectile struct {
	Type      ProjectileType
	X, Y      float64
	VX, VY    float64
	Angle     float64
	Damage    float64
	Life      float64
	MaxLife   float64
	Radius    float64
	Penetrate int
	Active    bool
}

func NewSpear(x, y, targetX, targetY float64, damage float64) *Projectile {
	dx := targetX - x
	dy := targetY - y
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist < 0.001 {
		dist = 1.0
		dx = 1.0
	}
	speed := 440.0
	vx := (dx / dist) * speed
	vy := (dy / dist) * speed
	angle := math.Atan2(dy, dx)

	return &Projectile{
		Type:      ProjSpear,
		X:         x,
		Y:         y,
		VX:        vx,
		VY:        vy,
		Angle:     angle,
		Damage:    damage,
		Life:      1.2,
		MaxLife:   1.2,
		Radius:    7.0,
		Penetrate: 1,
		Active:    true,
	}
}

func NewOwlDart(x, y float64, damage float64) *Projectile {
	return &Projectile{
		Type:      ProjOwlDart,
		X:         x,
		Y:         y,
		Damage:    damage,
		Life:      2.0,
		MaxLife:   2.0,
		Radius:    6.0,
		Penetrate: 1,
		Active:    true,
	}
}

func (p *Projectile) Update(dt float64, targetX, targetY float64) {
	if !p.Active {
		return
	}
	p.Life -= dt
	if p.Life <= 0 {
		p.Active = false
		return
	}

	// Owl homing behavior
	if p.Type == ProjOwlDart {
		dx := targetX - p.X
		dy := targetY - p.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		speed := 310.0
		if dist > 1.0 {
			p.VX = (dx / dist) * speed
			p.VY = (dy / dist) * speed
		}
	}

	p.X += p.VX * dt
	p.Y += p.VY * dt
}

func (p *Projectile) Draw(screen *ebiten.Image, camX, camY float64) {
	if !p.Active {
		return
	}
	op := &ebiten.DrawImageOptions{}

	switch p.Type {
	case ProjSpear:
		op.GeoM.Translate(-9, -2)
		op.GeoM.Rotate(p.Angle)
		op.GeoM.Translate(p.X-camX, p.Y-camY)
		screen.DrawImage(render.Atlas.SpearProj, op)
	case ProjOwlDart:
		op.GeoM.Translate(p.X-camX-5, p.Y-camY-5)
		screen.DrawImage(render.Atlas.OwlProj, op)
	}
}
