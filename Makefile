
start:
	@go run ./cmd/worker

deploy:
	@go run ./cmd/deploy

upgrade_deps:
	@GONOPROXY=buf.build/gen/go/dgroux/newsteam go get -u ./...
	@go mod tidy