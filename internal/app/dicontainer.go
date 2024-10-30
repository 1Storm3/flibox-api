package app

import (
	"kinopoisk-api/database/postgres"
	"kinopoisk-api/internal/config"
	externalservice "kinopoisk-api/pkg/external-service"
)

type diContainer struct {
	config          *config.Config
	storage         *postgres.Storage
	externalService ExternalService

	filmModule        *filmModule
	userModule        *userModule
	filmSequelModule  *filmSequelModule
	filmSimilarModule *filmSimilarModule
	userFilmModule    *userFilmModule
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

func (d *diContainer) UserFilmModule() (*userFilmModule, error) {
	if d.userFilmModule == nil {
		storage, err := d.Storage()
		if err != nil {
			return nil, err
		}
		filmModule, err := d.FilmModule()
		if err != nil {
			return nil, err
		}
		d.userFilmModule = NewUserFilmModule(storage, filmModule)
	}
	return d.userFilmModule, nil
}

func (d *diContainer) FilmSimilarModule() (*filmSimilarModule, error) {
	if d.filmSimilarModule == nil {
		storage, err := d.Storage()
		if err != nil {
			return nil, err
		}
		filmModule, err := d.FilmModule()
		if err != nil {
			return nil, err
		}
		d.filmSimilarModule = NewFilmSimilarModule(storage, d.Config(), filmModule)
	}
	return d.filmSimilarModule, nil
}

func (d *diContainer) FilmModule() (*filmModule, error) {
	if d.filmModule == nil {
		storage, err := d.Storage()
		if err != nil {
			return nil, err
		}
		externalService, err := d.ExternalService()
		if err != nil {
			return nil, err
		}
		d.filmModule = NewFilmModule(storage, externalService)
	}
	return d.filmModule, nil
}

func (d *diContainer) UserModule() (*userModule, error) {
	if d.userModule == nil {
		storage, err := d.Storage()
		if err != nil {
			return nil, err
		}
		d.userModule = NewUserModule(storage)
	}
	return d.userModule, nil
}

func (d *diContainer) FilmSequelModule() (*filmSequelModule, error) {
	if d.filmSequelModule == nil {
		storage, err := d.Storage()
		if err != nil {
			return nil, err
		}
		filmModule, err := d.FilmModule()
		if err != nil {
			return nil, err
		}
		d.filmSequelModule = NewFilmSequelModule(storage, d.Config(), filmModule)
	}
	return d.filmSequelModule, nil
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

func (d *diContainer) ExternalService() (ExternalService, error) {
	if d.externalService == nil {
		d.externalService = externalservice.NewExternalService(d.config)
	}
	return d.externalService, nil
}
