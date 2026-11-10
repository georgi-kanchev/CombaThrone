package game

import (
	"game/code/constants/icon"
	"game/code/constants/popup"
	"game/code/constants/screen"
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/gui"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/motion"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/easing"
	"pure-game-kit/packages/utility/text"
)

var View *graphics.View

var Screen screen.Enum = screen.Title
var Popup popup.Enum = popup.MainMenu
var Screens [screen.COUNT]func()

var Tags = [123]string{
	"⓪", "①", "②", "③", "④", "⑤", "⑥", "⑦", "⑧", "⑨", "⑩", "⑪", "⑫", "⑬", "⑭", "⑮", "⑯", "⑰", "⑱", "⑲", "⑳",
	"⑴", "⑵", "⑶", "⑷", "⑸", "⑹", "⑺", "⑻", "⑼", "⑽", "⑾", "⑿", "⒀", "⒁", "⒂", "⒃", "⒄", "⒅", "⒆", "⒇",
	"⒈", "⒉", "⒊", "⒋", "⒌", "⒍", "⒎", "⒏", "⒐", "⒑", "⒒", "⒓", "⒔", "⒕", "⒖", "⒗", "⒘", "⒙", "⒚", "⒛",
	"ⓐ", "ⓑ", "ⓒ", "ⓓ", "ⓔ", "ⓕ", "ⓖ", "ⓗ", "ⓘ", "ⓙ", "ⓚ", "ⓛ", "ⓜ", "ⓝ", "ⓞ", "ⓟ", "ⓠ", "ⓡ", "ⓢ", "ⓣ", "ⓤ",
	"ⓥ", "ⓦ", "ⓧ", "ⓨ", "ⓩ", "Ⓐ", "Ⓑ", "Ⓒ", "Ⓓ", "Ⓔ", "Ⓕ", "Ⓖ", "Ⓗ", "Ⓘ", "Ⓙ", "Ⓚ", "Ⓛ", "Ⓜ", "Ⓝ", "Ⓞ", "Ⓟ",
	"Ⓠ", "Ⓡ", "Ⓢ", "Ⓣ", "Ⓤ", "Ⓥ", "Ⓦ", "Ⓧ", "Ⓨ", "Ⓩ", "⓫", "⓬", "⓭", "⓮", "⓯", "⓰", "⓱", "⓲", "⓳", "⓴",
}
var ThemeGUI assets.GUIThemeId
var GUI assets.AtlasId
var SlotId, PanelId, ButtonUpId, ButtonDownId assets.ImageId
var ScreenFade = motion.NewTween(255)

func Init() {
	View = new(graphics.NewView(1))
	ScreenFade.GoTo(0.35, easing.SineIn, 0)

	InitAudio()
	InitGUI()

	InitMainMenu()
	InitSettings()
	InitCredits()

	InitTitleScreen()
	InitStoryMap()
	InitScene()
}
func InitGUI() {
	GUI = assets.LoadAtlas(assets.LoadImage("data/gui.png"), "data/gui.xml")

	for i := range icon.COUNT {
		assets.FontId(0).EmbedImage(text.At(Tags[i], 0), GUI.Crops("icons-text")[i])
	}

	ThemeGUI = assets.LoadGUITheme("data/gui-theme.xml")
	SlotId = assets.LoadImage9Patch(GUI.Crops("slot")[0], 8, 8, 8, 8).DuplicateId(-12345_3)

	var btn = GUI.Crops("button")
	ButtonUpId = assets.LoadImage9Patch(btn[0], 8, 8, 8, 8).DuplicateId(-12345_1)   // ensuring the ImageIds are known after
	ButtonDownId = assets.LoadImage9Patch(btn[1], 8, 8, 8, 8).DuplicateId(-12345_2) // load/crop for GUI themes
	PanelId = assets.LoadImage9Patch(GUI.Crops("panel")[0], 8, 8, 8, 8)
}

func Update() {
	switch Screen {
	case screen.Title:
		UpdateTitleScreen()
	case screen.StoryMap:
		UpdateStoryMap()
	case screen.Scene:
		UpdateScene()
	}

	switch Popup {
	case popup.MainMenu:
		UpdateMainMenu()
	case popup.Settings:
		UpdateSettings()
	case popup.Credits:
		UpdateCredits()
	}

	handleEscape()
	UpdateMusic()
	UpdateDebug()

	View.DrawColor(color.RGBA(0, 0, 0, uint8(ScreenFade.CurrentValues()[0])))
}

func AlignView() {
	var _, bly = CurrentZone.Ground.PointFromEdge(0.5, 1)
	View.FitSize(CurrentZone.Ground.Width, 0)
	var _, h = View.Size()
	View.X, View.Y = 0, (bly-h/2)-2
	// View.Zoom = 4 // left for debug purposes
}

func DimScreenUI() {
	gui.Rectangle(color.RGBA(0, 0, 0, 180), gui.AreaHUD(0.5, 0.5, 1, 1), geometry.Area{}, true) // dim popup background
}

// private ========================================================

func handleEscape() {
	if Screen == screen.Title && Popup == popup.MainMenu {
		return
	}

	if keyboard.IsKeyJustPressed(key.Escape) {
		if Popup != popup.MainMenu {
			Popup = popup.MainMenu
			TimeScale = 0
			PauseSounds()
		} else {
			Popup = popup.NONE
			TimeScale = 1
			ResumeSounds()
		}
	}
	if keyboard.IsKeyJustPressed(key.Backtick) {
		if TimeScale == 1 {
			TimeScale = 0
			PauseSounds()
		} else {
			TimeScale = 1
			ResumeSounds()
		}
	}
}
