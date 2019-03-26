build: vendor bin
	go build -o bin/longrun .

bin:
	mkdir -p bin

vendor:
	dep ensure

test: vendor
	go test -v github.com/factorysh/go-longrun/longrun
	go test -v github.com/factorysh/go-longrun/run

clean:
	rm -rf vendor bin