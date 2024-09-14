package external

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"kinopoisk-api/internal/config"
	"net/http"
)

type Film struct {
	ExternalId      int     `json:"kinopoiskId"`
	NameRu          string  `json:"nameRu"`
	NameOriginal    string  `json:"nameOriginal"`
	Year            int     `json:"year"`
	PosterUrl       string  `json:"posterUrl"`
	RatingKinopoisk float64 `json:"ratingKinopoisk"`
	Description     string  `json:"description"`
	LogoUrl         string  `json:"logoUrl"`
	Type            string  `json:"type"`
}

func GetFilm(ctx *fiber.Ctx, cfg *config.Config) error {
	filmId := ctx.Params("id")

	apiKey := cfg.DB.ApiKey
	baseUrlForAllFilms := "https://kinopoiskapiunofficial.tech/api/v2.2/films/"

	urlAllFilms := fmt.Sprintf("%s%s", baseUrlForAllFilms, filmId)

	req, err := http.NewRequest("GET", urlAllFilms, nil)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	req.Header.Add("X-API-KEY", apiKey)

	client := &http.Client{}
	resAllFilms, err := client.Do(req)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}
	defer resAllFilms.Body.Close()

	bodyAllFilms, err := io.ReadAll(resAllFilms.Body)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	var film Film
	err = json.Unmarshal(bodyAllFilms, &film)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to parse JSON")
	}

	return ctx.JSON(film)
}
