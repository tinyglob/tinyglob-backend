run:
	go run main.go

run prod:
	railway run go run main.go

build:
	go build -o bin/main main.go