package main

type Object interface {
	Draw(dt float64)
	GetObjectType() ObjectType
	Delete()
	GetID() int
	Update(float64)
	GetX() float64
	GetY() float64
}

type ObjectType int

const (
	ObjectTypeTileSet = iota + 1
	ObjectTypeSprite
)
