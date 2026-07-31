package main

import (
	"math/rand/v2"
	"time"
)

// Dice allow you to roll dices
type Dice interface {
	roll(max int) int
}

type dice struct {
	rng *rand.Rand
}

// newDice returns a live time-seeded dice.
func newDice() *dice {
	return &dice{rng: rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), 0))}
}

// newDiceWithSource returns a dice backed by the given source, useful for
// deterministic tests.
func newDiceWithSource(source rand.Source) *dice {
	return &dice{rng: rand.New(source)}
}

func (d *dice) roll(max int) int {
	return d.rng.IntN(max)
}
