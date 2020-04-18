package models

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

//func init() {
//	Salt = config.DefaultConnectToServerKey
//}

//var Salt string

type TokenClaims struct {
	RunId      string
	Host       string
	TcpPort    int
	KcpPort    int
	TlsPort    int
	P2PApiPort int
	Permission int
	jwt.StandardClaims
}

func GetToken(salt, id, host string, tcpPort, kcpPort, tlsPort, p2pApiPort, permission int, expiresecd int64) (token string, err error) {
	tokenModel := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaims{
		id,
		host,
		tcpPort,
		kcpPort,
		tlsPort,
		p2pApiPort,
		permission,
		jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 8*60*60,
			ExpiresAt: time.Now().Unix() + expiresecd,
		},
	})
	tokenStr, err := tokenModel.SignedString([]byte(salt))
	if err != nil {
		fmt.Printf(err.Error())
		return "", err
	}
	return tokenStr, nil
}

func DecodeToken(salt, tokenStr string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(salt), nil
	})
	if err != nil {
		log.Println("错误")
		return &TokenClaims{}, err
	}
	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		//log.Println(claims["foo"], claims["nbf"])
		return claims, nil
	} else {
		return &TokenClaims{}, fmt.Errorf("jwt decode err")
	}
}

func DecodeUnverifiedToken(tokenStr string) (*TokenClaims, error) {
	token, _ := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(""), nil
	})
	//不校验是否是正确加密的signature
	if token == nil || token.Claims == nil {
		return &TokenClaims{}, fmt.Errorf("token or token.Claims is nil")
	}
	if claims, ok := token.Claims.(*TokenClaims); ok {
		//log.Println(claims["foo"], claims["nbf"])
		return claims, nil
	} else {
		return &TokenClaims{}, fmt.Errorf("jwt decode err，not TokenClaims")
	}
}
