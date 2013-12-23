GOPATH = $(shell pwd)

all:
	GOPATH=$(GOPATH) go build -o riddick src/main.go

depends:
	GOPATH=$(GOPATH) go get github.com/codegangsta/martini
	GOPATH=$(GOPATH) go get github.com/nu7hatch/gouuid

run:
	GOPATH=$(GOPATH) ./docdb
