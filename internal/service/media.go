package service

import (
	"MusicList_MarkII/internal/domain"
	"MusicList_MarkII/internal/repository"
	"math/rand"
)

type MediaService struct {
	repo repository.MediaRepository
}

func NewMediaService(repo repository.MediaRepository) *MediaService {
	return &MediaService{repo}
}

func (s *MediaService) GetAll() ([]domain.Media, error) {
	return s.repo.FindAll()
}

func (s *MediaService) Search(query string) ([]domain.Media, error) {
	return s.repo.Search(query)
}

func (s *MediaService) FilterByDate(from, to string) ([]domain.Media, error) {
	return s.repo.FilterByDate(from, to)
}

func (s *MediaService) Add(m *domain.Media) error {
	return s.repo.Create(m)
}

func (s *MediaService) Update(id string, m *domain.Media) error {
	return s.repo.Update(id, m)
}

func (s *MediaService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *MediaService) GetStatistics() map[string]any {
	tracks, _ := s.repo.FindAll()
	artistCounts := make(map[string]int)
	for _, t := range tracks {
		artistCounts[t.Artist]++
	}
	maxCount := 0
	popularArtist := ""
	for artist, count := range artistCounts {
		if count > maxCount {
			maxCount = count
			popularArtist = artist
		}
	}
	return map[string]any{
		"total_tracks":   len(tracks),
		"artist_counts":  artistCounts,
		"popular_artist": popularArtist,
	}
}

func (s *MediaService) Shuffle() ([]domain.Media, error) {
	tracks, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	rand.Shuffle(len(tracks), func(i, j int) {
		tracks[i], tracks[j] = tracks[j], tracks[i]
	})
	return tracks, nil
}
