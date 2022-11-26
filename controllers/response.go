package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
	返回前端 code 200 正常 异常都是 -1
*/

func Res(c *gin.Context, args ...any) {
	// 1. args 长度 1 为报错信息 其他为正常返回
	if len(args) == 1 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  args[0],
			"data": "",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  args[0],
		"data": args[1],
	})
}
