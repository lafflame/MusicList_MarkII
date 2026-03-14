package repository

import (
	"MusicList_MarkII/internal/domain"

	"gorm.io/gorm"
)

type MediaRepository interface {
	FindAll() ([]domain.Media, error)
	Search(query string) ([]domain.Media, error)
	FilterByDate(from, to string) ([]domain.Media, error)
	Create(m *domain.Media) error
	Update(id string, m *domain.Media) error
	Delete(id string) error
}

type mediaRepo struct {
	db *gorm.DB
}

func NewMediaRepo(db *gorm.DB) MediaRepository {
	return &mediaRepo{db}
}

func (r *mediaRepo) FindAll() ([]domain.Media, error) {
	var tracks []domain.Media
	err := r.db.Find(&tracks).Error
	return tracks, err
}

func (r *mediaRepo) Search(query string) ([]domain.Media, error) {
	var tracks []domain.Media
	err := r.db.Where("artist ILIKE ? OR track ILIKE ?", "%"+query+"%", "%"+query+"%").Find(&tracks).Error
	return tracks, err
}

func (r *mediaRepo) FilterByDate(from, to string) ([]domain.Media, error) {
	var tracks []domain.Media
	q := r.db.Model(&domain.Media{})
	if from != "" {
		q = q.Where("created_at >= ?", from)
	}
	if to != "" {
		q = q.Where("created_at <= ?", to)
	}
	err := q.Find(&tracks).Error
	return tracks, err
}

func (r *mediaRepo) Create(m *domain.Media) error {
	return r.db.Create(m).Error
}

func (r *mediaRepo) Update(id string, m *domain.Media) error {
	return r.db.Model(&domain.Media{}).Where("track_id = ?", id).Updates(m).Error
}

func (r *mediaRepo) Delete(id string) error {
	return r.db.Delete(&domain.Media{}, id).Error
}
