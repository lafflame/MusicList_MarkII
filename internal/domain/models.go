package domain

import (
	"time"
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
