package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	proxy "github.com/kgoulding1/segment-redis-proxy/redisproxy"
	"github.com/mediocregopher/radix.v2/pool"
)

// Command-line flags.
var (
	redisAddr  = flag.String("redisAddr", "localhost:7001", "Address of the backing Redis")
	expiryTime = flag.Duration("expiry", 60*time.Second, "Cache expiry time")
	capacity   = flag.Int("capacity", 5000, "Capacity (number of keys)")
	port       = flag.String("port", ":8080", "TCP/IP port number the proxy listens on")
)

func main() {
	flag.Parse()
	fmt.Printf("%s, %v, %d, %s\n", *redisAddr, *expiryTime, *capacity, *port)

	connPool, err := pool.New("tcp", *redisAddr, 10)
	if err != nil {
		log.Fatalf("Could not connect: %v\n", err)
	}
	defer connPool.Empty()

	http.Handle("/", proxy.NewGetServer(*connPool, *capacity, *expiryTime))
	log.Fatal(http.ListenAndServe(*port, nil))
}
