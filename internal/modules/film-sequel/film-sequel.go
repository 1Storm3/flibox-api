package filmsequel

import (
	"kbox-api/database/postgres"
	"kbox-api/internal/config"
	"kbox-api/internal/modules/film"
	"kbox-api/internal/modules/film-sequel/handler"
	"kbox-api/internal/modules/film-sequel/repository"
	"kbox-api/internal/modules/film-sequel/service"
)

type Module struct {
	storage              *postgres.Storage
	config               *config.Config
	filmSequelRepository service.FilmSequelRepository
	filmSequelService    handler.FilmSequelService
	filmModule           *film.Module
	filmSequelHandler    *handler.FilmSequelHandler
}

func NewFilmSequelModule(
	storage *postgres.Storage,
	config *config.Config,
	filmModule *film.Module,
) *Module {
	return &Module{
		storage:    storage,
		config:     config,
		filmModule: filmModule,
	}
}

func (f *Module) FilmSequelService() (handler.FilmSequelService, error) {
	if f.filmSequelService == nil {
		repo, err := f.FilmSequelRepository()
		if err != nil {
			return nil, err
		}
		filmService, err := f.filmModule.FilmService()
		if err != nil {

		}
		f.filmSequelService = service.NewFilmsSequelService(repo, f.config, filmService)
	}
	return f.filmSequelService, nil
}

func (f *Module) FilmSequelRepository() (service.FilmSequelRepository, error) {
	if f.filmSequelRepository == nil {
		f.filmSequelRepository = repository.NewFilmSequelRepository(f.storage)
	}
	return f.filmSequelRepository, nil
}

func (f *Module) FilmSequelHandler() (*handler.FilmSequelHandler, error) {
	if f.filmSequelHandler == nil {
		filmSequelService, err := f.FilmSequelService()
		if err != nil {
			return nil, err
		}
		f.filmSequelHandler = handler.NewFilmSequelHandler(filmSequelService)
	}
	return f.filmSequelHandler, nil
}
