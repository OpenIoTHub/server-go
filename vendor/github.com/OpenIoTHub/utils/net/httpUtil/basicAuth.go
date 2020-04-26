package httpUtil

import (
	"encoding/base64"
	"net/http"
	"strings"
)

//http basic Auth
func unAuth(w http.ResponseWriter) bool {
	w.Header().Set("WWW-Authenticate", `Basic realm="Dotcoo User Login"`)
	w.WriteHeader(http.StatusUnauthorized)
	return false
}

func Auth(w http.ResponseWriter, r *http.Request, uname string, pwd string) bool {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return unAuth(w)
	}
	auths := strings.SplitN(auth, " ", 2)
	if len(auths) != 2 {
		return unAuth(w)
	}
	authMethod := auths[0]
	authB64 := auths[1]
	switch authMethod {
	case "Basic":
		authstr, err := base64.StdEncoding.DecodeString(authB64)
		if err != nil {
			return unAuth(w)
		}
		userPwd := strings.SplitN(string(authstr), ":", 2)
		if len(userPwd) != 2 {
			return unAuth(w)
		}
		username := userPwd[0]
		userpwd := userPwd[1]
		if username == uname && userpwd == pwd {
			return true
		}
	default:
		return unAuth(w)
	}
	return unAuth(w)
}
