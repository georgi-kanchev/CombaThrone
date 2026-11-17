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
func UpdateStoryMap() {
	UpdateEvents()

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
		ProgressOnMap()
	}

	SceneHUD.View.Zoom = View.Zoom
	gui.Scale = View.Zoom
	gui.ScaleNinePatches = gui.Scale
	gui.Image(PanelId, geometry.NewArea(-TileSize*5-8, -TileSize*3.5-8, TileSize*8.6, TileSize*3), false)

	SceneHUD.UpdateBack()
	SceneHUD.UpdateFront()

	var hovered, cur = geometry.Shape{}, progress.Start
	for i := progress.Tutorial; i <= progress.Underground; i++ {
		var p = storyMapProgressPoints[i]
		p.Width, p.Height = TileSize+4, TileSize+4
		var bg, img = GUI.Crops("progress")[i-1], storyProgressImages[i]
		var bgCrop, imgCrop = bg.CropArea(), img.CropArea()
		var tint = color.FadeOut(palette.White, 0.5)
		if Player.Progress == i {
			tint = palette.White
		}
		if p.ContainsPoint(View.MousePosition()) {
			hovered = p
			cur = i
		}
		View.DrawImage(p.X, p.Y, bgCrop.Width, bgCrop.Height, 0, bg, tint, geometry.Area{})
		View.DrawImage(p.X, p.Y, imgCrop.Width, imgCrop.Height, 0, img, tint, geometry.Area{})
	}

	storyMapMoveTween.TimeScale = TimeScale
	var pos = storyMapMoveTween.CurrentValues()
	View.DrawImage(pos[0], pos[1], TileSize, TileSize, 0, Characters[char.Explorer].Icon, palette.White, geometry.Area{})

	if hovered != (geometry.Shape{}) && Player.Progress == cur {
		SceneHUD.Highlight(View, hovered, palette.White)
		mouse.SetCursor(cursor.Hand)

		if mouse.IsAnyButtonJustPressed() {
			MoveOnMap(hovered.X, hovered.Y)
		}
	} else if hovered != (geometry.Shape{}) {
		SceneHUD.Highlight(View, hovered, palette.Red)
		mouse.SetCursor(cursor.NotAllowed)
	}
}

func ProgressOnMap() {
	storyMapFogFade[3+Player.Progress].GoTo(1.0, nil, 0)
	Player.Progress++
}
func MoveOnMap(targetX, targetY float32) {
	var pos = storyMapMoveTween.CurrentValues()
	var x, y = pos[0], pos[1]
	storyMapMoveTween.ForceStop()
	storyMapCollisionGrid.FindPathDiagonally(x, y, targetX, targetY, true, &storyMapPath)
	for i := 2; i < len(storyMapPath); i += 2 {
		var x1, y1, x2, y2 = storyMapPath[i-2], storyMapPath[i-1], storyMapPath[i], storyMapPath[i+1]
		var dist = point.DistanceToPoint(x1, y1, x2, y2)
		storyMapMoveTween.GoTo(dist/50, nil, x2, y2)
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
