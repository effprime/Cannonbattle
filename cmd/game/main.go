package main

import "github.com/effprime/cannonbattle/pkg/game"

func main() {
	err := game.Run()
	if err != nil {
		panic(err)
	}
}
