package game

import (
	"encoding/json"
	"os"
	"slices"
)

var (
	// use this for colors and special asciis
	ClearAscii = "\033[2J"
	green      = "\033[32m"
	yellow     = "\033[33m"
	blue       = "\033[34m"
	red        = "\033[31m"
	reset      = "\033[0m"
	white      = "\033[37m"
)

var (
	DefaultState = State{
		P1State: PlayerState{
			Position: Point{40, 10},
			Action:   Idle,
		},
		P2State: PlayerState{
			Position: Point{30, 10},
			Action:   Idle,
		},
	}
)

func NewGame(dataPath string, char1Name string, char2Name string) (Game, error) {
	// composing the game
	cfgFilePath := dataPath + "/config.json"
	data, err := os.ReadFile(cfgFilePath)

	if err != nil {
		return Game{}, err
	}

	var c Config
	err = json.Unmarshal(data, &c)
	if err != nil {
		return Game{}, err
	}

	metaDataPath := dataPath + "/metaData.json"
	data, err = os.ReadFile(metaDataPath)

	if err != nil {
		return Game{}, err
	}
	var metaData MetaData
	err = json.Unmarshal(data, &metaData)
	if err != nil {
		return Game{}, err
	}

	if !slices.Contains(metaData.AvailableCharecters, char1Name) || !slices.Contains(metaData.AvailableCharecters, char2Name) {
		return Game{}, ErrCharFilesMissing
	}

	ap := NewAssets()
	err = ap.LoadCharecter(dataPath+"/assets/charecters/"+char1Name+".json", p1ID)
	if err != nil {
		return Game{}, err
	}
	err = ap.LoadCharecter(dataPath+"/assets/charecters/"+char2Name+".json", p2ID)
	if err != nil {
		return Game{}, err
	}
	DefaultCharecter, err := ap.GetChar(p1ID)
	if err != nil {
		return Game{}, err
	}
	err = ap.LoadAnimationSet(DefaultCharecter.getAnimationSetPath(dataPath), p1ID)
	if err != nil {
		return Game{}, err
	}
	err = ap.LoadAnimationSet(DefaultCharecter.getAnimationSetPath(dataPath), p2ID)
	if err != nil {
		return Game{}, err
	}
	currState := DefaultState
	animS, err := ap.GetAnimData(p1ID, Idle)
	animS2, err := ap.GetAnimData(p2ID, Idle)
	if err != nil {
		return Game{}, err
	}
	currState.P1State.AnimState = animS
	currState.P1State.OwnerID = p1ID
	currState.P2State.AnimState = animS2
	currState.P2State.OwnerID = p2ID
	g := Game{
		Config:         c,
		State:          currState,
		AssetsProvider: ap,
	}

	return g, nil

}

//
// state should tell the renderer everything it needs to render the scene

func (game Game) GenerateNextState(s State, inputs Input) State {
	// NewState, err := game.GetBindFunc(inputs, Player1)(s)
	// use the new way of handling inputs and events , for now simply we just return the same state
	return s
}
