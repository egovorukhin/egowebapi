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
		SetParameters(ewa.NewPathParam("/{id}", "Id пользователя")).
		InitParametersByModel(models.ModelUser).
		SetSummary("Get user").
		SetResponse(422, "", nil, "Return parse parameter error").
		SetResponse(200, models.ModelUser, nil, "Return user struct").
		SetEmptyParam("Get users").SetResponseArray(200, models.ModelUser, nil, "Return array users")

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

	route.SetSecurity(security.BasicAuth).
		SetParameters(ewa.NewBodyParam(true, models.ModelUser, false, "Must have request body")).
		SetSummary("Create user").
		SetResponse(200, models.ModelResponse, nil, "OK").
		SetResponse(400, "", nil, "Parse body error")

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
}

func (User) Put(route *ewa.Route) {

	route.SetSecurity(security.BasicAuth).
		InitParametersByModel(models.ModelUser).
		SetParameters(ewa.NewBodyParam(true, models.ModelUser, false, "Must have request body")).
		SetSummary("Update user").
		SetResponse(400, "", nil, "Parse body error").
		SetResponse(422, "", nil, "Return query error").
		SetResponse(200, models.ModelResponse, nil, "OK")

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
		return c.JSON(200, models.Response{
			Id:       user.Id,
			Message:  "Updated",
			Datetime: time.Now(),
		})
	}
}

func (User) Delete(route *ewa.Route) {

	route.SetSecurity(security.BasicAuth).
		SetParameters(ewa.NewPathParam("/{id}", "ID user")).
		SetSummary("Delete user").
		SetResponse(422, "", nil, "Return parse parameter error").
		SetResponse(200, models.ModelResponse, nil, "OK")

	route.Handler = func(c *ewa.Context) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.SendString(422, err.Error())
		}
		user := models.User{
			Id: id,
		}
		user.Delete()
		return c.JSON(200, models.Response{
			Id:       user.Id,
			Message:  "Deleted",
			Datetime: time.Now(),
		})
	}
}
