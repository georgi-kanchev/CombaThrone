package game

func BehaviorMan(self *Unit) {
}

func BehaviorWoman(self *Unit) {
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
