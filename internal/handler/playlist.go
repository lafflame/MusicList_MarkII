package handler

import (
	"MusicList_MarkII/internal/domain"
	"MusicList_MarkII/internal/service"

	"github.com/gin-gonic/gin"
)

type PlaylistHandler struct {
	service *service.PlaylistService
}

func NewPlaylistHandler(s *service.PlaylistService) *PlaylistHandler {
	return &PlaylistHandler{s}
}

func (h *PlaylistHandler) Register(r *gin.Engine) {
	api := r.Group("/api")
	api.GET("/playlists", h.GetAll)
	api.POST("/playlists", h.Create)
	api.PUT("/playlists/:id", h.Rename)
	api.DELETE("/playlists/:id", h.Delete)
	api.GET("/playlists/:id/tracks", h.GetTracks)
	api.POST("/playlists/:id/tracks/:track_id", h.AddTrack)
	api.DELETE("/playlists/:id/tracks/:track_id", h.RemoveTrack)
}

func (h *PlaylistHandler) GetAll(c *gin.Context) {
	playlists, err := h.service.GetAll()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, playlists)
}

func (h *PlaylistHandler) Create(c *gin.Context) {
	var p domain.Playlist
	if err := c.BindJSON(&p); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Create(&p); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, p)
}

func (h *PlaylistHandler) Rename(c *gin.Context) {
	var body struct {
		Name string `json:"name"`
	}
	c.BindJSON(&body)
	if err := h.service.Rename(c.Param("id"), body.Name); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "переименовано"})
}

func (h *PlaylistHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Param("id")); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "удалено"})
}

func (h *PlaylistHandler) GetTracks(c *gin.Context) {
	playlist, err := h.service.GetTracks(c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, playlist.Tracks)
}

func (h *PlaylistHandler) AddTrack(c *gin.Context) {
	if err := h.service.AddTrack(c.Param("id"), c.Param("track_id")); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "трек добавлен"})
}

func (h *PlaylistHandler) RemoveTrack(c *gin.Context) {
	if err := h.service.RemoveTrack(c.Param("id"), c.Param("track_id")); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "трек удалён"})
}
