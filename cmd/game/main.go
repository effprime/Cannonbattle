package main

import "github.com/effprime/cannonbattle/pkg/cannonbattle"

func main() {
	err := cannonbattle.Run()
	if err != nil {
		panic(err)
	}
}
