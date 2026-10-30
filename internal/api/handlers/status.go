package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sebx771/monitor-bot/internal/status"
)

// StatusHandler manages status queries and manual worker triggers.
type StatusHandler struct {
	store *status.Store
}

// NewStatusHandler creates a new StatusHandler instance.
func NewStatusHandler(store *status.Store) *StatusHandler {
	return &StatusHandler{
		store: store,
	}
}

// GetStatus returns the current snapshot of all registered workers.
// Endpoint: GET /api/v1/status
func (h *StatusHandler) GetStatus(c *gin.Context) {
	workers := h.store.GetAll()
	c.JSON(http.StatusOK, workers)
}

// TriggerWorker triggers an immediate run of the specified worker.
// Endpoint: POST /api/v1/workers/:id/trigger
func (h *StatusHandler) TriggerWorker(c *gin.Context) {
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
