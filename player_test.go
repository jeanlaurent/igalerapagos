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

func TestFoodStatusAsString(t *testing.T) {
	player := newPlayer(1)
	assert.Equal(t, player.getFoodStatusAsString(), "fed")
	player.foodStatus = hungry
	assert.Equal(t, player.getFoodStatusAsString(), "hungry")
	player.foodStatus = starving
	assert.Equal(t, player.getFoodStatusAsString(), "starving")
	player.foodStatus = -1
	assert.Equal(t, player.getFoodStatusAsString(), "unknown")
}
