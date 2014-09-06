.PHONY: get clean

nw: nw.go
	go build -o nw

get:
	go get -v ./...

clean:
	rm -f nw
