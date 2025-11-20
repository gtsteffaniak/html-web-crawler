test:
	go test -v ./...

build:
	docker build -t html-web-crawler .

lint:
	go tool golangci-lint run ./...