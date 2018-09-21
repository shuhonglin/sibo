package component

import (
	log "github.com/sirupsen/logrus"
	"reflect"
	"sibo/entity"
	"strings"
	"github.com/jmoiron/sqlx"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

type SkillComponent struct {
	MapComponent

	skillMap map[int]*entity.Skill
}

func (s *SkillComponent) InitComponent(playerId int64) {
	s.dbSaveProxy = s
	s.playerId = playerId
	s.keyPrefix = "skill_"
	s.selectSql = "SELECT * FROM tb_skill WHERE playerId=?"
	s.deleteSql = "DELETE FROM tb_skill WHERE playerId="+s.Key()+" AND skillId IN (?)"
	if s.init == false {
		s.MapComponent.InitComponent(playerId)
		s.skillMap = make(map[int]*entity.Skill)
		s.init = true
	}
}

func (s SkillComponent) GetType() reflect.Type {
	return reflect.TypeOf(s)
}

/*func (s SkillComponent) ID() int64 {
	return s.playerId
}*/

/*func (s *SkillComponent) Save2DB() error {
	for _, v := range s.skillMap {
		v.GetStructMap(v)
	}
	log.Println("save component to database")
	if s.addSet.Intersect(s.delSet).Cardinality() > 0 {
		log.Warn("addSet 与 delSet冲突")
	}
	s.deleteEntityFromDB()
	s.saveNewEntityToDB()
	s.saveUpdateEntityToDB()
	return nil
}*/

func (s *SkillComponent) save2SqlDB() error {
	if len(s.skillMap)>0 {
		tx := SQL_DB.MustBegin()
		for _, v := range s.skillMap {
			if s.insertSql == "" {
				nameValues := v.GetStructMap()
				names := make([]string, 0, len(nameValues))
				values := make([]string, 0, len(nameValues))
				for k := range nameValues {
					val := strings.ToLower(k)
					names = append(names, val)
					values = append(values, ":"+val)
				}
				s.insertSql = "REPLACE INTO tb_skill ("+strings.Join(names, ",")+") VALUES (" + strings.Join(values,",")+")"
			}
			_,err:= tx.NamedExec(s.insertSql, v)
			if err!=nil {
				log.Errorln(err)
			}
		}
		if s.delSet.Cardinality() > 0 {
			query, args, err := sqlx.In(s.deleteSql, s.delSet.ToSlice())
			if err != nil {
				log.Errorln(err)
			}
			query = SQL_DB.Rebind(query)
			tx.Exec(query, args)
		}
		tx.Commit()
		log.Println("save component to database")
	}

	/*if s.addSet.Intersect(s.delSet).Cardinality() > 0 {
		log.Warn("addSet 与 delSet冲突")
	}
	s.deleteEntityFromDB()
	s.saveNewEntityToDB()
	s.saveUpdateEntityToDB()*/
	return nil
}

func (s *SkillComponent) save2NoSqlDB() error {
	log.Info("save skill to nosql db")
	jsonData, _ := json.Marshal(s.skillMap)
	reply, err := REDIS_DB.Get().Do("SET", s.Key(), jsonData)
	log.Println(reply, err)
	return nil
}

func (s *SkillComponent) initFromSqlDB() error {
	// execute selectSql
	log.Println(s.selectSql)

	var skills []entity.Skill
	err := SQL_DB.Select(&skills, s.selectSql, s.playerId)
	if err != nil {
		log.Errorf("初始化玩家{%d}技能错误", s.playerId)
	}
	for _, skill := range skills {
		s.skillMap[skill.SkillId] = &skill
	}
	log.Println("init skill map from sql db")
	return nil
}
func (s *SkillComponent) initFromNoSqlDB() error {
	jsonData, err := redis.Bytes(REDIS_DB.Get().Do("GET", s.Key()))
	if err != nil {
		return err
	}
	json.Unmarshal(jsonData, s.skillMap)
	return nil
}

/*func (s *SkillComponent) InitFromDB() error {
	if s.updateSet == nil {
		s.updateSet = mapset.NewSet()
	}
	if s.addSet == nil {
		s.addSet = mapset.NewSet()
	}
	if s.delSet == nil {
		s.delSet = mapset.NewSet()
	}
	s.loadAllEntityFromDB()
	log.Println("init from db")
	return nil
}*/

func (s *SkillComponent) GetSkill(skillId int) (*entity.Skill, bool) {
	skill, ok := s.skillMap[skillId]
	return skill, ok
}

func (s *SkillComponent) DeleteSkill(skillId int) {
	s.delSet.Add(skillId)
	if s.addSet.Contains(skillId) {
		s.addSet.Remove(skillId)
	}
	if s.updateSet.Contains(skillId) {
		s.updateSet.Remove(skillId)
	}
	delete(s.skillMap, skillId)
}

func (s *SkillComponent) AddSkill(skill *entity.Skill) {
	s.addSet.Add(skill.SkillId)
	if s.delSet.Contains(skill.SkillId) {
		s.delSet.Remove(skill.SkillId)
	}
	if s.updateSet.Contains(skill.SkillId) {
		s.updateSet.Remove(skill.SkillId)
	}
	s.skillMap[skill.SkillId] = skill
}

func (s *SkillComponent) SaveSkill(skill *entity.Skill) {
	_, ok := s.skillMap[skill.SkillId]
	if !ok {
		s.AddSkill(skill)
	} else {
		s.updateSet.Add(skill.SkillId)
	}
}

func (s SkillComponent) IsInit() bool {
	return s.init
}

func (s *SkillComponent) SetInit(init bool) {
	s.init = init
}

/*func (s *SkillComponent) loadAllEntityFromDB() {
	selectSql := "SELECT * FROM tb_skill where playerId = ?"
	if s.skillMap == nil {
		s.skillMap = make(map[int]*entity.Skill)
	}
	// execute selectSql
	log.Println(selectSql)

	skills := []entity.Skill{}
	err := SQL_DB.Select(&skills, selectSql, s.playerId)
	if err != nil {
		log.Errorf("初始化玩家{%d}技能错误", s.playerId)
	}
	for _, skill := range skills {
		s.skillMap[skill.SkillId] = &skill
	}
}*/

/*func (s *SkillComponent) saveUpdateEntityToDB() {
	//realUpdateSet := s.updateSet.Difference(s.delSet)
	if s.updateSet.Cardinality() > 0 {
		updateSql := "UPDATE tb_skill SET hole=:hole WHERE playerId=:playerId AND skillId=:skillId"

		tx := SQL_DB.MustBegin()
		for i := range s.updateSet.Iter() {
			skill, ok := s.skillMap[i.(int)]
			if ok { //存在内存中
				//tx.MustExec()
				tx.NamedExec(updateSql, skill)
			}
			s.updateSet.Remove(i)
		}
		tx.Commit()
		// execute updatesql
		log.Println(updateSql)
	}
}*/

/*func (s *SkillComponent) saveNewEntityToDB() {
	//realAddSet := s.addSet.Difference(s.delSet)
	// 可优化转化为批量插入
	if s.addSet.Cardinality() > 0 {
		insertSql := "REPLACE INTO tb_skill(skillId, playerId, hole) VALUES (:skillId, :playerId, :hole)"
		tx := SQL_DB.MustBegin()
		for i := range s.addSet.Iter() {
			skill, ok := s.skillMap[i.(int)]
			if ok { //存在内存中
				tx.NamedExec(insertSql, skill)
			}
			s.addSet.Remove(i)
		}
		tx.Commit()
		// execute insertsql
		log.Println(insertSql)
	}
}*/

/*func (s *SkillComponent) deleteEntityFromDB() {
	if s.delSet.Cardinality() > 0 {
		delSlice := s.delSet.ToSlice()
		delSql := "DELETE FROM tb_skill WHERE playerId= ? AND skillId IN(?)"
		for index, i := range delSlice {
			_, ok := s.skillMap[i.(int)]
			s.delSet.Remove(i)
			if !ok { //不在内存中
				delSlice[index] = -1
			}
		}
		_, err := SQL_DB.Exec(delSql, s.playerId, delSlice)
		if err != nil {
			log.Error("删除玩家技能失败, ", err)
		}
		// execute delSql
		log.Println(delSql)
	}
}*/
