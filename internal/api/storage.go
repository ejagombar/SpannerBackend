package api

import (
	bolt "go.etcd.io/bbolt"
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
	db *bolt.DB
}

func NewSpannerStorage(db *bolt.DB) *SpannerStorage {
	return &SpannerStorage{db: db}
}

func (s *SpannerStorage) SaveToken(token Token) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		var err error

		b := tx.Bucket([]byte("SpotifyToken"))
		if b == nil {
			b, err = tx.CreateBucket([]byte("SpotifyToken"))
			if err != nil {
				return err
			}
		}
		err = b.Put([]byte("AccessToken"), []byte(token.Access))
		err = b.Put([]byte("RefreshToken"), []byte(token.Refresh))
		err = b.Put([]byte("TokenExpiry"), []byte(token.Expiry))
		return err
	})

	return err
}

func (s *SpannerStorage) GetToken() (token Token, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("SpotifyToken"))
		if b == nil {
			return nil
		}
		token.Access = string(b.Get([]byte("AccessToken")))
		token.Refresh = string(b.Get([]byte("RefreshToken")))
		token.Expiry = string(b.Get([]byte("TokenExpiry")))
		return nil
	})
	return token, err
}
