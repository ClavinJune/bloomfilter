lint:
	@go vet ./... & go fmt ./... && goimports -w .
test:
	@mkdir -p out/ && go test ./... -race -coverprofile out/coverage.out