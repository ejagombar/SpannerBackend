package storage

import (
    "time"
    bolt "go.etcd.io/bbolt"
)

func LoadBbolt(name string, timeout time.Duration) (*bolt.DB, error) {
    return bolt.Open(name, 0600, &bolt.Options{Timeout: timeout})
}

func CloseBbolt(db *bolt.DB) error {
    return db.Close()
}
