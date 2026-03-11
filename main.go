// @title Music Library API
// @version 1.0
// @description API для управления личной музыкальной коллекцией
// @host localhost:8080
// @BasePath /api
package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	_ "MusicList_MarkII/docs"

	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Media struct {
	TrackID   uint      `gorm:"primary_key" json:"track_id"`
	Artist    string    `json:"artist"`
	Track     string    `json:"track"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}

type Playlist struct {
	PlaylistID uint    `gorm:"primary_key" json:"playlist_id"`
	Name       string  `json:"name"`
	Tracks     []Media `gorm:"many2many:playlist_tracks;" json:"tracks,omitempty"`
}

func main() {
	dsn := "host=localhost user=postgres password=0311 dbname=mydb port=5432 sslmode=disable client_encoding=UTF8"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database", err)
	}
	db.AutoMigrate(&Media{}, &Playlist{})
	err = db.AutoMigrate(&Media{}, &Playlist{})
	if err != nil {
		log.Fatal("Ошибка миграции:", err)
	}
	// Запускаем HTTP сервер со Swagger в отдельной горутине
	go startAPI(db)

	// CLI работает как раньше
	menu(db)
}

func startAPI(db *gorm.DB) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Сначала статические маршруты
	r.GET("/api/tracks/search", func(c *gin.Context) { searchTracks(c, db) })
	r.GET("/api/tracks/filter", func(c *gin.Context) { filterTracks(c, db) })
	r.GET("/api/statistics", func(c *gin.Context) { getStatistics(c, db) })

	// Потом параметрические
	r.GET("/api/tracks", func(c *gin.Context) { getTracks(c, db) })
	r.POST("/api/tracks", func(c *gin.Context) { addTrack(c, db) })
	r.DELETE("/api/tracks/:id", func(c *gin.Context) { deleteTrack(c, db) })
	r.PUT("/api/tracks/:id", func(c *gin.Context) { updateTrack(c, db) })

	// Плейлисты — тоже сначала статические
	r.GET("/api/playlists", func(c *gin.Context) { getPlaylists(c, db) })
	r.POST("/api/playlists", func(c *gin.Context) { createPlaylist(c, db) })
	r.DELETE("/api/playlists/:id", func(c *gin.Context) { deletePlaylist(c, db) })
	r.PUT("/api/playlists/:id", func(c *gin.Context) { renamePlaylist(c, db) })
	r.GET("/api/playlists/:id/tracks", func(c *gin.Context) { getPlaylistTracks(c, db) })
	r.POST("/api/playlists/:id/tracks/:track_id", func(c *gin.Context) { addTrackToPlaylist(c, db) })
	r.DELETE("/api/playlists/:id/tracks/:track_id", func(c *gin.Context) { removeTrackFromPlaylist(c, db) })

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	fmt.Println("Swagger UI: http://localhost:8080/swagger/index.html")
	r.Run(":8080")
}

// @Summary Получить все треки
// @Tags tracks
// @Produce json
// @Success 200 {array} Media
// @Router /tracks [get]
func getTracks(c *gin.Context, db *gorm.DB) {
	var tracks []Media
	db.Find(&tracks)
	c.JSON(200, tracks)
}

// @Summary Добавить трек
// @Tags tracks
// @Accept json
// @Produce json
// @Param track body Media true "Данные трека"
// @Success 201 {object} Media
// @Router /tracks [post]
func addTrack(c *gin.Context, db *gorm.DB) {
	var media Media
	if err := c.BindJSON(&media); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	db.Create(&media)
	c.JSON(201, media)
}

// @Summary Удалить трек
// @Tags tracks
// @Param id path int true "ID трека"
// @Success 200
// @Router /tracks/{id} [delete]
func deleteTrack(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	db.Delete(&Media{}, id)
	c.JSON(200, gin.H{"message": "удалено"})
}

// @Summary Обновить трек
// @Tags tracks
// @Accept json
// @Produce json
// @Param id path int true "ID трека"
// @Param track body Media true "Новые данные"
// @Success 200 {object} Media
// @Router /tracks/{id} [put]
func updateTrack(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var media Media
	if err := c.BindJSON(&media); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	db.Model(&Media{}).Where("track_id = ?", id).Updates(media)
	c.JSON(200, media)
}

// @Summary Поиск треков
// @Tags tracks
// @Param query query string true "Название или исполнитель"
// @Produce json
// @Success 200 {array} Media
// @Router /tracks/search [get]
func searchTracks(c *gin.Context, db *gorm.DB) {
	query := c.Query("query")
	var tracks []Media
	db.Where("artist ILIKE ? OR track ILIKE ?", "%"+query+"%", "%"+query+"%").Find(&tracks)
	c.JSON(200, tracks)
}

// @Summary Статистика медиатеки
// @Tags statistics
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /statistics [get]
func getStatistics(c *gin.Context, db *gorm.DB) {
	var tracks []Media
	db.Find(&tracks)

	artistCounts := make(map[string]int)
	for _, track := range tracks {
		artistCounts[track.Artist]++
	}

	maxCount := 0
	popularArtist := ""
	for artist, count := range artistCounts {
		if count > maxCount {
			maxCount = count
			popularArtist = artist
		}
	}

	c.JSON(200, gin.H{
		"total_tracks":   len(tracks),
		"artist_counts":  artistCounts,
		"popular_artist": popularArtist,
	})
}

// Все твои старые CLI функции остаются без изменений ниже
func menu(db *gorm.DB) {
	var choice int
	fmt.Println("Выберите пункт:\n1.Вывести все треки\n2.Добавить трек\n3.Удалить трек" +
		"\n4.Изменить трек по ID\n5.Перемешивание треков\n6.Поиск клипа на YouTube" +
		"\n7.Показать статистику\n8.Поиск треков\n")
	fmt.Scan(&choice)

	switch choice {
	case 1:
		output(db)
	case 2:
		add(db)
	case 3:
		del(db)
	case 4:
		changing(db)
	case 5:
		shuffleAndOutput(db)
	case 6:
		playYouTubeClip()
	case 7:
		showStatistics(db)
	case 8:
		searching(db)
	case 10:
		fmt.Println("Хорошего дня!")
		return
	case 228:
		gettingInfo()
	default:
		fmt.Println("Неправильный выбор")
	}
	menu(db)
}

func add(db *gorm.DB) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Введите исполнителя:")
	artist, _ := reader.ReadString('\n')
	artist = strings.TrimSpace(artist)
	fmt.Println("Введите название трека:")
	track, _ := reader.ReadString('\n')
	track = strings.TrimSpace(track)
	fmt.Println("Введите ссылку на клип:")
	url, _ := reader.ReadString('\n')
	url = strings.TrimSpace(url)
	db.Create(&Media{Track: track, Artist: artist, URL: url})
}

func del(db *gorm.DB) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Введите ID исполнителя:")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)
	db.Delete(&Media{}, id)
}

func output(db *gorm.DB) {
	var tracks []Media
	db.Find(&tracks)
	fmt.Println("\nID | Исполнитель | Название трека | Ссылка на клип\n")
	for _, track := range tracks {
		fmt.Printf("%d: %s - %s  |  %s\n", track.TrackID, track.Artist, track.Track, track.URL)
	}
	fmt.Print("\n\n")
}

func changing(db *gorm.DB) {
	var id int
	fmt.Println("Введите ID для изменения трека:")
	fmt.Scan(&id)
	var artist, track string
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Введите новое имя исполнителя:")
	artist, _ = reader.ReadString('\n')
	artist = strings.TrimSpace(artist)
	fmt.Println("Введите новое название трека:")
	track, _ = reader.ReadString('\n')
	track = strings.TrimSpace(track)
	fmt.Println("Введите новую ссылку на клип:")
	url, _ := reader.ReadString('\n')
	url = strings.TrimSpace(url)
	db.Model(&Media{}).Where("track_id = ?", id).Updates(Media{Artist: artist, Track: track, URL: url})
}

func shuffleAndOutput(db *gorm.DB) {
	var tracks []Media
	db.Find(&tracks)
	rand.Shuffle(len(tracks), func(i, j int) {
		tracks[i], tracks[j] = tracks[j], tracks[i]
	})
	fmt.Println("\nПеремешанные треки:\nID | Исполнитель | Название трека | Ссылка на клип\n")
	for _, track := range tracks {
		fmt.Printf("%d: %s - %s  |  %s\n", track.TrackID, track.Artist, track.Track, track.URL)
	}
	fmt.Print("\n\n")
}

func openBrowser(url string) error {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func playYouTubeClip() {
	fmt.Println("Введите название трека или исполнителя для поиска на YouTube:")
	reader := bufio.NewReader(os.Stdin)
	query, _ := reader.ReadString('\n')
	query = strings.TrimSpace(query)
	if query == "" {
		fmt.Println("Запрос не может быть пустым!")
		return
	}
	searchURL := "https://www.youtube.com/results?search_query=" + strings.ReplaceAll(query, " ", "+")
	fmt.Println("Открываю YouTube...")
	err := openBrowser(searchURL)
	if err != nil {
		fmt.Println("Ошибка при открытии браузера:", err)
	}
}

func showStatistics(db *gorm.DB) {
	var tracks []Media
	db.Find(&tracks)
	trackCount := len(tracks)
	artistCounts := make(map[string]int)
	for _, track := range tracks {
		artistCounts[track.Artist]++
	}
	fmt.Println("\nСтатистика медиатеки:")
	fmt.Printf("Общее количество треков: %d\n", trackCount)
	fmt.Println("\nКоличество треков по исполнителям:")
	for artist, count := range artistCounts {
		fmt.Printf("%s: %d треков\n", artist, count)
	}
	maxCount := 0
	popularArtist := ""
	for artist, count := range artistCounts {
		if count > maxCount {
			maxCount = count
			popularArtist = artist
		}
	}
	if popularArtist != "" {
		fmt.Printf("\nСамый популярный исполнитель: %s (%d треков)\n", popularArtist, maxCount)
	} else {
		fmt.Println("\nНет данных о самом популярном исполнителе.")
	}
}

func searching(db *gorm.DB) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Введите название трека или имя исполнителя для поиска:")
	query, _ := reader.ReadString('\n')
	query = strings.TrimSpace(query)
	var tracks []Media
	result := db.Where("artist ILIKE ? OR track ILIKE ?", "%"+query+"%", "%"+query+"%").Find(&tracks)
	if result.Error != nil {
		fmt.Println("Ошибка при поиске треков:", result.Error)
		return
	}
	if len(tracks) == 0 {
		fmt.Println("Треки не найдены.")
		return
	}
	fmt.Println("\nНайденные треки:")
	fmt.Println("ID | Исполнитель | Название трека | Ссылка на клип")
	for _, track := range tracks {
		fmt.Printf("%d: %s - %s  |  %s\n", track.TrackID, track.Artist, track.Track, track.URL)
	}
}

func gettingInfo() {
	url := "https://music.yandex.ru/artist/1426524"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Ошибка при выполнении запроса:", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Ошибка: статус код %d", resp.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal("Ошибка при парсинге HTML:", err)
	}
	trackName := doc.Find("h1.page-artist__title").Text()
	fmt.Println(trackName)
}

func getPlaylists(c *gin.Context, db *gorm.DB) {
	var playlists []Playlist
	db.Find(&playlists)
	c.JSON(200, playlists)
}

func createPlaylist(c *gin.Context, db *gorm.DB) {
	var playlist Playlist
	if err := c.BindJSON(&playlist); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	db.Create(&playlist)
	c.JSON(201, playlist)
}

func deletePlaylist(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	db.Delete(&Playlist{}, id)
	c.JSON(200, gin.H{"message": "удалено"})
}

func renamePlaylist(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var body struct {
		Name string `json:"name"`
	}
	c.BindJSON(&body)
	db.Model(&Playlist{}).Where("playlist_id = ?", id).Update("name", body.Name)
	c.JSON(200, gin.H{"message": "переименовано"})
}

func addTrackToPlaylist(c *gin.Context, db *gorm.DB) {
	playlistID := c.Param("id")
	trackID := c.Param("track_id")
	var playlist Playlist
	var track Media
	db.First(&playlist, playlistID)
	db.First(&track, trackID)
	db.Model(&playlist).Association("Tracks").Append(&track)
	c.JSON(200, gin.H{"message": "трек добавлен"})
}

func removeTrackFromPlaylist(c *gin.Context, db *gorm.DB) {
	playlistID := c.Param("id")
	trackID := c.Param("track_id")
	var playlist Playlist
	var track Media
	db.First(&playlist, playlistID)
	db.First(&track, trackID)
	db.Model(&playlist).Association("Tracks").Delete(&track)
	c.JSON(200, gin.H{"message": "трек удалён из плейлиста"})
}

func getPlaylistTracks(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var playlist Playlist
	db.Preload("Tracks").First(&playlist, id)
	c.JSON(200, playlist.Tracks)
}

func filterTracks(c *gin.Context, db *gorm.DB) {
	from := c.Query("from")
	to := c.Query("to")
	var tracks []Media
	query := db.Model(&Media{})
	if from != "" {
		query = query.Where("created_at >= ?", from)
	}
	if to != "" {
		query = query.Where("created_at <= ?", to)
	}
	query.Find(&tracks)
	c.JSON(200, tracks)
}
