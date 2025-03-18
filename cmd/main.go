package main

import (
	"REST_API_Songs/internal/config"
	"REST_API_Songs/internal/db"
	"REST_API_Songs/internal/repository"
	"REST_API_Songs/internal/structure"
	"context"
	"log"
	"time"
)

import (
	"fmt"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Sprintf("Ай ай")
	}
	dab, err := db.NewDataBase(cfg.DataBaseURL)
	if err != nil {
		fmt.Errorf("Провалилось подклюение к БД")
	}
	defer dab.CloseDatabase()

	repo := repository.NewRepository(dab)

	song := &structure.Song{
		GroupTitle:  "Coldplay",
		SongTitle:   "Yellow",
		ReleaseDate: time.Date(2000, 7, 26, 0, 0, 0, 0, time.UTC),
		SongText:    "Look at the stars",
		Link:        "http://example.com/yellow",
	}
	id, err := repo.CreateSong(context.Background(), song)
	if err != nil {
		log.Fatalf("Failed to create song: %v", err)
	}
	fmt.Printf("Created song with ID: %d\n", id)
	id1, err1 := repo.CreateSong(context.Background(), song)
	if err1 != nil {
		log.Fatalf("Failed to create song: %v", err1)
	}
	fmt.Printf("Created song with ID: %d\n", id1)

	filters := map[string]string{
		"group_title": "Coldplay",
	}
	page, limit := 1, 10
	songs, err := repo.GetSongs(context.Background(), filters, page, limit)
	if err != nil {
		log.Fatalf("Failed to fetch songs: %v", err)
	}
	fmt.Printf("Retrieved songs: %+v\n", songs)

	songByID, err := repo.GetSongById(context.Background(), id)
	if err != nil {
		log.Fatalf("Failed to fetch song by ID: %v", err)
	}
	fmt.Printf("Retrieved song by ID: %+v\n", songByID)

	updatedSong := &structure.Song{
		GroupTitle:  "Coldplay",
		SongTitle:   "Updated Yellow",
		ReleaseDate: time.Date(2000, 7, 26, 0, 0, 0, 0, time.UTC),
		SongText:    "Updated text",
		Link:        "http://example.com/updated-yellow",
	}
	err = repo.UpdateSong(context.Background(), id, updatedSong)
	if err != nil {
		log.Fatalf("Failed to update song: %v", err)
	}
	fmt.Println("Updated song successfully")

	err = repo.DeleteSong(context.Background(), id)
	if err != nil {
		log.Fatalf("Failed to delete song: %v", err)
	}
	fmt.Println("Deleted song successfully")
}
