package rest

import (
	"github.com/1Storm3/flibox-api/internal/delivery/middleware"
	"github.com/1Storm3/flibox-api/internal/modules/auth"
	"github.com/1Storm3/flibox-api/internal/modules/collection"
	"github.com/1Storm3/flibox-api/internal/modules/collection-film"
	"github.com/1Storm3/flibox-api/internal/modules/comment"
	"github.com/1Storm3/flibox-api/internal/modules/user"
	"github.com/gofiber/fiber/v2"
)

func (r *Router) setCommentRoutes(commentRoute fiber.Router) {
	commentRoute.Get("by/:filmId", r.commentHandler.GetAllByFilmID)
	commentRoute.Post("", middleware.ValidateMiddleware[comment.CreateCommentDTO](), r.commentHandler.Create)
	commentRoute.Delete(":id", r.commentHandler.Delete)
	commentRoute.Patch(":id", middleware.ValidateMiddleware[comment.UpdateCommentDTO](), r.commentHandler.Update)
}

func (r *Router) setHistoryFilmsRoutes(historyFilmsRoute fiber.Router) {
	historyFilmsRoute.Post(":Id", r.historyFilmsHandler.Add)
}

func (r *Router) setExternalRoutes(externalRoute fiber.Router) {
	externalRoute.Put("", r.externalHandler.UploadFile)
}
func (r *Router) setSequelRoutes(sequelRoute fiber.Router) {
	sequelRoute.Get(":id", r.filmSequelHandler.GetAll)
}

func (r *Router) setSimilarRoutes(similarRoute fiber.Router) {
	similarRoute.Get(":id", r.filmSimilarHandler.GetAll)
}

func (r *Router) setAuthRoutes(authRoute fiber.Router, authMiddleware fiber.Handler) {
	authRoute.Post("login", middleware.ValidateMiddleware[auth.LoginDTO](), r.authHandler.Login)
	authRoute.Post("register", middleware.ValidateMiddleware[user.CreateUserDTO](), r.authHandler.Register)
	authRoute.Put("me", authMiddleware, r.authHandler.Me)
	authRoute.Post("verify/:token", r.authHandler.Verify)
}

func (r *Router) setUserRoutes(userRoute fiber.Router) {
	userRoute.Get(":nickName", r.userHandler.GetOneByNickName)
	userRoute.Patch(":id", r.userHandler.Update)
}

func (r *Router) setFilmRoutes(filmRoute fiber.Router) {
	filmRoute.Get(":id", r.filmHandler.GetOneByID)
	filmRoute.Get("", r.filmHandler.Search)
}

func (r *Router) setCollectionRoutes(collectionRoute fiber.Router) {
	collectionRoute.Get("", r.collectionHandler.GetAll)
	collectionRoute.Get("my", r.collectionHandler.GetAllMy)
	collectionRoute.Get(":id", r.collectionHandler.GetOne)

	collectionRoute.Post("", middleware.ValidateMiddleware[collection.CreateCollectionDTO](), r.collectionHandler.Create)
	collectionRoute.Delete(":id", r.collectionHandler.Delete)
	collectionRoute.Patch(":id", middleware.ValidateMiddleware[collection.UpdateCollectionDTO](), r.collectionHandler.Update)
}

func (r *Router) setCollectionFilmRoutes(collectionFilmRoute fiber.Router) {
	collectionFilmRoute.Post(":id/film", middleware.ValidateMiddleware[collectionfilm.CreateCollectionFilmDTO](), r.collectionFilmHandler.Add)
	collectionFilmRoute.Delete(":id/film", middleware.ValidateMiddleware[collectionfilm.DeleteCollectionFilmDTO](), r.collectionFilmHandler.Delete)
	collectionFilmRoute.Get(":id/films", r.collectionFilmHandler.GetFilmsByCollectionId)
}
