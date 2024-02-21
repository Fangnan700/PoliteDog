package PoliteDog

import (
	"encoding/base64"
	"net/http"
)

/**
Basic 认证
*/

type Account struct {
	Users         map[string]string
	UnAuthHandler func(ctx *Context)
}

func (a *Account) BasicAuthHandler(ctx *Context) {
	username, password, ok := ctx.r.BasicAuth()
	if !ok {
		a.unAuthHandler(ctx)
		return
	}

	pwd, exist := a.Users[username]
	if !exist {
		a.unAuthHandler(ctx)
		return
	}

	if pwd != password {
		a.unAuthHandler(ctx)
		return
	}

	ctx.Next()
}

func (a *Account) unAuthHandler(ctx *Context) {
	if a.UnAuthHandler != nil {
		a.UnAuthHandler(ctx)
	} else {
		ctx.Data(http.StatusUnauthorized, nil)
		ctx.Abort()
	}
}

// GetBasicAuth 获取Basic认证加密结果
func (a *Account) GetBasicAuth(username string, password string) string {
	encodeStr := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	return "Basic " + encodeStr
}
