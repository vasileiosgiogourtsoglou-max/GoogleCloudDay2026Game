package entities

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"olympus2/internal/render"
)

type Pickup struct {
	X, Y      float64
	VX, VY    float64
	Value     int
	Radius    float64
	Collected bool
}

func NewPickup(x, y float64, value int) *Pickup {
	return &Pickup{
		X:      x,
		Y:      y,
		Value:  value,
		Radius: 6.0,
	}
}

func (p *Pickup) Update(dt float64, playerX, playerY float64) bool {
	if p.Collected {
		return false
	}
	dx := playerX - p.X
	dy := playerY - p.Y
	dist := math.Sqrt(dx*dx + dy*dy)

	magnetRange := 90.0
	pickupRange := 16.0

	// Magnet attraction to Athena
	if dist <= magnetRange {
		speed := 270.0
		p.VX = (dx / dist) * speed
		p.VY = (dy / dist) * speed
		p.X += p.VX * dt
		p.Y += p.VY * dt
	}

	if dist <= pickupRange {
		p.Collected = true
		return true // Collected!
	}
	return false
}

func (p *Pickup) Draw(screen *ebiten.Image, camX, camY float64) {
	if p.Collected {
		return
	}
	op := &ebiten.DrawImageOptions{}

	var sprite *ebiten.Image
	offset := 3.0
	if p.Value >= 10 {
		sprite = render.Atlas.DrachmaLarge
		offset = 5.0
	} else if p.Value >= 4 {
		sprite = render.Atlas.DrachmaMed
		offset = 4.0
	} else {
		sprite = render.Atlas.DrachmaSmall
		offset = 3.0
	}

	op.GeoM.Translate(p.X-camX-offset, p.Y-camY-offset)
	screen.DrawImage(sprite, op)
}
