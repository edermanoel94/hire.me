build:
	docker-compose build --no-cache

run:
	docker-compose up -d

down:
	docker-compose down

test:
	docker-compose run api go test -v ./...