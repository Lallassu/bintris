package main

type Effect float32

// Not using iota as it's easier to see values as we
// need to use them raw in the shaders.
const (
	EffectNone       = 0
	EffectMetaballs  = 1
	EffectTileTop    = 2
	EffectStats      = 3
	EffectNumber     = 4
	EffectGameOver   = 5
	EffectGameOver1  = 6
	EffectGameOver2  = 7
	EffectStatsBlink = 8
	EffectBg         = 9
	EffectMenu       = 10
	EffectMenuStill  = 11
)
