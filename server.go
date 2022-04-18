package egowebapi

import (
	"fmt"
	"github.com/egovorukhin/egowebapi/security"
	"github.com/invopop/jsonschema"
	p "path"
	"regexp"
	"strings"
)

const (
	Name    = "EgoWebApi"
	Version = "v0.2.6"
)

type Server struct {
	Config      Config
	IsStarted   bool
	WebServer   IServer
	Controllers []*Controller
	Swagger     *Swagger
}

type IServer interface {
	Start(addr string) error
	StartTLS(addr, cert, key string) error
	Stop() error
	Static(prefix, root string)
	Any(path string, handler interface{})
	Use(params ...interface{})
	Add(method, path string, handler Handler)
	GetApp() interface{}
	NotFoundPage(path, page string)
	ConvertParam(param string) string
}

type Suffix struct {
	Index       int
	Value       string
	isParam     bool
	Description string
}

func NewSuffix(suffix ...Suffix) (s []Suffix) {
	for _, item := range suffix {
		s = append(s, item)
	}
	return
}

func New(server IServer, config Config) *Server {

	// Устанавливаем статические файлы
	if config.Static != nil {
		server.Static(config.Static.Prefix, config.Static.Root)
	}

	s := &Server{
		Config:    config,
		WebServer: server,
		Swagger: &Swagger{
			Swagger:             "2.0",
			Host:                fmt.Sprintf("localhost:%d", config.Port),
			BasePath:            "/",
			SecurityDefinitions: SecurityDefinitions{},
			Paths:               Paths{},
			Definitions:         map[string]*jsonschema.Schema{},
		},
	}

	return s
}

// GetWebServer вернуть интерфейс веб сервера
func (s *Server) GetWebServer() interface{} {
	return s.WebServer.GetApp()
}

// Start запуск сервера
func (s *Server) Start() (err error) {

	for _, c := range s.Controllers {

		c.initialize(s.Swagger.BasePath)

		// Добавляем тэги контроллера
		if c.IsShow {
			s.Swagger.Tags = append(s.Swagger.Tags, c.Tag)
		}

		// Проверка интерфейса на соответствие
		if i, ok := c.Interface.(IGet); ok {
			err = s.get(i, c)
			if err != nil {
				return
			}
		}
		if i, ok := c.Interface.(IPost); ok {
			err = s.post(i, c)
			if err != nil {
				return
			}
		}
		if i, ok := c.Interface.(IPut); ok {
			err = s.put(i, c)
			if err != nil {
				return
			}
		}
		if i, ok := c.Interface.(IDelete); ok {
			err = s.delete(i, c)
			if err != nil {
				return
			}
		}
		if i, ok := c.Interface.(IOptions); ok {
			err = s.options(i, c)
			if err != nil {
				return
			}
		}
		if i, ok := c.Interface.(IPatch); ok {
			err = s.patch(i, c)
			if err != nil {
				return
			}
		}
		if i, ok := c.Interface.(IHead); ok {
			err = s.head(i, c)
			if err != nil {
				return
			}
		}
		if i, ok := c.Interface.(IConnect); ok {
			err = s.connect(i, c)
			if err != nil {
				return
			}
		}
		if i, ok := c.Interface.(ITrace); ok {
			err = s.trace(i, c)
			if err != nil {
				return
			}
		}
	}

	//Флаг старта
	s.IsStarted = true
	// Получение адреса
	addr := fmt.Sprintf(":%d", s.Config.Port)
	// Установка порта в swagger
	s.Swagger.setPort(addr)
	// Если флаг для безопасности true, то запускаем механизм с TLS
	if s.Config.Secure != nil {
		// Добавляем схему в Swagger
		s.Swagger.SetSchemes("https")
		// Возвращаем данные по сертификату
		cert, key := s.Config.Secure.Get()
		// Запускаем слушатель с TLS настройкой
		return s.WebServer.StartTLS(addr, cert, key)
	}

	// Добавляем схему в Swagger
	s.Swagger.SetSchemes("http")

	// Запуск слушателя веб сервера
	return s.WebServer.Start(addr)
}

// Stop Остановка сервера
func (s *Server) Stop() error {
	s.IsStarted = false
	return s.WebServer.Stop()
}

// Устанавливаем глобальные настройки для маршрутов
func (s *Server) newRoute() *Route {

	route := &Route{
		Operation: Operation{
			Responses: map[string]Response{
				"default": {
					Description: "successful operation",
				},
			},
		},
	}
	if s.Config.Permission != nil {
		route.isPermission = s.Config.Permission.AllRoutes
	}
	if s.Config.Authorization.AllRoutes != security.NoAuth {
		route.SetSecurity(s.Config.Authorization.AllRoutes)
	}

	return route
}

// Обрабатываем метод GET
func (s *Server) get(i IGet, c *Controller) error {
	route := s.newRoute()
	i.Get(route)
	return s.add(MethodGet, c, route)
}

// Обрабатываем метод POST
func (s *Server) post(i IPost, c *Controller) error {
	route := s.newRoute()
	i.Post(route)
	return s.add(MethodPost, c, route)
}

// Обрабатываем метод PUT
func (s *Server) put(i IPut, c *Controller) error {
	route := s.newRoute()
	i.Put(route)
	return s.add(MethodPut, c, route)
}

// Обрабатываем метод DELETE
func (s *Server) delete(i IDelete, c *Controller) error {
	route := s.newRoute()
	i.Delete(route)
	return s.add(MethodDelete, c, route)
}

// Обрабатываем метод OPTIONS
func (s *Server) options(i IOptions, c *Controller) error {
	route := s.newRoute()
	i.Options(route)
	return s.add(MethodOptions, c, route)
}

// Обрабатываем метод PATCH
func (s *Server) patch(i IPatch, c *Controller) error {
	route := s.newRoute()
	i.Patch(route)
	return s.add(MethodPatch, c, route)
}

// Обрабатываем метод HEAD
func (s *Server) head(i IHead, c *Controller) error {
	route := s.newRoute()
	i.Head(route)
	return s.add(MethodHead, c, route)
}

// Обрабатываем метод CONNECT
func (s *Server) connect(i IConnect, c *Controller) error {
	route := s.newRoute()
	i.Connect(route)
	return s.add(MethodConnect, c, route)
}

// Обрабатываем метод TRACE
func (s *Server) trace(i ITrace, c *Controller) error {
	route := s.newRoute()
	i.Trace(route)
	return s.add(MethodTrace, c, route)
}

// Добавить маршрут в веб сервер
func (s *Server) add(method string, c *Controller, route *Route) error {

	// Если нет ни одного handler, то выходим
	if route.Handler == nil {
		return nil
	}

	params := route.Operation.getParams()

	if params == nil || route.isEmptyParam {
		params = append(params, "")
	}

	/*var view *View
	// Проверка на view
	if s.Config.Views != nil {
		files, _ := filepath.Glob(filepath.Join(s.Config.Views.Root, strings.ToLower(name)+s.Config.Views.Engine))
		for _, file := range files {
			view = &View{
				Filename: strings.Replace(filepath.Base(file), s.Config.Views.Engine, "", -1),
				Filepath: file,
				Layout:   s.Config.Views.Layout,
			}
		}
	}*/

	// Авторизация в swagger
	for _, sec := range route.Security {
		for key := range sec {
			s.Swagger.setSecurityDefinition(key, s.Config.Authorization.Get(key).Definition())
		}
	}

	// Добавляем в swagger параметр указанный в суффиксе
	for _, suffix := range c.Suffix {
		if suffix.isParam {
			continue
		}
		route.Operation.Parameters = append(route.Operation.Parameters, NewInPath(suffix.Value, true, suffix.Description))
	}

	// Добавляем ссылку на тэг в контроллере
	route.Operation.addTag(c.Tag.Name)

	// Получаем handler маршрута
	h := route.getHandler(s.Config, nil, *s.Swagger)

	// Перебираем параметры адресной строки
	for _, param := range params {

		// Объединяем путь и параметры
		fullPath := p.Join(c.Path, param)

		// Проверка на соответствие базового пути
		ok, l := s.Swagger.compareBasePath(c.Path)
		if (param != "" || (param == "" && !route.isEmptyParam)) && (ok && c.IsShow) {
			// Добавляем пути и методы в swagger
			s.Swagger.setPath(fullPath[l:], strings.ToLower(method), route.Operation)
		}

		// Проверка на пустые пути
		if param != "" {
			matches := regexp.MustCompile(`{(\w+)}`).FindStringSubmatch(fullPath)
			if len(matches) == 2 {
				fullPath = strings.ReplaceAll(fullPath, matches[0], s.WebServer.ConvertParam(matches[1]))
			}
		}

		// Добавляем метод, путь и обработчик
		s.WebServer.Add(method, fullPath, h)
	}

	return nil
}

// Register Регистрация контроллера
func (s *Server) Register(i interface{}) *Controller {
	controller := &Controller{
		Interface: i,
		IsShow:    true,
	}
	s.Controllers = append(s.Controllers, controller)
	return controller
}

// Функция вернет Имя и версию
func (s *Server) String() string {
	return fmt.Sprintf("%s %s", Name, Version)
}
