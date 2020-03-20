.PHONY: clean build

clean:
	rm -rf ./build

build:
	GOOS=linux GOARCH=amd64 go build -o build/main ./cmd/otp-lambda/main.go
