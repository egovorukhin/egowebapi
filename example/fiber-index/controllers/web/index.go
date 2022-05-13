package web

import ewa "github.com/egovorukhin/egowebapi"

type Index struct{}

func (Index) Get(route *ewa.Route) {
	//route.Session()
	route.Handler = func(c *ewa.Context) error {
		return c.Render("index", nil)
	}
}
