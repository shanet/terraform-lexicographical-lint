BIN = bin/terraform-lexicographical-lint
SRC = *.go

.PHONY: all lint clean

all:
	go get -d ./...
	go build -o $(BIN) $(SRC)

fmt: all
	go fmt $(SRC)

clean:
	rm $(BIN)
