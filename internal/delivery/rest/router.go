package rest

import "github.com/gofiber/fiber/v2"

type Router struct {
	filmHandler        *FilmHandler
	filmSequelHandler  *FilmSequelHandler
	filmSimilarHandler *FilmSimilarHandler
	userHandler        *UserHandler
}

func NewRouter(filmHandler *FilmHandler,
	filmSequelHandler *FilmSequelHandler,
	userHandler *UserHandler,
	filmSimilarHandler *FilmSimilarHandler,
) *Router {
	return &Router{
		filmHandler:        filmHandler,
		filmSequelHandler:  filmSequelHandler,
		userHandler:        userHandler,
		filmSimilarHandler: filmSimilarHandler,
	}
}

func (r *Router) LoadRoutes(app fiber.Router) {
	filmRoute := app.Group("api/films")
	filmRoute.Get(":id", r.filmHandler.GetOneByID)

	sequelRoute := app.Group("api/sequels")
	sequelRoute.Get(":id", r.filmSequelHandler.GetAll)

	similarRoute := app.Group("api/similars")
	similarRoute.Get(":id", r.filmSimilarHandler.GetAll)

	userRoute := app.Group("api/users")
	userRoute.Get(":user_token", r.userHandler.GetOne)
}
