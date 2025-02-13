package main

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type media struct {
	TrackID uint `gorm:"primary_key"`
	Artist  string
	Track   string
	URL     string
}

func main() {
	dsn := "host=localhost user=postgres password=YOURPASSWORD dbname=mydb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&media{}) // Создание БД, если её ещё нет
	if err != nil {
		log.Fatal("failed to connect to database", err)
	}
	menu(db)
}

func menu(db *gorm.DB) {
	var choice int
	fmt.Println("Выберите пункт:\n1.Вывести все треки\n2.Добавить трек\n3.Удалить трек" +
		"\n4.Изменить трек по ID\n5.Перемешивание треков\n6.Поиск клипа на YouTube" +
		"\n7.Показать статистику\n8.Поиск треков\n")
	fmt.Scan(&choice)

	var null string
	fmt.Scanln(&null) // Очистка буфера

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
		gettingInfo() // Скрытая фича
	default:
		fmt.Println("Неправильный выбор")
		menu(db)
	}
	shutDown(db)
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

	db.Create(&media{Track: track, Artist: artist, URL: url})
}

func del(db *gorm.DB) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Введите ID исполнителя:")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)
	db.Delete(&media{}, id)
}

func output(db *gorm.DB) {
	var tracks []media
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

	db.Model(&media{}).Where("track_id = ?", id).Updates(media{Artist: artist, Track: track, URL: url})
}

func shuffleAndOutput(db *gorm.DB) {
	var tracks []media
	db.Find(&tracks)

	// Перемешивание треков
	rand.Shuffle(len(tracks), func(i, j int) {
		tracks[i], tracks[j] = tracks[j], tracks[i]
	})

	// Вывод перемешанных треков
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
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)

	return exec.Command(cmd, args...).Start()
}

// Поиск трека на ютубе
func playYouTubeClip() {
	fmt.Println("Введите название трека или исполнителя для поиска на YouTube:")
	reader := bufio.NewReader(os.Stdin)
	query, _ := reader.ReadString('\n')
	query = strings.TrimSpace(query)

	if query == "" {
		fmt.Println("Запрос не может быть пустым!")
		return
	}

	// Формируем URL для поиска на YouTube
	searchURL := "https://www.youtube.com/results?search_query=" + strings.ReplaceAll(query, " ", "+")

	fmt.Println("Открываю YouTube...")
	err := openBrowser(searchURL)
	if err != nil {
		fmt.Println("Ошибка при открытии браузера:", err)
	}
}

func showStatistics(db *gorm.DB) {
	var tracks []media
	db.Find(&tracks) // Получаем все треки из базы данных

	// Подсчет статистики
	trackCount := len(tracks)            // Общее количество треков
	artistCounts := make(map[string]int) // Мапа для подсчета треков по исполнителям

	for _, track := range tracks {
		artistCounts[track.Artist]++
	}

	// Вывод общей статистики
	fmt.Println("\nСтатистика медиатеки:")
	fmt.Printf("Общее количество треков: %d\n", trackCount)

	// Вывод количества треков по каждому исполнителю
	fmt.Println("\nКоличество треков по исполнителям:")
	for artist, count := range artistCounts {
		fmt.Printf("%s: %d треков\n", artist, count)
	}

	// Поиск самого популярного исполнителя
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

	// Ищем треки, где название трека или имя исполнителя содержат запрос
	var tracks []media
	result := db.Where("artist ILIKE ? OR track ILIKE ?", "%"+query+"%", "%"+query+"%").Find(&tracks)

	// Проверяем, есть ли ошибки
	if result.Error != nil {
		fmt.Println("Ошибка при поиске треков:", result.Error)
		return
	}

	// Проверяем, найдены ли треки
	if len(tracks) == 0 {
		fmt.Println("Треки не найдены.")
		return
	}

	// Выводим найденные треки
	fmt.Println("\nНайденные треки:")
	fmt.Println("ID | Исполнитель | Название трека | Ссылка на клип")
	for _, track := range tracks {
		fmt.Printf("%d: %s - %s  |  %s\n", track.TrackID, track.Artist, track.Track, track.URL)
	}
}

// Парсинг
func gettingInfo() {
	url := "https://music.yandex.ru/artist/1426524"
	// Выполняем HTTP-запрос
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Ошибка при выполнении запроса:", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Ошибка: статус код %d", resp.StatusCode)
	}

	// Парсим HTML-документ
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal("Ошибка при парсинге HTML:", err)
	}
	// Извлекаем название трека
	trackName := doc.Find("h1.page-artist__title").Text()
	fmt.Println(trackName)
}

func shutDown(db *gorm.DB) {
	var choice string
	fmt.Println("Вы хотите продолжить? y/n")
	fmt.Scan(&choice)
	if choice == "y" {
		menu(db)
	} else {
		fmt.Println("Хорошего дня!")
		return
	}
}
