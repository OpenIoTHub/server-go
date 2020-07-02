package models

import (
	"fmt"
	"github.com/OpenIoTHub/utils/net/ip"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"log"
	"time"
)

type TokenClaims struct {
	RunId      string
	Host       string
	TcpPort    int
	KcpPort    int
	TlsPort    int
	GrpcPort   int
	UDPApiPort int
	KCPApiPort int
	Permission int
	jwt.StandardClaims
}

func GetToken(gatewayConfig *GatewayConfig, permission int, expiresecd int64) (token string, err error) {
	fmt.Println("Get Token:")
	fmt.Println(
		gatewayConfig.LastId,
		gatewayConfig.Server.ServerHost,
		gatewayConfig.Server.TcpPort,
		gatewayConfig.Server.KcpPort,
		gatewayConfig.Server.TlsPort,
		gatewayConfig.Server.GrpcPort,
		gatewayConfig.Server.UdpApiPort,
		gatewayConfig.Server.KcpApiPort,
	)
	tokenModel := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaims{
		gatewayConfig.LastId,
		gatewayConfig.Server.ServerHost,
		gatewayConfig.Server.TcpPort,
		gatewayConfig.Server.KcpPort,
		gatewayConfig.Server.TlsPort,
		gatewayConfig.Server.GrpcPort,
		gatewayConfig.Server.UdpApiPort,
		gatewayConfig.Server.KcpApiPort,
		permission,
		jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 8*60*60,
			ExpiresAt: time.Now().Unix() + expiresecd,
		},
	})
	tokenStr, err := tokenModel.SignedString([]byte(gatewayConfig.Server.LoginKey))
	if err != nil {
		fmt.Printf(err.Error())
		return "", err
	}
	return tokenStr, nil
}

func GetTokenByServerConfig(serverConfig *ServerConfig, permission int, expiresecd int64) (gatewayToken, openIoTHubToken string, err error) {
	uuidStr := uuid.Must(uuid.NewV4()).String()
	myPublicIp, err := ip.GetMyPublicIpInfo()
	if err != nil {
		return "", "", err
	}
	gatewayConfig := &GatewayConfig{
		ConnectionType: "tcp",
		LastId:         uuidStr,
		GrpcPort:       1082,
		Server: &Srever{
			ServerHost: myPublicIp,
			TcpPort:    serverConfig.Common.TcpPort,
			KcpPort:    serverConfig.Common.KcpPort,
			UdpApiPort: serverConfig.Common.UdpApiPort,
			KcpApiPort: serverConfig.Common.KcpApiPort,
			TlsPort:    serverConfig.Common.TlsPort,
			GrpcPort:   serverConfig.Common.GrpcPort,
			LoginKey:   serverConfig.Security.LoginKey,
		},
	}
	gatewayToken, err = GetToken(gatewayConfig, 1, expiresecd)
	if err != nil {
		return "", "", err
	}
	openIoTHubToken, err = GetToken(gatewayConfig, 2, expiresecd)
	if err != nil {
		return "", "", err
	}
	return gatewayToken, openIoTHubToken, err
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
