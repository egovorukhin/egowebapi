package main

import (
	"fmt"
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/fiber-api/controllers"
	"github.com/egovorukhin/egowebapi/example/fiber-api/controllers/api"
	f "github.com/egovorukhin/egowebapi/fiber"
	"github.com/egovorukhin/egowebapi/swagger"
	"github.com/egovorukhin/egowebapi/swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	//BasicAuth
	basicAuthHandler := func(user string, pass string) bool {
		if user == "user" && pass == "Qq123456" {
			return true
		}
		return false
	}

	// Fiber
	app := fiber.New()
	// Cors
	app.Use(cors.New())
	server := &f.Server{App: app}
	// Конфиг
	cfg := ewa.Config{
		Port: 8070,
		Secure: &ewa.Secure{
			Path: "./cert",
			Key:  "key.pem",
			Cert: "cert.pem",
		},
		Authorization: ewa.Authorization{
			Basic: basicAuthHandler,
		},
		Swagger: &swagger.Config{
			Host: "10.28.0.73:8070",
			Info: v2.Info{
				Description:    "Описание приложения",
				Version:        "1.0.0",
				Title:          "FiberApi",
				TermsOfService: "",
				Contact: v2.Contact{
					Email: "user@mail.ru",
				},
				License: v2.License{
					Name: "Пользуйся на здоровье",
				},
			},
		},
	}
	//Инициализируем сервер
	ws := ewa.New(server, cfg)
	ws.Register(new(api.User), "")
	// Swagger
	ws.Register(new(controllers.Api), "")

	// Канал для получения ошибки, если таковая будет
	errChan := make(chan error, 2)
	go func() {
		errChan <- ws.Start()
	}()

	// Ждем сигнал от ОС
	go getSignal(errChan)

	fmt.Println("Старт приложения")
	fmt.Printf("Остановка приложения. %s", <-errChan)
}

func getSignal(errChan chan error) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT)
	errChan <- fmt.Errorf("%s", <-c)
}
