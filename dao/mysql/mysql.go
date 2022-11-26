package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" // 注意驱动导入
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var db *sqlx.DB

func Init() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.dbname"))
	// 也可以使用MustConnect连接不成功就panic
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Debug("connect DB failed, err:%v\n", zap.Error(err))
		return
	}
	db.SetMaxOpenConns(viper.GetInt("mysql.MaxOpenConns"))
	db.SetMaxIdleConns(viper.GetInt("mysql.MaxIdleConns"))
	return
}

func Close() {
	_ = db.Close()
}
