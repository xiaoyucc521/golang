package model

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"column:id;type:int unsigned;not null;primaryKey;autoIncrement;comment:主键ID" json:"id"`
	Username  string    `gorm:"column:username;type:varchar(20);not null;unique;comment:用户名" json:"username"`
	Nickname  string    `gorm:"column:nickname;type:varchar(20);not null;default:'';comment:昵称" json:"nickname"`
	Password  string    `gorm:"column:password;type:varchar(100);not null;comment:密码"  json:"password"`
	Mobile    string    `gorm:"column:mobile;type:varchar(20);not null;unique;default:'';comment:手机号"  json:"mobile"`
	Avatar    string    `gorm:"column:avatar;type:varchar(255);not null;default:'';comment:头像"  json:"avatar"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间"  json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
}

func (*User) TableName() string {
	return "user"
}
