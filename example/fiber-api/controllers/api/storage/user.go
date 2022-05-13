package storage

import (
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/fiber-api/models"
	"github.com/egovorukhin/egowebapi/security"
	"strconv"
	"time"
)

type User struct{}

func (User) Get(route *ewa.Route) {

	route.SetSecurity(security.BasicAuth).
		SetEmptyParam(ewa.NewEmptyPathParam("Get users").SetResponse(200, ewa.NewResponse(ewa.NewSchemaArray(route.ParameterModel()), "Return array users"))).
		SetParameters(ewa.NewPathParam("/{id}", "Id пользователя")).
		InitParametersByModel().
		SetSummary("Get user")
	route.SetResponse(422, ewa.NewResponse(nil, "Return parse parameter error")).
		SetResponse(200, ewa.NewResponse(ewa.NewSchema(route.ParameterModel()), "Return user struct"))

	route.Handler = func(c *ewa.Context) error {
		id, err := strconv.Atoi(c.Params("id", "0"))
		if err != nil {
			return c.SendString(422, err.Error())
		}
		if id > 0 {
			user := models.GetUser(id)
			return c.JSON(200, user)
		}
		users := models.GetUsers()
		return c.JSON(200, users)
	}
}

func (User) Post(route *ewa.Route) {
	route.SetSecurity(security.BasicAuth)
	route.SetParameters(ewa.NewBodyParam(true, ewa.NewSchema(route.ParameterModel()), "Must have request body"))
	route.SetSummary("Create user")
	route.Handler = func(c *ewa.Context) error {
		user := models.User{}
		err := c.BodyParser(&user)
		if err != nil {
			return c.SendString(400, err.Error())
		}
		user.Set()
		return c.JSON(200, models.Response{
			Id:       user.Id,
			Message:  "Created",
			Datetime: time.Now(),
		})
	}
	route.SetResponse(200, ewa.NewResponse(ewa.NewSchema(route.ResponseModel()), "OK"))
	route.SetResponse(400, ewa.NewResponse(nil, "Parse body error"))
}

func (User) Put(route *ewa.Route) {
	route.SetSecurity(security.BasicAuth)
	route.SetParameters(
		ewa.NewQueryParam("id", false, "id user"),
		ewa.NewBodyParam(true, ewa.NewSchema(route.ParameterModel()), "Must have request body"),
	)
	route.SetSummary("Update user")
	route.Handler = func(c *ewa.Context) error {

		id, err := strconv.Atoi(c.QueryParam("id"))
		if err != nil {
			return c.SendString(422, err.Error())
		}
		user := models.User{}
		err = c.BodyParser(&user)
		if err != nil {
			return c.SendString(400, err.Error())
		}
		err = user.Update(id)
		if err != nil {
			return c.SendString(400, err.Error())
		}
		return c.SendString(200, "OK")
	}
	route.SetResponse(400, ewa.NewResponse(nil, "Parse body error"))
	route.SetResponse(422, ewa.NewResponse(nil, "Return query error"))
	route.SetResponse(200, ewa.NewResponse(ewa.NewSchema(route.ResponseModel()), "OK"))
}

func (User) Delete(route *ewa.Route) {
	route.SetSecurity(security.BasicAuth)
	route.SetParameters(ewa.NewPathParam("/{id}", "ID user"))
	route.SetSummary("Delete user")
	route.Handler = func(c *ewa.Context) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.SendString(422, err.Error())
		}
		user := models.User{
			Id: id,
		}
		user.Delete()
		return c.SendString(200, "OK")
	}
	route.SetResponse(422, ewa.NewResponse(nil, "Return parse parameter error"))
	route.SetResponse(200, ewa.NewResponse(ewa.NewSchema(route.ResponseModel()), "OK"))
}
