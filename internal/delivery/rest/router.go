package rest

import "github.com/gofiber/fiber/v2"

type Router struct {
	filmHandler       *FilmHandler
	filmSequelHandler *FilmSequelHandler
	userHandler       *UserHandler
}

func NewRouter(filmHandler *FilmHandler, filmSequelHandler *FilmSequelHandler, userHandler *UserHandler) *Router {
	return &Router{filmHandler: filmHandler, filmSequelHandler: filmSequelHandler, userHandler: userHandler}
}

func (r *Router) LoadRoutes(app fiber.Router) {
	filmRoute := app.Group("api/films")
	filmRoute.Get(":id", r.filmHandler.GetOneByID)

	sequelRoute := app.Group("api/sequels")
	sequelRoute.Get(":id", r.filmSequelHandler.GetAll)

	userRoute := app.Group("api/users")
	userRoute.Get(":user_token", r.userHandler.GetOne)

	// actors
	// serials
	// etc
}
