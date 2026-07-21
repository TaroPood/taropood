package handler

import (
	"log/slog"
	"net/http"

	"github.com/TaroPood/taropood/internal/dto"
	ruleuc "github.com/TaroPood/taropood/internal/usecases/rule"
	"github.com/gin-gonic/gin"
)

type RuleHandler struct {
	uc *ruleuc.UseCase
}

func NewRuleHandler(uc *ruleuc.UseCase) *RuleHandler {
	return &RuleHandler{uc: uc}
}

func (h *RuleHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/rules", h.Create)
	r.GET("/rules/:id", h.GetByID)
}

func (h *RuleHandler) Create(c *gin.Context) {
	var req dto.CreateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON: " + err.Error()})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "name is required"})
		return
	}

	resp, err := h.uc.Create(c.Request.Context(), &req)
	if err != nil {
		if ruleuc.IsDuplicate(err) {
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
			return
		}
		slog.Error("failed to create rule", "err", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to create rule"})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *RuleHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	resp, err := h.uc.GetByID(c.Request.Context(), id)
	if err != nil {
		if ruleuc.IsNotFound(err) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		slog.Error("failed to get rule", "id", id, "err", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to get rule"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
