default:
	make gotest

gotest:
	go test ./...

install:
	go install

build:
	go build -o ./out/gsh ./
