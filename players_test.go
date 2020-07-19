package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLength(t *testing.T) {
	players := newPlayers()
	assert.Equal(t, 0, len(players))
}

func TestAddGet(t *testing.T) {
	players := newPlayers()
	player := Player{name: "foo"}
	players = append(players, &player)
	assert.Equal(t, 1, len(players))
	anotherPlayer := players[0]
	assert.Equal(t, "foo", anotherPlayer.name)
}

func TestExist(t *testing.T) {
	players := newPlayers()
	players = append(players, &Player{name: "foo"})
	players = append(players, &Player{name: "foo2"})
	players = append(players, &Player{name: "foo3"})
	assert.True(t, players.exist("foo"))
	assert.False(t, players.exist("bar"))
}

func TestListNames(t *testing.T) {
	players := newPlayers()
	players = append(players, &Player{name: "foo"})
	players = append(players, &Player{name: "bar"})
	players = append(players, &Player{name: "qix"})
	assert.Equal(t, "foo, bar and qix", players.listNames())
}

func TestListSingleName(t *testing.T) {
	players := newPlayers()
	players = append(players, &Player{name: "foo"})
	assert.Equal(t, "foo", players.listNames())
}

func TestList2Names(t *testing.T) {
	players := newPlayers()
	players = append(players, &Player{name: "foo"})
	players = append(players, &Player{name: "bar"})
	assert.Equal(t, "foo and bar", players.listNames())
}
