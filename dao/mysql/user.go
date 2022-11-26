package mysql

import (
	"github.com/ylinyang/blog-demo/models"
	"go.uber.org/zap"
)

func CheckUserExist(username string) bool {
	var c int
	if err := db.Get(&c, `select count(*) from user where username = ? `, username); err != nil {
		zap.L().Error("get user failed: ", zap.Error(err))
	}
	return c > 0
}

func InsertUser(p *models.User) (err error) {
	if _, err := db.Exec(`insert into user(user_id,username,password) values (?,?,?)`, p.UserId, p.UserName, p.PassWord); err != nil {
		zap.L().Error("insert user failed: ", zap.Error(err))
	}
	return
}
