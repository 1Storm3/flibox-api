package rest

import (
	"github.com/gofiber/fiber/v2"

	"kbox-api/internal/delivery/middleware"
	"kbox-api/internal/modules/auth"
	"kbox-api/internal/modules/collection"
	"kbox-api/internal/modules/collection-film"
	"kbox-api/internal/modules/comment"
	"kbox-api/internal/modules/external"
	"kbox-api/internal/modules/film"
	"kbox-api/internal/modules/film-sequel"
	"kbox-api/internal/modules/film-similar"
	"kbox-api/internal/modules/history-films"
	"kbox-api/internal/modules/user"
	"kbox-api/internal/modules/user-film"
	"kbox-api/pkg/constant"
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

	authRoute.Post("login",
		middleware.ValidateMiddleware[auth.LoginDTO](),
		r.authHandler.Login)

	authRoute.Post("register",
		middleware.ValidateMiddleware[user.CreateUserDTO](),
		r.authHandler.Register)

	authRoute.Put("me",
		authMiddleware, r.authHandler.Me)

	authRoute.Post("verify/:token",
		r.authHandler.Verify)

	userFilmFavouriteRoute := apiRoute.Group("user/my",
		authMiddleware,
		middleware.RoleMiddleware(constant.Admin, constant.User),
	)
	userFilmFavouriteRoute.Get("/", r.userFilmHandler.GetAll)
	userFilmFavouriteRoute.Post("/:filmId", r.userFilmHandler.Add)
	userFilmFavouriteRoute.Delete("/:filmId", r.userFilmHandler.Delete)

	userRoute := apiRoute.Group("user", authMiddleware)
	userRoute.Get(":nickName", r.userHandler.GetOneByNickName)
	userRoute.Patch(":id", r.userHandler.Update)

	filmRoute := apiRoute.Group("film")
	filmRoute.Get(":id", r.filmHandler.GetOneByID)
	filmRoute.Get("", r.filmHandler.Search)

	sequelRoute := apiRoute.Group("sequel")
	sequelRoute.Get(":id", r.filmSequelHandler.GetAll)

	similarRoute := apiRoute.Group("similar")
	similarRoute.Get(":id", r.filmSimilarHandler.GetAll)

	externalRoute := apiRoute.Group("upload", authMiddleware)
	externalRoute.Put("", r.externalHandler.UploadFile)

	commentRoute := apiRoute.Group("comment", authMiddleware)
	commentRoute.Get("by/:filmId", r.commentHandler.GetAllByFilmID)

	commentRoute.Post("",
		middleware.ValidateMiddleware[comment.CreateCommentDTO](),
		r.commentHandler.Create)

	commentRoute.Delete(":id", r.commentHandler.Delete)

	commentRoute.Patch(":id",
		middleware.ValidateMiddleware[comment.UpdateCommentDTO](),
		r.commentHandler.Update)

	collectionRoute := apiRoute.Group("collection",
		authMiddleware,
		middleware.RoleMiddleware(constant.Admin, constant.User),
	)
	collectionRoute.Get("", r.collectionHandler.GetAll)
	collectionRoute.Get("my", r.collectionHandler.GetAllMy)
	collectionRoute.Get(":id", r.collectionHandler.GetOne)

	collectionRoute.Post("",
		middleware.ValidateMiddleware[collection.CreateCollectionDTO](),
		r.collectionHandler.Create)

	collectionRoute.Delete(":id", r.collectionHandler.Delete)

	collectionRoute.Patch(":id",
		middleware.ValidateMiddleware[collection.UpdateCollectionDTO](),
		r.collectionHandler.Update)

	collectionFilmRoute := apiRoute.Group("collection",
		authMiddleware,
		middleware.RoleMiddleware(constant.Admin, constant.User),
	)

	collectionFilmRoute.Post(":id/film",
		middleware.ValidateMiddleware[collectionfilm.CreateCollectionFilmDTO](),
		r.collectionFilmHandler.Add)

	collectionRoute.Delete(":id/film",
		middleware.ValidateMiddleware[collectionfilm.DeleteCollectionFilmDTO](),
		r.collectionFilmHandler.Delete)

	collectionRoute.Get(":id/films", r.collectionFilmHandler.GetFilmsByCollectionId)

	historyFilmsRoute := apiRoute.Group("film/history", authMiddleware)
	historyFilmsRoute.Post(":Id", r.historyFilmsHandler.Add)
}
