package main

import (
	"fmt"
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/fiber-api/controllers"
	"github.com/egovorukhin/egowebapi/example/fiber-api/controllers/api/storage"
	"github.com/egovorukhin/egowebapi/example/fiber-api/models"
	f "github.com/egovorukhin/egowebapi/fiber"
	"github.com/egovorukhin/egowebapi/security"
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
		Authorization: security.Authorization{
			Basic: &security.Basic{
				Handler: basicAuthHandler,
			},
		},
	}

	info := ewa.Info{
		Description: "Description",
		Version:     "1.0.0",
		Title:       "FiberApi",
		Contact: &ewa.Contact{
			Email: "user@mail.ru",
		},
		License: &ewa.License{
			Name: "License",
		},
	}

	hostname := ewa.Suffix{
		Index:       3,
		Value:       "{hostname}",
		Description: "Set hostname device",
	}

	response := models.Response{}

	//Инициализируем сервер
	ws := ewa.New(server, cfg)
	ws.Register(new(storage.User)).SetSuffix(hostname).SetModel(models.User{}, response).SetDescription("Users")
	ws.Register(new(controllers.Home)).SetPath("/")
	// Swagger
	ws.Register(new(controllers.Api)).NotShow()

	// Описываем swagger
	ws.Swagger.SetInfo(fmt.Sprintf("10.28.0.73:%d", cfg.Port), &info, nil).SetBasePath("/api")
	//ws.Swagger.SetDefinitions(models.User{})

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
