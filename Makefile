default: build

NAME := gosea
VERSION := v0.0.4

clean:
	rm -rf ./build

build: clean build-darwin build-linux build-windows

build-darwin:
	@echo Version: $(VERSION) $(BUILDDATE)
	CGO_ENABLED=0 \
	GOARCH=amd64 \
	GOOS=darwin \
	go build -ldflags "-X main.Version=${VERSION}" -trimpath -o build/darwin-amd64/${NAME} *.go

build-linux:
	@echo Version: $(VERSION) $(BUILDDATE)
	CGO_ENABLED=0 \
	GOARCH=amd64 \
	GOOS=linux \
	go build -ldflags "-X main.Version=${VERSION}" -trimpath -o build/linux-amd64/${NAME} *.go

build-windows:
	@echo Version: $(VERSION) $(BUILDDATE)
	CGO_ENABLED=0 \
	GOARCH=amd64 \
	GOOS=windows \
	go build -ldflags "-X main.Version=${VERSION}" -trimpath -o build/windows-amd64/${NAME} *.go

build-docker:
	@echo Version: $(VERSION) $(BUILDDATE)
	docker build -t gosea:latest -f Dockerfile ./
