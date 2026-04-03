package models

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"time"
)

type UuidTokenClaims struct {
	Uuid       string
	Role       string
	Permission []string
	Txts       map[string]string
	jwt.RegisteredClaims
}

func (t *UuidTokenClaims) IfContainPermission(permission string) bool {
	for _, p := range t.Permission {
		if p == permission {
			return true
		}
	}
	return false
}

// 列表内的权限是否都包括
func (t *UuidTokenClaims) IfContainPermissions(permissions []string) bool {
	for _, p := range permissions {
		if t.IfContainPermission(p) {
			continue
		} else {
			return false
		}
	}
	return true
}

func GetUuidToken(key, uuid, role string, permission []string, txts map[string]string, expiresecd int64) (token string, err error) {
	tokenModel := jwt.NewWithClaims(jwt.SigningMethodHS256, UuidTokenClaims{
		uuid,
		role,
		permission,
		txts,
		jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Hour * -24)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	})
	tokenStr, err := tokenModel.SignedString([]byte(key))
	if err != nil {
		fmt.Printf(err.Error())
		return "", err
	}
	return tokenStr, nil
}

func DecodeUuidToken(salt, tokenStr string) (*UuidTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UuidTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(salt), nil
	})
	if err != nil {
		log.Println("错误")
		return &UuidTokenClaims{}, err
	}
	if claims, ok := token.Claims.(*UuidTokenClaims); ok && token.Valid {
		//log.Println(claims["foo"], claims["nbf"])
		return claims, nil
	} else {
		return &UuidTokenClaims{}, fmt.Errorf("jwt decode err")
	}
}

func DecodeUnverifiedUuidToken(tokenStr string) (*UuidTokenClaims, error) {
	token, _ := jwt.ParseWithClaims(tokenStr, &UuidTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(""), nil
	})
	//不校验是否是正确加密的signature
	if token == nil || token.Claims == nil {
		return &UuidTokenClaims{}, fmt.Errorf("token or token.Claims is nil")
	}
	if claims, ok := token.Claims.(*UuidTokenClaims); ok {
		//log.Println(claims["foo"], claims["nbf"])
		return claims, nil
	} else {
		return &UuidTokenClaims{}, fmt.Errorf("jwt decode err，not UuidTokenClaims")
	}
}
