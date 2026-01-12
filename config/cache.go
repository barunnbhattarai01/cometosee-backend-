package config

import (
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	c                      *cache.Cache
	DefaultCacheExpiration = 5 * time.Minute
)

func InitCache() {
	//so here 5 min expiration,10 min cleanup interval
	c = cache.New(5*time.Minute, 10*time.Minute)
}

// setcache stores an items in the cahce
// save value with key,think as label
func SetCache(key string, value interface{}) {
	c.Set(key, value, DefaultCacheExpiration)
}

// getcache retrives an data from the cahcee
func GetCache(key string) (interface{}, bool) {
	return c.Get(key)
}

// removecachekey deletes a single cache key
func RemoveCacheKey(key string) {
	c.Delete(key)
}

// removebyprefix delets all keys starting with given prefix
// ex:there is teacher_page and studnt_page so you update student_page ,so it only remove studentpage and uppdate it
// "student_page_1_limit_10" this is k and student_page is prefix
func RemoveByPrefix(prefix string) {
	for k := range c.Items() {
		if strings.HasPrefix(k, prefix) {
			c.Delete(k)
		}
	}
}

// flushcache clears all items from the cache
func FlushCache() {
	c.Flush()
}
