build:
	go build -o bin/server cmd/server/main.go
	go build -o bin/client cmd/client/main.go

release:
	go build -ldflags="-s -w" -o bin/server cmd/server/main.go
	go build -ldflags="-s -w" -o bin/client cmd/client/main.go
	cargo build --release

benchmark: release
	./scripts/go-benchmark.sh
	./scripts/rust-benchmark.sh