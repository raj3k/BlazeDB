build:
	@go build -o bin/blazedb cmd/main.go

run:	build
	@./bin/blazedb

test:
	@go test -v ./...

clean:
	@rm -rf ./bin