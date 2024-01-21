package api 

import (
	bolt "go.etcd.io/bbolt"
)

type spannerDB struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenExpiry  string `json:"token_expiry"`
}

type SpannerStorage struct {
	db *bolt.DB
}

func NewSpannerStorage(db *bolt.DB) *SpannerStorage {
	return &SpannerStorage{db: db}
}

func (s *SpannerStorage) SaveToken(AccessToken, RefreshToken, TokenExpiry string) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		var err error

		b := tx.Bucket([]byte("SpotifyToken"))
		if b == nil {
			b, err = tx.CreateBucket([]byte("SpotifyToken"))
			if err != nil {
				return err
			}
		}
		err = b.Put([]byte("AccessToken"), []byte(AccessToken))
		err = b.Put([]byte("RefreshToken"), []byte(RefreshToken))
		err = b.Put([]byte("TokenExpiry"), []byte(TokenExpiry))
		return err
	})

	return err
}

func (s *SpannerStorage) GetToken() (AccessToken, RefreshToken, TokenExpiry string, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("SpotifyToken"))
		if b == nil {
			return nil
		}
		AccessToken = string(b.Get([]byte("AccessToken")))
		RefreshToken = string(b.Get([]byte("RefreshToken")))
		TokenExpiry = string(b.Get([]byte("TokenExpiry")))
		return nil
	})
	return
}
