package proxy

import (
	"log"
	"net/http"
	"time"

	"github.com/karlseguin/ccache"
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
)

// GetServer acts as a server for any incoming GET requests for the redis proxy. It has a cache,
// a pool of connections to the redis, and how long an entry should be kept before it expires.
type GetServer struct {
	clientPool pool.Pool
	cache      ccache.Cache
	expiryTime time.Duration
}

// Create a new Get Server.
func NewGetServer(clientPool pool.Pool, capacity int, expiryTime time.Duration) *GetServer {
	cache := ccache.New(ccache.Configure().MaxSize(int64(capacity)))
	c := &GetServer{clientPool: clientPool, cache: *cache, expiryTime: expiryTime}
	return c
}

// ServeHTTP implements the HTTP user interface.
func (c *GetServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]
	item := c.cache.Get(key)

	if item == nil || item.Expired() {
		conn, err := c.clientPool.Get()
		if err != nil {
			log.Printf("Error connection from the pool: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		// Call to the redis, return the connection to the pool, handle any errors.
		resp := conn.Cmd("GET", key)
		c.clientPool.Put(conn)
		if resp.Err != nil {
			log.Printf("Error getting from the redis: %v\n", resp.Err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		c.cache.Set(key, resp, c.expiryTime)
		resp.WriteTo(w)
	} else {
		item.Value().(*redis.Resp).WriteTo(w)
	}

}
