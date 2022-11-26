package models

type ParamsUser struct {
	UserName string `json:"username" binding:"required"` // binding:"required" 表示必传参数 binding为固定key
	PassWord string `json:"password" binding:"required"`
}
