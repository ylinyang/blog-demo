package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ylinyang/blog-demo/logic"
	"github.com/ylinyang/blog-demo/models"
	"github.com/ylinyang/blog-demo/pkg"
	"go.uber.org/zap"
	"net/http"
	"strings"
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
	aToken, rToken, err := logic.Login(u)
	if err != nil {
		zap.L().Error("用户名或密码错误, ", zap.Error(err))
		Res(c, "用户名或者密码错误")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":          200,
		"message":       "登录成功",
		"token":         aToken,
		"refresh_token": rToken,
	})
}

func RefreshTokenHandler(c *gin.Context) {
	// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
	// 这里假设Token放在Header的Authorization中，并使用Bearer开头
	// 这里的具体实现方式要依据你的实际业务情况决定
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		Res(c, "请求头缺少Auth Token")
		c.Abort()
		return
	}
	// 按空格分割
	parts := strings.SplitN(authHeader, "|", 2)
	//if !(len(parts) == 2 && parts[0] == "Bearer") {
	//	Res(c, "Token格式不对")
	//	c.Abort()
	//	return
	//}
	aToken, _, err := pkg.RefreshToken(parts[0], parts[1])
	if err == pkg.RTokenExpire {
		zap.L().Error("RefreshToken过期, ", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 9999,
			"msg":  "RefreshToken过期",
		})
		return
	}
	if err != nil {
		zap.L().Error("刷新access_token失败, ", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "刷新access_token失败",
		})
		return
	}
	zap.L().Info("新获取到的access_token是, " + aToken)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取新的access_token成功",
		"data": map[string]interface{}{
			"token": aToken,
		},
	})
}
