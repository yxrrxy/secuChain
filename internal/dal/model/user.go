package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Username    string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password    string         `gorm:"size:100;not null" json:"-"` // 不返回密码
	Email       string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Status      string         `gorm:"size:20;default:active" json:"status"` // active, inactive, banned
	LastLoginAt *time.Time     `json:"last_login_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
