package game

import (
	"game/code/constants/music"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/gui"
)

func InitCredits() {
}

func UpdateCredits() {
	const scale = 0.2
	var panel, _, _ = settingsLayout.Box(0, scale)
	panel.Width /= 2
	gui.Image(PanelId, panel, false)

	gui.Text("🔈Graphics🔉\n🌘🟧🟤CraftPix⬜⚫\n\n"+
		"🔈Music🔉\n🌘🟧🟤Joel Steudler⬜⚫\n\n"+
		"🔈Sound Effects🔉\n🌘🟧🟤Joel Francis Burford\nOvani\nPMSFX\nSonnis",
		panel, geometry.Area{}, ThemeGUI, false, 0.6)

	//=================================================================

	if CurrentMusic() != nil {
		ShowCurrentMusic()
	} else {
		PlayMusic(music.Menu)
	}
}
