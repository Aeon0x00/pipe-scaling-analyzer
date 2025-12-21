run:
	go run ./cmd/app/main.go

docker-up:
	docker compose up -d

docker-down:
	docker compose down
