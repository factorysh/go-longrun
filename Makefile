build: bin
	go build -o bin/longrun .

bin:
	mkdir -p bin

test:
	go test -v github.com/factorysh/go-longrun/longrun
	go test -v github.com/factorysh/go-longrun/longrun/sse
	go test -v github.com/factorysh/go-longrun/run

clean:
	rm -rf bin