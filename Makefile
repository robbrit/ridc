GOPATH = $(shell pwd)

all:
	GOPATH=$(GOPATH) go build -o ridc src/main.go

depends:
	GOPATH=$(GOPATH) go get github.com/codegangsta/martini
	GOPATH=$(GOPATH) go get github.com/nu7hatch/gouuid
