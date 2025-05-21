package handlers

import (
	_ "controlDeviceServer/internal/storage/sqlite"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Tables(c *gin.Context) {
	logger := c.MustGet("logger").(*slog.Logger)
	logger.Debug("handling Tables request")
	c.JSON(http.StatusOK, "pong")
}

func (h *Handler) InsertCurrentAnalyzer(c *gin.Context) {
	h.insertGeneric(c, h.storage.InsertCurrentAnalyzer, "current analyzer")
}

func (h *Handler) InsertController(c *gin.Context) {
	h.insertGeneric(c, h.storage.InsertController, "controller")
}

func (h *Handler) InsertInputModule(c *gin.Context) {
	h.insertGeneric(c, h.storage.InsertInputModule, "input module")
}

func (h *Handler) InsertLCD(c *gin.Context) {
	h.insertGeneric(c, h.storage.InsertLCD, "lcd")
}

func (h *Handler) InsertShield(c *gin.Context) {
	h.insertGeneric(c, h.storage.InsertShield, "shield")
}

func (h *Handler) insertGeneric(c *gin.Context, insertFunc func(map[string]interface{}) error, entityName string) {
	logger := c.MustGet("logger").(*slog.Logger)

	var rows []map[string]interface{}
	if err := c.BindJSON(&rows); err != nil {
		logger.Error("failed to bind JSON array", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON array:" + err.Error()})
		return
	}

	var errors []string
	var successCount int
	var errorCount int
	for _, row := range rows {
		if err := insertFunc(row); err != nil {
			logger.Error("failed to insert "+entityName, "error", err)
			errors = append(errors, fmt.Sprintf("%v: ", row["№"])+err.Error())
			errorCount++
			continue
		}
		successCount++
	}
	if errorCount > 0 {
		c.JSON(http.StatusMultiStatus, gin.H{
			"inserted": successCount,
			"total":    len(rows),
			"errors":   errors,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"inserted": successCount,
		"total":    len(rows),
	})
}
