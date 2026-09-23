package game

var Behaviors = map[CharacterKind]func(self *Unit){
	CharDummy: func(self *Unit) {
		switch self.State {
		case StateDyingStart:
			self.State = StateIdling
			self.Health = 1
		case StateSummoned:
			self.X, self.Y = TileSize*5, -TileSize*6
		case StateWalking:
			self.State = StateIdling
		}
	},
	CharMiner: func(self *Unit) {

	},
	CharCook: func(self *Unit) {
	},
	CharBowyer: func(self *Unit) {
		if self.State == StateSummoned {
			if !self.IsGarrisoner() {
				self.Stats.ActRange += 2
			}
		}
	},
}
