run:
	go run main.go

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

test:
	go test -v ./...

test-cover-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

tailwind-build:
	tailwindcss -i ./static/css/styles.css -o ./static/css/output.css --minify

tailwind-watch:
	tailwindcss -i ./static/css/styles.css -o ./static/css/output.css --watch

.PHONY: run lint lint-fix test test-cover-html tailwind-build tailwind-watch
