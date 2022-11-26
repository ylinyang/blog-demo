package logic

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/ylinyang/blog-demo/dao/mysql"
	"github.com/ylinyang/blog-demo/models"
	"github.com/ylinyang/blog-demo/pkg"
	"go.uber.org/zap"
)

// SignUp 相关的业务逻辑代码 数据库操作也是业务逻辑
func SignUp(u *models.ParamsUser) (err error) {
	//	 1. 判断用户是否存在
	if mysql.CheckUserExist(u.Username) {
		zap.L().Info("用户已经存在")
		return errors.New("用户已经存在")
	}

	// 2. 生成uid
	uid := pkg.GetUserID()

	// 3. 入库
	user := &models.User{
		UserId:   uid,
		UserName: u.Username,
		PassWord: fmt.Sprintf("%x", md5.Sum([]byte(u.Password))),
	}
	if err = mysql.InsertUser(user); err != nil {
		zap.L().Error("用户注册失败, ", zap.Error(err))
		return errors.New("用户注册失败,请重试")
	}
	return nil
}
