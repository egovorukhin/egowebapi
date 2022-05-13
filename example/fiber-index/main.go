package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	//BasicAuth
	/*basicAuthHandler := func(user string, pass string) bool {
		if user == "user" && pass == "Qq123456" {
			return true
		}
		return false
	}
	//Session
	checkSession := func(key string) (string, error) {
		if value, ok := storage.GetStorage(key); ok {
			return value, nil
		}
		return "", errors.New("Элемент не найден")
	}
	//Обработчик ошибок
	errorHandler := func(c *ewa.Context, code int, err interface{}) error {
		return c.Render("error", map[string]interface{}{"Code": code, "Text": err})
	}
	//Permission
	checkPermission := func(username, path string) bool {
		if username == "user" {
			switch path {
			case "/":
				return true
			}
		}
		return false
	}
	contextHandler := func(handler ewa.Handler) interface{} {
		return func(ctx *fiber.Ctx) error {
			return handler(ewa.NewContext(&fewa.Context{Ctx: ctx}))
		}
	}*/

	//root := "./views"

	// Fiber
	app := fiber.New(fiber.Config{
		//Views: html.New("/", ".html"),
	})
	app.Static("/views", "dist")
	app.Add(fiber.MethodGet, "/", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("./dist/index.html")
	})
	/*	app.Use(favicon.New(favicon.Config{
		File: "./views/favicon.ico",
	}))*/
	//app.Use(cors.New())
	app.Listen(":3005")

	/*server := &fewa.Server{App: app}
	// Конфиг
	cfg := ewa.Config{
		Port: 3005,
		Authorization: security.Authorization{
			Basic: &security.Basic{
				Handler: basicAuthHandler,
			},
		},
		Session: &session.Config{
			RedirectPath:   "/login",
			Expires:        1 * time.Minute,
			SessionHandler: checkSession,
		},
		Permission: &ewa.Permission{
			Handler: checkPermission,
		},
		ErrorHandler:   errorHandler,
		ContextHandler: contextHandler,
	}

	//Инициализируем сервер
	ws := ewa.New(server, cfg)
	ws.Register(new(web.Index))
	ws.Register(new(web.Auth)).SetPath("/login")
	ws.Register(new(web.Logout)).SetPath("/logout")

	// Канал для получения ошибки, если таковая будет
	errChan := make(chan error, 2)
	go func() {
		errChan <- ws.Start()
	}()

	// Ждем сигнал от ОС
	go getSignal(errChan)

	fmt.Println("Старт приложения")
	fmt.Printf("Остановка приложения. %s", <-errChan)*/

}

func getSignal(errChan chan error) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT)
	errChan <- fmt.Errorf("%s", <-c)
}
