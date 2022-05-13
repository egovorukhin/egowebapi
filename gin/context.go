package gin

/*type Context struct {
	Ctx *gin.Context
}

func (c *Context) Render(name string, data interface{}, layouts ...string) error {
	c.Ctx.HTML(200, name, data)
	return nil
}

func (c *Context) Params(key string, defaultValue ...string) string {
	value := c.Ctx.Param(key)
	if value == "" && defaultValue != nil {
		return defaultValue[0]
	}
	return value
}

func (c *Context) Get(key string, defaultValue ...string) string {
	value := c.Ctx.GetHeader(key)
	if value == "" && defaultValue != nil {
		return defaultValue[0]
	}
	return value
}

func (c *Context) Set(key, value string) {
	c.Ctx.Header(key, value)
}

func (c *Context) SendStatus(code int) error {
	c.Ctx.Status(code)
	return nil
}

func (c *Context) Cookies(key string) string {
	value, _ := c.Ctx.Cookie(key)
	return value
}

func (c *Context) SetCookie(cookie *consts.Cookie) {
	c.Ctx.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
}

// TODO ClearCookie
func (c *Context) ClearCookie(key string) {
	for _, cookie := range c.Ctx..Cookies() {
		if cookie.Name == key {

		}
	}
}

func (c *Context) Redirect(location string, status int) error {
	c.Ctx.Redirect(status, location)
	return nil
}

func (c *Context) Path() string {
	return c.Ctx.FullPath()
}

func (c *Context) SendString(code int, s string) error {
	c.Ctx.Data(code, "", s)
	return nil
}

func (c *Context) Send(code int, contentType string, b []byte) error {
	c.Ctx.Data(code, contentType, b)
	return nil
}

func (c *Context) SendFile(file string) error {
	c.Ctx.File(file)
	return nil
}

func (c *Context) SaveFile(fileHeader *multipart.FileHeader, path string) error {
	return c.Ctx.SaveFile(fileHeader, path)
}

func (c *Context) SendStream(code int, contentType string, stream io.Reader) error {
	return c.Ctx.Stream(code, contentType, stream)
}

func (c *Context) JSON(code int, data interface{}) error {
	c.Ctx.JSON(code, data)
	return nil
}

func (c *Context) Body() []byte {
	body := c.Ctx.Request().Body
	b, _ := ioutil.ReadAll(body)
	defer body.Close()
	return b
}

func (c *Context) BodyParser(out interface{}) error {
	return c.Ctx.Bind(out)
}

func (c *Context) QueryParam(name string, defaultValue ...string) string {
	value := c.Ctx.Query(name)
	if value == "" && defaultValue != nil {
		return defaultValue[0]
	}
	return value
}

func (c *Context) QueryValues() url.Values {
	return c.Ctx.QueryParams()
}

func (c *Context) QueryParams(h func(key, value string)) {
	for k, v := range c.Ctx.QueryParams() {
		s := ""
		if len(v) > 0 {
			s = v[0]
		}
		h(k, s)
	}
}

func (c *Context) Hostname() string {
	return c.Ctx.Host
}

func (c *Context) FormValue(name string) string {
	return c.Ctx.PostForm(name)
}

func (c *Context) FormFile(name string) (*multipart.FileHeader, error) {
	return c.Ctx.FormFile(name)
}

func (c *Context) Scheme() string {
	return c.Ctx.Scheme()
}

func (c *Context) MultipartForm() (*multipart.Form, error) {
	return c.Ctx.MultipartForm()
}
*/
