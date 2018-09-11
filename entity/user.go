package entity

type User struct {
	Base
	UserId    int64
	UserToken string
	Players   []int64
}

func (u User) GetStructMap() map[string]interface{} {
	return u.Base.GetStructMap(u)
}
