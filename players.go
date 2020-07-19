package main

// Players is a simple type alias on an array of Player
type Players []*Player

func (p *Players) getRandom(d Dice) *Player {
	index := d.roll(len(*p))
	return (*p)[index]
}

func (p *Players) exist(playerName string) bool {
	for _, candidate := range *p {
		if candidate.name == playerName {
			return true
		}
	}
	return false
}

func (p *Players) listNames() string {
	if len(*p) == 1 {
		return (*p)[0].name
	}
	names := (*p)[0].name
	for i := 1; i < len(*p); i++ {
		if len(*p)-1 == i {
			names += " and "
		} else {
			names += ", "
		}
		names += (*p)[i].name
	}
	return names
}

func newPlayers() Players {
	return []*Player{}
}
