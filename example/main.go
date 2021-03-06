package main

import (
	"errors"
	"fmt"
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/controllers/api"
	"github.com/egovorukhin/egowebapi/example/controllers/web"
	"github.com/egovorukhin/egowebapi/example/controllers/web/section1"
	__1 "github.com/egovorukhin/egowebapi/example/controllers/web/section1/1_1"
	"github.com/egovorukhin/egowebapi/example/src/storage"
	"github.com/gofiber/fiber/v2"
	"os"
	"strings"
)

func main() {

	/*	app := fiber.New()
		app.Get("/", func(c *fiber.Ctx) error {
			// Create cookie
			cookie := new(fiber.Cookie)
			cookie.Name = "john"
			cookie.Value = "doe"
			cookie.Expires = time.Now().Add(24 * time.Hour)

			// Set cookie
			c.Cookie(cookie)

			return nil
		})

		err := app.Listen(":3005")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Scanln()*/

	//BasicAuth
	authorizer := func(user string, pass string) bool {
		if user == "user" && pass == "Qq123456" {
			return true
		}
		return false
	}
	//Session
	checkSession := func(key string) (string, string, error) {
		if value, ok := storage.GetStorage(key); ok {
			return value, "", nil
		}
		return "", "", errors.New("Элемент не найден")
	}
	//Обработчик ошибок
	errorHandler := func(ctx *fiber.Ctx, code int, err string) error {
		return ctx.Render("error", fiber.Map{"Code": code, "Text": err})
	}
	//Permission
	checkPermission := func(key, route string) bool {
		user, _ := storage.GetStorage(key)
		if user == "user" && strings.Contains(route, "/section1/1_1") {
			return true
		}
		return false
	}
	//WEB
	cfg := ewa.Config{
		Port:    3005,
		Timeout: ewa.NewTimeout(30, 30, 30),
		Views: &ewa.Views{
			Directory: "views",
			Extension: ewa.Html,
			Engine: &ewa.Engine{
				Reload: false,
			},
		},
		Static:    "views",
		BasicAuth: ewa.NewBasicAuth(authorizer, nil),
		Session: &ewa.Session{
			RedirectPath:      "/login",
			SessionHandler:    checkSession,
			PermissionHandler: checkPermission,
			ErrorHandler:      errorHandler,
		},
	}
	//Инициализируем сервер
	system := ewa.Suffix{
		Index: 2,
		Value: ":system",
	}
	version := ewa.Suffix{
		Index: 3,
		Value: ":version",
	}
	ws, _ := ewa.New("Example", cfg)
	ws.RegisterWeb(new(web.Home), "/")
	ws.RegisterWeb(new(web.Login), "/login")
	ws.RegisterWeb(new(web.Logout), "/logout")
	ws.RegisterWeb(new(__1.Document), "/section1/1_1/document")
	ws.RegisterWeb(new(__1.List), "/section1/1_1/list")
	ws.RegisterWeb(new(section1.Section_1_2), "/section1/1_2")
	ws.RegisterRest(new(api.User), "", "person", system, version)
	//ws.SetBasicAuth(ba)
	//Cors = nil - DefaultConfig
	ws.SetCors(nil)
	//ws.SetStore(nil)
	ws.Start()

	for {
		var input string
		_, err := fmt.Fscan(os.Stdin, &input)
		if err != nil {
			os.Exit(1)
		}
		switch strings.ToLower(input) {
		case "exit":
			fmt.Println(ws.Stop())
			os.Exit(0)
		}
	}
}
