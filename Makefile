run: build
	@./bin/gokv

build:
	@go build -o bin/gokv .

