

all:
	go build ./...

check:
	gometalinter -D gotype ./tx ./ ./sql
