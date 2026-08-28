package game

import (
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
)

func BehaviorMan(self *Unit) {
}

func BehaviorWoman(self *Unit) {
	if keyboard.IsKeyJustPressed(key.A) {
		self.VelocityX = 0
		self.State = StateActionStart
	}
}

func BehaviorHunter(self *Unit) {
	if self.State == StateSummoned {
		if !self.IsGarrisoner() {
			self.Stats.ActRange += 2
		}
	}
}

func BehaviorDummy(self *Unit) {
	if self.State == StateDyingStart {
		self.State = StateIdling
		self.Health = 1
	}
	if self.State == StateSummoned {
		self.X, self.Y = TileSize*5, -TileSize*6
	}
}
