package rest

import (
	"github.com/gofiber/fiber/v2"
	"kinopoisk-api/internal/delivery/middleware"
	"kinopoisk-api/internal/modules/auth/handler"
	"kinopoisk-api/internal/modules/auth/service"
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
	authHandler        *handler.AuthHandler
}

func NewRouter(
	filmHandler *filmhandler.FilmHandler,
	filmSequelHandler *filmsequelhandler.FilmSequelHandler,
	userHandler *userhandler.UserHandler,
	filmSimilarHandler *filmsimilarhandler.FilmSimilarHandler,
	userFilmHandler *userfilmhandler.UserFilmHandler,
	authHandler *handler.AuthHandler,
) *Router {
	return &Router{
		filmHandler:        filmHandler,
		filmSequelHandler:  filmSequelHandler,
		userHandler:        userHandler,
		filmSimilarHandler: filmSimilarHandler,
		userFilmHandler:    userFilmHandler,
		authHandler:        authHandler,
	}
}

func (r *Router) LoadRoutes(app fiber.Router) {
	authRoute := app.Group("api/auth")
	authRoute.Post("login", middleware.ValidateMiddleware[handler.RequestLogin](), r.authHandler.Login)
	authRoute.Post("register", middleware.ValidateMiddleware[service.RequestUser](), r.authHandler.Register)
	authRoute.Put("me", middleware.AuthMiddleware, r.authHandler.Me)

	userRoute := app.Group("api/user")
	userRoute.Get(":nickName", middleware.AuthMiddleware, r.userHandler.GetOneByNickName)

	filmRoute := app.Group("api/film")
	filmRoute.Get(":id", middleware.AuthMiddleware, r.filmHandler.GetOneByID)
	filmRoute.Get("", middleware.AuthMiddleware, r.filmHandler.Search)

	sequelRoute := app.Group("api/sequel")
	sequelRoute.Get(":id", middleware.AuthMiddleware, r.filmSequelHandler.GetAll)

	similarRoute := app.Group("api/similar")
	similarRoute.Get(":id", middleware.AuthMiddleware, r.filmSimilarHandler.GetAll)

	userFilmRoute := app.Group("api/user/:user_id/films")
	userFilmRoute.Get("", middleware.AuthMiddleware, r.userFilmHandler.GetAll)
	userFilmRoute.Post(":film_id", middleware.AuthMiddleware, r.userFilmHandler.Add)
	userFilmRoute.Delete(":film_id", middleware.AuthMiddleware, r.userFilmHandler.Delete)
}
