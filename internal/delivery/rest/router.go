package rest

import (
	"github.com/gofiber/fiber/v2"

	"kbox-api/internal/delivery/middleware"
	dtoAuth "kbox-api/internal/modules/auth/dto"
	authHandler "kbox-api/internal/modules/auth/handler"
	externalHandler "kbox-api/internal/modules/external/handler"
	filmSequelHandler "kbox-api/internal/modules/film-sequel/handler"
	filmSimilarHandler "kbox-api/internal/modules/film-similar/handler"
	filmHandler "kbox-api/internal/modules/film/handler"
	userFilmHandler "kbox-api/internal/modules/user-film/handler"
	dtoUser "kbox-api/internal/modules/user/dto"
	userHandler "kbox-api/internal/modules/user/handler"
)

type Router struct {
	filmHandler        *filmHandler.FilmHandler
	filmSequelHandler  *filmSequelHandler.FilmSequelHandler
	filmSimilarHandler *filmSimilarHandler.FilmSimilarHandler
	userHandler        *userHandler.UserHandler
	userFilmHandler    *userFilmHandler.UserFilmHandler
	authHandler        *authHandler.AuthHandler
	externalHandler    *externalHandler.ExternalHandler
}

func NewRouter(
	filmHandler *filmHandler.FilmHandler,
	filmSequelHandler *filmSequelHandler.FilmSequelHandler,
	userHandler *userHandler.UserHandler,
	filmSimilarHandler *filmSimilarHandler.FilmSimilarHandler,
	userFilmHandler *userFilmHandler.UserFilmHandler,
	authHandler *authHandler.AuthHandler,
	externalHandler *externalHandler.ExternalHandler,
) *Router {
	return &Router{
		filmHandler:        filmHandler,
		filmSequelHandler:  filmSequelHandler,
		userHandler:        userHandler,
		filmSimilarHandler: filmSimilarHandler,
		userFilmHandler:    userFilmHandler,
		authHandler:        authHandler,
		externalHandler:    externalHandler,
	}
}

func (r *Router) LoadRoutes(app fiber.Router) {
	authRoute := app.Group("api/auth")
	authRoute.Post("login", middleware.ValidateMiddleware[dtoAuth.LoginDTO](), r.authHandler.Login)
	authRoute.Post("register", middleware.ValidateMiddleware[dtoUser.CreateUserDTO](), r.authHandler.Register)
	authRoute.Put("me", middleware.AuthMiddleware, r.authHandler.Me)

	userFilmRoute := app.Group("api/user/favourites")
	userFilmRoute.Get("", middleware.AuthMiddleware, r.userFilmHandler.GetAll)
	userFilmRoute.Post(":film_id", middleware.AuthMiddleware, r.userFilmHandler.Add)
	userFilmRoute.Delete(":film_id", middleware.AuthMiddleware, r.userFilmHandler.Delete)

	userRoute := app.Group("api/user")
	userRoute.Get(":nickName", middleware.AuthMiddleware, r.userHandler.GetOneByNickName)
	userRoute.Patch(":id", middleware.AuthMiddleware, r.userHandler.Update)

	filmRoute := app.Group("api/film")
	filmRoute.Get(":id", middleware.AuthMiddleware, r.filmHandler.GetOneByID)
	filmRoute.Get("", middleware.AuthMiddleware, r.filmHandler.Search)

	sequelRoute := app.Group("api/sequel")
	sequelRoute.Get(":id", middleware.AuthMiddleware, r.filmSequelHandler.GetAll)

	similarRoute := app.Group("api/similar")
	similarRoute.Get(":id", middleware.AuthMiddleware, r.filmSimilarHandler.GetAll)

	externalRoute := app.Group("api/upload")
	externalRoute.Put("", middleware.AuthMiddleware, r.externalHandler.UploadFile)
}
