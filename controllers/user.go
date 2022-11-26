package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ylinyang/blog-demo/logic"
	"github.com/ylinyang/blog-demo/models"
	"go.uber.org/zap"
	"net/http"
)

// SignUp 用户注册
func SignUp(c *gin.Context) {
	// 1. 入参 + 校验
	u := new(models.ParamsUser)
	if err := c.ShouldBindJSON(u); err != nil {
		zap.L().Error("signUp with params err: ", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "请求参数错误",
		})
		return
	}
	// 2. 业务处理
	if err := logic.SignUp(u); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
	})
}
