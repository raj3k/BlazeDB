build:
	@go build -o bin/blazedb cmd/main.go

run:	build
	@./bin/blazedb

clean:
	@rm -rf ./bin