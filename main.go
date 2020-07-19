package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	reader := bufio.NewReader(os.Stdin)

	game := newGame(12) // put a boolean to wether or not we need to output something
	for !game.isOver() {
		game.runDay()
		fmt.Print("Press any key to continue")
		reader.ReadString('\n')
	}
	if game.succesfullEscape {
		fmt.Println("You win, you manage to save", len(game.players), "persons in ", game.dayCount, "days.")
	} else {
		fmt.Println("You lose, everybody died after ", game.dayCount, "days.")
	}
}
