.DEFAULT_GOAL := test

.PHONY: test
test: clean
	go test

.PHONY: benchmark
benchmark:
	go test -bench=.

.PHONY: coverage
coverage: clean
	go test -cover -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

.PHONY: clean
clean:
	$(RM) *.log *.log.gz *.out *.html

.PHONY: lint
lint:
	golangci-lint run -v
