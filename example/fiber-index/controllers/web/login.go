package web

import (
	"errors"
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/fiber/src/storage"
	"time"
)

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

/*func (Login) Get(route *ewa.Route) {
	route.Handler = func(c *ewa.Context) error {
		return c.Render("login", nil)
	}
}*/

func (a Auth) Post(route *ewa.Route) {
	route.Session(ewa.On)
	route.Handler = func(c *ewa.Context) error {

		err := c.BodyParser(&a)
		if err != nil {
			return c.JSON(400, ewa.Map{
				"message": err.Error(),
			})
		}

		if a.Username == "user" && a.Password == "Qq123456" {
			if c.Session != nil {
				storage.SetStorage(c.Session.Value, a.Username)
				return c.JSON(200, ewa.Map{
					c.Session.Key: c.Session.Value,
					"created":     c.Session.Created,
					"last_time":   c.Session.LastTime.Format(time.RFC3339),
				})
			}
		}

		return errors.New("Не верное имя пользователя или пароль!")
	}
}
