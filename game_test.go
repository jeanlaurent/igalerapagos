package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type DiceStub struct {
	values []int
	next   int
}

func (d *DiceStub) roll(max int) int {
	result := d.values[d.next]
	d.next++
	return result
}

func TestGameAsXPlayers(t *testing.T) {
	game := newGame(12)
	assert.Equal(t, 12, len(game.players))
}

func TestStartFoodIsGreaterThanPlayerNumber(t *testing.T) {
	assert.LessOrEqual(t, 12, newGame(12).foodStock)
}

func TestDayCountStartAt0(t *testing.T) {
	assert.Equal(t, 0, newGame(12).dayCount)
}

func TestWhenDayPassDayCountProgress(t *testing.T) {
	game := newGame(12)
	game.runDay()
	assert.Equal(t, 1, game.dayCount)
	game.runDay()
	assert.Equal(t, 2, game.dayCount)
	game.runDay()
	assert.Equal(t, 3, game.dayCount)
}

func TestWhenDayPassDayFoodGetsDown(t *testing.T) {
	game := newGame(12)
	game.foodStock = 14
	game.runLunchPhase()
	assert.Equal(t, 2, game.foodStock)
}

func TestWhenFoodStockCantBeNegative(t *testing.T) {
	game := newGame(12)
	game.foodStock = 0
	game.runLunchPhase()
	assert.Equal(t, 0, game.foodStock)
}

func TestEverybodyGotToBeHungry(t *testing.T) {
	playerCount := 12
	game := newGame(playerCount)
	game.foodStock = 0
	game.runLunchPhase()
	for _, player := range game.players {
		assert.Equal(t, hungry, player.foodStatus)
	}
}

func TestFind3Players(t *testing.T) {
	playerCount := 12
	game := newGame(playerCount)
	players := game.findXPlayers(3)
	assert.Equal(t, 3, len(players))
}

func TestFind200Players(t *testing.T) {
	playerCount := 12
	game := newGame(playerCount)
	players := game.findXPlayers(200)
	assert.Equal(t, playerCount, len(players))
}

func TestIsGameOver(t *testing.T) {
	game := newGame(12)
	assert.False(t, game.isOver())
	for _, player := range game.players {
		player.alive = false
	}
	assert.True(t, game.isOver())
}

func TestIsGameOverBySuccesfullEscape(t *testing.T) {
	game := newGame(12)
	game.succesfullEscape = true
	assert.True(t, game.isOver())
}

func TestRemovedDeadPlayerNone(t *testing.T) {
	game := newGame(12)
	game.removeDeadPlayers()
	assert.Equal(t, 12, len(game.players))
}

func TestRemovedDeadPlayerThree(t *testing.T) {
	game := newGame(12)
	game.players[0].alive = false
	game.players[1].alive = false
	game.players[2].alive = false
	game.removeDeadPlayers()
	assert.Equal(t, 9, len(game.players))
}

func TestRemovedDeadPlayerAll(t *testing.T) {
	game := newGame(12)
	for _, player := range game.players {
		player.alive = false
	}
	game.removeDeadPlayers()
	assert.Equal(t, 0, len(game.players))
}

func TestKillPlayerOn(t *testing.T) {
	game := newGame(4)
	game.dice = &DiceStub{[]int{10, 40, 60, 90}, 0}
	game.killPlayersOn(50)
	assert.False(t, game.players[0].alive)
	assert.False(t, game.players[1].alive)
	assert.True(t, game.players[2].alive)
	assert.True(t, game.players[3].alive)
}
