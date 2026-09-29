package game

import (
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/direction"
	"pure-game-kit/packages/utility/number"
)

var Behaviors = map[CharacterKind]func(self *Unit){
	CharDummy: func(self *Unit) {
		switch self.State {
		case StateDyingStart:
			self.State = StateIdling
			self.Health = 1
		case StateSummoned:
			self.X, self.Y = TileSize*5, -TileSize*6
		case StateMoving:
			self.State = StateIdling
		}
	},
	CharMiner: func(self *Unit) {
		var attackable, entrance = self.EnemyEntrance()
		if attackable && entrance != nil && self.UnitFront == nil {
			self.AddEffect(EffectMiner)
		} else {
			self.RemoveEffect(EffectMiner)
		}
	},
	CharCook: func(self *Unit) {
		var canAct = number.IsNaN(self.ActTimer) || self.ActTimer < 0
		var targetUnit *Unit
		for _, u := range Units {
			var isClose = number.IsWithin(u.X, self.X, TileSize*0.5)
			var isMaxHp = u.Health == u.Values.MaxHealth
			if u.Team == self.Team && u != self && u.Lane == self.Lane-1 && isClose && !isMaxHp && u.IsAlive() {
				targetUnit = u
				if canAct {
					self.State = StateActStart
				}
				break
			}
		}
		if self.State == StateActTrigger && targetUnit != nil {
			if targetUnit.IsAlive() && targetUnit.Health != targetUnit.Values.MaxHealth {
				targetUnit.Heal(4)
			}
		}
	},
	CharBowyer: func(self *Unit) {
		if self.State == StateSummoned && !self.IsGarrisoner() {
			self.AddEffect(EffectBowyer)
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
	CharKid: func(self *Unit) {
		if self.Carrying != nil {
			self.RemoveEffect(EffectKid)
		} else {
			self.AddEffect(EffectKid)
		}
	},
	CharHorse: func(self *Unit) {
		for _, u := range Units {
			if self == u || self.Team != u.Team {
				continue
			}

			var affectedLane = self.Lane == u.Lane || self.Lane == u.Lane+1
			if affectedLane && number.IsWithin(u.X, self.X, float32(self.Values.ActRange)*TileSize) {
				u.AddEffect(EffectHorse)
			} else {
				u.RemoveEffect(EffectHorse)
			}
		}
	},
	CharFisherman: func(self *Unit) {
		if self.State == StateSummoned {
			self.ActTimer = self.Values.ActTimer
		}

		var offset float32 = TileSize / 2
		if self.IsReturning {
			offset = -TileSize / 2
		}
		if self.LastState == StateActTrigger { // is fishing
			for _, u := range Units {
				if self == u || self.Team == u.Team {
					continue
				}
				var hb = u.Hitbox()
				var closeEnough = number.IsWithin(self.X+offset, hb.X, hb.Width/2)
				if !closeEnough {
					continue
				}

				var correctLane = self.Lane == u.Lane+3
				if closeEnough && correctLane {
					u.VelocityY = -150
					u.Lane += 2
					self.ActTimer = -1 // trigger act
					break
				}
			}
			var height float32 = TileSize * 1.35
			var x, y = self.X + 18.5, self.Y + height
			if (self.Team == TeamAlly && self.IsReturning) || (self.Team == TeamEnemy && !self.IsReturning) {
				x = self.X - 18.5
			}
			View.DrawShape(x, y, 3, height, 0, 0, color.RGB(37, 19, 26), geometry.Area{})
			View.DrawShape(x, y, 1, height, 0, 0, color.RGB(167, 172, 186), geometry.Area{})
			View.DrawShape(x, y+height/2, 5, 5, 0, 0.8, color.RGB(37, 19, 26), geometry.Area{})
			View.DrawShape(x, y+height/2, 3, 3, 0, 0.8, color.RGB(208, 25, 80), geometry.Area{})
		}

		if self.ActTimer < 0 {
			if self.LastState == StateActTrigger {
				self.ActTimer = self.Values.ActTimer // consume act to hook up
				self.State = StateActRecovering      // and keep moving
				self.LastState = StateActRecovering
			} else {
				self.State = StateActStart
			}
		}
		if self.LastState == StateActTrigger {
			self.State = StateActTrigger // force keep StateActEnd
		}
	},
	CharTroll: func(self *Unit) {
		if self.State == StateActTrigger {
			for _, u := range Units {
				if self.Team == u.Team {
					continue
				}
				if number.IsWithin(self.X, u.X, TileSize*2) {
					var dirX, dirY = direction.BetweenPoints(self.X, self.Y+self.Height/2, u.X, u.Y)
					u.VelocityX, u.VelocityY = dirX*100, dirY*100
					u.Values.SleepTimer = 1.5
				}
			}
		}
	},
	CharSkirmisherGoblin:  behaviorGoblin,
	CharArbalestierGoblin: behaviorGoblin,
}

func behaviorGoblin(self *Unit) {
	if self.UnitFront != nil && self.UnitFront.Values.Role == RoleDefender {
		self.AddEffect(EffectGoblin)
	} else {
		self.RemoveEffect(EffectGoblin)
	}

	if self.ClosestEnemyInRange != nil && self.ClosestEnemyInRange.Values.Role == RoleDefender {
		self.AddEffect(EffectGoblin)
	} else {
		self.RemoveEffect(EffectGoblin)
	}
}
