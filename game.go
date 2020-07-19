package main

import "fmt"

// Game hold the game data context
type Game struct {
	players          Players // All alive players current status
	foodStock        int     // The food stock level
	woodStock        int     // The wood stock
	dayCount         int     // Number of days elapsed since the start
	weather          weather // Current Weather
	campLevel        int     // the current camp level
	succesfullEscape bool    // Did the team successfully escaped the island ?
	dice             Dice    // A dice to get random numbers from
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
	fmt.Println()
	skipActionPhase := g.escapePhase()
	fmt.Println()
	if !g.succesfullEscape {
		if !skipActionPhase {
			g.runActionPhase()
			fmt.Println()
		}
		g.runLunchPhase()
		fmt.Println()
		g.runCampfirePhase()
		fmt.Println()
	}
	g.endPhase()
}

func (g *Game) startPhase() {
	g.dayCount++
	fmt.Println("Start of day", g.dayCount, ".")
	g.weather.changeWeather(g.dice)
	fmt.Println("\tIt is going to be a", g.weather.weatherAsString(), "day.")
}

func (g *Game) escapePhase() bool {
	if g.campLevel < 25 {
		return false
	}
	if g.dice.roll(100)+g.campLevel < 50 {
		fmt.Println("\tThe group decide that even if the raft is ready, the condition for escaping are not met.")
		return false
	}
	fmt.Println("\tThe group decide to try to use the raft and escape the island.")
	return g.escape()
}

func (g *Game) escape() bool {
	roll := g.dice.roll(50) + g.campLevel - 25 + len(g.players)
	if roll < 10 {
		g.campLevel -= 20
		drown := g.killPlayersOn(20)
		fmt.Println("\tThe raft is crushed almost immediatly on some rocks by huge waves.")
		if len(drown) > 0 {
			fmt.Println("A huge wave fall on the raft", drown.listNames(), "have been thrown into the sea, we will never see them back.")
		}
		fmt.Println("\tThe group is back to the island with a rather broken raft.")
		fmt.Println("\tWe still have time to get some work done.")
		g.removeDeadPlayers()
		return false
	} else if roll < 30 {
		g.campLevel -= 10
		fmt.Println("\tAfter some long hours, the raft break apparts.")
		fmt.Println("\tThe group managed to get back on the origin island with a slightly damaged raft.")
		fmt.Println("\tWe still have time to get some work done.")
		return false
	} else if roll < 40 {
		drown := g.killPlayersOn(10)
		fmt.Println("\tAs the group enters open water. A huge wave swipe the raft.")
		if len(drown) > 0 {
			fmt.Println("\t", drown.listNames(), "have been swiped by the wave, and died.")
		} else {
			fmt.Println("\tbut everyone survived")
		}
		g.removeDeadPlayers()
		return g.escape()
	} else if roll < 50 {
		fmt.Println("\tThe raft is within range of a rescue ship.")
		roll = g.dice.roll(100)
		if roll < 30 {
			fmt.Println("\tThe passengers waves furiously but they fail to be noticed, and are back at the camp.")
			fmt.Println("\tThis took so long, that they reach the island by night, and can't work today")
		} else {
			fmt.Println("\tThe passengers waves furiously. After an exhausting session of shouting and waving, the boat notice them. The team has been rescued succesfully.")
			g.succesfullEscape = true
		}
		return true
	} else {
		fmt.Println("\tAlmost when all hopes are lost. The group reach another island with some civilization and are saved.")
		g.succesfullEscape = true
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
		} else if actionRoll > 33 && actionRoll < 66 {
			foodGroup = append(foodGroup, player)
		} else {
			campGroup = append(campGroup, player)
		}
	}
	// Perform actions
	if len(starvingGroup) > 0 {
		fmt.Println("\t", starvingGroup.listNames(), "are starving, and are too weak to work on anything.")
	}
	if len(woodGroup) > 0 {
		woodGathered := g.dice.roll(6) + 1 + len(woodGroup) //+ g.woodGatherindBonus()
		if woodGathered < 0 {
			woodGathered = 0
		}
		g.woodStock += woodGathered
		fmt.Println("\t", "A group made of", woodGroup.listNames(), "gathered", woodGathered, "wood pieces.")
	} else {
		fmt.Println("\t", "No one wanted to pickup wood today.")
	}
	if len(foodGroup) > 0 {
		foodGathered := g.dice.roll(6) + 1 + g.weather.foodGatheringBonus() + len(foodGroup)
		if foodGathered < 0 {
			foodGathered = 0
		}
		g.foodStock += foodGathered
		fmt.Println("\t", "A group made of", foodGroup.listNames(), "gathered", foodGathered, "fruits and other food.")
	} else {
		fmt.Println("\t", "No one wanted to pickup food today.")
	}
	if len(campGroup) > 0 {
		campImprovement := g.dice.roll(6) + 1 + len(campGroup)
		if g.woodStock < campImprovement {
			campImprovement = g.woodStock
		}
		g.woodStock -= campImprovement
		g.campLevel += campImprovement
		fmt.Println("\t", "A group made of", campGroup.listNames(), "worked the camp they raised the camp level to", g.campLevel, ". They used wood for that, there are", g.woodStock, "wood left")
	} else {
		fmt.Println("\t", "No one wanted to work the camp today.")
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
		fmt.Println("\t", "Everybody ate some food today.")
	} else {
		if len(playerEatingLunch) > 0 {
			fmt.Println("\t", playerEatingLunch.listNames(), "managed to get some food.")
		}
		fmt.Println("\t", playerMissingLunch.listNames(), "did not eat tonight")
		if len(playerDiedOfHunger) > 0 {
			fmt.Println("\t", playerDiedOfHunger.listNames(), "died of hunger")
		}
	}
	g.removeDeadPlayers()
}

func (g *Game) runCampfirePhase() {
	if len(g.players) == 0 {
		fmt.Println("\tnoone is alive, so no campfire tonight")
		return
	}
	firepower := g.dice.roll(6) + 1
	if g.woodStock < firepower {
		fmt.Println("\tDuring the campfire tonight the group start to burn some", g.woodStock, "logs")
		firepower -= g.woodStock
		g.woodStock = 0
		g.campLevel -= firepower
		fmt.Println("\tSince they were not enough logs to keep the fire burning, we broke part of the camp.")
		if g.campLevel < 0 {
			g.campLevel = 0
		}
	} else {
		fmt.Println("\tDuring the campfire tonight the group burns", firepower, "logs")
		firepower = 0
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
	fmt.Println("At the end of day", g.dayCount, "there are", len(g.players), "persons alive.", g.foodStock, "meals are left, we got ", g.woodStock, "wood log left. The camp level is at", g.campLevel)
	fmt.Println()
	fmt.Println("====================")
}

func (g *Game) isOver() bool {
	if g.succesfullEscape {
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
	game.dice = &dice{}
	for i := 0; i < playerCount; i++ {
		newplayer := newPlayer(i)
		game.players = append(game.players, &newplayer)
	}
	game.foodStock = int(float32(playerCount) * 1.5)
	game.woodStock = playerCount / 2
	game.campLevel = 0
	game.dayCount = 0
	game.succesfullEscape = false
	game.weather = newWeather()
	return game
}
