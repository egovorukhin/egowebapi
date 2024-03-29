package controllers

import (
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/consts"
)

type Api struct{}

func (Api) Get(route *ewa.Route) {
	route.Handler = func(c *ewa.Context) error {

		b, err := c.Swagger.JSON()
		if err != nil {
			return c.SendString(422, err.Error())
		}

		return c.Send(200, consts.MIMEApplicationJSON, b)
		//return c.JSON(200, c.Swagger)
	}
}
