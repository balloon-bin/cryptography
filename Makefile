BINARY_DIR = bin
BINARIES = $(patsubst cmd/%/,%,$(wildcard cmd/*/))

.PHONY: all build test coverage validate clean purge $(BINARIES)

all: build
	

build: $(BINARIES)
	

$(BINARY_DIR):
	mkdir -p $(BINARY_DIR)

$(BINARIES): %: $(BINARY_DIR)
	go build -o $(BINARY_DIR)/$@ ./cmd/$@/

test:
	go test ./... -cover

coverage:
	mkdir -p reports/
	go test -coverprofile=reports/coverage.out ./... && go tool cover -html=reports/coverage.out

validate:
	@test -z "$(shell gofumpt -l .)" && echo "No files need formatting" || (echo "Incorrect formatting in:"; gofumpt  -l .; exit 1)
	go vet ./...

clean:
	rm -rf $(BINARY_DIR)
	go clean

purge: clean
	rm -rf reports
