package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ylinyang/blog-demo/logic"
	"github.com/ylinyang/blog-demo/models"
	"go.uber.org/zap"
)

// SignUp 用户注册
func SignUp(c *gin.Context) {
	// 1. 入参 + 校验
	u := new(models.ParamsUser)
	if err := c.ShouldBindJSON(u); err != nil {
		zap.L().Error("signUp with params err: ", zap.Error(err))
		Res(c, "请求参数错误")
		return
	}
	// 2. 业务处理
	if err := logic.SignUp(u); err != nil {
		Res(c, err.Error())
		return
	}

	Res(c, "success", "")
}

// Login 用户登录
func Login(c *gin.Context) {
	u := new(models.ParamsUser)
	if err := c.ShouldBindJSON(u); err != nil {
		zap.L().Error("login with params err: ", zap.Error(err))
		Res(c, "请求参数错误")
		return
	}
	if err := logic.Login(u); err != nil {
		Res(c, err.Error())
		return
	}
	Res(c, "success", "")
}
