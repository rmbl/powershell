GOFMT=gofmt -s -tabs=false -tabwidth=4

GOFILES=$(wildcard *.go **/*.go)

all:
	go build powershell.go

format:
	${GOFMT} -w ${GOFILES}
