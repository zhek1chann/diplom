package handler

import (
	"context"
	"net/http"
	"strconv"

	"diploma/modules/auth/jwt"
	contractConverter "diploma/modules/contract/handler/converter"
	apiModel "diploma/modules/contract/handler/model"
	serviceModel "diploma/modules/contract/model"
	contextkeys "diploma/pkg/context-keys"
	"github.com/gin-gonic/gin"
)

type Service interface {
	SignContract(ctx context.Context, contractID int64, role int, signature string) error
	GetContract(ctx context.Context, id int64) (*serviceModel.Contract, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Sign godoc
// @Summary Подписать контракт
// @Description Подпись контракта (клиент или поставщик)
// @Tags contracts
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body apiModel.SignRequest true "Contract ID и Подпись"
// @Success 200 {object} map[string]string "Подпись сохранена"
// @Failure 400 {object} map[string]string "Ошибка валидации"
// @Failure 401 {object} map[string]string "Неавторизован"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /api/contract/sign [post]
func (h *Handler) Sign(c *gin.Context) {
	claims := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)

	var req apiModel.SignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.SignContract(c.Request.Context(), req.ContractID, claims.Role, req.Signature)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "signature recorded"})
}

// Get godoc
// @Summary Получить контракт
// @Description Возвращает контракт по ID
// @Tags contracts
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "ID контракта"
// @Success 200 {object} apiModel.ContractResponse "Контракт"
// @Failure 400 {object} map[string]string "Неверный ID"
// @Failure 401 {object} map[string]string "Неавторизован"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /api/contract/{id} [get]

func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	contractID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid contract ID"})
		return
	}

	contract, err := h.service.GetContract(c.Request.Context(), contractID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, contractConverter.ToAPI(contract))
}
