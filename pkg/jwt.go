package pkg

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"time"
)

// aTokenExpireDuration 过期时间
const aTokenExpireDuration = time.Second * 5

var RTokenExpire = errors.New("rToken过期")

// secret 用于加盐的字符串
var secret = []byte("www.baidu.com")

// CustomClaims 自定义声明类型 并内嵌jwt.RegisteredClaims
// jwt包自带的jwt.RegisteredClaims只包含了官方字段
// 假设我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type CustomClaims struct {
	// 可根据需要自行添加字段
	UserId               int64 `json:"user_id"`
	jwt.RegisteredClaims       // 内嵌标准的声明
}

// GenToken 生成JWT
func GenToken(userId int64) (aToken, rToken string, err error) {
	// 创建一个我们自己的声明
	claims := CustomClaims{
		userId, // 自定义字段
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(aTokenExpireDuration)),
			Issuer:    "my-project", // 签发人
		},
	}
	// 加密并获得完整的编码后的字符串aToken
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)

	// refresh token 不需要存任何自定义数据
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 15)), // 过期时间
		Issuer:    "bluebell",                                           // 签发人
	}).SignedString(secret)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return
}

// ParseToken 解析access_token
func ParseToken(tokenString string) (claims *CustomClaims, err error) {
	// 解析token
	var token *jwt.Token
	claims = new(CustomClaims)
	token, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return secret, nil
	})
	if err != nil {
		return
	}
	if !token.Valid { // 校验token
		err = errors.New("invalid token")
	}
	return
}

// RefreshToken 刷新access_token
func RefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
	// refresh token无效直接返回
	if _, err = jwt.Parse(rToken, func(token *jwt.Token) (i interface{}, err error) {
		// 直接使用标准的Claim则可以直接使用Parse方法
		return secret, nil
	}); err != nil {
		zap.L().Error("rToken解析失败, ", zap.Error(err))
		return "", "", RTokenExpire
	}

	// 从旧access token中解析出claims数据	解析出payload负载信息
	var claims CustomClaims
	_, err = jwt.ParseWithClaims(aToken, &claims, func(token *jwt.Token) (i interface{}, err error) {
		return secret, nil
	})
	// aToken过期
	zap.L().Error("从旧access token中解析出claims数据解析失败, ", zap.Error(err))
	v, _ := err.(*jwt.ValidationError)

	// 当access token是过期错误 并且 refresh token没有过期时就创建一个新的access token
	if v.Errors == jwt.ValidationErrorExpired {
		return GenToken(claims.UserId)
	}
	return
}
