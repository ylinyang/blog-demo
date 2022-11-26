package settings

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Init() (err error) {
	// 1. 直接指定文件路径
	//viper.SetConfigFile("./config/config.yaml")

	// 2. 指定配置文件的名和配置文件的位置 viper自行查找
	viper.SetConfigName("config")    // 配置文件名称(无扩展名)
	viper.AddConfigPath("./config/") // 配置文件存储路径
	//viper.SetConfigType("yaml")      // 用于远程获取配置文件的时候时候指定文件类型的  远程配置中心

	err = viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {            // 处理读取配置文件的错误
		return err
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config file changed: ", in.Name)
	})
	return
}
