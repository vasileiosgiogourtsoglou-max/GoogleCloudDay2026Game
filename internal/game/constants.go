package game

const (
	VirtualWidth  = 640
	VirtualHeight = 360
	TargetTPS      = 60
)

type StateType int

const (
	StateTitle StateType = iota
	StateGameplay
	StateLevelUp
	StateGameOver
	StateVictory
)
