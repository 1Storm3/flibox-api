package historyfilms

import (
	"github.com/1Storm3/flibox-api/internal/modules/recommendation/adapter"
	"github.com/1Storm3/flibox-api/internal/shared/logger"
	"github.com/1Storm3/flibox-api/pkg/token"
	"github.com/gofiber/fiber/v2"
)

type HandlerInterface interface {
	Add(c *fiber.Ctx) error
}

type Handler struct {
	service          ServiceInterface
	recommendService adapter.RecommendService
}

func NewHistoryFilmsHandler(
	service ServiceInterface,
	recommendService adapter.RecommendService,
) *Handler {
	return &Handler{
		service:          service,
		recommendService: recommendService,
	}
}

func (h *Handler) Add(c *fiber.Ctx) error {
	userID := c.Locals("userClaims").(*token.Claims).UserID
	filmID := c.Params("Id")
	ctx := c.Context()
	err := h.service.Add(ctx, filmID, userID)
	if err != nil {
		return err
	}
	go func() {
		err := h.recommendService.CreateRecommendations(adapter.RecommendationsParams{
			UserID: userID,
		})
		if err != nil {
			logger.Info("Произошла ошибка при создании рекомендаций")
			logger.Error(err.Error())
		}
	}()

	return c.JSON(fiber.Map{
		"data": "Фильм добавлен в историю просмотра",
	})
}
