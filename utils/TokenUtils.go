package utils

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var (
	TokenUtils = &tokenUtils{}
	jwtkey     = []byte("jelly")
)

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

type tokenUtils struct {
}

//生成token
func (tokenUtils) CreateToken(userId uint) string {
	expireTime := time.Now().Add(10 * 365 * 24 * time.Hour)
	claims := &Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间
			IssuedAt:  time.Now().Unix(),
			Issuer:    "127.0.0.1",  // 签名颁发者
			Subject:   "user token", //签名主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtkey)
	if err != nil {
		fmt.Println(err)
	}
	return tokenString
}

//解析token
func (tokenUtils) GetUserId(tokenString string) (userId uint, err error) {
	token, claims, err := ParseToken(tokenString)
	if err != nil || !token.Valid {
		err = errors.New("解析失败")
		return
	}
	userId = claims.UserId
	return
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	Claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, Claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtkey, nil
	})
	return token, Claims, err
}
