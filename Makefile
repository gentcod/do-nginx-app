run:
	go run .

build:
	go build -o /bin/do-nginx .

docker-build:
	docker build -t do-nginx .

.PHONY: run build docker-build