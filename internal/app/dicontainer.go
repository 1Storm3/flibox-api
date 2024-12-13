package app

import (
	"kbox-api/database/postgres"
	"kbox-api/internal/config"
	"kbox-api/internal/delivery/grpc"
	"kbox-api/internal/modules/auth"
	"kbox-api/internal/modules/collection"
	"kbox-api/internal/modules/collection-film"
	"kbox-api/internal/modules/comment"
	"kbox-api/internal/modules/external"
	"kbox-api/internal/modules/film"
	"kbox-api/internal/modules/film-sequel"
	"kbox-api/internal/modules/film-similar"
	"kbox-api/internal/modules/history-films"
	"kbox-api/internal/modules/recommendation"
	"kbox-api/internal/modules/recommendation/adapter"
	"kbox-api/internal/modules/user"
	"kbox-api/internal/modules/user-film"
	"kbox-api/internal/shared/logger"
)

type diContainer struct {
	config  *config.Config
	storage *postgres.Storage

	grpcClient grpc.ClientConnInterface

	filmModule           film.ModuleInterface
	historyFilmsModule   historyfilms.ModuleInterface
	commentModule        comment.ModuleInterface
	externalModule       external.ModuleInterface
	userModule           user.ModuleInterface
	filmSequelModule     filmsequel.ModuleInterface
	filmSimilarModule    filmsimilar.ModuleInterface
	userFilmModule       userfilm.ModuleInterface
	authModule           auth.ModuleInterface
	collectionModule     collection.ModuleInterface
	collectionFilmModule collectionfilm.ModuleInterface
	recommendationModule adapter.ModuleInterface
}

func newDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) Config() *config.Config {
	if d.config == nil {
		d.config = config.MustLoad()
	}
	return d.config
}
func (d *diContainer) RecommendationModule() (adapter.ModuleInterface, error) {
	if d.recommendationModule == nil {
		grpcClient, err := d.GrpcClient()
		if err != nil {
			return nil, err
		}
		filmModule, err := d.FilmModule()
		if err != nil {
			return nil, err
		}
		historyFilmsModule, err := d.HistoryFilmsModule()
		if err != nil {
			return nil, err
		}
		userFilmModule, err := d.UserFilmModule()
		if err != nil {
			return nil, err
		}
		d.recommendationModule = recommendation.NewRecommendationModule(
			filmModule, historyFilmsModule, userFilmModule, grpcClient,
		)
	}
	return d.recommendationModule, nil
}

func (d *diContainer) UserFilmModule() (userfilm.ModuleInterface, error) {
	if d.userFilmModule == nil {
		storage, err := d.Storage()
		if err != nil {
			return nil, err
		}
		filmModule, err := d.FilmModule()
		if err != nil {
			return nil, err
		}
		recommendModuleFactory := func() (adapter.ModuleInterface, error) {
			return d.RecommendationModule()
		}
		d.userFilmModule = userfilm.NewUserFilmModule(storage, filmModule, recommendModuleFactory)
	}
	return d.userFilmModule, nil
}

func (d *diContainer) HistoryFilmsModule() (historyfilms.ModuleInterface, error) {
	if d.historyFilmsModule == nil {
		storage, err := d.Storage()
		if err != nil {
			return nil, err
		}
		recommendModuleFactory := func() (adapter.ModuleInterface, error) {
			return d.RecommendationModule()
		}
		d.historyFilmsModule = historyfilms.NewHistoryFilmsModule(storage, recommendModuleFactory)
	}
	return d.historyFilmsModule, nil
}

func (d *diContainer) CollectionFilmModule() (collectionfilm.ModuleInterface, error) {
	if d.collectionFilmModule == nil {
		storage, err := d.Storage()
		if err != nil {
			return nil, err
		}
		d.collectionFilmModule = collectionfilm.NewCollectionFilmModule(storage)
	}
	return d.collectionFilmModule, nil
}

func (d *diContainer) CollectionModule() (collection.ModuleInterface, error) {
	if d.collectionModule == nil {
		storage, err := d.Storage()
		if err != nil {
			return nil, err
		}
		d.collectionModule = collection.NewCollectionModule(storage)
	}
	return d.collectionModule, nil
}

func (d *diContainer) CommentModule() (comment.ModuleInterface, error) {
	if d.commentModule == nil {
		storage, err := d.Storage()
		if err != nil {
			return nil, err
		}
		d.commentModule = comment.NewCommentModule(storage)
	}
	return d.commentModule, nil
}

func (d *diContainer) AuthModule() (auth.ModuleInterface, error) {
	if d.authModule == nil {
		userModule, err := d.UserModule()
		if err != nil {
			return nil, err
		}
		externalModule, err := d.ExternalModule()
		if err != nil {
			return nil, err
		}
		d.authModule = auth.NewAuthModule(userModule, externalModule, d.Config())
	}
	return d.authModule, nil
}

func (d *diContainer) ExternalModule() (external.ModuleInterface, error) {
	if d.externalModule == nil {
		d.externalModule = external.NewExternalModule(d.Config())
	}
	return d.externalModule, nil
}

func (d *diContainer) FilmSimilarModule() (filmsimilar.ModuleInterface, error) {
	if d.filmSimilarModule == nil {
		storage, err := d.Storage()
		if err != nil {
			return nil, err
		}
		filmModule, err := d.FilmModule()
		if err != nil {
			return nil, err
		}
		d.filmSimilarModule = filmsimilar.NewFilmSimilarModule(storage, d.Config(), filmModule)
	}
	return d.filmSimilarModule, nil
}

func (d *diContainer) FilmModule() (film.ModuleInterface, error) {
	if d.filmModule == nil {
		storage, err := d.Storage()
		if err != nil {
			return nil, err
		}
		externalModule, err := d.ExternalModule()
		if err != nil {
			return nil, err
		}

		d.filmModule = film.NewFilmModule(storage, externalModule)
	}
	return d.filmModule, nil
}

func (d *diContainer) UserModule() (user.ModuleInterface, error) {
	if d.userModule == nil {
		storage, err := d.Storage()
		if err != nil {
			return nil, err
		}
		externalModule, err := d.ExternalModule()
		if err != nil {
			return nil, err
		}
		d.userModule = user.NewUserModule(storage, externalModule)
	}
	return d.userModule, nil
}

func (d *diContainer) FilmSequelModule() (filmsequel.ModuleInterface, error) {
	if d.filmSequelModule == nil {
		storage, err := d.Storage()
		if err != nil {
			return nil, err
		}
		filmModule, err := d.FilmModule()
		if err != nil {
			return nil, err
		}
		d.filmSequelModule = filmsequel.NewFilmSequelModule(storage, d.Config(), filmModule)
	}
	return d.filmSequelModule, nil
}

func (d *diContainer) GrpcClient() (grpc.ClientConnInterface, error) {
	if d.grpcClient == nil {
		var err error
		d.grpcClient, err = grpc.NewClient(d.Config())
		if err != nil {
			logger.Info("Клиент gRPC не создан: " + err.Error())
		}
	}
	return d.grpcClient, nil
}

func (d *diContainer) Storage() (*postgres.Storage, error) {
	if d.storage == nil {
		var err error
		d.storage, err = postgres.NewStorage(d.Config().DB.URL)
		if err != nil {
			return nil, err
		}
	}
	return d.storage, nil
}
