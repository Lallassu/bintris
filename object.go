package main

type Object interface {
	Draw(dt float64)
	Delete()
	GetID() int
	Update(float64)
	GetX() float64
	GetY() float64
}
