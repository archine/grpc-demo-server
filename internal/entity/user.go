package entity

type User struct {
	Id        int    // 用户ID
	Age       int    // 用户年龄
	Name      string // 用户名
	Email     string // 用户邮箱
	CreatedAt int
}

func (u *User) TableName() string {
	return "users"
}
