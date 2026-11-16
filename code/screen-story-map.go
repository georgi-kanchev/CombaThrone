package game

import (
	char "game/code/constants/character"
	"game/code/constants/progress"
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/gui"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/input/mouse"
	"pure-game-kit/packages/input/mouse/cursor"
	"pure-game-kit/packages/motion"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/point"
)

func InitStoryMap() {
	var layers, _ = assets.LoadTileLayersFromTiled("data/story.tmx")
	for i, l := range layers {
		var tmap = new(graphics.NewTilemap(l))
		storyMapLayers = append(storyMapLayers, tmap)

		if i < 3 {
			storyMapFogFade = append(storyMapFogFade, nil)
			continue
		}

		storyMapFogFade = append(storyMapFogFade, motion.NewTween(FOG_OPACITY))
	}
	storyMapLayers[0].Details.Tint = palette.White

	storyMapCollisionGrid = new(geometry.NewShapeGrid(TileSize))
	storyMapCollisionGrid.AddShapes(storyMapLayers[len(storyMapLayers)-2].TilemapShapes()...)
	storyMapProgressPoints = storyMapLayers[len(storyMapLayers)-1].TilemapShapes()

	var start = storyMapProgressPoints[progress.Start]
	storyMapMoveTween = motion.NewTween(start.X, start.Y)

	var diff, info = GUI.Crops("difficulty"), GUI.Crops("info")
	storyProgressImages = [progress.COUNT]assets.ImageId{0, info[0], diff[0], diff[1], diff[2], diff[3]}
}
func Progress() {
	storyMapFogFade[3+Player.Progress].GoTo(1.0, nil, 0)
	Player.Progress++
}
func UpdateStoryMap() {
	View.FillSize(TileSize*18, TileSize*10)
	mouse.SetCursor(cursor.Arrow)

	View.X, View.Y = -TileSize*9, -TileSize*10
	for i, l := range storyMapLayers {
		if i >= 3 {
			var alpha = byte(storyMapFogFade[i].CurrentValues()[0])
			l.Details.Tint = color.RGBA(255, 255, 255, alpha)
		}
		if i < len(storyMapLayers)-2 {
			View.DrawObject(l)
		}
	}

	if keyboard.IsKeyJustPressed(key.F) {
		Progress()
	}

	SceneHUD.View.Zoom = View.Zoom
	gui.Scale = View.Zoom
	gui.ScaleNinePatches = gui.Scale
	gui.Image(PanelId, geometry.NewArea(-TileSize*5-8, -TileSize*3.5-8, TileSize*8.6, TileSize*3), false)

	SceneHUD.UpdateBack()
	SceneHUD.UpdateFront()

	for i := progress.Tutorial; i <= progress.Underground; i++ {
		var point = storyMapProgressPoints[i]
		var bg = GUI.Crops("progress")[i-1] // ignore Start
		var bgCrop = bg.CropArea()
		var img = storyProgressImages[i]
		var imgCrop = img.CropArea()
		View.DrawImage(point.X, point.Y, bgCrop.Width, bgCrop.Height, 0, bg, palette.White, geometry.Area{})
		View.DrawImage(point.X, point.Y, imgCrop.Width, imgCrop.Height, 0, img, palette.White, geometry.Area{})

		point.Width, point.Height = TileSize*1.5, TileSize*1.5
		if point.ContainsPoint(View.MousePosition()) {
			SceneHUD.Highlight(View, point, palette.White)
			mouse.SetCursor(cursor.Hand)
		}
	}

	var pos = storyMapMoveTween.CurrentValues()
	var x, y = pos[0], pos[1]
	View.DrawImage(x, y, TileSize, TileSize, 0, Characters[char.Explorer].Icon, palette.White, geometry.Area{})

	var mx, my = View.MousePosition()
	storyMapCollisionGrid.FindPathDiagonally(x, y, mx, my, true, &storyMapPath)

	if mouse.IsAnyButtonJustPressed() {
		storyMapMoveTween.ForceStop()
		for i := 2; i < len(storyMapPath); i += 2 {
			var x1, y1, x2, y2 = storyMapPath[i-2], storyMapPath[i-1], storyMapPath[i], storyMapPath[i+1]
			var dist = point.DistanceToPoint(x1, y1, x2, y2)
			storyMapMoveTween.GoTo(dist/50, nil, x2, y2)
		}
	}
}

// private ========================================================

const FOG_OPACITY, COLLISION_TILE_ID = 128, 144

var storyMapProgressPoints []geometry.Shape
var storyMapLayers []*graphics.Object
var storyMapFogFade []*motion.Tween
var storyMapCollisionGrid *geometry.ShapeGrid
var storyMapMoveTween *motion.Tween
var storyMapPath []float32
var storyProgressImages [progress.COUNT]assets.ImageId
