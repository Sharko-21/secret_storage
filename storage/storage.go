package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"

	"main/cipher"
)

const maxCounter = 3
const maxLifetime = 3 * 5 * time.Second //3 day

type Storage interface {
	Get(id string, key string) (secret string)
	Save(secret string, key string) (id string)
}

type storage struct {
	sync.Map
}

type item struct {
	counter   int8
	createdAt time.Time
	value     []byte
}

var secretStorage storage

func (s *storage) Get(id string, key string) (secret string) {
	elem, ok := s.Load(id)
	if !ok {
		return ""
	}
	i, ok := elem.(*item)
	if !ok {
		return ""
	}
	if i.counter >= maxCounter {
		return ""
	}
	if time.Since(i.createdAt) > maxLifetime {
		return ""
	}
	i.counter++
	if i.counter >= maxCounter {
		s.Delete(id)
	}
	return string(cipher.Decrypt(i.value, []byte(key)))
}

func (s *storage) Save(secret string, key string) (id string) {
	v := sha256.Sum256([]byte(secret))
	id = hex.EncodeToString(v[:])
	s.Store(id, &item{
		createdAt: time.Now(),
		value:     cipher.Encrypt([]byte(secret), []byte(key)),
	})
	return id
}

func (s *storage) Clean() {
	s.Range(func(key, value any) bool {
		v, ok := value.(*item)
		if !ok {
			s.Delete(key)
		}
		if time.Since(v.createdAt) > maxLifetime {
			s.Delete(key)
		}
		return true
	})
}

var once sync.Once

func GetStorage() Storage {
	once.Do(func() {
		secretStorage = storage{}
		go secretStorage.Clean()
	})
	return &secretStorage
}
