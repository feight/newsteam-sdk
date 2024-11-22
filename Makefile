
start:
	@go run ./cmd/worker

deploy:
	@go run github.com/feight/deploy@v1.0.4

upgrade_deps:
	@GONOPROXY=buf.build/gen/go/dgroux/newsteam go get -u ./...
	@go mod tidy