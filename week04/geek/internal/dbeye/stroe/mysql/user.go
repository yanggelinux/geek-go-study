package mysql

import (
	"time"
)

func NewUser() *User {
	return &User{}
}

type User struct {
	ID       uint32 `json:"id" gorm:"column:id"`
	UserName string `json:"user_name" gorm:"column:user_name"`
	FullName string `json:"full_name" gorm:"column:full_name"`
	*Model
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) GetUserByID(id uint32) (error, *User) {
	user := &User{}
	err = DB.First(&user, id).Error
	return err, user
}

func (u *User) CreateUser(user *User) error {
	err := DB.Create(user).Error
	return err
}

func (u *User) UpdateUser(id uint32, user *User) error {
	//通过结构体变量更新字段值, gorm库会忽略零值字段。就是字段值等于0, nil, "", false这些值会被忽略掉，不会更新。
	// 如果想更新零值，可以使用map类型替代结构体,也可以结构体的类型是指针类型。
	condition := make(map[string]interface{})
	condition["id"] = id
	condition["is_deleted"] = 0
	err := DB.Model(&User{}).Where(condition).Updates(user).Error
	return err

}

func (u *User) UpdateUse2(user *User) error {
	err := DB.Save(user).Error
	return err
}

func (u *User) DeleteUser(user *User) error {
	nowTime := uint64(time.Now().UnixNano() / 1e6)
	user.IsDeleted = 1
	user.DeletedTime = nowTime
	err := DB.Save(user).Error
	return err
}

func (u *User) DeleteUser2(id uint32) error {
	condition := make(map[string]interface{})
	condition["is_deleted"] = 0
	condition["id"] = id
	nowTime := uint64(time.Now().UnixNano() / 1e6)
	user := User{Model: &Model{IsDeleted: 1, DeletedTime: nowTime}}
	err := DB.Model(&User{}).Where(condition).Updates(user).Error
	return err
}
