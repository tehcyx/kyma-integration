default: build

build: test cover
	go build -i -o bin/app

test:
	go test ./...

cover:
	go test ./... -cover

clean:
	rm -rf bin