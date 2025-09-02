build:
	go build -o out

run:
	./out

dev:
	air

build-run:
	go build -o out && ./out
