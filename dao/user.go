package dao

import (
	"errors"
	"gorm.io/gorm"
	"livestreamall/model"
	"livestreamall/util"
	"strconv"
)

// SearchUser 用于按照具有唯一性的字段来查找用户信息(不包括密码)
// attributes: 用于查找user的字段,该字段必须有唯一性[uid,mail,nickname]
// value：字段的值
func SearchUser(attribute string, value string) (user model.User, err error) {
	query := DB.Model(&model.User{})
	switch attribute {
	case "uID":
		query = query.Where("id = ?", value)
	case "mail":
		query = query.Where("mail = ?", value)
	case "nickname":
		query = query.Where("nickname = ?", value)
	default:
		return model.User{}, util.FieldsError // 返回自定义错误
	}

	err = query.First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.User{}, gorm.ErrRecordNotFound // 用户不存在时返回空值
	}
	return
}

// SearchUserPassword 用于按照具有唯一性的字段来查找用户密码
// attributes: 用于查找user的字段,该字段必须有唯一性[uid,mail,nickname]
// value：字段的值
func SearchUserPassword(attribute string, value string) (uID string, password string, err error) {
	var user model.User

	query := DB.Model(&model.User{})
	switch attribute {
	case "uID":
		query = query.Where("id = ?", value)
	case "mail":
		query = query.Where("mail = ?", value)
	case "nickname":
		query = query.Where("nickname = ?", value)
	default:
		return "", "", util.FieldsError
	}

	err = query.Select("id", "password").First(&user).Error
	if err != nil {
		return "", "", err
	}

	return strconv.Itoa(int(user.ID)), user.Password, nil
}

func CreateUser(user model.User) (err error) {
	err = DB.Create(&user).Error
	return
}

// UserProfile 更新用户信息
func UserProfile(user model.User) (err error) {
	// 仅更新指定字段
	err = DB.Model(&model.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"nickname": user.Nickname,
	}).Error
	return
}
