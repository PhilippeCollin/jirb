rm -rf ./build
env GOOS=darwin GOARCH=amd64 go build -ldflags "-w" -o ./build/jirb-darwin-amd64 -v
env GOOS=linux GOARCH=amd64 go build -ldflags "-w" -o ./build/jirb-linux-amd64 -v
env GOOS=windows GOARCH=amd64 go build -ldflags "-w" -o ./build/jirb-windows-amd64.exe -v
