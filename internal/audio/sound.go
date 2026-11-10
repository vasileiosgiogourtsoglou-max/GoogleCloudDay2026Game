package audio

import (
	"bytes"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const SampleRate = 44100

type SoundManager struct {
	ctx        *audio.Context
	bgmPlayer  *audio.Player
	spearData  []byte
	dashData   []byte
	coinData   []byte
	levelData  []byte
	boomData   []byte
	hitData    []byte
	lastPlayed map[string]time.Time
}

var GlobalManager *SoundManager

func InitAudio() {
	if GlobalManager != nil {
		return
	}
	ctx := audio.NewContext(SampleRate)
	sm := &SoundManager{
		ctx:        ctx,
		lastPlayed: make(map[string]time.Time),
	}

	sm.spearData = generateSpearThrow()
	sm.dashData = generateDash()
	sm.coinData = generateDrachmaChime()
	sm.levelData = generateOlympianFanfare()
	sm.boomData = generateTitanSmite()
	sm.hitData = generateBronzeHit()

	// Heroic Ancient Greek Olympian March BGM (Dorian mode, brass, marching percussion)
	bgmData := generateGreekMarchBGM()
	loop := audio.NewInfiniteLoop(bytes.NewReader(bgmData), int64(len(bgmData)))
	if p, err := ctx.NewPlayer(loop); err == nil {
		p.SetVolume(0.55)
		sm.bgmPlayer = p
		sm.bgmPlayer.Play()
	}

	GlobalManager = sm
}

func (sm *SoundManager) playCue(name string, data []byte, cooldown time.Duration, vol float64) {
	if sm == nil || sm.ctx == nil || len(data) == 0 {
		return
	}
	now := time.Now()
	if last, ok := sm.lastPlayed[name]; ok && now.Sub(last) < cooldown {
		return
	}
	sm.lastPlayed[name] = now

	p := sm.ctx.NewPlayerFromBytes(data)
	if p != nil {
		p.SetVolume(vol)
		p.Play()
	}
}

func PlayLaser()     { GlobalManager.playCue("spear", GlobalManager.spearData, 60*time.Millisecond, 0.45) }
func PlayDash()      { GlobalManager.playCue("dash", GlobalManager.dashData, 150*time.Millisecond, 0.5) }
func PlayXP()        { GlobalManager.playCue("coin", GlobalManager.coinData, 40*time.Millisecond, 0.45) }
func PlayLevelUp()   { GlobalManager.playCue("level", GlobalManager.levelData, 300*time.Millisecond, 0.7) }
func PlayExplosion() { GlobalManager.playCue("smite", GlobalManager.boomData, 80*time.Millisecond, 0.6) }
func PlayHit()       { GlobalManager.playCue("hit", GlobalManager.hitData, 50*time.Millisecond, 0.4) }

func encodePCM(samplesL, samplesR []float64) []byte {
	n := len(samplesL)
	buf := make([]byte, n*4)
	for i := 0; i < n; i++ {
		l := math.Max(-1.0, math.Min(1.0, samplesL[i]))
		r := math.Max(-1.0, math.Min(1.0, samplesR[i]))
		valL := int16(l * 32767)
		valR := int16(r * 32767)

		idx := i * 4
		buf[idx] = byte(valL)
		buf[idx+1] = byte(valL >> 8)
		buf[idx+2] = byte(valR)
		buf[idx+3] = byte(valR >> 8)
	}
	aligned := len(buf) &^ 3
	return buf[:aligned]
}

func generateSpearThrow() []byte {
	dur := 0.14
	samples := int(float64(SampleRate) * dur)
	l := make([]float64, samples)
	r := make([]float64, samples)

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(SampleRate)
		env := (1.0 - t/dur) * (1.0 - t/dur)
		freq := 700.0 * math.Exp(-t*20.0)
		sig := math.Sin(2.0*math.Pi*freq*t) + 0.3*(rand.Float64()*2.0-1.0)
		l[i] = sig * env * 0.65
		r[i] = sig * env * 0.65
	}
	return encodePCM(l, r)
}

func generateDash() []byte {
	dur := 0.22
	samples := int(float64(SampleRate) * dur)
	l := make([]float64, samples)
	r := make([]float64, samples)

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(SampleRate)
		env := math.Sin((t / dur) * math.Pi)
		noise := (rand.Float64()*2.0 - 1.0)
		sig := noise * env * 0.6
		l[i] = sig
		r[i] = sig
	}
	return encodePCM(l, r)
}

func generateDrachmaChime() []byte {
	dur := 0.18
	samples := int(float64(SampleRate) * dur)
	l := make([]float64, samples)
	r := make([]float64, samples)
	freq := 1760.0 // A6 Greek Lyre harmonic

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(SampleRate)
		env := math.Exp(-t * 22.0)
		sig := math.Sin(2.0*math.Pi*freq*t) + 0.35*math.Sin(2.0*math.Pi*freq*2.0*t)
		l[i] = sig * env * 0.45
		r[i] = sig * env * 0.45
	}
	return encodePCM(l, r)
}

func generateOlympianFanfare() []byte {
	notes := []float64{440.0, 554.37, 659.25, 880.0} // A4, C#5, E5, A5 (Heroic Major Triad)
	noteDur := 0.11
	totalDur := noteDur * 4.0
	samples := int(float64(SampleRate) * totalDur)
	l := make([]float64, samples)
	r := make([]float64, samples)

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(SampleRate)
		idx := int(t / noteDur)
		if idx >= len(notes) {
			idx = len(notes) - 1
		}
		freq := notes[idx]
		subT := math.Mod(t, noteDur)
		env := math.Exp(-subT * 7.0)
		sig := math.Sin(2.0*math.Pi*freq*t) + 0.3*math.Sin(2.0*math.Pi*freq*2.0*t) + 0.15*math.Sin(2.0*math.Pi*freq*3.0*t)
		l[i] = sig * env * 0.6
		r[i] = sig * env * 0.6
	}
	return encodePCM(l, r)
}

func generateTitanSmite() []byte {
	dur := 0.35
	samples := int(float64(SampleRate) * dur)
	l := make([]float64, samples)
	r := make([]float64, samples)

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(SampleRate)
		env := math.Exp(-t * 9.0)
		subBass := math.Sin(2.0 * math.Pi * 75.0 * (1.0 - t/dur) * t)
		noise := (rand.Float64()*2.0 - 1.0)
		sig := (subBass*0.75 + noise*0.35) * env
		l[i] = sig
		r[i] = sig
	}
	return encodePCM(l, r)
}

func generateBronzeHit() []byte {
	dur := 0.08
	samples := int(float64(SampleRate) * dur)
	l := make([]float64, samples)
	r := make([]float64, samples)

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(SampleRate)
		env := math.Exp(-t * 30.0)
		sig := (math.Sin(2.0*math.Pi*420.0*t) + 0.5*math.Sin(2.0*math.Pi*840.0*t)) * env
		l[i] = sig * 0.5
		r[i] = sig * 0.5
	}
	return encodePCM(l, r)
}

func generateGreekMarchBGM() []byte {
	bpm := 128.0
	secondsPerBeat := 60.0 / bpm
	totalBeats := 32
	totalDur := secondsPerBeat * float64(totalBeats)
	samples := int(float64(SampleRate) * totalDur)
	l := make([]float64, samples)
	r := make([]float64, samples)

	// Authentic Ancient Greek Dorian / Phrygian Mode: D, E, F, G, A, Bb, C
	bassNotes := []float64{146.83, 146.83, 174.61, 164.81} // D3, D3, F3, E3
	melodyNotes := []float64{293.66, 329.63, 349.23, 392.00, 440.00, 392.00, 349.23, 293.66}

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(SampleRate)
		beat := t / secondsPerBeat
		subBeat := math.Mod(beat, 1.0)

		// 1. Spartan Tympanon (War Drum) on beats 1 and 3
		drum := 0.0
		barPos := math.Mod(beat, 2.0)
		if barPos < 0.5 {
			drumEnv := math.Exp(-subBeat * 16.0)
			drum = math.Sin(2.0*math.Pi*80.0*subBeat) * drumEnv * 0.7
		}

		// 2. Bronze Crotales / Shaker
		crotale := (rand.Float64()*2.0 - 1.0) * math.Exp(-subBeat*30.0) * 0.15

		// 3. Greek Lyre Bass Plucks
		bassStep := int(beat * 2.0)
		subStep := math.Mod(beat*2.0, 1.0)
		bassIdx := (bassStep / 4) % len(bassNotes)
		baseFreq := bassNotes[bassIdx]
		bassEnv := math.Exp(-subStep * 9.0)
		bass := (math.Sin(2.0*math.Pi*baseFreq*t) + 0.4*math.Sin(2.0*math.Pi*baseFreq*2.0*t)) * bassEnv * 0.35

		// 4. Aulos Horn & Lyre Melody
		melStep := int(beat * 2.0)
		subMel := math.Mod(beat*2.0, 1.0)
		melFreq := melodyNotes[melStep%len(melodyNotes)]
		melEnv := math.Exp(-subMel * 5.0)
		mel := (math.Sin(2.0*math.Pi*melFreq*t) + 0.25*math.Sin(2.0*math.Pi*melFreq*2.0*t) + 0.15*math.Sin(2.0*math.Pi*melFreq*3.0*t)) * melEnv * 0.28

		mixL := drum + crotale*0.6 + bass*0.85 + mel*0.8
		mixR := drum + crotale*0.8 + bass*0.85 + mel*0.85

		l[i] = mixL
		r[i] = mixR
	}
	return encodePCM(l, r)
}
