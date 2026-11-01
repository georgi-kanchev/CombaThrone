package game

import (
	"game/code/constants/music"
	"pure-game-kit/packages/gui"
)

func InitCredits() {
}

func UpdateCredits() {
	const scale = 0.2
	var panel, _, _ = settingsLayout.Box(0, scale)
	gui.Image(PanelId, panel, false)

	//=================================================================

	if CurrentMusic() != nil {
		ShowCurrentMusic()
	} else {
		PlayMusic(music.Menu)
	}
}
