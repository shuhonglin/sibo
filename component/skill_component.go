package component

import (
	"sibo/entity"
	"reflect"
	"log"
)

type SkillComponent struct {
	SkillMap map[int]entity.Skill
}

func (s SkillComponent) GetType() reflect.Type {
	return reflect.TypeOf(s)
}

func (s SkillComponent) Save2DB() error {
	log.Println("save component to database")
	return nil
}

func (s SkillComponent) GetSkillById(skillId int) (entity.Skill,bool) {
	skill,ok := s.SkillMap[skillId]
	return skill, ok
}
