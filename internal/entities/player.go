package entities

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"olympus2/internal/audio"
	"olympus2/internal/render"
)

type UpgradeType int

const (
	UpgradeLaser UpgradeType = iota
	UpgradeOrbitalAegis
	UpgradeCyberOwl
	UpgradeLightning
	UpgradeSpeed
	UpgradeMaxHP
)

type Player struct {
	X, Y           float64
	VX, VY         float64
	Speed          float64
	MaxHP          float64
	HP             float64
	MaxShield      float64
	Shield         float64
	ShieldRegenCD  float64
	IsDashing      bool
	DashDuration   float64
	DashCooldown   float64
	DashTimer      float64
	InvincibleTime float64

	Level        int
	XP           int
	XPToNext     int
	TotalKills   int
	SurvivalTime float64

	// Weapon Levels
	LaserLvl       int
	LaserTimer     float64
	OrbitalLvl     int
	OrbitalAngle   float64
	CyberOwlLvl    int
	OwlTimer       float64
	LightningLvl   int
	LightningTimer float64
}

func NewPlayer(x, y float64) *Player {
	return &Player{
		X:             x,
		Y:             y,
		Speed:         145.0,
		MaxHP:         100.0,
		HP:            100.0,
		MaxShield:     50.0,
		Shield:        50.0,
		Level:         1,
		XP:            0,
		XPToNext:      10,
		LaserLvl:      1,
		OrbitalLvl:    1,
		LaserTimer:    0.2,
		DashCooldown:  1.2,
		ShieldRegenCD: 0.0,
	}
}

func (p *Player) Update(dt float64, inputX, inputY float64, dashTrigger bool) {
	p.SurvivalTime += dt

	if p.InvincibleTime > 0 {
		p.InvincibleTime -= dt
	}

	if p.DashTimer > 0 {
		p.DashTimer -= dt
	}
	if p.IsDashing {
		p.DashDuration -= dt
		if p.DashDuration <= 0 {
			p.IsDashing = false
		}
	} else if dashTrigger && p.DashTimer <= 0 {
		p.IsDashing = true
		p.DashDuration = 0.22
		p.DashTimer = p.DashCooldown
		p.InvincibleTime = 0.25
		audio.PlayDash()
	}

	currentSpeed := p.Speed
	if p.IsDashing {
		currentSpeed *= 2.8
	}

	mag := math.Sqrt(inputX*inputX + inputY*inputY)
	if mag > 0.001 {
		p.VX = (inputX / mag) * currentSpeed
		p.VY = (inputY / mag) * currentSpeed
	} else if !p.IsDashing {
		p.VX *= 0.8
		p.VY *= 0.8
	}

	p.X += p.VX * dt
	p.Y += p.VY * dt

	if p.ShieldRegenCD > 0 {
		p.ShieldRegenCD -= dt
	} else if p.Shield < p.MaxShield {
		p.Shield = math.Min(p.MaxShield, p.Shield+12.0*dt)
	}

	// Orbital shield rotation
	p.OrbitalAngle += 3.5 * dt
	if p.OrbitalAngle > math.Pi*2 {
		p.OrbitalAngle -= math.Pi * 2
	}
}

func (p *Player) TakeDamage(dmg float64) bool {
	if p.InvincibleTime > 0 || p.IsDashing {
		return false
	}
	p.InvincibleTime = 0.4
	p.ShieldRegenCD = 3.0
	audio.PlayHit()

	if p.Shield > 0 {
		if p.Shield >= dmg {
			p.Shield -= dmg
			return false
		}
		remaining := dmg - p.Shield
		p.Shield = 0
		p.HP -= remaining
	} else {
		p.HP -= dmg
	}
	return p.HP <= 0
}

func (p *Player) AddXP(amount int) bool {
	p.XP += amount
	audio.PlayXP()
	if p.XP >= p.XPToNext {
		p.XP -= p.XPToNext
		p.Level++
		p.XPToNext = int(float64(p.XPToNext) * 1.35)
		audio.PlayLevelUp()
		return true
	}
	return false
}

func (p *Player) Draw(screen *ebiten.Image, camX, camY float64) {
	// Shadow
	shOp := &ebiten.DrawImageOptions{}
	shOp.GeoM.Translate(p.X-camX-11, p.Y-camY+9)
	screen.DrawImage(render.Atlas.Shadow, shOp)

	// Athena Sprite (flicker when invincible)
	if p.InvincibleTime > 0 && int(p.InvincibleTime*20)%2 == 0 {
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.X-camX-15, p.Y-camY-15)
	if p.IsDashing {
		screen.DrawImage(render.Atlas.AthenaDash, op)
	} else {
		screen.DrawImage(render.Atlas.AthenaIdle, op)
	}

	// Whirling Bronze Aegis Shields
	if p.OrbitalLvl > 0 {
		count := p.OrbitalLvl + 1
		radius := 36.0
		for i := 0; i < count; i++ {
			ang := p.OrbitalAngle + (float64(i) * (2.0 * math.Pi / float64(count)))
			ox := p.X + math.Cos(ang)*radius - camX - 8
			oy := p.Y + math.Sin(ang)*radius - camY - 8

			orbOp := &ebiten.DrawImageOptions{}
			orbOp.GeoM.Translate(ox, oy)
			screen.DrawImage(render.Atlas.BronzeAegis, orbOp)
		}
	}
}
