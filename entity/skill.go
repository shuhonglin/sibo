package entity

type Skill struct {
	playerId int64
	skillId int
	Hole int
}

func NewSkill(playerId int64, skillId int) *Skill {
	return &Skill{playerId:playerId, skillId:skillId}
}

func (skill Skill) PlayerId() int64 {
	return skill.playerId
}

func (skill Skill) SkillId() int {
	return skill.skillId
}
