package models

type ParamsUser struct {
	Username string `json:"username" binding:"required"` // binding:"required" 表示必传参数 binding为固定key
	Password string `json:"password" binding:"required"`
}
