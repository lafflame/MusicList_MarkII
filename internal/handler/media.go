package handler

import (
	"MusicList_MarkII/internal/domain"
	"MusicList_MarkII/internal/service"

	"github.com/gin-gonic/gin"
)

type MediaHandler struct {
	service *service.MediaService
}

func NewMediaHandler(s *service.MediaService) *MediaHandler {
	return &MediaHandler{s}
}

func (h *MediaHandler) Register(r *gin.Engine) {
	api := r.Group("/api")
	api.GET("/tracks/search", h.Search)
	api.GET("/tracks/filter", h.Filter)
	api.GET("/tracks", h.GetAll)
	api.POST("/tracks", h.Add)
	api.PUT("/tracks/:id", h.Update)
	api.DELETE("/tracks/:id", h.Delete)
	api.GET("/statistics", h.Statistics)
}

func (h *MediaHandler) GetAll(c *gin.Context) {
	tracks, err := h.service.GetAll()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, tracks)
}

func (h *MediaHandler) Add(c *gin.Context) {
	var m domain.Media
	if err := c.BindJSON(&m); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Add(&m); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, m)
}

func (h *MediaHandler) Update(c *gin.Context) {
	var m domain.Media
	if err := c.BindJSON(&m); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Update(c.Param("id"), &m); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, m)
}

func (h *MediaHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Param("id")); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "удалено"})
}

func (h *MediaHandler) Search(c *gin.Context) {
	tracks, err := h.service.Search(c.Query("query"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, tracks)
}

func (h *MediaHandler) Filter(c *gin.Context) {
	tracks, err := h.service.FilterByDate(c.Query("from"), c.Query("to"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, tracks)
}

func (h *MediaHandler) Statistics(c *gin.Context) {
	c.JSON(200, h.service.GetStatistics())
}
