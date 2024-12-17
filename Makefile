lint-go:
	golangci-lint run --max-issues-per-linter=0 --max-same-issues=0 --sort-results

lint-fix:
	golangci-lint run --fix

lint-imports:
	./make-lint-imports.sh

lint-imports-fix:
	goimports-reviser -company-prefixes "go.farcloser.world" ./...

tidy:
	go mod tidy

up:
	go get -u ./...