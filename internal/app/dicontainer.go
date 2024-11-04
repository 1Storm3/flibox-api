package app

import (
	"kbox-api/database/postgres"
	"kbox-api/internal/config"
	"kbox-api/internal/modules/auth"
	"kbox-api/internal/modules/external"
	"kbox-api/internal/modules/film"
	"kbox-api/internal/modules/film-sequel"
	"kbox-api/internal/modules/film-similar"
	"kbox-api/internal/modules/user"
	"kbox-api/internal/modules/user-film"
)

type diContainer struct {
	config  *config.Config
	storage *postgres.Storage

	filmModule        *film.Module
	externalModule    *external.Module
	userModule        *user.Module
	filmSequelModule  *filmsequel.Module
	filmSimilarModule *film_similar.Module
	userFilmModule    *user_film.Module
	authModule        *auth.Module
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

func (d *diContainer) AuthModule() (*auth.Module, error) {
	if d.authModule == nil {
		userModule, err := d.UserModule()
		if err != nil {
			return nil, err
		}
		d.authModule = auth.NewAuthModule(userModule, d.Config())
	}
	return d.authModule, nil
}

func (d *diContainer) ExternalModule() (*external.Module, error) {
	if d.externalModule == nil {
		d.externalModule = external.NewExternalModule(d.Config())
	}
	return d.externalModule, nil
}

func (d *diContainer) UserFilmModule() (*user_film.Module, error) {
	if d.userFilmModule == nil {
		storage, err := d.Storage()
		if err != nil {
			return nil, err
		}
		filmModule, err := d.FilmModule()
		if err != nil {
			return nil, err
		}
		d.userFilmModule = user_film.NewUserFilmModule(storage, filmModule)
	}
	return d.userFilmModule, nil
}

func (d *diContainer) FilmSimilarModule() (*film_similar.Module, error) {
	if d.filmSimilarModule == nil {
		storage, err := d.Storage()
		if err != nil {
			return nil, err
		}
		filmModule, err := d.FilmModule()
		if err != nil {
			return nil, err
		}
		d.filmSimilarModule = film_similar.NewFilmSimilarModule(storage, d.Config(), filmModule)
	}
	return d.filmSimilarModule, nil
}

func (d *diContainer) FilmModule() (*film.Module, error) {
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

func (d *diContainer) UserModule() (*user.Module, error) {
	if d.userModule == nil {
		storage, err := d.Storage()
		if err != nil {
			return nil, err
		}
		d.userModule = user.NewUserModule(storage)
	}
	return d.userModule, nil
}

func (d *diContainer) FilmSequelModule() (*filmsequel.Module, error) {
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
