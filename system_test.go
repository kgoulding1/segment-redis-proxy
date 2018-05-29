package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/mediocregopher/radix.v2/redis"
)

func redisProxyGet(t *testing.T, key string) string {
	resp, err := http.Get(fmt.Sprintf("%s/%s", "http://localhost:8080", key))
	if err != nil {
		t.Errorf("Unable to connect to redisProxy: %s\n", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Unable to read response body: %s\n", err)
	}
	resp.Body.Close()
	return string(body)
}

func TestExpectString(t *testing.T) {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Errorf("Unable to connect to the redis client %s\n", err)
	}
	defer conn.Close()

	conn.Cmd("SET", "foo", "bar")
	body := redisProxyGet(t, "foo")
	if strings.Compare("$3\r\nbar\r\n", body) != 0 {
		t.Fail()
	}
}

func TestExpectEmptyString(t *testing.T) {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Errorf("Unable to connect to the redis client %s\n", err)
	}
	defer conn.Close()

	conn.Cmd("SET", "empty", "")
	body := redisProxyGet(t, "empty")
	if strings.Compare("$0\r\n\r\n", body) != 0 {
		t.Fail()
	}
}

func TestExpectNil(t *testing.T) {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Errorf("Unable to connect to the redis client %s\n", err)
	}
	defer conn.Close()

	conn.Cmd("DEL", "nil")
	body := redisProxyGet(t, "nil")
	if strings.Compare("$-1\r\n", body) != 0 {
		t.Fail()
	}
}

func TestUsesCache(t *testing.T) {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Errorf("Unable to connect to the redis client %s\n", err)
	}
	defer conn.Close()

	conn.Cmd("SET", "change", "one")

	body := redisProxyGet(t, "change")
	if strings.Compare("$3\r\none\r\n", body) != 0 {
		t.Fail()
	}

	conn.Cmd("SET", "change", "two")

	body = redisProxyGet(t, "change")
	if strings.Compare("$3\r\none\r\n", body) != 0 {
		t.Fail()
	}
}

func TestCacheExpires(t *testing.T) {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Errorf("Unable to connect to the redis client %s\n", err)
	}
	defer conn.Close()

	conn.Cmd("SET", "change", "one")

	body := redisProxyGet(t, "change")
	if strings.Compare("$3\r\none\r\n", body) != 0 {
		t.Fail()
	}

	time.Sleep(11 * time.Second)
	conn.Cmd("SET", "change", "two")

	body = redisProxyGet(t, "change")
	if strings.Compare("$3\r\ntwo\r\n", body) != 0 {
		t.Fail()
	}

}

func TestCacheLRUFixedKeySize(t *testing.T) {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Errorf("Unable to connect to the redis client %s\n", err)
	}
	defer conn.Close()

	// The cache size is 9. To test the fixed key size and it's LRU-ness we create
	// 10 new entries in the redis and get them into the cache, then change the
	// value in the redis. Then we get them all again from back to front. 0 should
	// have the new entry, everything else should be the old entry.

	for i := 0; i < 10; i++ {
		conn.Cmd("SET", fmt.Sprint(i), fmt.Sprint(i))
		redisProxyGet(t, fmt.Sprint(i))
		conn.Cmd("SET", fmt.Sprint(i), fmt.Sprintf("%d prime", i))
	}

	var body string
	for i := 9; i > 0; i-- {
		body = redisProxyGet(t, fmt.Sprint(i))
		if strings.Compare(fmt.Sprintf("$1\r\n%d\r\n", i), body) != 0 {
			log.Printf("Expected %s but got %s\n", fmt.Sprintf("$1\r\n%d\r\n", i), body)
			t.Fail()
		}
	}

	body = redisProxyGet(t, fmt.Sprint(0))
	if strings.Compare("$7\r\n0 prime\r\n", body) != 0 {
		log.Printf("Expected %s but got %s\n", "$7\r\n0 prime\r\n", body)
		t.Fail()
	}

}

func worker(t *testing.T, i int) {
	body := redisProxyGet(t, fmt.Sprint(i))
	if strings.Compare(fmt.Sprintf("$1\r\n%d\r\n", i), body) != 0 {
		log.Printf("Expected %s but got %s\n", fmt.Sprintf("$1\r\n%d\r\n", i), body)
		t.Fail()
	}
}

func TestProcessingConcurrent(t *testing.T) {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Errorf("Unable to connect to the redis client %s\n", err)
	}
	defer conn.Close()

	for i := 0; i < 20; i++ {
		conn.Cmd("SET", fmt.Sprint(i), fmt.Sprint(i))
	}

	for i := 0; i < 20; i++ {
		go worker(t, i)
	}

}
