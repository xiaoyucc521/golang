package util

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var jwtSecret = []byte("server.jwtSecret")

// Claims 自定义声明类型 并内嵌jwt.RegisteredClaims
// jwt包自带的jwt.RegisteredClaims只包含了官方字段
// 假设我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	// 内嵌标准的声明
	//jwt.RegisteredClaims
	jwt.MapClaims
}

// GenerateToken 生成 token
func GenerateToken(UserID int64, Username string) (string, error) {

	claims := Claims{
		UserID:   UserID,
		Username: Username,
		//RegisteredClaims: jwt.RegisteredClaims{
		//	// 定义过期时间
		//	ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		//	// 发布时间
		//	IssuedAt: jwt.NewNumericDate(time.Now()),
		//	// 签发人
		//	Issuer: "XiaoYu",
		//},
		MapClaims: jwt.MapClaims{
			// 定义过期时间
			"exp": time.Now().Add(24 * time.Hour).Unix(),
			// 发布时间
			"iat": time.Now().Unix(),
			// 签发人
			"iss": "XiaoYu",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

// ParseToken token 解析
func ParseToken(tokenString string) (*Claims, error) {
	// 解析token
	//var claims = new(Claims)

	tokenClaims, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
