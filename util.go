package fstask

import (
	"crypto/md5"
	"encoding/hex"
	"sync"
	"time"
)

var (
	debounceKey string
	mutex       sync.Mutex
)

//MD5 encryption
func MD5(data []byte) string {
	crypto := md5.New()
	crypto.Write(data)
	return hex.EncodeToString(crypto.Sum(nil))
}

// Include Is in the array
func Include(name string, list []string) bool {
	for _, v := range list {
		if name == v {
			return true
		}
	}
	return false
}

// Debounce prevents jitter and ensures that it is executed only once in a given period of time
func Debounce(key string, fn func(), delay time.Duration) {
	mutex.Lock()
	debounceKey = key
	mutex.Unlock()

	go func() {
		time.Sleep(delay)
		if key == debounceKey {
			fn()
		}
	}()
}
