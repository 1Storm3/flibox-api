package app

import (
	"kinopoisk-api/internal/config"
	"kinopoisk-api/internal/repository"
	"kinopoisk-api/internal/service"
)

type diContainer struct {
	config         *config.Config
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
		d.filmRepository = repository.NewFilmRepository()
	}
	return d.filmRepository, nil

}
