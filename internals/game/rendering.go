package game

import (
	"fmt"
	"strings"
)

func (f frame) renderSprite(s []string, color string, renderPoint Point) frame {
	outputFrame := f
	midRowIndex := len(s) / 2
	for rowNumber, row := range s {
		relativeY := rowNumber - midRowIndex
		midCharecterIndex := len(row) / 2
		finalY := relativeY + renderPoint.Y
		if (finalY < 0) || (finalY > len(f)-1) {
			continue
		}
		for charecterIndex, charecter := range row {
			relativeX := charecterIndex - midCharecterIndex
			finalX := relativeX + renderPoint.X
			if (finalX > len(f[0])-1) || (finalX < 0) {
				continue
			}
			outputFrame[finalY][finalX] = color + string(charecter) + reset

		}
	}

	return outputFrame
}

// we have a bucket of events we want to use them in order to pretty much alter the state every single time
// so we need something that can construct this events list
func (game Game) MakeFrame(s State) frame {
	// background filling first layer
	// make this into a function btw , generateBasePlate()
	cfg := game.Config    // abbreviation
	res := cfg.Resulotion // abbreviation
	var f frame

	for range res.Height {
		var rowToAdd []string
		for range res.Width {
			rowToAdd = append(rowToAdd, string(cfg.BackgroundTile))
		}
		f = append(f, rowToAdd)
	}
	// second layer for now : rendering the charecters
	position := game.State.P1State.Position
	animationFrame, err := game.AssetsProvider.GetFrame(game.State.P1State.OwnerID, game.State.P1State.Action, game.State.P1State.AnimState.CurrentFrame)
	if err != nil {
		fmt.Println(err) // improve the error handling here
	}
	f = f.renderSprite(animationFrame, green, position)
	// next step would be to render the actual charecter instead of this

	return f

}

// optimize this later along with the frame gimmick
func (f frame) RenderFrame() {
	frameToRender := ClearAscii
	for _, row := range f {
		frameToRender += strings.Join(row, "") + "\r\n"
	}
	fmt.Println(frameToRender)
}
