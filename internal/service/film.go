package service

type FilmService struct {
	filmRepo FilmRepository
}

func NewFilmService(filmRepo FilmRepository) *FilmService {
	return &FilmService{
		filmRepo: filmRepo,
	}
}
