BINARY := gobackup
CONFIG := config.json

.PHONY: all build run clean

all: build

build:
	go build -o $(BINARY) .

run: build
	@if [ ! -f "$(CONFIG)" ]; then \
		echo "Error: $(CONFIG) not found. Copy config.json.example to config.json and fill in your values."; \
		exit 1; \
	fi
	./$(BINARY)

clean:
	rm -f $(BINARY)
