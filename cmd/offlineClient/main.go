package main

import (
	"getLost/internals/game"
	"log"
	"time"
)

var (
	configFilePath = "./config.json"
)

func main() {
	chars := game.Charecters{
		Char1: game.DefaultCharecter,
		Char2: game.DefaultCharecter,
	}
	g, err := game.NewGame(configFilePath, chars)
	if err != nil {
		log.Fatal(err)
	}
	anim := chars.Char1.AnimationSet[game.Walking]
	g.State.P1State.CurrentAnimation = &anim
	// game loop
	// what's the plan , init the game , make a rendering loop and yeah that's it
	for {
		position := &g.State.P1State.Position
		f := g.MakeFrame(g.State)
		f.RenderFrame()
		position.X += 1
		if position.X >= g.Config.Resulotion.Width {
			position.X = 0
		}
		time.Sleep(time.Millisecond * 30)
	}
}
