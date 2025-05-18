package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type PlaylistAlbum struct {
	ID                   int64
	SpotifyURI           string
	AlbumName            string
	ArtistName           string
	ReleaseDate          string
	ReleaseDatePrecision string
}

type PlaylistRepository struct {
	db *sql.DB
}

func NewPlaylistRepository(dbPath string) (*PlaylistRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS playlist_albums (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			spotify_uri TEXT NOT NULL UNIQUE,
			artist_name TEXT NOT NULL,
			album_name TEXT NOT NULL,
			release_date TEXT NOT NULL,
			release_date_precision TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return nil, err
	}

	return &PlaylistRepository{db: db}, nil
}

func (r *PlaylistRepository) StorePlaylistAlbums(albums []PlaylistAlbum) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO playlist_albums (spotify_uri, artist_name, album_name, release_date, release_date_precision)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(spotify_uri) DO UPDATE SET
			artist_name = excluded.artist_name,
			album_name = excluded.album_name,
			release_date = excluded.release_date,
			release_date_precision = excluded.release_date_precision
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, album := range albums {
		_, err = stmt.Exec(album.SpotifyURI, album.ArtistName, album.AlbumName, album.ReleaseDate, album.ReleaseDatePrecision)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PlaylistRepository) Close() error {
	return r.db.Close()
}
