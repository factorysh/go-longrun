build: vendor bin
	go build -o bin/longrun .

bin:
	mkdir -p bin

vendor:
	dep ensure

clean:
	rm -rf vendor bin