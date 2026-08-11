package game

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func (AS *AssetsChunks) LoadCharecter(filepath string, ID string) error {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	var char Charecter
	err = json.Unmarshal(file, &char)
	if err != nil {
		return err
	}

	AS.chunksMap[FormatCharKey(ID)] = char

	return nil

}

func (AS *AssetsChunks) LoadAnimationSet(filePath string, ID string) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	var cfgs []AnimationConfig
	err = json.Unmarshal(file, &cfgs)
	if err != nil {
		return err
	}
	for _, cfg := range cfgs {
		AS.loadAnimation(cfg, ID)
	}

	return nil
}

func parseAnimCfg(AC AnimationConfig) (AnimationData, panels) {
	return AnimationData{
		CurrentFrame:       1,
		RepeatingAnimation: AC.RepeatingAnimation,
	}, AC.Panels

}

func FormatFrameKey(ownerID string, action Action, frameNumber int) string {
	return "frame" + action.String() + ownerID + strconv.Itoa(frameNumber)
}

func FormatCharKey(ownerID string) string {
	return "char" + ownerID
}

func FormatAnimDataKey(ownerID string, action Action) string {
	return "animData" + action.String() + ownerID
}

// make it go through the animation add the animation and add it's data saperately
func (AS *AssetsChunks) loadAnimation(AC AnimationConfig, ID string) {
	action, err := ParseAction(AC.Action)
	if err != nil {
		panic("this action does not exist , refer to example animation set")
	}

	data, pnls := parseAnimCfg(AC)
	data.FramesCount = len(pnls)
	AS.chunksMap[FormatAnimDataKey(ID, action)] = data
	for frameNum, frame := range pnls {
		AS.chunksMap[FormatFrameKey(ID, action, (frameNum+1))] = frame
	}
}

func (c Charecter) getAnimationSetPath(dataPath string) string {
	formattedName := strings.ReplaceAll(c.Name, " ", "_")
	return fmt.Sprintf("%s/assets/animations/%s.json", dataPath, formattedName)
}

// we need to loop over the list of animations
// and we also need to validate the final map
// maybve we can do this by uhmm a small program inside the assets to validate all the assets

// returns the next frame in the animation sequence

type assetsProvider interface {
	GetFrame(ID string, action Action, frameIndex int) ([]string, error)
	GetAnimData(ID string, action Action) (AnimationData, error)
}

// this assets provider stores needed assets for the match in memory
type AssetsChunks struct {
	chunksMap map[string]any
}

// the idea is to not give direct access to the user for asset chunks ,
func NewAssets() *AssetsChunks {
	return &AssetsChunks{
		chunksMap: map[string]any{},
	}
}

// getAnimation uses the prefix thing for animations
func (AC *AssetsChunks) GetFrame(ID string, action Action, frameIndex int) ([]string, error) {
	anim, ok := AC.chunksMap[FormatFrameKey(ID, action, frameIndex)]
	if !ok {
		return []string{}, AssetNotFound
	}

	if reflect.TypeOf(anim) != reflect.TypeOf([]string{}) {
		return []string{}, TypeMismatchErr
	}

	return anim.([]string), nil
}

func (AC *AssetsChunks) GetChar(ID string) (Charecter, error) {
	char, ok := AC.chunksMap[FormatCharKey(ID)]
	if !ok {
		return Charecter{}, AssetNotFound
	}

	if reflect.TypeOf(char) != reflect.TypeOf(Charecter{}) {
		return Charecter{}, TypeMismatchErr
	}

	return char.(Charecter), nil
}

func (AC *AssetsChunks) GetAnimData(ID string, action Action) (AnimationData, error) {
	animData, ok := AC.chunksMap[FormatAnimDataKey(ID, action)]
	if !ok {
		return AnimationData{}, AssetNotFound
	}

	if reflect.TypeOf(animData) != reflect.TypeOf(AnimationData{}) {
		return AnimationData{}, TypeMismatchErr
	}

	return animData.(AnimationData), nil
}

func (AS AnimationData) progressAnimState() (AnimationData, error) {
	AS.CurrentFrame++
	if AS.CurrentFrame > AS.FramesCount {
		if !AS.RepeatingAnimation {
			return AnimationData{}, AnimationIsOver
		}
		AS.CurrentFrame = 1
	}

	return AS, nil
}
