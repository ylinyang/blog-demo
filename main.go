package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"github.com/ylinyang/blog-demo/dao/mysql"
	"github.com/ylinyang/blog-demo/dao/redis"
	"github.com/ylinyang/blog-demo/logger"
	"github.com/ylinyang/blog-demo/pkg"
	"github.com/ylinyang/blog-demo/routers"
	"github.com/ylinyang/blog-demo/settings"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	/*
		1. 配置加载  本地或者远程加载
		2. 日志初始化
		3. 初始化第三方服务连接 mysql、redis等服务
		4. 注册路由
		5. 启动服务  优雅关机
	*/
	if err := settings.Init(); err != nil {
		log.Println("init config failed:", err)
		return
	}

	if err := logger.Init(); err != nil {
		log.Println("init logger failed:", err)
		return
	}
	// 将延迟注册
	defer zap.L().Sync()

	// 释放连接
	defer mysql.Close()
	if err := mysql.Init(); err != nil {
		log.Println("init mysql failed:", err)
		return
	}

	defer redis.Close()
	if err := redis.Init(); err != nil {
		log.Println("init redis failed:", err)
		return
	}

	if err := pkg.Init(); err != nil {
		log.Println("init snowflake failed:", err)
		return
	}

	zap.L().Debug("init success!!!")

	// 路由注册
	r := routers.SetUp()

	//	优雅关机
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s", viper.GetString("app.port")),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Error("listen: %s\n", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Error("Server Shutdown: ", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}
