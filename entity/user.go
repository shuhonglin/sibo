package entity

type User struct {
	Base
	PlayerId  int64
	UserId    int64
	UserToken string
	//Players   string
}

func (u User) GetStructMap() map[string]interface{} {
	return u.Base.GetStructMap(u)
}

/*
func (u User) GetStructFieldNames() []string {
	return u.Base.GetStructFieldNames(u)
}*/
