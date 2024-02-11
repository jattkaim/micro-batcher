USER_TESTS_DIR := tests

test: ## Run unit tests
	@echo "mode: count" > coverage-all.out
	@go test -v -p=1 -cover -coverpkg=./pkg/... -covermode=count -coverprofile=coverage.out ${USER_TESTS_DIR}/*.go
	@tail -n +2 coverage.out >> coverage-all.out
