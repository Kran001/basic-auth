PATH_BUILD = build/
PATH_BUILD_BIN = build/bin/
PATH_BUILD_TEST = build_test/

.PHONY: all clean build install uninstall dependencies
all: dependencies clean build

dependencies:
	go mod download

build:
	mkdir -p $(PATH_BUILD)
	mkdir -p $(PATH_BUILD_BIN)
	go build -v ./cmd/basic-auth && mv basic-auth $(PATH_BUILD_BIN)
	cp -r configs $(PATH_BUILD)
	cp -r scripts/sql $(PATH_BUILD)

clean:
	rm -rf $(PATH_BUILD)