package egowebapi

import (
	"errors"
	"fmt"
	"github.com/egovorukhin/egowebapi/consts"
	"github.com/egovorukhin/egowebapi/security"
	"github.com/mustan989/jsonschema"
	p "path"
	"regexp"
	"strings"
)

const (
	Name    = "EgoWebApi"
	Version = "v0.2.26"
)

type Server struct {
	Config      Config
	IsStarted   bool
	WebServer   IServer
	Controllers []*Controller
	Swagger     Swagger
}

type IServer interface {
	Start(addr string) error
	StartTLS(addr, cert, key string) error
	Stop() error
	Static(prefix, root string)
	Any(path string, handler interface{})
	Use(params ...interface{})
	Add(method, path string, handler interface{})
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

	if config.Session != nil {
		config.Session.Default()
	}

	s := &Server{
		Config:    config,
		WebServer: server,
		Swagger: Swagger{
			Swagger:             "2.0",
			Host:                fmt.Sprintf("localhost:%d", config.Port),
			BasePath:            "/",
			SecurityDefinitions: SecurityDefinitions{},
			Paths:               Paths{},
			Definitions:         jsonschema.Definitions{},
			models:              Models{},
		},
	}

	// Глобальная переменная для указания ссылки на объект
	jsonschema.SetReferencePrefix(RefDefinitions)

	return s
}

// GetWebServer вернуть интерфейс веб сервера
func (s *Server) GetWebServer() interface{} {
	return s.WebServer.GetApp()
}

// Start запуск сервера
func (s *Server) Start() (err error) {

	if s.Config.ContextHandler == nil {
		return errors.New("Specify the handler - ContextHandler")
	}

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
			Produces: []string{
				consts.MIMEApplicationJSON,
				consts.MIMEApplicationXML,
			},
			Responses: map[string]Response{
				"default": {
					Description: "successful operation",
				},
			},
		},
		models: s.Swagger.models,
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
	return s.add(consts.MethodGet, c, route)
}

// Обрабатываем метод POST
func (s *Server) post(i IPost, c *Controller) error {
	route := s.newRoute()
	i.Post(route)
	return s.add(consts.MethodPost, c, route)
}

// Обрабатываем метод PUT
func (s *Server) put(i IPut, c *Controller) error {
	route := s.newRoute()
	i.Put(route)
	return s.add(consts.MethodPut, c, route)
}

// Обрабатываем метод DELETE
func (s *Server) delete(i IDelete, c *Controller) error {
	route := s.newRoute()
	i.Delete(route)
	return s.add(consts.MethodDelete, c, route)
}

// Обрабатываем метод OPTIONS
func (s *Server) options(i IOptions, c *Controller) error {
	route := s.newRoute()
	i.Options(route)
	return s.add(consts.MethodOptions, c, route)
}

// Обрабатываем метод PATCH
func (s *Server) patch(i IPatch, c *Controller) error {
	route := s.newRoute()
	i.Patch(route)
	return s.add(consts.MethodPatch, c, route)
}

// Обрабатываем метод HEAD
func (s *Server) head(i IHead, c *Controller) error {
	route := s.newRoute()
	i.Head(route)
	return s.add(consts.MethodHead, c, route)
}

// Обрабатываем метод CONNECT
func (s *Server) connect(i IConnect, c *Controller) error {
	route := s.newRoute()
	i.Connect(route)
	return s.add(consts.MethodConnect, c, route)
}

// Обрабатываем метод TRACE
func (s *Server) trace(i ITrace, c *Controller) error {
	route := s.newRoute()
	i.Trace(route)
	return s.add(consts.MethodTrace, c, route)
}

// Добавить маршрут в веб сервер
func (s *Server) add(method string, c *Controller, route *Route) error {

	// Если нет ни одного handler, то выходим
	if route.Handler == nil {
		return nil
	}

	pathParams := route.Operation.getPathParams()
	params := []string{pathParams}
	if pathParams == "" || route.emptyPathParam != nil {
		params = append(params, "")
	}

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
		route.Operation.Parameters = append(route.Operation.Parameters, NewPathParam(suffix.Value, suffix.Description))
	}

	// Добавляем ссылку на тэг в контроллере
	route.Operation.addTag(c.Tag.Name)

	// Получаем handler маршрута
	h := s.Config.ContextHandler(route.getHandler(s.Config, s.Swagger))

	// Перебираем параметры адресной строки
	for _, param := range params {

		// Объединяем путь и параметры
		fullPath := p.Join(c.Path, param)

		// Проверка на соответствие базового пути
		ok, l := s.Swagger.compareBasePath(c.Path)
		if ok && c.IsShow {

			operation := route.Operation
			// Если пустой путь, то применяем некоторые настройки из основного
			if param == "" && route.emptyPathParam != nil {
				operation.Responses = make(map[string]Response)
				for key, value := range route.Responses {
					operation.Responses[key] = value
				}
				for key, value := range route.emptyPathParam.Responses {
					operation.Responses[key] = value
				}
				operation.Description = route.emptyPathParam.Description
				operation.Summary = route.emptyPathParam.Summary
				operation.Parameters = route.Operation.getParams(params...)
			}

			lowerMethod := strings.ToLower(method)
			// Установка ID операции
			operation.ID = lowerMethod + strings.ReplaceAll(fullPath[l:], "/", "-")

			// Добавляем пути и методы в swagger
			s.Swagger.setPath(fullPath[l:], lowerMethod, operation)
		}

		// Корректировка параметров пути
		fullPath = s.convertParams(fullPath)

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

// convertParams Корректировка параметров адресной строки
func (s *Server) convertParams(path string) string {
	matches := regexp.MustCompile(`{(\w+)}`).FindAllStringSubmatch(path, -1)
	for _, match := range matches {
		if len(match) == 2 {
			path = strings.ReplaceAll(path, match[0], s.WebServer.ConvertParam(match[1]))
		}
	}
	return path
}
