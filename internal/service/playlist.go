package service

import (
	"MusicList_MarkII/internal/domain"
	"MusicList_MarkII/internal/repository"
)

type PlaylistService struct {
	repo repository.PlaylistRepository
}

func NewPlaylistService(repo repository.PlaylistRepository) *PlaylistService {
	return &PlaylistService{repo}
}

func (s *PlaylistService) GetAll() ([]domain.Playlist, error) {
	return s.repo.FindAll()
}

func (s *PlaylistService) GetTracks(id string) (*domain.Playlist, error) {
	return s.repo.FindByIDWithTracks(id)
}

func (s *PlaylistService) Create(p *domain.Playlist) error {
	return s.repo.Create(p)
}

func (s *PlaylistService) Rename(id, name string) error {
	return s.repo.Rename(id, name)
}

func (s *PlaylistService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *PlaylistService) AddTrack(playlistID, trackID string) error {
	return s.repo.AddTrack(playlistID, trackID)
}

func (s *PlaylistService) RemoveTrack(playlistID, trackID string) error {
	return s.repo.RemoveTrack(playlistID, trackID)
}
