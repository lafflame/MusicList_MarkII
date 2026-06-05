package service

import (
	"MusicList_MarkII/internal/domain"
	"errors"
	"testing"
)

// mockMediaRepo — заглушка репозитория для изоляции от БД
type mockMediaRepo struct {
	tracks []domain.Media
	err    error
}

func (m *mockMediaRepo) FindAll() ([]domain.Media, error) {
	return m.tracks, m.err
}

func (m *mockMediaRepo) Search(query string) ([]domain.Media, error) {
	return m.tracks, m.err
}

func (m *mockMediaRepo) FilterByDate(from, to string) ([]domain.Media, error) {
	return m.tracks, m.err
}

func (m *mockMediaRepo) Create(t *domain.Media) error {
	return m.err
}

func (m *mockMediaRepo) Update(id string, t *domain.Media) error {
	return m.err
}

func (m *mockMediaRepo) Delete(id string) error {
	return m.err
}

// TestGetAll — проверяет, что GetAll возвращает все треки без ошибок
func TestGetAll(t *testing.T) {
	repo := &mockMediaRepo{
		tracks: []domain.Media{
			{TrackID: 1, Artist: "Travis Scott", Track: "Antidote"},
			{TrackID: 2, Artist: "Михаил Круг", Track: "Девочка-Пай"},
		},
	}
	svc := NewMediaService(repo)

	result, err := svc.GetAll()
	if err != nil {
		t.Fatalf("GetAll вернул ошибку: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("ожидалось 2 трека, получено %d", len(result))
	}
}

// TestSearch_ByArtist — проверяет поиск по исполнителю
func TestSearch_ByArtist(t *testing.T) {
	repo := &mockMediaRepo{
		tracks: []domain.Media{
			{TrackID: 1, Artist: "Travis Scott", Track: "Antidote"},
		},
	}
	svc := NewMediaService(repo)

	result, err := svc.Search("Travis")
	if err != nil {
		t.Fatalf("Search вернул ошибку: %v", err)
	}
	if len(result) == 0 {
		t.Error("Search не вернул результатов по запросу 'Travis'")
	}
}

// TestSearch_ByTrack — проверяет поиск по названию трека
func TestSearch_ByTrack(t *testing.T) {
	repo := &mockMediaRepo{
		tracks: []domain.Media{
			{TrackID: 2, Artist: "Михаил Круг", Track: "Девочка-Пай"},
		},
	}
	svc := NewMediaService(repo)

	result, err := svc.Search("Девочка")
	if err != nil {
		t.Fatalf("Search вернул ошибку: %v", err)
	}
	if len(result) == 0 {
		t.Error("Search не вернул результатов по запросу 'Девочка'")
	}
}

// TestAdd — проверяет добавление трека без ошибок
func TestAdd(t *testing.T) {
	repo := &mockMediaRepo{}
	svc := NewMediaService(repo)

	err := svc.Add(&domain.Media{Artist: "Ludacris", Track: "Act A Fool"})
	if err != nil {
		t.Fatalf("Add вернул ошибку: %v", err)
	}
}

// TestDelete — проверяет удаление трека без ошибок
func TestDelete(t *testing.T) {
	repo := &mockMediaRepo{}
	svc := NewMediaService(repo)

	err := svc.Delete("1")
	if err != nil {
		t.Fatalf("Delete вернул ошибку: %v", err)
	}
}

// TestGetStatistics_PopularArtist — проверяет, что статистика корректно определяет популярного исполнителя
func TestGetStatistics_PopularArtist(t *testing.T) {
	repo := &mockMediaRepo{
		tracks: []domain.Media{
			{TrackID: 1, Artist: "Travis Scott", Track: "Antidote"},
			{TrackID: 2, Artist: "Travis Scott", Track: "Meltdown"},
			{TrackID: 3, Artist: "Михаил Круг", Track: "Девочка-Пай"},
		},
	}
	svc := NewMediaService(repo)

	stats := svc.GetStatistics()
	if stats["popular_artist"] != "Travis Scott" {
		t.Errorf("ожидался 'Travis Scott', получено '%v'", stats["popular_artist"])
	}
	if stats["total_tracks"] != 3 {
		t.Errorf("ожидалось 3 трека, получено %v", stats["total_tracks"])
	}
}

// TestShuffle — проверяет, что Shuffle возвращает все треки без потерь
func TestShuffle(t *testing.T) {
	original := []domain.Media{
		{TrackID: 1, Artist: "Travis Scott", Track: "Antidote"},
		{TrackID: 2, Artist: "Михаил Круг", Track: "Девочка-Пай"},
		{TrackID: 3, Artist: "Ludacris", Track: "Act A Fool"},
	}
	repo := &mockMediaRepo{tracks: original}
	svc := NewMediaService(repo)

	result, err := svc.Shuffle()
	if err != nil {
		t.Fatalf("Shuffle вернул ошибку: %v", err)
	}
	if len(result) != len(original) {
		t.Errorf("Shuffle изменил количество треков: было %d, стало %d", len(original), len(result))
	}
}

// TestGetAll_Error — проверяет обработку ошибки репозитория
func TestGetAll_Error(t *testing.T) {
	repo := &mockMediaRepo{err: errors.New("db connection failed")}
	svc := NewMediaService(repo)

	_, err := svc.GetAll()
	if err == nil {
		t.Error("ожидалась ошибка, но GetAll вернул nil")
	}
}

// TestAdd_EmptyArtist — добавление трека без исполнителя должно вернуть ошибку
func TestAdd_EmptyArtist(t *testing.T) {
	repo := &mockMediaRepo{}
	svc := NewMediaService(repo)

	err := svc.Add(&domain.Media{Artist: "", Track: "Some Track"})
	if err == nil {
		t.Error("ожидалась ошибка при пустом исполнителе")
	}
}

// TestAdd_EmptyTrack — добавление трека без названия должно вернуть ошибку
func TestAdd_EmptyTrack(t *testing.T) {
	repo := &mockMediaRepo{}
	svc := NewMediaService(repo)

	err := svc.Add(&domain.Media{Artist: "Travis Scott", Track: ""})
	if err == nil {
		t.Error("ожидалась ошибка при пустом названии трека")
	}
}

// TestAdd_WhitespaceOnly — пробелы приравниваются к пустой строке
func TestAdd_WhitespaceOnly(t *testing.T) {
	repo := &mockMediaRepo{}
	svc := NewMediaService(repo)

	err := svc.Add(&domain.Media{Artist: "   ", Track: "   "})
	if err == nil {
		t.Error("ожидалась ошибка при строке из пробелов")
	}
}

// TestUpdate_EmptyFields — обновление с пустыми полями должно вернуть ошибку
func TestUpdate_EmptyFields(t *testing.T) {
	repo := &mockMediaRepo{}
	svc := NewMediaService(repo)

	err := svc.Update("1", &domain.Media{Artist: "", Track: ""})
	if err == nil {
		t.Error("ожидалась ошибка при обновлении с пустыми полями")
	}
}
