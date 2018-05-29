package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/mediocregopher/radix.v2/redis"
)

func redisProxyGet(t *testing.T, key string) string {
	resp, err := http.Get(fmt.Sprintf("%s/%s", "http://localhost:8080/", key))
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
	conn, err := redis.Dial("tcp", "localhost:7001")
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
	conn, err := redis.Dial("tcp", "localhost:7001")
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
	conn, err := redis.Dial("tcp", "localhost:7001")
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

}

func TestCacheExpires(t *testing.T) {

}

func TestCacheLRU(t *testing.T) {

}

func TestCacheFixedKeySize(t *testing.T) {

}

func TestProcessingConcurrent(t *testing.T) {

}
