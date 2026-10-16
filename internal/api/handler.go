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

func (h *Handler) TriggerWorker(c *gin.Context) {
	id := c.Param("id")
	w, ok := h.store.GetWorker(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "worker not found or not active"})
		return
	}
	if !w.Trigger() {
		c.JSON(http.StatusConflict, gin.H{"error": "worker is busy or a trigger is already pending"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "trigger sent", "worker": id})
}