package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"time"
)

// TokenExpireDuration token 过期时间
const TokenExpireDuration = time.Hour * 2

// mySecret JWT 签名密钥（生产环境建议走配置）
var mySecret = []byte("老子加盐让你无法解密")

type MyClaims struct {
	// 自定义 JWT 载荷
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// GenToken 生成 JWT token。
func GenToken(userID int64, username, role string) (string, error) {
	c := MyClaims{
		userID,   //自定义字段
		username, //自定义字段
		role,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), //过期时间
			Issuer:    "bluebell",                                 //签发人
		},
	}
	//使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	//使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}

// ParseToken 解析 JWT，并返回自定义 Claims。
func ParseToken(tokenString string) (*MyClaims, error) {
	var myClaims = new(MyClaims) //指针
	token, err := jwt.ParseWithClaims(tokenString, myClaims, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})

	if err != nil {
		return nil, err
	}

	if token.Valid {
		return myClaims, nil
	}
	return nil, errors.New("invalid token")
}

//
