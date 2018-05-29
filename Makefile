build:
	go build main.go

linuxbuild:
	GOARCH=amd64 GOOS=linux make build

docker: 
	docker build -f ./Dockerfile -t redisproxy ./

docker-run:
	docker run --net="host" redisproxy -redisAddr localhost:6379 -expiry 10s -capacity 9 -port :8080

redis:
	docker run --net="host" --name my-redis-container  -d redis

run:
	go run main.go -redisAddr redis:7001 -expiry 10s -capacity 9 -port :8080

test:
	# For this end to end test we stand up a redis, stand up a proxy, and then
	go test 