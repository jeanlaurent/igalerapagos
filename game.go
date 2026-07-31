package main

import (
	"fmt"
	"io"
	"os"
)

// Game hold the game data context
type Game struct {
	players          Players   // All alive players current status
	foodStock        int       // The food stock level
	woodStock        int       // The wood stock
	dayCount         int       // Number of days elapsed since the start
	weather          weather   // Current Weather
	campLevel        int       // the current camp level
	successfulEscape bool      // Did the team successfully escape the island?
	dice             Dice      // A dice to get random numbers from
	writer           io.Writer // Destination for all game output
}

func (g *Game) findXPlayers(max int) Players {
	found := newPlayers()
	if max > len(g.players) {
		return g.players
	}
	for len(found) < max {
		candidate := g.players.getRandom(g.dice)
		if !(found.exist(candidate.name)) {
			found = append(found, candidate)
		}
	}
	return found
}

func (g *Game) decrementFood() int {
	numberOfPlayerMissingALunch := 0
	if g.foodStock < len(g.players) {
		numberOfPlayerMissingALunch = len(g.players) - g.foodStock
		g.foodStock = 0
	} else {
		g.foodStock -= len(g.players)
	}
	return numberOfPlayerMissingALunch
}

func (g *Game) runDay() {
	g.startPhase()
	fmt.Fprintln(g.writer)
	skipActionPhase := g.escapePhase()
	fmt.Fprintln(g.writer)
	if !g.successfulEscape {
		if !skipActionPhase {
			g.runActionPhase()
			fmt.Fprintln(g.writer)
		}
		g.runLunchPhase()
		fmt.Fprintln(g.writer)
		g.runCampfirePhase()
		fmt.Fprintln(g.writer)
	}
	g.endPhase()
}

func (g *Game) startPhase() {
	g.dayCount++
	fmt.Fprintln(g.writer, "Start of day", g.dayCount, ".")
	g.weather.changeWeather(g.dice)
	fmt.Fprintln(g.writer, "\tIt is going to be a", g.weather.state, "day.")
}

func (g *Game) escapePhase() bool {
	if g.campLevel < 25 {
		return false
	}
	if g.dice.roll(100)+g.campLevel < 50 {
		fmt.Fprintln(g.writer, "\tThe group decide that even if the raft is ready, the condition for escaping are not met.")
		return false
	}
	fmt.Fprintln(g.writer, "\tThe group decide to try to use the raft and escape the island.")
	return g.escape()
}

func (g *Game) escape() bool {
	roll := g.dice.roll(50) + g.campLevel - 25 + len(g.players)
	if roll < 10 {
		g.campLevel -= 20
		drown := g.killPlayersOn(20)
		fmt.Fprintln(g.writer, "\tThe raft is crushed almost immediately on some rocks by huge waves.")
		if len(drown) > 0 {
			fmt.Fprintln(g.writer, "A huge wave fall on the raft", drown.listNames(), "have been thrown into the sea, we will never see them back.")
		}
		fmt.Fprintln(g.writer, "\tThe group is back to the island with a rather broken raft.")
		fmt.Fprintln(g.writer, "\tWe still have time to get some work done.")
		g.removeDeadPlayers()
		return false
	} else if roll < 30 {
		g.campLevel -= 10
		fmt.Fprintln(g.writer, "\tAfter some long hours, the raft breaks apart.")
		fmt.Fprintln(g.writer, "\tThe group managed to get back on the origin island with a slightly damaged raft.")
		fmt.Fprintln(g.writer, "\tWe still have time to get some work done.")
		return false
	} else if roll < 40 {
		drown := g.killPlayersOn(10)
		fmt.Fprintln(g.writer, "\tAs the group enters open water. A huge wave sweeps the raft.")
		if len(drown) > 0 {
			fmt.Fprintln(g.writer, "\t", drown.listNames(), "have been swept by the wave, and died.")
		} else {
			fmt.Fprintln(g.writer, "\tbut everyone survived")
		}
		g.removeDeadPlayers()
		return g.escape()
	} else if roll < 50 {
		fmt.Fprintln(g.writer, "\tThe raft is within range of a rescue ship.")
		roll = g.dice.roll(100)
		if roll < 30 {
			fmt.Fprintln(g.writer, "\tThe passengers waves furiously but they fail to be noticed, and are back at the camp.")
			fmt.Fprintln(g.writer, "\tThis took so long, that they reach the island by night, and can't work today")
		} else {
			fmt.Fprintln(g.writer, "\tThe passengers waves furiously. After an exhausting session of shouting and waving, the boat notice them. The team has been rescued successfully.")
			g.successfulEscape = true
		}
		return true
	} else {
		fmt.Fprintln(g.writer, "\tAlmost when all hopes are lost. The group reach another island with some civilization and are saved.")
		g.successfulEscape = true
		return true
	}

}

func (g *Game) killPlayersOn(percentage int) Players {
	killed := newPlayers()
	for _, player := range g.players {
		roll := g.dice.roll(100)
		if roll <= percentage {
			player.alive = false
			killed = append(killed, player)
		}
	}
	return killed
}

func (g *Game) runActionPhase() { // create interface for phase, move phase into dedicated func
	starvingGroup := newPlayers()
	woodGroup := newPlayers()
	foodGroup := newPlayers()
	campGroup := newPlayers()
	// dispatch players randomly in groups
	for _, player := range g.players {
		if player.foodStatus == starving {
			starvingGroup = append(starvingGroup, player)
			continue
		}
		actionRoll := g.dice.roll(100)
		if actionRoll < 33 {
			woodGroup = append(woodGroup, player)
		} else if actionRoll >= 33 && actionRoll <= 66 {
			foodGroup = append(foodGroup, player)
		} else {
			campGroup = append(campGroup, player)
		}
	}
	// Perform actions
	if len(starvingGroup) > 0 {
		fmt.Fprintln(g.writer, "\t", starvingGroup.listNames(), "are starving, and are too weak to work on anything.")
	}
	if len(woodGroup) > 0 {
		woodGathered := g.dice.roll(6) + 1 + len(woodGroup) //+ g.woodGatheringBonus()
		if woodGathered < 0 {
			woodGathered = 0
		}
		g.woodStock += woodGathered
		fmt.Fprintln(g.writer, "\t", "A group made of", woodGroup.listNames(), "gathered", woodGathered, "wood pieces.")
	} else {
		fmt.Fprintln(g.writer, "\t", "No one wanted to pickup wood today.")
	}
	if len(foodGroup) > 0 {
		foodGathered := g.dice.roll(6) + 1 + g.weather.foodGatheringBonus() + len(foodGroup)
		if foodGathered < 0 {
			foodGathered = 0
		}
		g.foodStock += foodGathered
		fmt.Fprintln(g.writer, "\t", "A group made of", foodGroup.listNames(), "gathered", foodGathered, "fruits and other food.")
	} else {
		fmt.Fprintln(g.writer, "\t", "No one wanted to pickup food today.")
	}
	if len(campGroup) > 0 {
		campImprovement := g.dice.roll(6) + 1 + len(campGroup)
		if g.woodStock < campImprovement {
			campImprovement = g.woodStock
		}
		g.woodStock -= campImprovement
		g.campLevel += campImprovement
		fmt.Fprintln(g.writer, "\t", "A group made of", campGroup.listNames(), "worked the camp they raised the camp level to", g.campLevel, ". They used wood for that, there are", g.woodStock, "wood left")
	} else {
		fmt.Fprintln(g.writer, "\t", "No one wanted to work the camp today.")
	}
}

func (g *Game) runLunchPhase() {
	// Lunch time
	playerDiedOfHunger := newPlayers()
	playerEatingLunch := newPlayers()

	numberOfPlayerMissingALunch := g.decrementFood()
	playerMissingLunch := g.findXPlayers(numberOfPlayerMissingALunch)
	for _, player := range g.players {
		if playerMissingLunch.exist(player.name) {
			player.lunch(false)
			if !player.alive {
				playerDiedOfHunger = append(playerDiedOfHunger, player)
			}
		} else {
			player.lunch(true)
			playerEatingLunch = append(playerEatingLunch, player)
		}
	}
	if len(playerMissingLunch) == 0 {
		fmt.Fprintln(g.writer, "\t", "Everybody ate some food today.")
	} else {
		if len(playerEatingLunch) > 0 {
			fmt.Fprintln(g.writer, "\t", playerEatingLunch.listNames(), "managed to get some food.")
		}
		fmt.Fprintln(g.writer, "\t", playerMissingLunch.listNames(), "did not eat tonight")
		if len(playerDiedOfHunger) > 0 {
			fmt.Fprintln(g.writer, "\t", playerDiedOfHunger.listNames(), "died of hunger")
		}
	}
	g.removeDeadPlayers()
}

func (g *Game) runCampfirePhase() {
	if len(g.players) == 0 {
		fmt.Fprintln(g.writer, "\tNo one is alive, so no campfire tonight")
		return
	}
	firepower := g.dice.roll(6) + 1
	if g.woodStock < firepower {
		fmt.Fprintln(g.writer, "\tDuring the campfire tonight the group start to burn some", g.woodStock, "logs")
		firepower -= g.woodStock
		g.woodStock = 0
		g.campLevel -= firepower
		fmt.Fprintln(g.writer, "\tSince they were not enough logs to keep the fire burning, we broke part of the camp.")
		if g.campLevel < 0 {
			g.campLevel = 0
		}
	} else {
		fmt.Fprintln(g.writer, "\tDuring the campfire tonight the group burns", firepower, "logs")
		g.woodStock -= firepower
	}
}

func (g *Game) removeDeadPlayers() {
	alivePlayers := newPlayers()
	for _, player := range g.players {
		if player.alive {
			alivePlayers = append(alivePlayers, player)
		}
	}
	g.players = alivePlayers
}

func (g *Game) endPhase() {
	fmt.Fprintln(g.writer, "At the end of day", g.dayCount, "there are", len(g.players), "persons alive.", g.foodStock, "meals are left, we got ", g.woodStock, "wood log left. The camp level is at", g.campLevel)
	fmt.Fprintln(g.writer)
	fmt.Fprintln(g.writer, "====================")
}

func (g *Game) isOver() bool {
	if g.successfulEscape {
		return true
	}
	for _, player := range g.players {
		if player.alive {
			return false
		}
	}
	return true
}

func newGame(playerCount int) Game {
	initialPlayers := newPlayers()
	game := Game{players: initialPlayers}
	game.dice = newDice()
	game.writer = os.Stdout
	for i := 0; i < playerCount; i++ {
		newplayer := newPlayer(i)
		game.players = append(game.players, &newplayer)
	}
	game.foodStock = int(float32(playerCount) * 1.5)
	game.woodStock = playerCount / 2
	game.campLevel = 0
	game.dayCount = 0
	game.successfulEscape = false
	game.weather = newWeather()
	return game
}
