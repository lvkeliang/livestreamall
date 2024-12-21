package model

import "time"

type User struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"ID"`               // 主键，自增
	Mail      string    `gorm:"type:varchar(255);unique;not null" json:"mail"`    // 邮箱，唯一且非空
	Nickname  string    `gorm:"type:varchar(50);unique;not null" json:"nickname"` // 昵称，唯一且非空
	Password  string    `gorm:"type:varchar(255);not null" json:"password"`       // 密码，非空
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`                 // 创建时间
}
