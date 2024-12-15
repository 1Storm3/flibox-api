package rest

import (
	"github.com/1Storm3/flibox-api/internal/delivery/middleware"
	"github.com/1Storm3/flibox-api/internal/modules/auth"
	"github.com/1Storm3/flibox-api/internal/modules/collection"
	"github.com/1Storm3/flibox-api/internal/modules/collection-film"
	"github.com/1Storm3/flibox-api/internal/modules/comment"
	"github.com/1Storm3/flibox-api/internal/modules/external"
	"github.com/1Storm3/flibox-api/internal/modules/film"
	"github.com/1Storm3/flibox-api/internal/modules/film-sequel"
	"github.com/1Storm3/flibox-api/internal/modules/film-similar"
	"github.com/1Storm3/flibox-api/internal/modules/history-films"
	"github.com/1Storm3/flibox-api/internal/modules/user"
	"github.com/1Storm3/flibox-api/internal/modules/user-film"
	"github.com/gofiber/fiber/v2"
)

const (
	Admin = "admin"
	User  = "user"
)

type Router struct {
	filmHandler           film.HandlerInterface
	filmSequelHandler     filmsequel.HandlerInterface
	filmSimilarHandler    filmsimilar.HandlerInterface
	userHandler           user.HandlerInterface
	userFilmHandler       userfilm.HandlerInterface
	authHandler           auth.HandlerInterface
	externalHandler       external.HandlerInterface
	commentHandler        comment.HandlerInterface
	collectionHandler     collection.HandlerInterface
	collectionFilmHandler collectionfilm.HandlerInterface
	historyFilmsHandler   historyfilms.HandlerInterface
}

func NewRouter(
	filmHandler film.HandlerInterface,
	filmSequelHandler filmsequel.HandlerInterface,
	userHandler user.HandlerInterface,
	filmSimilarHandler filmsimilar.HandlerInterface,
	userFilmHandler userfilm.HandlerInterface,
	authHandler auth.HandlerInterface,
	externalHandler external.HandlerInterface,
	commentHandler comment.HandlerInterface,
	collectionHandler collection.HandlerInterface,
	collectionFilmHandler collectionfilm.HandlerInterface,
	historyFilmsHandler historyfilms.HandlerInterface,
) *Router {
	return &Router{
		filmHandler:           filmHandler,
		filmSequelHandler:     filmSequelHandler,
		userHandler:           userHandler,
		filmSimilarHandler:    filmSimilarHandler,
		userFilmHandler:       userFilmHandler,
		authHandler:           authHandler,
		externalHandler:       externalHandler,
		commentHandler:        commentHandler,
		collectionHandler:     collectionHandler,
		collectionFilmHandler: collectionFilmHandler,
		historyFilmsHandler:   historyFilmsHandler,
	}
}

func (r *Router) LoadRoutes(app fiber.Router, authMiddleware fiber.Handler) {
	apiRoute := app.Group("api")

	authRoute := apiRoute.Group("auth")
	r.setAuthRoutes(authRoute, authMiddleware)

	userRoute := apiRoute.Group("user", authMiddleware, middleware.RoleMiddleware(Admin, User))
	r.setUserRoutes(userRoute)

	filmRoute := apiRoute.Group("film")
	r.setFilmRoutes(filmRoute)

	sequelRoute := apiRoute.Group("sequel")
	r.setSequelRoutes(sequelRoute)

	similarRoute := apiRoute.Group("similar")
	r.setSimilarRoutes(similarRoute)

	externalRoute := apiRoute.Group("upload", authMiddleware, middleware.RoleMiddleware(Admin, User))
	r.setExternalRoutes(externalRoute)

	commentRoute := apiRoute.Group("comment", authMiddleware, middleware.RoleMiddleware(Admin, User))
	r.setCommentRoutes(commentRoute)

	collectionRoute := apiRoute.Group("collection", authMiddleware, middleware.RoleMiddleware(Admin, User))
	r.setCollectionRoutes(collectionRoute)

	collectionFilmRoute := apiRoute.Group("collection", authMiddleware, middleware.RoleMiddleware(Admin, User))
	r.setCollectionFilmRoutes(collectionFilmRoute)

	historyFilmsRoute := apiRoute.Group("film/history", authMiddleware, middleware.RoleMiddleware(Admin, User))
	r.setHistoryFilmsRoutes(historyFilmsRoute)
}
