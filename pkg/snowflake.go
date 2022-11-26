package pkg

import (
	"github.com/bwmarrin/snowflake"
	"github.com/spf13/viper"
	"time"
)

// 雪花算法
var node *snowflake.Node

func Init() (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", viper.GetString("app.startTime"))
	if err != nil {
		return
	}
	snowflake.Epoch = st.UnixNano() / 1e6
	node, err = snowflake.NewNode(viper.GetInt64("app.machineID"))
	return
}

// GetUserID 生成用户唯一ID
func GetUserID() int64 {
	return node.Generate().Int64()
}
