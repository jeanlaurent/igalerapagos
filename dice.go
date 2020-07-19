package main

import "math/rand"

// Dice allow you to roll dices
type Dice interface {
	roll(max int) int
}

type dice struct{}

func (d *dice) roll(max int) int {
	return rand.Intn(max)
}
