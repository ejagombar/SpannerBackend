package api

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type spannerDB struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenExpiry  string `json:"token_expiry"`
}

type Token struct {
	Access  string
	Refresh string
	Expiry  string
}

type SpannerStorage struct {
	db *sql.DB
}

func NewSpannerStorage(db *sql.DB) *SpannerStorage {
	return &SpannerStorage{db: db}
}

// func (s *SpannerStorage) NewSQLStorage(dataSourceName string) (*SpannerStorage, error) {
// 	db, err := sql.Open("mysql", dataSourceName)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	// Create the necessary table if not exists
// 	_, err = db.Exec(`
// 		CREATE TABLE IF NOT EXISTS SpotifyToken (
// 			AccessToken VARCHAR(255) PRIMARY KEY,
// 			RefreshToken VARCHAR(255),
// 			TokenExpiry VARCHAR(255)
// 		);
// 	`)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &SpannerStorage{db: db}, nil
// }

func (s *SpannerStorage) SaveToken(token Token) error {
	_, err := s.db.Exec(`
		INSERT INTO SpotifyToken (AccessToken, RefreshToken, TokenExpiry)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE
			AccessToken = VALUES(AccessToken),
			RefreshToken = VALUES(RefreshToken),
			TokenExpiry = VALUES(TokenExpiry);
	`, token.Access, token.Refresh, token.Expiry)
	return err
}

func (s *SpannerStorage) GetToken() (token Token, err error) {
	row := s.db.QueryRow(`
		SELECT AccessToken, RefreshToken, TokenExpiry
		FROM SpotifyToken
		LIMIT 1;
	`)

	err = row.Scan(&token.Access, &token.Refresh, &token.Expiry)
	if err == sql.ErrNoRows {
		return Token{}, nil
	}

	return token, err
}
