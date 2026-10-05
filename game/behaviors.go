package game

import (
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/direction"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/text"
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
				targetUnit.AffectHealth(4)
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
		if len(self.Carrying) == 0 {
			self.AddEffect(EffectKid)
		} else {
			self.RemoveEffect(EffectKid)
		}
	},
	CharHorse: func(self *Unit) {
		for _, u := range Units {
			if self != u && self.Team == u.Team {
				var affectedLane = self.Lane == u.Lane || self.Lane == u.Lane+1
				if affectedLane && number.IsWithin(u.X, self.X, float32(self.Values.ActRange)*TileSize) {
					u.AddEffect(EffectHorse)
				} else {
					u.RemoveEffect(EffectHorse)
				}
			}
		}
	},
	CharFisherman: func(self *Unit) {
		if self.State == StateSummoned {
			self.ActTimer = self.Values.ActTimer
		}

		var offset float32 = TileSize / 2.5
		if self.IsReturning {
			offset *= -1
		}
		if self.LastState == StateActTrigger { // is fishing
			for _, u := range Units {
				if self == u || self.Team == u.Team {
					continue
				}
				var hb = u.Hitbox()
				if self.Lane == u.Lane+1 && number.IsWithin(self.X+offset, hb.X, hb.Width/2) {
					u.VelocityY = -150
					u.Lane += 2
					self.ActTimer = -1 // trigger act
					break
				}
			}
			var height float32 = TileSize * 0.7
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
				if self.Team != u.Team && number.IsWithin(self.X, u.X, TileSize*2) {
					var dirX, dirY = direction.BetweenPoints(self.X, self.Y+self.Height/2, u.X, u.Y)
					u.VelocityX, u.VelocityY = dirX*100, dirY*100
					u.Values.SleepTimer = 1.5
				}
			}
		}
	},
	CharSkirmisherGoblin:  behaviorGoblin,
	CharArbalestierGoblin: behaviorGoblin,
	CharStabberBandit: func(self *Unit) {
		if self.State != StateActTrigger || self.UnitFront == nil || self.UnitFront.ProjectHealth(-self.Values.ActPoints) > 0 {
			return // i have not attacked or enemy will live from my attack, so bail
		}

		var eff, hasEffect = self.Effects[EffectStabberBandit] // accumulated effect from all kills
		if !hasEffect {
			self.AddEffect(EffectStabberBandit)
			eff = self.Effects[EffectStabberBandit]
		} else { // dmg accumulates on the unit effect itself for tooltip display but its dmg value should remain unchanged
			self.Effects[EffectStabberBandit] = Effects[EffectStabberBandit] // so swap for original
			self.TickEffect(EffectStabberBandit)                             // tick original
			eff.ActPoints += Effects[EffectStabberBandit].ActPoints          // accumulate the unit one
		}
		eff.EffectInfo = text.New("🟩", Tags[IconPlus], "🟧", eff.ActPoints, " ", Tags[IconSword], "damage ⬜(Self)")
		self.Effects[EffectStabberBandit] = eff // return back to unit
	},
	CharGunnerBandit: func(self *Unit) {
		if self.UnitFront == nil {
			self.RemoveEffect(EffectGunnerBandit1)
			self.RemoveEffect(EffectGunnerBandit2)
		} else {
			self.AddEffect(EffectGunnerBandit1)
			self.AddEffect(EffectGunnerBandit2)
		}
	},
	CharVulture: func(self *Unit) {
		for _, u := range Units {
			if self != u && self.Team == u.Team && u.State == StateDecaying {
				if self.Lane == u.Lane+1 && number.IsWithin(u.X, self.X, float32(self.Values.ActRange)*TileSize) {
					u.AddEffect(EffectVulture)                // effect is just for tooltip, no stats
					u.Values.ReviveTimer -= DeltaTimeScaled() // applied here instead
				} else {
					u.RemoveEffect(EffectVulture)
				}
			}
		}
	},
	CharRaven: func(self *Unit) {
		for _, u := range Units {
			if self != u && self.Team != u.Team {
				if self.Lane == u.Lane+1 && number.IsWithin(u.X, self.X, float32(self.Values.ActRange)*TileSize) {
					u.AddEffect(EffectRaven)
				} else {
					u.RemoveEffect(EffectRaven)
				}
			}
		}
	},
	CharRedSnake: func(self *Unit) {
		for _, u := range Units {
			if self != u && self.Team != u.Team {
				if self.Lane == u.Lane+1 && number.IsWithin(u.X, self.X, float32(self.Values.ActRange)*TileSize) {
					u.AddEffect(EffectRedSnake)
				} else {
					u.RemoveEffect(EffectRedSnake)
				}
			}
		}
	},
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
