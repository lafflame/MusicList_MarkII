package main

import (
	"fmt"
	"log"
	"time"

	_ "MusicList_MarkII/docs"
	"MusicList_MarkII/internal/config"
	"MusicList_MarkII/internal/domain"
	"MusicList_MarkII/internal/handler"
	"MusicList_MarkII/internal/repository"
	"MusicList_MarkII/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load()
	cfg := config.Load()
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable client_encoding=UTF8",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)
	db := connectDB(dsn)

	if err := db.AutoMigrate(&domain.Media{}, &domain.Playlist{}); err != nil {
		log.Fatal("migration error:", err)
	}

	mediaRepo := repository.NewMediaRepo(db)
	playlistRepo := repository.NewPlaylistRepo(db)

	mediaService := service.NewMediaService(mediaRepo)
	playlistService := service.NewPlaylistService(playlistRepo)

	mediaHandler := handler.NewMediaHandler(mediaService)
	playlistHandler := handler.NewPlaylistHandler(playlistService)

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

	mediaHandler.Register(r)
	playlistHandler.Register(r)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	fmt.Println("Swagger UI: http://localhost:8080/swagger/index.html")
	r.Run(":8080")
}

func connectDB(dsn string) *gorm.DB {
	var db *gorm.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			return db
		}
		log.Printf("waiting for database... (%d/10)", i+1)
		time.Sleep(2 * time.Second)
	}
	log.Fatal("failed to connect to database:", err)
	return nil
}
