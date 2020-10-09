
test:
	go test -timeout 30s -cover \
	github.com/factorysh/go-longrun/sse \
	github.com/factorysh/go-longrun/run \
	github.com/factorysh/go-longrun/rest \
	github.com/factorysh/go-longrun/client
