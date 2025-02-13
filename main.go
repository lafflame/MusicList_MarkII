package main

import (
	"bufio"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"strings"
)

type media struct {
	TrackID uint `gorm:"primary_key"`
	Artist  string
	Track   string
	URL     string
}

func main() {
	dsn := "host=localhost user=postgres password= dbname=mydb port=5432 sslmode=disable"
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
		"\n4.Изменить трек по ID")
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
	case 10:
		fmt.Println("Хорошего дня!")
		return
	default:
		fmt.Println("Неправильный выбор")
		menu(db)
	}
	shut(db)
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

// TODO Исправить вывод, чтобы он был адекватный, А не в виде массива
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

func shut(db *gorm.DB) {
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
