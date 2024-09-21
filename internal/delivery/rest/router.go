package rest

import "github.com/gofiber/fiber/v2"

type Router struct {
	filmHandler   *FilmHandler
	sequelHandler *SequelHandler
}

func NewRouter(filmHandler *FilmHandler, sequelHandler *SequelHandler) *Router {
	return &Router{filmHandler: filmHandler, sequelHandler: sequelHandler}
}

func (r *Router) LoadRoutes(app fiber.Router) {
	filmRoute := app.Group("api/films")
	filmRoute.Get(":id", r.filmHandler.GetOneByID)

	actorRoute := app.Group("api/sequels")
	actorRoute.Get(":id", r.sequelHandler.GetAll)

	// actors
	// serials
	// etc
}
