package system_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

const (
	redisURI      = "redis:6379"
	rproxyHTTPURI = "http://rproxy:9090/v1/redis"
	rproxyTCPURI  = "rproxy:6379"
)

func redisGet(key, uri string) (string, error) {
	client := redis.NewClient(&redis.Options{
		Addr: uri,
	})

	cmd := client.Get(key)
	return cmd.Result()
}

func redisSet(key, val, uri string) error {
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	_, err := client.Set(key, val, 0).Result()
	return err
}

func redisDelete(key, uri string) error {
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	_, err := client.Del(key).Result()
	return err
}

func rproxyGet(key, uri string) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(uri+"?key=%s", key), nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func TestTrivialGet(t *testing.T) {
	key := "5"
	val := "abc"
	if err := redisSet(key, val, redisURI); err != nil {
		t.Errorf("unexpected err happened: %v", err)
	}

	rpVal, err := rproxyGet(key, rproxyHTTPURI)
	if err != nil {
		t.Errorf("unexpected err happened: %v", err)
	}

	if rpVal != val {
		t.Errorf("Val was %v but we expected %v", rpVal, val)
	}
}

func TestCacheExpiry(t *testing.T) {
	key := "5"
	val := "abc"
	if err := redisSet(key, val, redisURI); err != nil {
		t.Errorf("unexpected err happened: %v", err)
	}

	// first get - establish in cache
	_, err := rproxyGet(key, rproxyHTTPURI)
	if err != nil {
		t.Errorf("unexpected err happened: %v", err)
	}

	// delete in actual redis
	redisDelete(key, redisURI)

	// Wait for  global expiry
	time.Sleep(5 * time.Second)

	// Should be gone from cache and redis
	rpVal, err := rproxyGet(key, rproxyHTTPURI)
	if err != nil {
		t.Errorf("unexpected err happened: %v", err)
	}

	if rpVal != "" {
		t.Errorf("Val was %v but we expected %v", rpVal, "")
	}
}

func TestLRU(t *testing.T) {
	// Assign 25 key/vals in redis.
	for i := 0; i < 25; i++ {
		if err := redisSet(strconv.Itoa(i), strconv.Itoa(i), redisURI); err != nil {
			t.Errorf("unexpected err happened: %v", err)
		}
	}

	// Get all 25 keys, putting them in the cache as we go.
	for i := 0; i < 25; i++ {
		_, err := rproxyGet(strconv.Itoa(i), rproxyHTTPURI)
		if err != nil {
			t.Errorf("unexpected err happened: %v", err)
		}
	}

	// Delete all 25 key/vals in redis, so we're depending on the cache.
	for i := 0; i < 25; i++ {
		if err := redisDelete(strconv.Itoa(i), redisURI); err != nil {
			t.Errorf("unexpected err happened: %v", err)
		}
	}

	// Get all 25 values. Since default capacity is 10, we should only get 10 elements.
	for i := 0; i < 14; i++ {
		val, err := rproxyGet(strconv.Itoa(i), rproxyHTTPURI)
		if err != nil {
			t.Errorf("unexpected err happened: %v", err)
		}

		if val != "" {
			t.Errorf("LRU cache is not limiting capacity")
		}
	}
	for i := 15; i < 25; i++ {
		val, err := rproxyGet(strconv.Itoa(i), rproxyHTTPURI)
		if err != nil {
			t.Errorf("unexpected err happened: %v", err)
		}

		if val != strconv.Itoa(i) {
			t.Errorf("Failed to get cached response")
		}
	}
}

func TestConcurrency(t *testing.T) {
	// set 3 keys 0-2
	for i := 0; i < 3; i++ {
		if err := redisSet(strconv.Itoa(i), strconv.Itoa(i), redisURI); err != nil {
			t.Errorf("unexpected err happened: %v", err)
		}
	}

	// set 3 keys 10-12
	for i := 10; i < 13; i++ {
		if err := redisSet(strconv.Itoa(i), strconv.Itoa(i), redisURI); err != nil {
			t.Errorf("unexpected err happened: %v", err)
		}
	}

	// set 3 keys 20-22
	for i := 20; i < 23; i++ {
		if err := redisSet(strconv.Itoa(i), strconv.Itoa(i), redisURI); err != nil {
			t.Errorf("unexpected err happened: %v", err)
		}
	}

	// Get all values in separate goroutines. Should receive all values back correctly.
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		for i := 0; i < 3; i++ {
			val, err := rproxyGet(strconv.Itoa(i), rproxyHTTPURI)
			if err != nil {
				t.Errorf("unexpected err happened: %v", err)
			}

			if val != strconv.Itoa(i) {
				t.Errorf("LRU cache is not limiting capacity")
			}
		}
		wg.Done()
	}()
	go func() {
		for i := 10; i < 13; i++ {
			val, err := rproxyGet(strconv.Itoa(i), rproxyHTTPURI)
			if err != nil {
				t.Errorf("unexpected err happened: %v", err)
			}

			if val != strconv.Itoa(i) {
				t.Errorf("LRU cache is not limiting capacity")
			}
		}
		wg.Done()
	}()
	go func() {
		for i := 20; i < 23; i++ {
			val, err := rproxyGet(strconv.Itoa(i), rproxyHTTPURI)
			if err != nil {
				t.Errorf("unexpected err happened: %v", err)
			}

			if val != strconv.Itoa(i) {
				t.Errorf("LRU cache is not limiting capacity")
			}
		}
		wg.Done()
	}()

	// Make sure all have finished
	wg.Wait()
}
