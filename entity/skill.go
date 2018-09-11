package entity

type Skill struct {
	Base
	PlayerId int64
	SkillId  int
	Hole     int
}

func NewSkill(playerId int64, skillId int) *Skill {
	return &Skill{PlayerId: playerId, SkillId: skillId}
}
