package game

import "pure-game-kit/packages/utility/number"

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
		var attackable, entrance = self.EnemyEntrance()
		if attackable && entrance != nil && self.UnitFront == nil {
			self.AddEffect(EffectMoreDmgVsEntrances)
		} else {
			self.RemoveEffect(EffectMoreDmgVsEntrances)
		}
	},
	CharCook: func(self *Unit) {
		var canAct = number.IsNaN(self.ActTimer) || self.ActTimer < 0
		var targetUnit *Unit
		for _, u := range Units {
			var isClose = number.IsWithin(u.X, self.X, TileSize*0.5)
			var isAlive, isMaxHp = u.Health > 0, u.Health == u.Values.MaxHealth
			if u.Team == self.Team && u != self && u.Lane == self.Lane-1 && isClose && !isMaxHp && isAlive {
				targetUnit = u
				if canAct {
					self.State = StateActStart
				}
				break
			}
		}
		if self.State == StateActTrigger && targetUnit != nil {
			var isAlive, isMaxHp = targetUnit.Health > 0, targetUnit.Health == targetUnit.Values.MaxHealth
			if isAlive && !isMaxHp {
				targetUnit.Heal(4)
			}
		}
	},
	CharBowyer: func(self *Unit) {
		if self.State == StateSummoned && !self.IsGarrisoner() {
			self.AddEffect(EffectMoreRangeOnGround)
		}
	},
	CharSmith: func(self *Unit) {
		var canAct = number.IsNaN(self.ActTimer) || self.ActTimer < 0
		if self.State == StateIdling && self.UnitFront != nil && canAct {
			self.State = StateActStart
		} else if self.State == StateActTrigger && self.UnitFront != nil {
			self.UnitFront.VelocityY = -10
			self.UnitFront.Values.SleepTimer = 1.5
			if self.Team == TeamAlly {
				self.UnitFront.VelocityX = 50
			} else {
				self.UnitFront.VelocityX = -50
			}
		}
	},
}
