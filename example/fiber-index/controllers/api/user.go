package api

import (
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/gofiber/fiber/v2"
)

var users = Users{}

type User struct {
	Id        string
	Lastname  string
	Firstname string
}

type Users []User

func (u *User) Get(route *ewa.Route) {
	route.SetParams("", "/:id").
		SetDescription("Возвращаем всех пользователей либо по id").
		Auth(ewa.BasicAuth).Handler =
		func(c *fiber.Ctx) error {

			c.Set("System", c.Params("system"))
			c.Set("Version", c.Params("version"))

			id := c.Params("id")
			if id != "" {
				_, user := GetUser(id)
				return c.JSON(user)
			}

			return c.JSON(GetUsers())
		}
}

func (u *User) Post(route *ewa.Route) {
	route.Auth(ewa.BasicAuth).
		SetDescription("Добавляем пользователя").Handler =
		func(c *fiber.Ctx) error {
			user := &User{}
			user.Id = c.Query("id")
			user.Lastname = c.Query("lastname")
			user.Firstname = c.Query("firstname")
			SetUser(*user)
			return nil
		}
}

func (u *User) Put(route *ewa.Route) {
	route = nil
	/*route.SetParams("/:id").
	SetDescription("Изменяем пользователя по id").
	SetHandler(
		func(c *fiber.Ctx) error {
			u.Id = c.Params("id")
			u.Update()
			return nil
		})*/
}

func (u *User) Delete(route *ewa.Route) {
	route.SetParams("/:id").
		SetDescription("Удаляем пользователя по id").
		Handler =
		func(c *fiber.Ctx) error {
			u.Id = c.Params("id")
			u.Remove()
			return nil
		}
}

func (u *User) Options(swagger *s.Swagger) ewa.EmptyHandler {
	return func(ctx *fiber.Ctx) error {
		//ctx.Append("Allow", "GET, POST, PUT, DELETE, OPTIONS")
		swagger.Allow(ctx)
		return ctx.JSON(swagger)
	}
}

func GetUsers() Users {
	return users
}

func GetUser(id string) (int, *User) {
	for i, user := range users {
		if user.Id == id {
			return i, &user
		}
	}
	return -1, nil
}

func SetUser(u User) {
	users = append(users, u)
}

func (u *User) Update() {
	i, _ := GetUser(u.Id)
	users[i] = *u
}

func (u *User) Remove() {
	i, _ := GetUser(u.Id)
	users = append(users[:i], users[i+1:]...)
}
