test:
	go test -timeout 30s -race -v ./...

build:
	docker build -t html-web-crawler .

lint:
	go tool golangci-lint run ./...