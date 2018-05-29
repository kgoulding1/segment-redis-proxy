FROM golang:latest 
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 


run go get "github.com/kgoulding1/segment-redis-proxy/redisproxy"
run go get "github.com/mediocregopher/radix.v2/..."
run go get "github.com/karlseguin/ccache"

RUN go build main.go 

EXPOSE 8080


ENTRYPOINT ["/app/main"]

