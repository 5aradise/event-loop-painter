build:
	go build -C cmd/painter/ -o ../../bin/

run: build
	./bin/painter
