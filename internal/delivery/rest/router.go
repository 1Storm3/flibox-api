package rest

import (
	"github.com/gofiber/fiber/v2"
	filmsequelhandler "kinopoisk-api/internal/modules/film-sequel/handler"
	filmsimilarhandler "kinopoisk-api/internal/modules/film-similar/handler"
	filmhandler "kinopoisk-api/internal/modules/film/handler"
	userfilmhandler "kinopoisk-api/internal/modules/user-film/handler"
	userhandler "kinopoisk-api/internal/modules/user/handler"
)

type Router struct {
	filmHandler        *filmhandler.FilmHandler
	filmSequelHandler  *filmsequelhandler.FilmSequelHandler
	filmSimilarHandler *filmsimilarhandler.FilmSimilarHandler
	userHandler        *userhandler.UserHandler
	userFilmHandler    *userfilmhandler.UserFilmHandler
}

func NewRouter(
	filmHandler *filmhandler.FilmHandler,
	filmSequelHandler *filmsequelhandler.FilmSequelHandler,
	userHandler *userhandler.UserHandler,
	filmSimilarHandler *filmsimilarhandler.FilmSimilarHandler,
	userFilmHandler *userfilmhandler.UserFilmHandler,
) *Router {
	return &Router{
		filmHandler:        filmHandler,
		filmSequelHandler:  filmSequelHandler,
		userHandler:        userHandler,
		filmSimilarHandler: filmSimilarHandler,
		userFilmHandler:    userFilmHandler,
	}
}

func (r *Router) LoadRoutes(app fiber.Router) {
	filmRoute := app.Group("api/films")
	filmRoute.Get(":id", r.filmHandler.GetOneByID)
	filmRoute.Get("", r.filmHandler.Search)

	sequelRoute := app.Group("api/sequels")
	sequelRoute.Get(":id", r.filmSequelHandler.GetAll)

	similarRoute := app.Group("api/similars")
	similarRoute.Get(":id", r.filmSimilarHandler.GetAll)

	userRoute := app.Group("api/users")
	userRoute.Get(":user_token", r.userHandler.GetOne)

	userFilmRoute := app.Group("api/users/:user_id/films")
	userFilmRoute.Get("", r.userFilmHandler.GetAll)
	userFilmRoute.Post(":film_id", r.userFilmHandler.Add)
	userFilmRoute.Delete(":film_id", r.userFilmHandler.Delete)
}
