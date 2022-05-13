package echo

import (
	"context"
	"github.com/labstack/echo/v4"
)

type Server struct {
	App *echo.Echo
}

func (s *Server) Start(addr string) error {
	return s.App.Start(addr)
}

func (s *Server) StartTLS(addr, cert, key string) error {
	return s.App.StartTLS(addr, cert, key)
}

func (s *Server) Stop() error {
	return s.App.Shutdown(context.Background())
}

func (s *Server) Static(prefix, root string) {
	s.App.Static(prefix, root)
}

func (s *Server) Any(path string, handler interface{}) {
	if h, ok := handler.(echo.HandlerFunc); ok {
		s.App.Any(path, h)
	}
}

func (s *Server) Use(params ...interface{}) {
	for _, param := range params {
		if h, ok := param.(echo.MiddlewareFunc); ok {
			s.App.Use(h)
		}
	}
}

func (s *Server) Add(method, path string, handler interface{}) {
	s.App.Add(method, path, handler.(echo.HandlerFunc))
}

func (s *Server) GetApp() interface{} {
	return s.App
}

func (s *Server) NotFoundPage(path, page string) {
	s.App.Any(path, func(c echo.Context) error {
		return c.Render(200, page, nil)
	})
}

func (s *Server) ConvertParam(param string) string {
	return "{" + param + "}"
}
