package main

var names = []string{"Alice", "Bob", "Chloe", "Daphne", "Eric", "Franck", "George", "Hans", "Ines", "Jules", "Klaus", "Lawrence"}

const fed = 0
const hungry = 1
const starving = 2

// Player represent a person playing in the game
type Player struct {
	name       string // The player name, used as identifier
	foodStatus int    // Wether this player is fed, hungry or starving
	alive      bool   // wether this player is dead or not
}

func (p *Player) lunch(hasLunch bool) {
	if hasLunch {
		p.foodStatus = fed
	} else {
		p.foodStatus++
	}
	if p.foodStatus > starving {
		p.alive = false
	}
}

func (p *Player) getFoodStatusAsString() string {
	switch p.foodStatus {
	case fed:
		return "fed"
	case hungry:
		return "hungry"
	case starving:
		return "starving"
	}
	return "unknown"
}

func newPlayer(rank int) Player {
	return Player{name: names[rank], foodStatus: fed, alive: true}
}
