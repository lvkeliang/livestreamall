package api

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"livestreamall/dao"
	"livestreamall/model"
	"livestreamall/service"
	"livestreamall/util"
	"net/http"
	"regexp"
)

// 用于匹配mail的正则表达式
var mailreg = ".+@.+[.].+"
var mailregexp *regexp.Regexp = regexp.MustCompile(mailreg)

func Register(c *gin.Context) {

	mail := c.PostForm("mail")
	password := c.PostForm("password")
	nickname := c.PostForm("nickname")

	//mailregexp := regexp.MustCompile(mailreg)

	//判断提交的表单中的内容的格式是否符合要求
	//1.正则表达式判断是否是邮箱
	//2.昵称不应大于60字,要么就不填
	//3.密码不应小于6个字符
	if mailregexp.FindAllStringSubmatch(mail, 1) == nil || len(nickname) > 60 || len(password) < 6 {
		util.RespFormatError(c)
		return
	}

	//确保邮箱不重复
	_, err := service.SearchUser("mail", mail)

	//处理除了没查询到以外的错误
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		if err == util.FieldsError {
			//处理数据库查询字段不符的错误
			util.RespFieldsMatchError(c)
		} else if err == nil {
			//处理查询出数据的结果(即邮箱重复)
			util.RespMailRepeated(c)
		} else if err != sql.ErrNoRows {
			//处理意料之外的错误
			util.RespUnexceptedError(c)
			fmt.Println(err)
		}
		return
	}

	//如果输入了昵称，确保昵称不重复
	if len(nickname) > 0 {
		_, err := service.SearchUser("nickname", nickname)

		//处理除了没查询到以外的错误
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			if err == util.FieldsError {
				//处理数据库查询字段不符的错误
				util.RespFieldsMatchError(c)
			} else if err == nil {
				//处理查询出数据的结果(即昵称重复)
				util.RespNicknameRepeated(c)
			} else if err != sql.ErrNoRows {
				//处理意料之外的错误
				util.RespUnexceptedError(c)
			}
			return
		}
	} else {
		//若没输入昵称，则将"用户" + mail作为昵称
		nickname = "用户" + mail
	}

	err = service.CreateUser(model.User{
		Mail:     mail,
		Nickname: nickname,
		Password: password,
	})

	if err != nil {
		fmt.Println(err)
		util.RespCreateUserError(c)
		return
	}
	util.RespOK(c)
}

func Login(c *gin.Context) {
	mail := c.PostForm("mail")
	password := c.PostForm("password")

	//判断提交的表单中的内容的格式是否符合要求
	//1.正则表达式判断是否是邮箱
	//3.密码不应为空
	if mailregexp.FindAllStringSubmatch(mail, 1) == nil || len(password) < 1 {
		util.RespFormatError(c)
		return
	}

	_, _, userpassword, err := service.SearchUserPassword("mail", mail)

	if err != nil {
		if err == sql.ErrNoRows {
			//处理该用户不存在(回复邮箱或密码错误)
			util.RespIncorrectPassword(c)
		} else if err == util.FieldsError {
			//处理数据库查询字段不符的错误
			util.RespFieldsMatchError(c)
		} else if err != sql.ErrNoRows {
			//处理意料之外的错误
			util.RespUnexceptedError(c)
		}
		return
	}

	//处理密码不正确(回复邮箱或密码错误)
	if userpassword != password {
		//fmt.Println(userpassword)
		//fmt.Println(password)
		util.RespIncorrectPassword(c)
		return
	}

	//设置Cookie保持登录
	u, err := service.SearchUser("mail", mail)
	//fmt.Println(u.ID)
	//fmt.Println(strconv.Itoa(int(u.ID)))
	//c.SetCookie("uID", "uID"+strconv.Itoa(int(u.ID)), 604800, "", "/", false, false)

	//设置token
	SetToken(u.Mail, userpassword, c)

}

// GetUserInfo 获取用户信息接口
func GetUserInfo(c *gin.Context) {
	UID, exists := c.Get("uID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// 假设数据库中有一个方法可以根据用户ID获取用户信息
	var user model.User
	err := dao.DB.Where("id = ?", UID).First(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user data"})
		return
	}

	// 返回用户信息
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":         user.ID,
			"nickname":   user.Nickname,
			"email":      user.Mail,
			"created_at": user.CreatedAt,
		},
	})
}
