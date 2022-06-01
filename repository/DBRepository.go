package repository

import (
	bolt "go.etcd.io/bbolt"
)

type Bucket string

const (
	UserWords Bucket = "user_words"
)

type WordsStorage interface {
	Save(userId int64, word string, bucket Bucket) error
	Get(userId int64, bucket Bucket) (string, error)
}

func InitBolt() (*bolt.DB, error) {
	db, err := bolt.Open("bot.db", 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Batch(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(UserWords))
		if err != nil {
			return err
		}

		return err
	}); err != nil {
		return nil, err
	}

	return db, nil
}
