package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ylinyang/blog-demo/pkg"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

const ContextUserIDKey = "UserID"

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "请求头中auth为空",
			})
			c.Abort()
			return
		}
		// 按空格分割 由于前端直接返回的token无需处理
		parts := strings.SplitN(authHeader, " ", 2)
		//fmt.Println(parts)
		//if !(len(parts) == 2 && parts[0] == "Bearer") {
		//	c.JSON(http.StatusOK, gin.H{
		//		"code": -1,
		//		"msg":  "请求头中auth格式有误",
		//	})
		//	c.Abort()
		//	return
		//}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := pkg.ParseToken(parts[0])
		// 需要对过期的token返给前端401
		if err != nil {
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&(jwt.ValidationErrorExpired) != 0 {
					zap.L().Info("token过期")
					c.JSON(http.StatusUnauthorized, gin.H{
						"code": 8888,
						"msg":  "toKen过期",
					})
					c.Abort()
					return
				}
			}
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "无效的Token",
			})
			c.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set(ContextUserIDKey, mc.UserId)
		c.Next() // 后续的处理函数可以用过c.Get(ContextUserIDKey)来获取当前请求的用户信息
	}
}

// 后续可以在gin中获取到uid 进而获取用户信息

var ErrorUserNotLogin = errors.New("用户未登录")

func GetCurrentUser(c *gin.Context) (userId int64, err error) {
	v, _ := c.Get(ContextUserIDKey)
	userId, ok := v.(int64)
	if !ok {
		return userId, ErrorUserNotLogin
	}
	return
}
