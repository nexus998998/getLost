package game

import "fmt"

func MakePlayerActionID(a Action, p Player) string {
	return fmt.Sprintf("func%s%v", a.String(), p)
}

// we need a system to load in the events
// we need to do some validation of some sort right ?
// hmmmm

// so we need a function to map the strings to use cases
// what do we actually need from the asset provider ? we need to inject dependencies , mainly configs
// so we want to make it a function that returns an error
// we need a system to load in the configs first into this , or maybe we put them together ? idk what do you think
// we can also make a type big config thingy to make sure everything is configured
func (AS *AssetsChunks) LoadFuncs() error {

}

func MakePlayerActionEvent(a Action, p Player) event {
	return event{
		Catagory: ActivePlayerAction,
		FuncID:   MakePlayerActionID(a, p),
	}

}
