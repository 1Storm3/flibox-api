package rest

import "github.com/gofiber/fiber/v2"

type Router struct {
	filmHandler *FilmHandler
}

func NewRouter(filmHandler *FilmHandler) *Router {
	return &Router{filmHandler: filmHandler}
}

func (r *Router) LoadRoutes(app fiber.Router) {
	filmRoute := app.Group("/films")
	filmRoute.Get("/:id", r.filmHandler.GetOneByID)

	// actorRoute := app.Group("/actors")
	// actorRoute.Get("/:id", r.actorHandler.GetOneByID)

	// actors
	// serials
	// etc
}
