package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/ylinyang/blog-demo/controllers"
	_ "github.com/ylinyang/blog-demo/docs"
	"github.com/ylinyang/blog-demo/logger"
	"github.com/ylinyang/blog-demo/middlewares"
	"net/http"
)

func SetUp() (r *gin.Engine) {
	// gin默认的日志为debug模式 如果设置为release模式 将在不在打印日志  但是zap的日志还是会输出到文件中
	if viper.GetString("app.mode") == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r = gin.New()

	//	将zap日志中间间注册到gin中
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// swagger页面
	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//	 注册路由信息

	v1 := r.Group("/api/v1")
	// 用户模块
	v1.POST("/signup", controllers.SignUp)
	v1.POST("/login", controllers.Login)

	//v1.GET("/refresh_token", controllers.RefreshTokenHandler)
	v1.Use(middlewares.JWTAuthMiddleware()) // 应用JWT认证中间件

	{
		v1.GET("/community", func(c *gin.Context) {
			c.String(http.StatusOK, fmt.Sprintf("%s", viper.GetString("app.name")))
			fmt.Println("1")
		})
		v1.GET("/ping1", func(c *gin.Context) {
			c.String(http.StatusOK, fmt.Sprintf("%s", "xxxx"))
		})
	}

	return
}
