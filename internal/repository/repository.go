package repository

import (
	"REST_API_Songs/internal/db"
	"REST_API_Songs/internal/structure"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type Repository struct {
	db *db.DB
}

func NewRepository(db *db.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetSongs(ctx context.Context, filters map[string]string,
	page, limit int) ([]structure.Song, error) {
	logrus.Debugf("fetching songs with filters: %v, page %d, limit %d", filters, page, limit)
	var conditions []string
	var args []interface{}
	argIndex := 1
	for key, value := range filters {
		if value != "" {
			conditions = append(conditions, key+" LIKE $"+strconv.Itoa(argIndex))
			args = append(args, "%"+value+"%")
			argIndex++
		}
	}
	query := "SELECT id, group_title, song_title, release_date, song_text, link FROM songs"
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY id LIMIT $" + strconv.Itoa(argIndex) + " OFFSET $" + strconv.Itoa(argIndex+1)
	args = append(args, limit, (page-1)*limit)

	rows, err := r.db.Cnct.Query(ctx, query, args...)
	if err != nil {
		logrus.Errorf("Failed to query songs: %v", err)
		return nil, err
	}
	defer rows.Close()

	var songs []structure.Song
	for rows.Next() {
		var song structure.Song
		if err := rows.Scan(&song.Id,
			&song.GroupTitle,
			&song.SongTitle,
			&song.ReleaseDate,
			&song.SongText,
			&song.Link); err != nil {
			logrus.Errorf("Failed to scan song: %v", err)
			return nil, err
		}
		songs = append(songs, song)
	}
	logrus.Infof("Retrieved %d songs", len(songs))
	return songs, nil
}

func (r *Repository) GetSongById(ctx context.Context, id int64) (*structure.Song, error) {
	logrus.Debugf("Fetching song with ID: %d", id)
	var song structure.Song
	err := r.db.Cnct.
		QueryRow(ctx, "SELECT id, group_title, song_title, release_date, song_text, link FROM songs WHERE id = $1",
			id).Scan(&song.Id, &song.GroupTitle, &song.SongTitle, &song.ReleaseDate, &song.SongText, &song.Link)
	if errors.Is(err, pgx.ErrNoRows) {
		logrus.Warnf("Song with ID %d not found", id)
		return nil, nil
	}
	if err != nil {
		logrus.Errorf("Failed to fetch song: %v", err)
		return nil, err
	}
	logrus.Infof("Retrieved song with ID %d", id)
	return &song, nil
}

func (r *Repository) DeleteSong(ctx context.Context, id int64) error {
	logrus.Debugf("Deleting song with ID: %d", id)
	result, err := r.db.Cnct.Exec(ctx, "DELETE FROM songs WHERE id = $1", id)
	if err != nil {
		logrus.Errorf("Failed to delete song: %v", err)
		return err
	}
	if result.RowsAffected() == 0 {
		logrus.Warnf("Song with ID %d not found", id)
		return pgx.ErrNoRows
	}
	logrus.Infof("Deleted song with ID %d", id)
	return nil
}

func (r *Repository) UpdateSong(ctx context.Context, id int64, song *structure.Song) error {
	logrus.Debugf("Updating song with ID: %d", id)
	result, err := r.db.Cnct.Exec(ctx,
		"UPDATE songs SET group_title = $1, song_title = $2, release_date = $3, song_text = $4, link = $5 WHERE id = $6",
		song.GroupTitle, song.SongTitle, song.ReleaseDate, song.SongText, song.Link, id)
	if err != nil {
		logrus.Errorf("Failed to update song: %v", err)
		return err
	}
	if result.RowsAffected() == 0 {
		logrus.Warnf("Song with ID %d not found", id)
		return pgx.ErrNoRows
	}
	logrus.Infof("Updated song with ID %d", id)
	return nil
}

func (r *Repository) CreateSong(ctx context.Context, song *structure.Song) (int64, error) {
	logrus.Debugf("Creating new song: %s by %s", song.SongTitle, song.GroupTitle)
	var id int64
	err := r.db.Cnct.QueryRow(ctx,
		"INSERT INTO songs (group_title, song_title, release_date, song_text, link) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		song.GroupTitle, song.SongTitle, song.ReleaseDate, song.SongText, song.Link).Scan(&id)
	if err != nil {
		logrus.Errorf("Failed to create song: %v", err)
		return 0, err
	}
	logrus.Infof("Created song with ID %d", id)
	return id, nil
}
