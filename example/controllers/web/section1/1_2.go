package section1

import (
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/controllers/web/utils"
	"github.com/gofiber/fiber/v2"
)

type Section_1_2 struct {
	utils.NavBar
	Title string
}

func (h *Section_1_2) Get(route *ewa.Route) {
	route.SetDescription("Страница Home.html").Session().Permission()
	route.WebHandler = func(ctx *fiber.Ctx, identity *ewa.Identity) error {
		h.NavBar = utils.GetNavBar("section_1_2", identity.User)
		h.Title = "Привет раздел 1.2"
		return ctx.Render("section1/1_2", h, "layouts/base")
	}
}

func (h *Section_1_2) Post(route *ewa.Route) {
	route.Handler = nil
}
