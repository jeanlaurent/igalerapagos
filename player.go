package main

var names = []string{"Alice", "Bob", "Chloe", "Daphne", "Eric", "Franck", "George", "Hans", "Ines", "Jules", "Klaus", "Lawrence"}

// FoodStatus represents a player's hunger state.
type FoodStatus int

const (
	fed FoodStatus = iota
	hungry
	starving
)

// String implements fmt.Stringer for readable output and debugging.
func (f FoodStatus) String() string {
	switch f {
	case fed:
		return "fed"
	case hungry:
		return "hungry"
	case starving:
		return "starving"
	}
	return "unknown"
}

// Player represent a person playing in the game
type Player struct {
	name       string     // The player name, used as identifier
	foodStatus FoodStatus // Whether this player is fed, hungry or starving
	alive      bool       // whether this player is dead or not
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

func newPlayer(rank int) Player {
	return Player{name: names[rank], foodStatus: fed, alive: true}
}
