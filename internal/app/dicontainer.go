package app

import (
	"kinopoisk-api/database/postgres"
	"kinopoisk-api/internal/config"
	"kinopoisk-api/internal/repository"
	"kinopoisk-api/internal/service"
)

type diContainer struct {
	config                *config.Config
	storage               *postgres.Storage
	filmRepository        FilmRepository
	filmService           FilmService
	userRepository        UserRepository
	userService           UserService
	filmSequelRepository  FilmSequelRepository
	filmSequelService     FilmSequelService
	filmSimilarRepository FilmSimilarRepository
	filmSimilarService    FilmSimilarService
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

func (d *diContainer) FilmSequelRepository() (FilmSequelRepository, error) {
	if d.filmSequelRepository == nil {
		storage, err := d.Storage()
		if err != nil {
			return nil, err
		}
		d.filmSequelRepository = repository.NewFilmSequelRepository(storage)
	}
	return d.filmSequelRepository, nil
}

func (d *diContainer) FilmSimilarRepository() (FilmSimilarRepository, error) {
	if d.filmSimilarRepository == nil {
		storage, err := d.Storage()
		if err != nil {
			return nil, err
		}
		d.filmSimilarRepository = repository.NewFilmSimilarRepository(storage)
	}
	return d.filmSimilarRepository, nil
}

func (d *diContainer) FilmSequelService() (FilmSequelService, error) {
	if d.filmSequelService == nil {
		repo, err := d.FilmSequelRepository()
		if err != nil {
			return nil, err
		}
		filmService, err := d.FilmService()
		if err != nil {

		}
		d.filmSequelService = service.NewFilmsSequelService(repo, d.Config(), filmService)
	}
	return d.filmSequelService, nil
}

func (d *diContainer) FilmSimilarService() (FilmSimilarService, error) {
	if d.filmSimilarService == nil {
		repo, err := d.FilmSimilarRepository()
		if err != nil {
			return nil, err
		}
		filmService, err := d.FilmService()
		if err != nil {

		}
		d.filmSimilarService = service.NewFilmsSimilarService(repo, d.Config(), filmService)
	}
	return d.filmSimilarService, nil
}

func (d *diContainer) UserService() (UserService, error) {
	if d.userService == nil {
		repo, err := d.UserRepository()
		if err != nil {
			return nil, err
		}
		d.userService = service.NewUserService(repo)
	}
	return d.userService, nil
}

func (d *diContainer) UserRepository() (UserRepository, error) {
	if d.userRepository == nil {
		storage, err := d.Storage()
		if err != nil {
			return nil, err
		}
		d.userRepository = repository.NewUserRepository(storage)
	}
	return d.userRepository, nil
}
