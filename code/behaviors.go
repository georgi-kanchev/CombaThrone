package game

import (
	char "game/code/constants/character"
	"game/code/constants/effect"
	"game/code/constants/icon"
	"game/code/constants/lane"
	"game/code/constants/role"
	"game/code/constants/splash"
	"game/code/constants/state"
	"game/code/constants/team"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/utility/collection"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/direction"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/text"
)

var Behaviors = map[char.Enum]func(self *Unit){
	char.Dummy: func(self *Unit) {
		switch self.State {
		case state.DyingStart:
			self.State = state.Idling
			self.Health = 1
		case state.Summoned:
			self.X, self.Y = TileSize*5, -TileSize*6
		case state.Moving:
			self.State = state.Idling
		}
	},
	char.Miner: func(self *Unit) {
		var attackable, entrance = self.EnemyEntrance()
		if attackable && entrance != nil && self.UnitFront == nil {
			self.AddEffect(effect.Miner)
		} else {
			self.RemoveEffect(effect.Miner)
		}
	},
	char.Cook: func(self *Unit) {
		var canAct = number.IsNaN(self.ActTimer) || self.ActTimer < 0
		var targetUnit *Unit
		for _, u := range Units {
			var isClose = number.IsWithin(u.X, self.X, TileSize*0.5)
			var isMaxHp = u.Health == u.Values.MaxHealth
			if u.Team == self.Team && u != self && u.Lane == self.Lane-1 && isClose && !isMaxHp && u.IsAlive() {
				targetUnit = u
				if canAct {
					self.State = state.ActStart
				}
				break
			}
		}
		if self.State == state.ActTrigger && targetUnit != nil {
			if targetUnit.IsAlive() && targetUnit.Health != targetUnit.Values.MaxHealth {
				targetUnit.AffectHealth(self.Values.ActPoints, true, nil)
			}
		}
	},
	char.Bowyer: func(self *Unit) {
		if self.State == state.Summoned && !self.IsGarrisoner() {
			self.AddEffect(effect.Bowyer)
		}
	},
	char.Smith: func(self *Unit) {
		var canAct = number.IsNaN(self.ActTimer) || self.ActTimer < 0
		if self.State == state.Idling && self.UnitFront != nil && canAct {
			self.State = state.ActStart
		} else if self.State == state.ActTrigger && self.UnitFront != nil {
			self.UnitFront.VelocityY = -10
			self.UnitFront.Values.DisablesTimer += 1.5
			if self.Team == team.Ally {
				self.UnitFront.VelocityX = 50
			} else {
				self.UnitFront.VelocityX = -50
			}
		}
	},
	char.Kid: func(self *Unit) {
		if len(self.Carrying) == 0 {
			self.AddEffect(effect.Kid)
		} else {
			self.RemoveEffect(effect.Kid)
		}
	},
	char.Horse: func(self *Unit) {
		self.applyEffectInRange(effect.Horse, self.Team, []lane.Enum{self.Lane, self.Lane - 1}, nil)
	},
	char.Fisherman: func(self *Unit) {
		if self.State == state.Summoned {
			self.ActTimer = self.Values.ActTimer
		}

		var offset float32 = TileSize / 2.5
		if self.IsReturning {
			offset *= -1
		}
		if self.LastState == state.ActTrigger { // is fishing
			for _, u := range Units {
				if self == u || self.Team == u.Team {
					continue
				}
				var hb = u.Hitbox()
				if self.Lane == u.Lane+1 && number.IsWithin(self.X+offset, hb.X, hb.Width/2) {
					u.VelocityY = -150
					u.Lane = min(u.Lane+2, lane.Upper)
					self.ActTimer = -1 // trigger act
					break
				}
			}
			var height float32 = TileSize * 0.7
			var x, y = self.X + 18.5, self.Y + height
			if (self.Team == team.Ally && self.IsReturning) || (self.Team == team.Enemy && !self.IsReturning) {
				x = self.X - 18.5
			}
			View.DrawShape(x, y, 3, height, 0, 0, color.RGB(37, 19, 26), geometry.Area{})
			View.DrawShape(x, y, 1, height, 0, 0, color.RGB(167, 172, 186), geometry.Area{})
			View.DrawShape(x, y+height/2, 5, 5, 0, 0.8, color.RGB(37, 19, 26), geometry.Area{})
			View.DrawShape(x, y+height/2, 3, 3, 0, 0.8, color.RGB(208, 25, 80), geometry.Area{})
		}

		if self.ActTimer < 0 {
			if self.LastState == state.ActTrigger {
				self.ActTimer = self.Values.ActTimer // consume act to hook up
				self.State = state.ActRecovering     // and keep moving
				self.LastState = state.ActRecovering
			} else {
				self.State = state.ActStart
			}
		}
		if self.LastState == state.ActTrigger {
			self.State = state.ActTrigger // force keep StateActEnd
		}
	},
	char.Troll: func(self *Unit) {
		if self.State == state.ActTrigger {
			var hb = self.Hitbox()
			var offX float32 = 12
			if self.IsFacingLeft() {
				offX *= -1
			}
			Splashes = collection.Add(Splashes, NewSplash(splash.Dust, 0, hb.X+offX, hb.Y+hb.Height/2-14, self))
			for _, u := range Units {
				if self.Team != u.Team && number.IsWithin(self.X, u.X, TileSize*2) {
					var dirX, dirY = direction.BetweenPoints(self.X, self.Y+self.Height/2, u.X, u.Y)
					u.VelocityX, u.VelocityY = dirX*100, dirY*100
					u.Values.DisablesTimer += 1.5
				}
			}
		}
	},
	char.Spearlin:    behaviorGoblin,
	char.Arbalestier: behaviorGoblin,
	char.Stabber: func(self *Unit) {
		if self.State != state.ActTrigger || self.UnitFront == nil || self.UnitFront.ProjectHealth(-self.Values.ActPoints) > 0 {
			return // i have not attacked or enemy will live from my attack, so bail
		}

		var eff = self.Effects[effect.Stabber] // accumulated effect from all kills
		if eff == (Values{}) {
			self.AddEffect(effect.Stabber)
			eff = self.Effects[effect.Stabber]
		} else { // dmg accumulates on the unit effect itself for tooltip display but its dmg value should remain unchanged
			self.Effects[effect.Stabber] = Effects[effect.Stabber] // so swap for original
			self.TickEffect(effect.Stabber)                        // tick original
			eff.ActPoints += Effects[effect.Stabber].ActPoints     // accumulate the unit one
		}
		eff.EffectInfo = text.New("🟩", Tags[icon.Plus], "🟧", eff.ActPoints, " ", Tags[icon.Fighter], "damage ⬜(Self)")
		self.Effects[effect.Stabber] = eff // return back to unit
	},
	char.Gunner: func(self *Unit) {
		if self.UnitFront == nil {
			self.RemoveEffect(effect.Gunner1)
			self.RemoveEffect(effect.Gunner2)
		} else {
			self.AddEffect(effect.Gunner1)
			self.AddEffect(effect.Gunner2)
		}
	},
	char.Vulture: func(self *Unit) {
		for _, u := range Units { // special roll-out instead of func applyEffectInRange()
			if self != u && self.Team == u.Team && u.State == state.Decaying && self.Lane == u.Lane+1 {
				if self.IsUnitInRange(u) {
					u.AddEffect(effect.Vulture)               // effect is just for tooltip, no stats
					u.Values.ReviveTimer -= DeltaTimeScaled() // applied here instead
				} else {
					u.RemoveEffect(effect.Vulture)
				}
			}
		}
	},
	char.Raven: func(self *Unit) {
		self.applyEffectInRange(effect.RavenFighter, 1-self.Team, []lane.Enum{self.Lane - 1}, []role.Enum{role.Fighter})
		self.applyEffectInRange(effect.RavenBackliner, 1-self.Team, []lane.Enum{self.Lane - 1}, []role.Enum{role.Backliner})
	},
	char.Snake: func(self *Unit) {
		self.applyEffectInRange(effect.Snake, 1-self.Team, []lane.Enum{self.Lane, self.Lane - 1}, nil)
	},
	char.Healer: func(self *Unit) {
		var canAct = number.IsNaN(self.ActTimer) || self.ActTimer < 0
		var target *Unit
		for _, u := range Units {
			var isClose = self.IsUnitInRange(u)
			var isMaxHp = u.Health == u.Values.MaxHealth
			var inFront = (!self.IsFacingLeft() && u.X > self.X) || (self.IsFacingLeft() && u.X < self.X)
			if u.Team == self.Team && u != self && isClose && !isMaxHp && u.IsAlive() && inFront {
				target = u
				break
			}
		}
		if target != nil {
			if canAct {
				self.State = state.ActStart
			}
			if self.State == state.ActTrigger {
				self.shootAt(target)
			}
		}
	},
	char.Bear: func(self *Unit) {
		self.applyEffectInRange(effect.Bear, 1-self.Team, []lane.Enum{self.Lane, self.Lane - 1}, nil)
	},
	char.Boar: func(self *Unit) {
		if self.IsReturning {
			return
		}

		for _, u := range Units {
			if u != nil && self.Team != u.Team && self.Lane == u.Lane+1 && u.IsAlive() && self.Hitbox().Overlaps(u.Hitbox()) {
				if self.IsFacingLeft() {
					u.VelocityX = -150
				} else {
					u.VelocityX = 150
				}
				u.VelocityY = -80
				u.AddEffect(effect.Boar)
				self.IsReturning = true
			}
		}
	},
	char.Wolf: func(self *Unit) {
		if self.IsReturning {
			return
		}

		if self.Following == nil {
			for _, u := range Units {
				if u != self && u.Team == self.Team && u.IsAlive() && u.Lane == self.Lane-1 && self.Hitbox().Overlaps(u.Hitbox()) {
					self.Following = u
				}
			}
			return
		}

		if self.Hitbox().Overlaps(self.Following.Hitbox()) {
			self.Values.DisablesTimer = 0.1
			self.Values.MoveSpeed = self.Following.Values.MoveSpeed
			self.Following.AddEffect(effect.Wolf)              // effect is just for tooltip, no stats
			self.Following.ActTimer -= DeltaTimeScaled() * 0.5 // applied here instead
		}
		if !self.Following.IsSummoned() {
			self.IsReturning = true
			self.Following = nil
		}
	},
}

func behaviorGoblin(self *Unit) {
	if self.UnitFront != nil && self.UnitFront.Values.Role == role.Defender {
		self.AddEffect(effect.Goblin)
	} else {
		self.RemoveEffect(effect.Goblin)
	}

	if self.ClosestEnemyInRange != nil && self.ClosestEnemyInRange.Values.Role == role.Defender {
		self.AddEffect(effect.Goblin)
	} else {
		self.RemoveEffect(effect.Goblin)
	}
}
func (u *Unit) applyEffectInRange(effect effect.Enum, forTeam team.Enum, forLanes []lane.Enum, forRole []role.Enum) {
	if len(forRole) == 0 {
		forRole = []role.Enum{role.Fighter, role.Backliner, role.Defender, role.Griefer, role.Supplier, role.Collector}
	}

	for _, unit := range Units {
		var correctLane = collection.Contains(forLanes, unit.Lane)
		var correctRole = collection.Contains(forRole, unit.Values.Role)

		if u != unit && unit.Team == forTeam && unit.IsAlive() && correctLane && correctRole {
			if u.IsUnitInRange(unit) {
				unit.AddEffect(effect)
			} else {
				unit.RemoveEffect(effect)
			}
		}
	}
}
