package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlayerHasAName(t *testing.T) {
	assert.NotEmpty(t, newPlayer(1).name)
}

func TestPlayerisAlives(t *testing.T) {
	assert.True(t, newPlayer(1).alive)
}

func TestPlayerisFed(t *testing.T) {
	assert.Equal(t, newPlayer(1).foodStatus, fed)
}

func TestPlayerEatNormaly(t *testing.T) {
	player := newPlayer(1)
	player.lunch(true)
	assert.Equal(t, player.foodStatus, fed)
}

func TestPlayerMissALunch(t *testing.T) {
	player := newPlayer(1)
	player.lunch(false)
	assert.Equal(t, player.foodStatus, hungry)
}

func TestPlayerMiss2Lunches(t *testing.T) {
	player := newPlayer(1)
	player.lunch(false)
	player.lunch(false)
	assert.Equal(t, player.foodStatus, starving)
}

func TestPlayerMiss3Lunches(t *testing.T) {
	player := newPlayer(1)
	player.lunch(false)
	player.lunch(false)
	player.lunch(false)
	assert.False(t, player.alive)
}

func TestPlayeAfterMissing2LunchesEatAgain(t *testing.T) {
	player := newPlayer(1)
	player.lunch(false)
	player.lunch(false)
	player.lunch(true)
	assert.Equal(t, player.foodStatus, fed)
}

func TestFoodStatusString(t *testing.T) {
	player := newPlayer(1)
	assert.Equal(t, "fed", player.foodStatus.String())
	player.foodStatus = hungry
	assert.Equal(t, "hungry", player.foodStatus.String())
	player.foodStatus = starving
	assert.Equal(t, "starving", player.foodStatus.String())
	player.foodStatus = -1
	assert.Equal(t, "unknown", player.foodStatus.String())
}
