package filmsequel

import (
	"kbox-api/database/postgres"
	"kbox-api/internal/config"
	"kbox-api/internal/modules/film"
	"kbox-api/internal/modules/film-sequel/handler"
	"kbox-api/internal/modules/film-sequel/repository"
	"kbox-api/internal/modules/film-sequel/service"
)

var _ ModuleInterface = (*Module)(nil)

type ModuleInterface interface {
	FilmSequelService() (service.FilmSequelServiceInterface, error)
	FilmSequelRepository() (repository.FilmSequelRepositoryInterface, error)
	FilmSequelHandler() (handler.FilmSequelHandlerInterface, error)
}

type Module struct {
	storage              *postgres.Storage
	config               *config.Config
	filmSequelRepository repository.FilmSequelRepositoryInterface
	filmSequelService    service.FilmSequelServiceInterface
	filmModule           film.ModuleInterface
	filmSequelHandler    handler.FilmSequelHandlerInterface
}

func NewFilmSequelModule(
	storage *postgres.Storage,
	config *config.Config,
	filmModule film.ModuleInterface,
) *Module {
	return &Module{
		storage:    storage,
		config:     config,
		filmModule: filmModule,
	}
}

func (f *Module) FilmSequelService() (service.FilmSequelServiceInterface, error) {
	if f.filmSequelService == nil {
		repo, err := f.FilmSequelRepository()
		if err != nil {
			return nil, err
		}
		filmService, err := f.filmModule.FilmService()
		f.filmSequelService = service.NewFilmsSequelService(repo, f.config, filmService)
	}
	return f.filmSequelService, nil
}

func (f *Module) FilmSequelRepository() (repository.FilmSequelRepositoryInterface, error) {
	if f.filmSequelRepository == nil {
		f.filmSequelRepository = repository.NewFilmSequelRepository(f.storage)
	}
	return f.filmSequelRepository, nil
}

func (f *Module) FilmSequelHandler() (handler.FilmSequelHandlerInterface, error) {
	if f.filmSequelHandler == nil {
		filmSequelService, err := f.FilmSequelService()
		if err != nil {
			return nil, err
		}
		f.filmSequelHandler = handler.NewFilmSequelHandler(filmSequelService)
	}
	return f.filmSequelHandler, nil
}
