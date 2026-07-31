package main

import (
	"bytes"
	"io"
	"math/rand/v2"
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
	game.writer = io.Discard
	game.runDay()
	assert.Equal(t, 1, game.dayCount)
	game.runDay()
	assert.Equal(t, 2, game.dayCount)
	game.runDay()
	assert.Equal(t, 3, game.dayCount)
}

func TestWhenDayPassDayFoodGetsDown(t *testing.T) {
	game := newGame(12)
	game.writer = io.Discard
	game.foodStock = 14
	game.runLunchPhase()
	assert.Equal(t, 2, game.foodStock)
}

func TestWhenFoodStockCantBeNegative(t *testing.T) {
	game := newGame(12)
	game.writer = io.Discard
	game.foodStock = 0
	game.runLunchPhase()
	assert.Equal(t, 0, game.foodStock)
}

func TestEverybodyGotToBeHungry(t *testing.T) {
	playerCount := 12
	game := newGame(playerCount)
	game.writer = io.Discard
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

func TestIsGameOverBySuccessfulEscape(t *testing.T) {
	game := newGame(12)
	game.successfulEscape = true
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

// --- Bug 1: campfire wood consumption ---

// TestCampfireConsumesWoodWhenPlentiful verifies that runCampfirePhase deducts
// the correct number of logs from woodStock when there is plenty of wood.
// Before the fix, firepower was zeroed before the subtraction, so wood was
// never consumed.
func TestCampfireConsumesWoodWhenPlentiful(t *testing.T) {
	game := newGame(4)
	game.writer = io.Discard
	game.woodStock = 100
	// dice.roll(6) returns 3, so firepower = 3 + 1 = 4
	game.dice = &DiceStub{[]int{3}, 0}
	game.runCampfirePhase()
	assert.Equal(t, 96, game.woodStock) // 100 - 4 logs burned
	assert.Equal(t, 0, game.campLevel)  // camp untouched when wood is plentiful
}

// --- Bug 2: action dispatch off-by-one ---

// testActionDispatch is a helper that creates a 1-player game, sets up the
// dice so the first roll (the action roll) returns actionRoll followed by 5
// (for the food-gather roll), and then asserts the player was placed in the
// food bucket (not the camp bucket).
func testActionDispatch(t *testing.T, actionRoll int) {
	t.Helper()
	game := newGame(1)
	game.writer = io.Discard
	// Use woodStock=100 so any accidental camp routing is clearly visible.
	game.woodStock = 100
	game.foodStock = 0
	game.campLevel = 0
	// Dice sequence: [actionRoll, 5]
	// actionRoll is consumed by the per-player dispatch loop.
	// 5 is consumed by foodGathered = roll(6)+1+bonus+1 when food group fires.
	game.dice = &DiceStub{[]int{actionRoll, 5}, 0}
	game.runActionPhase()
	// Assertions that prove routing went to FOOD, not to CAMP:
	assert.Equal(t, 100, game.woodStock, "roll %d should not consume wood (not camp)", actionRoll)
	assert.Equal(t, 0, game.campLevel, "roll %d should not raise campLevel (not camp)", actionRoll)
	assert.Greater(t, game.foodStock, 0, "roll %d should increase foodStock (food bucket)", actionRoll)
}

// TestActionPhaseRoll33GoesToFood verifies that a dice roll of exactly 33 is
// routed to the food bucket. Before the fix the middle condition was
// `> 33 && < 66`, so 33 fell through to the camp (else) bucket.
func TestActionPhaseRoll33GoesToFood(t *testing.T) {
	testActionDispatch(t, 33)
}

// TestActionPhaseRoll66GoesToFood verifies that a dice roll of exactly 66 is
// routed to the food bucket. Before the fix the middle condition excluded 66,
// so it fell through to the camp (else) bucket.
func TestActionPhaseRoll66GoesToFood(t *testing.T) {
	testActionDispatch(t, 66)
}

// --- io.Writer: output assertion ---

// TestGameOutputContainsDayHeader verifies that game output is routed through
// the injected io.Writer: running runDay should produce a line beginning with
// "Start of day 1".
func TestGameOutputContainsDayHeader(t *testing.T) {
	buf := &bytes.Buffer{}
	game := newGame(4)
	game.writer = buf
	game.dice = newDiceWithSource(rand.NewPCG(1, 0))
	game.runDay()
	assert.Contains(t, buf.String(), "Start of day 1")
}

// --- RNG injection: determinism ---

// TestGameIsDeterministicWithSameSeed verifies that two Game instances
// initialised with the same RNG seed and the same writer produce identical
// output for a runDay call, proving that global rand state is not used.
func TestGameIsDeterministicWithSameSeed(t *testing.T) {
	buf1 := &bytes.Buffer{}
	game1 := newGame(4)
	game1.writer = buf1
	game1.dice = newDiceWithSource(rand.NewPCG(42, 0))
	game1.runDay()

	buf2 := &bytes.Buffer{}
	game2 := newGame(4)
	game2.writer = buf2
	game2.dice = newDiceWithSource(rand.NewPCG(42, 0))
	game2.runDay()

	assert.Equal(t, buf1.String(), buf2.String(), "games with the same seed must produce identical output")
}
