package repository

import (
	"MusicList_MarkII/internal/domain"

	"gorm.io/gorm"
)

type PlaylistRepository interface {
	FindAll() ([]domain.Playlist, error)
	FindByIDWithTracks(id string) (*domain.Playlist, error)
	Create(p *domain.Playlist) error
	Rename(id, name string) error
	Delete(id string) error
	AddTrack(playlistID, trackID string) error
	RemoveTrack(playlistID, trackID string) error
}

type playlistRepo struct {
	db *gorm.DB
}

func NewPlaylistRepo(db *gorm.DB) PlaylistRepository {
	return &playlistRepo{db}
}

func (r *playlistRepo) FindAll() ([]domain.Playlist, error) {
	var playlists []domain.Playlist
	err := r.db.Find(&playlists).Error
	return playlists, err
}

func (r *playlistRepo) FindByIDWithTracks(id string) (*domain.Playlist, error) {
	var playlist domain.Playlist
	err := r.db.Preload("Tracks").First(&playlist, id).Error
	return &playlist, err
}

func (r *playlistRepo) Create(p *domain.Playlist) error {
	return r.db.Create(p).Error
}

func (r *playlistRepo) Rename(id, name string) error {
	return r.db.Model(&domain.Playlist{}).Where("playlist_id = ?", id).Update("name", name).Error
}

func (r *playlistRepo) Delete(id string) error {
	return r.db.Delete(&domain.Playlist{}, id).Error
}

func (r *playlistRepo) AddTrack(playlistID, trackID string) error {
	var playlist domain.Playlist
	var track domain.Media
	r.db.First(&playlist, playlistID)
	r.db.First(&track, trackID)
	return r.db.Model(&playlist).Association("Tracks").Append(&track)
}

func (r *playlistRepo) RemoveTrack(playlistID, trackID string) error {
	var playlist domain.Playlist
	var track domain.Media
	r.db.First(&playlist, playlistID)
	r.db.First(&track, trackID)
	return r.db.Model(&playlist).Association("Tracks").Delete(&track)
}
