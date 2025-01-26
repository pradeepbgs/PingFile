echo "Building binaries..."

GOOS=windows GOARCH=amd64 go build -o bin/pingfile-windows.exe main.go

GOOS=darwin GOARCH=amd64 go build -o bin/pingfile-macos main.go

GOOS=darwin GOARCH=arm64 go build -o bin/pingfile-macos-arm main.go

GOOS=linux GOARCH=amd64 go build -o bin/pingfile-linux main.go

GOOS=linux GOARCH=arm64 go build -o bin/pingfile-linux-arm main.go

echo "Build complete! Binaries are in the bin directory."
