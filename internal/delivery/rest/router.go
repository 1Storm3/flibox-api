package rest

import "github.com/gofiber/fiber/v2"

type Router struct {
	filmHandler   *FilmHandler
	sequelHandler *SequelHandler
	userHandler   *UserHandler
}

func NewRouter(filmHandler *FilmHandler, sequelHandler *SequelHandler, userHandler *UserHandler) *Router {
	return &Router{filmHandler: filmHandler, sequelHandler: sequelHandler, userHandler: userHandler}
}

func (r *Router) LoadRoutes(app fiber.Router) {
	filmRoute := app.Group("api/films")
	filmRoute.Get(":id", r.filmHandler.GetOneByID)

	sequelRoute := app.Group("api/sequels")
	sequelRoute.Get(":id", r.sequelHandler.GetAll)

	userRoute := app.Group("api/users")
	userRoute.Get(":user_token", r.userHandler.GetOne)

	// actors
	// serials
	// etc
}
