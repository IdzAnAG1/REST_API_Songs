package service

import (
	"REST_API_Songs/internal/repository"
	"REST_API_Songs/internal/structure"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strings"
)

type Service struct {
	rep            *repository.Repository
	ExternalApiURL string
}

func NewService(rep *repository.Repository, ExternalApiUrl string) *Service {
	return &Service{
		rep:            rep,
		ExternalApiURL: ExternalApiUrl,
	}
}

func (s *Service) GetSongs(ctx context.Context, filters map[string]string, page, limit int) ([]structure.Song, error) {
	return s.rep.GetSongs(ctx, filters, page, limit)
}

func (s *Service) GetSongText(ctx context.Context, id int64, page, limit int) ([]string, error) {
	logrus.Debugf("Fetching song text for ID: %d, page: %d, limit: %d", id, page, limit)
	song, err := s.rep.GetSongById(ctx, id)
	if err != nil {
		return nil, err
	}
	if song == nil {
		return nil, fmt.Errorf("Song not found ")
	}
	lyrics := strings.Split(song.SongText, "\n\n")
	start := (page - 1) * limit
	end := start + limit
	if start > len(lyrics) {
		return []string{}, err
	}
	if end > len(lyrics) {
		end = len(lyrics)
	}
	logrus.Info("Retrieved %d verses for song ID %d", len(lyrics[start:end]), id)
	return lyrics[start:end], nil
}

func (s *Service) DeleteSong(ctx context.Context, id int64) error {
	return s.rep.DeleteSong(ctx, id)
}

func (s *Service) UpdateSong(ctx context.Context, id int64, song *structure.Song) error {
	return s.rep.UpdateSong(ctx, id, song)
}

func (s *Service) CreateSong(ctx context.Context, in *structure.NewSong) (int64, error) {
	logrus.Debugf("Creating song: %s by %s", in.SongTitle, in.GroupTitle)
	details, err := s.fetchSongDetails(in.GroupTitle, in.SongTitle)
	if err != nil {
		logrus.Errorf("Failed to fetch song details from external API: %v", err)
		return 0, err
	}
	song := &structure.Song{
		GroupTitle:  in.GroupTitle,
		SongTitle:   in.SongTitle,
		ReleaseDate: details.ReleaseDate,
		SongText:    details.SongText,
		Link:        details.Link,
	}
	id, err := s.rep.CreateSong(ctx, song)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Service) fetchSongDetails(groupTitle, songTitle string) (*structure.SupplementForSong, error) {
	logrus.Debugf("Fetching details for song: %s by %s from external API", songTitle, groupTitle)
	u, err := url.Parse(s.ExternalApiURL + "/info/")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("group_title", groupTitle)
	q.Set("song_title", songTitle)
	u.RawQuery = q.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external API returned status: %d", resp.StatusCode)
	}
	var details structure.SupplementForSong
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return nil, err
	}
	logrus.Infof("Fetched details for song: %s by %s", songTitle, groupTitle)
	return &details, nil
}
