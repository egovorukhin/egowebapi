package egowebapi

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

const (
	Md5Algorithm           = "MD5"
	Md5SessAlgorithm       = "MD5-sess"
	Sha256Algorithm        = "SHA-256"
	Sha256SessAlgorithm    = "SHA-256-sess"
	Sha512256Algorithm     = "SHA-512-256"
	Sha512256SessAlgorithm = "SHA-512-256-sess"
)

type Auth string

const (
	NoAuth     = "NoAuth"
	BasicAuth  = "Basic"
	DigestAuth = "Digest"
	ApiKeyAuth = "ApiKey"
	OAuth2Auth = "OAuth2"
)

type Authorization struct {
	AllRoutes    string
	Unauthorized ErrorHandler
	Basic        BasicAuthHandler
	Digest       DigestAuthHandler
	ApiKey       *ApiKey
}

type Basic struct {
	Handler      BasicAuthHandler
	Unauthorized ErrorHandler
}

func (b BasicAuthHandler) parseBasicAuth(auth string) (username, password string, ok bool) {
	const prefix = "Basic "
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	i := strings.IndexByte(cs, ':')
	if i < 0 {
		return
	}
	return cs[:i], cs[i+1:], true
}

func (b BasicAuthHandler) Do(c *Context) (err error) {

	err = errors.New(`Basic realm="Необходимо указать имя пользователя и пароль"`)
	auth := c.Get(HeaderAuthorization)
	if auth == "" {
		return
	}

	username, password, ok := b.parseBasicAuth(auth)
	if !ok || !b(username, password) {
		return
	}

	c.Identity = &Identity{
		Username: username,
		AuthName: BasicAuth,
	}

	return nil
}

/*func (b *Basic) Do(handler Handler, isPermission bool, permission *Permission) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		auth := ctx.Get(HeaderAuthorization)
		if auth == "" {
			return b.realm(ctx)
		}

		username, password, ok := b.parseBasicAuth(auth)
		if !ok || !b.Handler(username, password) {
			return b.realm(ctx)
		}

		// Получаем путь
		route := ctx.Route()
		// Проверяем на существование PermissionHandler
		if isPermission && permission != nil && route != nil {
			if !permission.Handler(username, route.Path) {
				if b.Unauthorized != nil {
					return b.Unauthorized(ctx, StatusForbidden)
				}
				return ctx.SendStatus(StatusForbidden)
			}
		}
		domain := ""
		a := strings.Split(username, `\`)
		if len(a) > 1 {
			domain = a[0]
			username = a[1]
		}

		c := &Context{
			Identity: &Identity{
				Username:   username,
				Domain: domain,
			},
		}

		// Возвращаем данные по пользователю и маршруту
		return handler(c)
	}
}*/

/*type Digest struct {
	Handler      DigestAuthHandler
	Unauthorized ErrorHandler
}*/

type Advanced struct {
	Realm       string
	Nonce       string
	Algorithm   string
	Qop         string
	NonceCount  string
	ClientNonce string
	Opaque      string
}

func (d DigestAuthHandler) Do(c *Context) (err error) {

	username := ""

	c.Identity = &Identity{
		Username: username,
	}

	return
}

type ApiKey struct {
	KeyName  string
	IsHeader bool
	Handler  ApiKeyAuthHandler
}

func (a ApiKey) Do(c *Context) (err error) {

	value := ""
	param := ""
	if a.IsHeader {
		// Пытаемся получить из заголовка токен
		value = c.Get(a.KeyName)
		param = "header"
	} else {
		// Если не нашли в заголовке, то ищем в переменных запроса адресной строки
		value = c.QueryParam(a.KeyName)
		param = "query"
	}

	if value == "" {
		return errors.New(fmt.Sprintf("Not found token by [%s]", param))
	}

	username := ""
	if a.Handler != nil {
		username, err = a.Handler(value)
	}

	c.Identity = &Identity{
		Username: username,
		AuthName: ApiKeyAuth,
	}

	return
}

func (a ApiKey) Get() (name, param string) {
	param = "query"
	if a.IsHeader {
		param = "header"
	}
	name = a.KeyName
	return
}
