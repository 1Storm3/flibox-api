package app

import (
	"kinopoisk-api/database/postgres"
	"kinopoisk-api/internal/config"
	"kinopoisk-api/internal/repository"
	"kinopoisk-api/internal/service"
)

type diContainer struct {
	config         *config.Config
	storage        *postgres.Storage
	filmRepository FilmRepository
	filmService    FilmService
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

func (d *diContainer) FilmService() (FilmService, error) {
	if d.filmService == nil {
		repo, err := d.FilmRepository()
		if err != nil {
			return nil, err
		}
		d.filmService = service.NewFilmService(repo, d.Config())
	}

	return d.filmService, nil
}

func (d *diContainer) FilmRepository() (FilmRepository, error) {
	if d.filmRepository == nil {
		storage, err := d.Storage()
		if err != nil {
			return nil, err
		}
		d.filmRepository = repository.NewFilmRepository(storage)
	}
	return d.filmRepository, nil

}
