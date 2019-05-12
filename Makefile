test:
	@echo > coverage.txt
	go test -race -coverprofile=coverage.txt -covermode=atomic 

.PHONY: test
