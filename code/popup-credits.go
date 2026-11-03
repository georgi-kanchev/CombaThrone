package game

import (
	"game/code/constants/music"
	"game/code/constants/popup"
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
		"🔈Sound Effects & Ambience🔉\n🌘🟧🟤Joel Francis Burford\nOvani\nPMSFX\nSonnis",
		panel, geometry.Area{}, ThemeGUI, false, 0.6)

	panel = geometry.NewArea(panel.X, panel.Y+panel.Height/2+10, 60, 20)
	gui.Button("Back", panel, geometry.Area{}, ThemeGUI.Variant(1), true, 0.7)
	if gui.IsJustClicked() {
		Popup = popup.MainMenu
	}

	//=================================================================

	if CurrentMusic() != nil {
		ShowCurrentMusic()
	} else {
		PlayMusic(music.Menu)
	}
}
