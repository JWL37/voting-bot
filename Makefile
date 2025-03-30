all:
	docker compose up -d
	go run ./cmd/voting-bot/main.go	
mattermost:
	docker run --name mattermost-preview -d --publish 8065:8065 mattermost/mattermost-preview
lint:
	golangci-lint run 