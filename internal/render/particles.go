package render

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Particle struct {
	X, Y     float64
	VX, VY   float64
	Life     float64
	MaxLife  float64
	Color    color.NRGBA
	Size     float64
	Active   bool
}

type DamageText struct {
	X, Y    float64
	Text    string
	Color   color.NRGBA
	Life    float64
	MaxLife float64
	Active  bool
}

type ParticleSystem struct {
	particles []Particle
	texts     []DamageText
	dotImg    *ebiten.Image
}

var GlobalVFX *ParticleSystem

func InitVFX() {
	if GlobalVFX != nil {
		return
	}
	dot := ebiten.NewImage(3, 3)
	dot.Fill(color.NRGBA{255, 255, 255, 255})

	GlobalVFX = &ParticleSystem{
		particles: make([]Particle, 250),
		texts:     make([]DamageText, 60),
		dotImg:    dot,
	}
}

func (ps *ParticleSystem) SpawnEliminationBurst(x, y float64, isBoss bool) {
	count := 12
	if isBoss {
		count = 35
	}
	gold := color.NRGBA{255, 215, 0, 255}
	terracotta := color.NRGBA{220, 75, 40, 255}
	divineWhite := color.NRGBA{255, 250, 220, 255}

	colors := []color.NRGBA{gold, terracotta, divineWhite}

	for i := 0; i < count; i++ {
		angle := rand.Float64() * math.Pi * 2.0
		speed := 50.0 + rand.Float64()*120.0
		if isBoss {
			speed *= 1.8
		}
		c := colors[rand.Intn(len(colors))]
		ps.addParticle(x, y, math.Cos(angle)*speed, math.Sin(angle)*speed, 0.45+rand.Float64()*0.3, c, 2.5)
	}
}

func (ps *ParticleSystem) SpawnDashTrail(x, y float64) {
	c := color.NRGBA{255, 225, 120, 180}
	for i := 0; i < 3; i++ {
		rx := x + (rand.Float64()*10.0 - 5.0)
		ry := y + (rand.Float64()*10.0 - 5.0)
		ps.addParticle(rx, ry, (rand.Float64()*20.0 - 10.0), (rand.Float64()*20.0 - 10.0), 0.25, c, 2.0)
	}
}

func (ps *ParticleSystem) SpawnDamageNumber(x, y float64, damage float64) {
	text := fmt.Sprintf("%d", int(damage+0.5))
	c := color.NRGBA{255, 235, 100, 255}
	if damage >= 30 {
		c = color.NRGBA{255, 70, 50, 255} // Critical hit
	}

	for i := range ps.texts {
		if !ps.texts[i].Active {
			ps.texts[i] = DamageText{
				X:       x + (rand.Float64()*12.0 - 6.0),
				Y:       y - 10.0,
				Text:    text,
				Color:   c,
				Life:    0.6,
				MaxLife: 0.6,
				Active:  true,
			}
			break
		}
	}
}

func (ps *ParticleSystem) addParticle(x, y, vx, vy, life float64, c color.NRGBA, size float64) {
	for i := range ps.particles {
		if !ps.particles[i].Active {
			ps.particles[i] = Particle{
				X:       x,
				Y:       y,
				VX:      vx,
				VY:      vy,
				Life:    life,
				MaxLife: life,
				Color:   c,
				Size:    size,
				Active:  true,
			}
			break
		}
	}
}

func (ps *ParticleSystem) Update(dt float64) {
	for i := range ps.particles {
		if !ps.particles[i].Active {
			continue
		}
		ps.particles[i].Life -= dt
		if ps.particles[i].Life <= 0 {
			ps.particles[i].Active = false
			continue
		}
		ps.particles[i].X += ps.particles[i].VX * dt
		ps.particles[i].Y += ps.particles[i].VY * dt
		ps.particles[i].VX *= 0.93 // Drag
		ps.particles[i].VY *= 0.93
	}

	for i := range ps.texts {
		if !ps.texts[i].Active {
			continue
		}
		ps.texts[i].Life -= dt
		if ps.texts[i].Life <= 0 {
			ps.texts[i].Active = false
			continue
		}
		ps.texts[i].Y -= 28.0 * dt // Float upward
	}
}

func (ps *ParticleSystem) Draw(screen *ebiten.Image, camX, camY float64) {
	// Draw particles
	for i := range ps.particles {
		p := &ps.particles[i]
		if !p.Active {
			continue
		}
		alpha := p.Life / p.MaxLife
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(p.Size/3.0, p.Size/3.0)
		op.GeoM.Translate(p.X-camX, p.Y-camY)
		op.ColorScale.Scale(
			float32(p.Color.R)/255.0,
			float32(p.Color.G)/255.0,
			float32(p.Color.B)/255.0,
			float32(alpha),
		)
		screen.DrawImage(ps.dotImg, op)
	}

	// Draw floating damage numbers
	if Fonts != nil && Fonts.CardTag != nil {
		for i := range ps.texts {
			t := &ps.texts[i]
			if !t.Active {
				continue
			}
			alpha := t.Life / t.MaxLife
			c := color.NRGBA{
				R: t.Color.R,
				G: t.Color.G,
				B: t.Color.B,
				A: uint8(alpha * 255.0),
			}
			DrawCentered(screen, t.Text, Fonts.CardTag, t.X-camX, t.Y-camY, c)
		}
	}
}
