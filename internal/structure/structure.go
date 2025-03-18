package structure

import "time"

type Sound struct {
	Id          int       `json:"id"`
	GroupTitle  string    `json:"group_title"`
	SongTitle   string    `json:"song_title"`
	ReleaseDate time.Time `json:"release_date"`
	SongText    string    `json:"song_text"`
	Link        string    `json:"link"`
}

type NewSong struct {
	GroupTitle string `json:"group_title"`
	SongTitle  string `json:"song_title"`
}

type SupplementForSong struct {
	ReleaseDate time.Time `json:"release_date"`
	SongText    string    `json:"song_text"`
	Link        string    `json:"link"`
}
