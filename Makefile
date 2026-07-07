BINARY := taskrunner
PKG := ./cmd/taskrunner

.PHONY: build test run lint fmt clean

build:
	go build -o $(BINARY) $(PKG)

test:
	go test ./...

run: build
	./$(BINARY) -file tasks.json -workers 3

lint:
	go vet ./...
	@test -z "$$(gofmt -l .)" || { echo "fichiers mal formatés:"; gofmt -l .; exit 1; }

fmt:
	gofmt -w .

clean:
	rm -f $(BINARY)
