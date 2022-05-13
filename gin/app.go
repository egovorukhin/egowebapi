package gin

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	App *gin.Engine
}

func (s *Server) Start(addr string) error {
	return s.App.Run(addr)
}

func (s *Server) StartTLS(addr, cert, key string) error {
	return s.App.RunTLS(addr, cert, key)
}

func (s *Server) Stop() error {
	return nil
	//return s.App.Shutdown(context.Background())
}

func (s *Server) Static(prefix, root string) {
	s.App.Static(prefix, root)
}

func (s *Server) Any(path string, handler interface{}) {
	if h, ok := handler.(gin.HandlerFunc); ok {
		s.App.Any(path, h)
	}
}

func (s *Server) Use(params ...interface{}) {
	for _, param := range params {
		if h, ok := param.(gin.HandlerFunc); ok {
			s.App.Use(h)
		}
	}
}

func (s *Server) Add(method, path string, handler interface{}) {
	s.App.Handle(method, path, handler.(gin.HandlerFunc))
}

func (s *Server) GetApp() interface{} {
	return s.App
}

func (s *Server) NotFoundPage(path, page string) {
	s.App.Any(path, func(c *gin.Context) {
		c.HTML(200, page, nil)
	})
}

func (s *Server) ConvertParam(param string) string {
	return "{" + param + "}"
}
