package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sebx771/monitor-bot/internal/status" 
)

type Handler struct {
	store *status.Store
}

func NewHandler(store *status.Store) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) GetStatus(c *gin.Context) {
	workers := h.store.GetAll()
	c.JSON(http.StatusOK, workers)
}