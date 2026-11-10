# Game Design Document (GDD)

> **Game Title**: Olympus 2.0: Aegis of Athena
> **Target 2D Genre**: Top-Down Bullet-Heaven / Action Roguelike
> **Camera Projection & Spatial Model**: 3/4 Top-Down Orthographic with smooth follow-camera
> **Simulation Model**: Real-Time Continuous (60 TPS fixed timestep in Ebitengine)
> **Target Platform**: Desktop (macOS/Windows/Linux) & WebAssembly (Browser / Game Dev Studio Playtest Tab / Cloud Run)
> **Virtual Canvas & Aspect Ratio**: 16:9 Widescreen (`640x360` virtual resolution, retro 16-bit pixel aesthetic)
> **Author / Lead Designer**: Vasilis Giogourtsogl & Antigravity (Google Cloud Day Athens 2026)

---

## 1. Executive Summary & Elevator Pitch
- **Elevator Pitch**: In a futuristic techno-mythological Greece, Mount Olympus has been overclocked into *Olympus 2.0*. As **Athena AI**, the digital goddess of wisdom and strategic warfare, players must fend off relentless waves of corrupted mythological automata and cyber-titans using autonomous orbital aegis lasers, smart cyber-drones, and tactical plasma shields.
- **Core Inspiration**: *Vampire Survivors*, *Hades*, *Enter the Gungeon*.
- **Target Mood & Pacing**: High-adrenaline neon cyber-arcade action, escalating intensity with hordes of mythological cyber-beasts, satisfying power fantasy build synergy, and heroic SEGA-style 16-bit Greek synthwave soundtrack.

---

## 2. Core Gameplay Loop, Simulation & Mechanics
- **Primary Gameplay Loop**: 
  1. **Evade & Position**: Move through neon cyber-marble ruins, dodging converging enemy swarms.
  2. **Auto-Attack**: Autonomous digital weapons (Aegis laser, Cyber-Owls, Lightning Surge) automatically acquire targets and destroy foes.
  3. **Harvest Aether**: Collect glowing neon XP orbs dropped by disintegrated cyber-monsters.
  4. **Overclock / Level Up**: Upon filling the XP bar, select 1 of 3 randomized blessings to evolve Athena AI's build.
  5. **Survive & Triumph**: Defeat the 5-minute Corrupted Titan Boss (*Chronos Core*) to restore systemic balance to Olympus 2.0.

- **Spatial & Simulation Rules**:
  - **Arena**: Continuous wrapping or bounded cyber-temple courtyard with neon columns and cyber-geometric mosaic tiles.
  - **Player Movement**: 8-directional top-down movement with smooth acceleration and Spacebar Dash (invulnerability frames + quick thrust).
  - **Enemy Archetypes**:
    - *Cyber-Scorpions / Drones*: Fast, low HP, swarming in packs.
    - *Minotaur Automata*: Heavy chargers that wind up high-speed line rushes.
    - *Gorgon Towers*: Stationary/slow turrets firing slow-moving digital petrification rings.
    - *Harpy Gliders*: Circling swoopers that sweep across the arena.
    - *Boss - Chronos Core*: Massive overclocked cyber-titan summoning orbital beam grids and shockwaves.
  - **Collision**: Spatial hash grid for fast broad-phase collision detection between hundreds of active projectiles, XP orbs, and enemies at a solid 60 FPS.

- **Win & Loss Conditions**:
  - **Win Condition**: Survive the escalating 5-minute onslaught and destroy the final boss (*Chronos Core*).
  - **Loss Condition**: Athena AI's shield/HP reaches 0.

---

## 3. Controls & Input Mapping Scheme
- **Primary Input Devices**: Keyboard, Mouse, Gamepad, Touch (Virtual Stick & Touch Dash).

| Logical Action | Keyboard / Mouse | Gamepad Button | Touch Gesture |
| :--- | :--- | :--- | :--- |
| `Move` | `WASD` / `Arrow Keys` | Left Stick / D-Pad | Virtual Joystick (drag) |
| `Dash / Thrust` | `Space` / `Shift` / Right Click | Button South (`A`) / Shoulder (`RB`) | Quick Double Tap / Dash Button |
| `Select / Confirm` | `Enter` / `Space` / Left Click | Button South (`A`) | Screen Tap |
| `Pause / Menu` | `Escape` / `P` | Start / Options | Pause Icon (Top-Right) |

---

## 4. Visual Art Direction & Sprite Pipeline (`sprite-design` + `nano-banana`)
- **Aesthetic Direction & Color Palette**: 
  - 16-bit cyber-Hellenic aesthetic: Cyber-white marble (`#F4F6FB`), gold photon trim (`#FFD700`, `#FFA500`), electric cyan lasers (`#00FFFF`), neon magenta corrupted core glow (`#FF007F`), and deep night obsidian shadows (`#0D0F18`).
- **Sprite Resolution Tier**:
  - **Retro Pixel Art (`32x32` - `64x64` per frame)**: Rendered on a `640x360` virtual canvas with integer scaling for razor-sharp pixel presentation.
- **Selected Visual Pipeline Mode**: Hybrid Pipeline
  - **Hero & Enemies**: Procedural `.sprite` DSL and optimized sprite sheets with `idle`, `move`, `dash`, `attack`, `hit`, and `destroy` states.
  - **VFX & Environment**: Ebitengine hardware-accelerated particle pools (plasma sparks, laser lines, expanding shockwaves, neon aether orb pulses).

---

## 5. Audio & Soundscape Pipeline (`procedural-sfx` + `procedural-music`)
- **Selected Audio Pipeline Mode**: Pure-Code Procedural (`.sfx` + `.music` DSL)
  - Zero external WAV/MP3 asset bloat; instant compilation, 4-byte stereo aligned 44.1 kHz PCM in memory.
- **Sound Effects (`procedural-sfx`)**:
  - `laser_aegis.sfx`: Crisp high-frequency FM laser chirp.
  - `dash_surge.sfx`: Filtered white-noise whoosh with resonant pitch sweep.
  - `xp_pickup.sfx`: Harmonically pleasing pentatonic chime on aether orb collect.
  - `level_up.sfx`: 4-note ascending arpeggio fanfare with polyphonic sustain.
  - `cyber_explosion.sfx`: Punchy low-end sub-bass transient with LFSR digital noise crunch.
  - `hit_impact.sfx`: Short gated square-wave click for snappy combat feedback.
- **Background Music (`procedural-music`)**:
  - **Theme**: *Hymn to Athena 2.0* - Fast-paced Greek synthwave (Tempo 136 BPM, D Dorian mode / Phrygian touches, pulsing bassline, FM brass leads, driving 16th-note arpeggiators).
  - **Synthesis Engine**: Pure-Go SEGA Mega Drive YM2612 4-Op FM + SN76489 PSG synthesis.

---

## 6. Scene Flow & UI / HUD Layout
- **Scene State Machine**:
  - `SceneTitle`: Title screen with pulsing "Press Space to Enter Olympus 2.0".
  - `SceneGame`: Active gameplay, spawning waves, particle updates, combat.
  - `SceneLevelUp`: Semi-transparent freeze-frame overlay showing 3 holographic upgrade cards with keyboard (1/2/3 or arrows) and mouse/touch selection.
  - `SceneGameOver`: Death screen with survival time, total kills, and quick retry button (`R`).
  - `SceneVictory`: Triumphant celebration of restoring Olympus 2.0, boss defeat stats.
- **HUD Layout (`640x360`)**:
  - **Top-Left**: Athena AI Portrait, HP Bar (Green/Cyan with white damage ghosting), Shield gauge.
  - **Top-Center**: Glowing XP bar spanning the top edge with Level indicator (`LVL 1`).
  - **Top-Right**: Mission Timer (`00:00 / 05:00`), Kill Counter (Skull icon + number).
  - **Bottom-Center**: Dash Cooldown indicator.

---

## 7. Ebitengine Architecture & Studio Workflow (`ebitengineer` + `game-dev-studio`)
- **Architecture**:
  - Standard Ebitengine v2 layout (`cmd/game/main.go`, `internal/game`, `internal/entities`, `internal/audio`, `internal/render`).
  - Pre-allocated entity slices and zero-allocation draw loops to maintain solid 60 FPS under 300+ simultaneous on-screen entities.
  - WebAssembly (`GOOS=js GOARCH=wasm`) and companion Go web server (`cmd/server/main.go`) with `wasm_exec.js` and responsive `index.html`.
- **Game Dev Studio Integration**:
  - Ready for immediate playtesting in `[Game Dev Studio](sidecar://agy-game-studio/studio/playtest)` with DVR rewind and visual annotations.
