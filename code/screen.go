package game

import (
	"game/code/constants/screen"
	"pure-game-kit/packages/graphics"
)

var View *graphics.View

var Screen screen.Enum
var Screens [screen.COUNT]func()

func Init() {
	var view = graphics.NewView(1)
	View = &view

	InitTitleScreen()

	InitAudio()
	InitUI() // needs to be before all other img loads since it relies on the asset IDs order
	InitCharacters()
	InitScene()
}

func Update() {
	UpdateAudio()

	switch Screen {
	case screen.TitleMenu:
		UpdateTitleScreen()
	case screen.Scene:
		UpdateScene()
	}

	UpdateDebug()
}
