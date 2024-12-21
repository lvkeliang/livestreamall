package model

import "time"

// LiveRoom 表示直播间模型
type LiveRoom struct {
	ID          uint      `gorm:"primaryKey"`                        // 主键
	Title       string    `gorm:"type:varchar(255);not null"`        // 直播间标题
	StreamName  string    `gorm:"type:varchar(255);unique;not null"` // 推流名称，唯一
	Description string    `gorm:"type:text"`                         // 直播间描述
	IsLive      bool      `gorm:"default:false"`                     // 是否正在直播
	CreatedAt   time.Time `gorm:"autoCreateTime"`                    // 创建时间
	UserID      string    `gorm:"not null"`                          // 用户ID，标识直播间所属的用户

	// 关联消息
	Messages []Message `gorm:"foreignKey:LiveRoomID"` // 直播间中的消息
}
