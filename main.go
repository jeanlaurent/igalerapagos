package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	game := newGame(12) // put a boolean to whether or not we need to output something
	for !game.isOver() {
		game.runDay()
		fmt.Print("Press any key to continue")
		_, _ = reader.ReadString('\n')
	}
	if game.successfulEscape {
		fmt.Println("You win, you manage to save", len(game.players), "persons in ", game.dayCount, "days.")
	} else {
		fmt.Println("You lose, everybody died after ", game.dayCount, "days.")
	}
}
