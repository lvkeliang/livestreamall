package model

import "time"

// Message 表示聊天记录模型
type Message struct {
	ID         uint      `gorm:"primaryKey"`         // 主键
	UserID     uint      `gorm:"not null"`           // 用户ID
	LiveRoomID uint      `gorm:"not null"`           // 直播间ID
	Content    string    `gorm:"type:text;not null"` // 消息内容
	CreatedAt  time.Time `gorm:"autoCreateTime"`     // 创建时间

	// 外键关系
	User     User     `gorm:"foreignKey:UserID"`     // 关联用户
	LiveRoom LiveRoom `gorm:"foreignKey:LiveRoomID"` // 关联直播间
}
