package bbolt

import (
	"dictionaryBot/log"
	repo "dictionaryBot/repository"
	bolt "go.etcd.io/bbolt"
	"strconv"
)

type WordsStorage struct {
	db *bolt.DB
}

func NewWordsStorage(db *bolt.DB) *WordsStorage {
	return &WordsStorage{db: db}
}

func (s *WordsStorage) Save(userId int64, word string, bucket repo.Bucket) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Put(intToBytes(userId), []byte(word))
	})
}

func (s *WordsStorage) Get(userId int64, bucket repo.Bucket) (string, error) {
	var word string

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		word = string(b.Get(intToBytes(userId)))
		return nil
	})

	if word == "" {
		log.Error("Words not found")
	}

	return word, err
}

func intToBytes(v int64) []byte {
	return []byte(strconv.FormatInt(v, 10))
}
