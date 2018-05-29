build:
	go build main.go

run:
	go run main.go -redisAddr localhost:7001 -expiry 10s -capacity 50 -port :8080

test:
	# For this end to end test we stand up a redis, stand up a proxy, and then
	go test 