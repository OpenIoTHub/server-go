package models

import (
	"fmt"
	"github.com/OpenIoTHub/utils/net/ip"
	"github.com/golang-jwt/jwt"
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
	Permission []string
	Txts       map[string]string
	jwt.StandardClaims
}

func (t *TokenClaims) IfContainPermission(permission string) bool {
	for _, p := range t.Permission {
		if p == permission {
			return true
		}
	}
	return false
}

// 列表内的权限是否都包括
func (t *TokenClaims) IfContainPermissions(permissions []string) bool {
	for _, p := range permissions {
		if t.IfContainPermission(p) {
			continue
		} else {
			return false
		}
	}
	return true
}

func GetToken(loginWithServer *LoginWithServer, permission []string, expiresecd int64) (token string, err error) {
	tokenModel := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaims{
		loginWithServer.LastId,
		loginWithServer.Server.ServerHost,
		loginWithServer.Server.TcpPort,
		loginWithServer.Server.KcpPort,
		loginWithServer.Server.TlsPort,
		loginWithServer.Server.GrpcPort,
		loginWithServer.Server.UdpApiPort,
		loginWithServer.Server.KcpApiPort,
		permission,
		map[string]string{},
		jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 8*60*60,
			ExpiresAt: time.Now().Unix() + expiresecd,
		},
	})
	tokenStr, err := tokenModel.SignedString([]byte(loginWithServer.Server.LoginKey))
	if err != nil {
		fmt.Printf(err.Error())
		return "", err
	}
	return tokenStr, nil
}

func GetTokenByServerConfig(serverConfig *ServerConfig, expiresecd int64) (gatewayToken, openIoTHubToken string, err error) {
	var myPublicIp string
	uuidStr := uuid.Must(uuid.NewV4()).String()
	if serverConfig.PublicIp != "" {
		myPublicIp = serverConfig.PublicIp
	} else {
		myPublicIp, err = ip.GetMyPublicIpv4()
		if err != nil {
			return "", "", err
		}
	}

	loginWithServer := &LoginWithServer{
		ConnectionType: "tcp",
		LastId:         uuidStr,
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
	gatewayToken, err = GetToken(loginWithServer, []string{PermissionGatewayLogin}, expiresecd)
	if err != nil {
		return "", "", err
	}
	openIoTHubToken, err = GetToken(loginWithServer, []string{PermissionOpenIoTHubLogin}, expiresecd)
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
