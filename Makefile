BIN := bin/medea

PHONY := all
all: build

PHONY += build
build:
	go build -ldflags="-s -w" -o $(BIN) src/main.go

PHONY += install
install: $(BIN)
	cp -f $(BIN) /usr/local/bin/

PHONY += clean
clean:
	rm -f $(BIN)

.PHONY: $(PHONY)
