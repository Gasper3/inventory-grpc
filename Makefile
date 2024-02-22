proto-gen:
	protoc --proto_path=proto proto/*.proto --go_out=. --go-grpc_out=.

server:
	go run ./cmd/server/server.go

client:
	go run ./cmd/client/client.go

compose-up:
	docker compose up -d

compose-down:
	docker compose down
