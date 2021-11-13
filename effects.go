package main

type Effect int

// Not using iota as it's easier to see values as we
// need to use them raw in the shaders.
const (
	EffectNone      = 0
	EffectMetaballs = 1
	EffectTileTop   = 2
	EffectStats     = 3
	EffectNumber    = 4
	EffectGameOver  = 5
	EffectGameOver1 = 6
	EffectGameOver2 = 7
)
