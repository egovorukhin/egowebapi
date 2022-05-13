package egowebapi

import (
	"github.com/egovorukhin/egowebapi/security"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
)

type Context struct {
	Identity *security.Identity
	Swagger  Swagger
	Session  *Session
	//View     *View
	IContext
}

type Session struct {
	Key      string
	Value    string
	Created  time.Time
	LastTime time.Time
}

type View struct {
	Filename string
	Filepath string
	Layout   string
}

type IContext interface {
	Render(name string, data interface{}, layouts ...string) error
	Params(key string, defaultValue ...string) string
	Get(key string, defaultValue ...string) string
	Set(key string, value string)
	SendStatus(code int) error
	Send(code int, contentType string, b []byte) error
	SendString(code int, s string) error
	SendFile(file string) error
	SaveFile(fileHeader *multipart.FileHeader, path string) error
	SendStream(code int, contentType string, stream io.Reader) error
	Cookies(key string) string
	SetCookie(cookie *http.Cookie)
	ClearCookie(key string)
	Redirect(location string, status int) error
	Path() string
	JSON(code int, data interface{}) error
	Body() []byte
	BodyParser(out interface{}) error
	QueryParam(name string, defaultValue ...string) string
	QueryValues() url.Values
	QueryParams(func(key, value string))
	Hostname() string
	FormValue(name string) string
	FormFile(name string) (*multipart.FileHeader, error)
	Scheme() string
	MultipartForm() (*multipart.Form, error)
}

func NewContext(c IContext) *Context {
	return &Context{IContext: c}
}
